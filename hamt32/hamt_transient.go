package hamt32

import (
	"fmt"
	"log"

	"github.com/lleo/go-hamt-key"
)

type HamtTransient struct {
	root     tableI
	nentries uint
	grade    bool
	compinit bool
}

func NewTransient(opt int) *HamtTransient {
	var h = new(HamtTransient)

	switch opt {
	case HybridTables:
		h.grade = true
		h.compinit = true
	case CompTablesOnly:
		h.grade = false
		h.compinit = true
	case FullTablesOnly:
		fallthrough
	default:
		h.grade = false
		h.compinit = false
	}

	return h
}

func (h *HamtTransient) IsEmpty() bool {
	return h.root == nil
}

func (h *HamtTransient) Nentries() uint {
	return h.nentries
}

func (h *HamtTransient) find(k key.Key) (tableStack, leafI, uint) {
	if h.IsEmpty() {
		return nil, nil, 0
	}

	var h30 = k.Hash30()
	var curTable = h.root

	var path = newTableStack()
	var leaf leafI
	var idx uint

	var depth uint
DepthIter:
	for depth = 0; depth <= maxDepth; depth++ {
		path.push(curTable)
		idx = h30.Index(depth)

		var curNode = curTable.get(idx)
		switch n := curNode.(type) {
		case nil:
			leaf = nil
			break DepthIter
		case leafI:
			leaf = n
			break DepthIter
		case tableI:
			if depth == maxDepth {
				log.Panicf("SHOULD NOT BE REACHED; depth,%d == maxDepth,%d & tableI entry found; %s", depth, maxDepth, n)
			}
			curTable = n
			// exit switch then loop for
		default:
			log.Panicf("SHOULD NOT BE REACHED: depth=%d; curNode unknown type=%T;", depth, curNode)
		}
	}

	return path, leaf, idx
}

// Get Hamt method looks up a given key in the Hamt data structure.
// BenchHamt32:
//func (h *HamtTransient) Get(k key.Key) (interface{}, bool) {
//	var _, leaf, _ = h.find(k)
//
//	if leaf == nil {
//		return nil, false
//	}
//
//	return leaf.get(k)
//}

func (h *HamtTransient) Get(k key.Key) (interface{}, bool) {
	if h.IsEmpty() {
		return nil, false
	}

	var val interface{}
	var found bool

	var h30 = k.Hash30()

	var curTable = h.root //ISA tableI

	for depth := uint(0); depth <= maxDepth; depth++ {
		var idx = h30.Index(depth)
		var curNode = curTable.get(idx) //nodeI

		if curNode == nil {
			return nil, false
		}

		if leaf, isLeaf := curNode.(leafI); isLeaf {
			val, found = leaf.get(k)
			return val, found
		}

		if depth == maxDepth {
			panic("SHOULD NOT HAPPEN")
		}
		curTable = curNode.(tableI)
	}

	panic("SHOULD NEVER BE REACHED")
}

func (h *HamtTransient) createRootTable(leaf leafI) tableI {
	if h.compinit {
		return createRootCompressedTable(leaf)
	}
	return createRootFullTable(leaf)
}

func (h *HamtTransient) createTable(depth uint, leaf1 leafI, leaf2 *flatLeaf) tableI {
	if h.compinit {
		return createCompressedTable(depth, leaf1, leaf2)
	}
	return createFullTable(depth, leaf1, leaf2)
}

// Put Hamt method inserts a given key/val pair into the Hamt data structure.
// It returns a boolean indicating if the key/val was inserted or whether or
// not the key already existed and the val was merely overwritten.
func (h *HamtTransient) Put(k key.Key, v interface{}) (Hamt, bool) {
	if h.IsEmpty() {
		h.root = h.createRootTable(newFlatLeaf(k, v))
		h.nentries++
		return h, true
	}

	var path, leaf, idx = h.find(k)

	var curTable = path.pop()
	var depth = uint(path.len())
	var added bool

	if leaf == nil {
		//check if upgrading allowed & if it is required
		if h.grade && (curTable.nentries()+1) == upgradeThreshold {
			var newTable tableI
			newTable = upgradeToFullTable(curTable.Hash30(), curTable.entries())
			if curTable == h.root {
				h.root = newTable
			} else {
				var parentTable = path.peek()
				var parentIdx = k.Hash30().Index(depth - 1)
				parentTable.replace(parentIdx, newTable)
			}
			curTable = newTable
		}
		curTable.insert(idx, newFlatLeaf(k, v))
		added = true
	} else {
		// This is the condition that allows collision leafs to exist at a level
		// less than maxDepth. I don't know if I want to allow this...
		if leaf.Hash30() == k.Hash30() {
			var newLeaf leafI
			// There are four possibilities here:
			// if leaf isa collision leaf
			//   k is identical to one of the kv pairs in collision leaf; hence
			//     we replace that ones value and added = false
			//   k is unique in the collision leaf and the kv pair is added;
			//     this is very rare; the underlying key basis is different but
			//     the Hash30 is identical.
			// if leaf isa flat leaf
			//   k is identical to the flat leaf's key; hence the value is
			//     replaced and added == false
			//   k is not identical to the flat leaf's key; and a collision leaf
			//     is created and added == true; again this is very rare; the
			//     underlying key basis is different but the Hash30 is identical
			newLeaf, added = leaf.put(k, v)
			curTable.replace(idx, newLeaf)
		} else {
			var tmpTable = h.createTable(depth+1, leaf, newFlatLeaf(k, v))
			curTable.replace(idx, tmpTable)
			added = true
		}
	}

	if added {
		h.nentries++
	}

	return h, added
}

// Del Hamt Method removes a given key from the Hamt data structure. It returns
// a pair of values: the value stored and a boolean indicating if the key was
// even found and deleted.
func (h *HamtTransient) Del(k key.Key) (Hamt, interface{}, bool) {
	if h.IsEmpty() {
		return h, nil, false
	}

	var path, leaf, idx = h.find(k)

	var curTable = path.pop()
	var depth = uint(path.len())

	var val interface{}
	var deleted bool

	if leaf == nil {
		return h, nil, false
	}

	var newLeaf leafI
	newLeaf, val, deleted = leaf.del(k)

	if !deleted {
		return h, nil, false
	}

	h.nentries--

	if newLeaf != nil { //leaf was a CollisionLeaf
		curTable.replace(idx, newLeaf)
	} else { //leaf was a FlatLeaf
		curTable.remove(idx)

		// Side-Effects of removing an KeyVal from the table
		switch {
		// if no entries left in table need to colapse down to parent
		case curTable != h.root && curTable.nentries() == 0:
			var parentTable = path.peek()
			var parentIdx = k.Hash30().Index(depth - 1)
			parentTable.remove(parentIdx)
			curTable = nil

			// else check if downgrade allowed and required
		case h.grade && curTable.nentries() == downgradeThreshold:
			//when nentries is decr'd it will be <downgradeThreshold
			var newTable = downgradeToCompressedTable(
				curTable.Hash30(), curTable.entries())
			if curTable == h.root { //aka path.len() == 0 or path.peek() == nil
				h.root = newTable
			} else {
				var parentTable = path.peek()
				var parentIdx = k.Hash30().Index(depth - 1)
				parentTable.replace(parentIdx, newTable)
			}
		}
	}

	return h, val, deleted
}

// String returns a string representation of the Hamt string.
func (h *HamtTransient) String() string {
	return fmt.Sprintf("HamtTransient{ nentries: %d, root: %s }", h.nentries, h.root.LongString("", 0))
}

// LongString returns a complete listing of the entire Hamt data structure.
func (h *HamtTransient) LongString(indent string) string {
	var str string
	if h.root != nil {
		str = indent + fmt.Sprintf("HamtTransient{ nentries: %d, root:\n", h.nentries)
		str += indent + h.root.LongString(indent, 0)
		str += indent + "} //HamtTransient"
	} else {
		str = indent + fmt.Sprintf("HamtTransient{ nentries: %d, root: nil }", h.nentries)
	}
	return str
}
