package hamt32

import (
	"log"

	"github.com/lleo/go-hamt-key"
)

type HamtFunctional struct {
	root     tableI
	nentries uint
	grade    bool
	compinit bool
}

func NewFunctional(opt int) *HamtFunctional {
	var h = new(HamtFunctional)

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

func (h *HamtFunctional) IsEmpty() bool {
	return h.root == nil
	//return h.nentries == 0
}

func (h *HamtFunctional) Nentries() uint {
	return h.nentries
}

// persist() is ONLY called on a fresh copy of the current Hamt.
// Hence, modifying it is allowed.
func (nh *HamtFunctional) persist(oldTable, newTable tableI, path tableStack) {
	if path.isEmpty() {
		nh.root = newTable
		return
	}

	var depth = uint(path.len())
	var parentDepth = depth - 1

	var parentIdx = oldTable.Hash30().Index(parentDepth)

	var oldParent = path.pop()
	var newParent tableI = oldParent.copy()

	if newTable == nil {
		newParent.remove(parentIdx)
	} else {
		newParent.replace(parentIdx, newTable)
	}

	nh.persist(oldParent, newParent, path) //recurses at most MaxDepth-1 times

	return
}

func (h *HamtFunctional) find(k key.Key) (tableStack, leafI, uint) {
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

//func (h *HamtTransient) Get(k key.Key) (interface{}, bool) {
//	var _, leaf, _ = h.find(k)
//
//	if leaf == nil {
//		return nil, false
//	}
//
//	return leaf.get(k)
//}

func (h *HamtFunctional) Get(k key.Key) (interface{}, bool) {
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

func (h *HamtFunctional) createRootTable(leaf leafI) tableI {
	if h.compinit {
		return createRootCompressedTable(leaf)
	}
	return createRootFullTable(leaf)
}

func (h *HamtFunctional) createTable(depth uint, leaf1 leafI, leaf2 *flatLeaf) tableI {
	if h.compinit {
		return createCompressedTable(depth, leaf1, leaf2)
	}
	return createFullTable(depth, leaf1, leaf2)
}

func (h *HamtFunctional) Put(k key.Key, v interface{}) (Hamt, bool) {
	var nh = new(HamtFunctional)
	*nh = *h

	if nh.IsEmpty() {
		nh.root = nh.createRootTable(newFlatLeaf(k, v))
		nh.nentries++
		return nh, true
	}

	var path, leaf, idx = nh.find(k)

	var curTable = path.pop()
	var depth = uint(path.len())
	var added bool

	var newTable tableI
	if leaf == nil {
		if nh.grade && (curTable.nentries()+1) == upgradeThreshold {
			newTable = upgradeToFullTable(curTable.Hash30(), curTable.entries())
		} else {
			newTable = curTable.copy()
			newTable.insert(idx, newFlatLeaf(k, v))
		}
		added = true
	} else {
		newTable = curTable.copy()
		if leaf.Hash30() == k.Hash30() {
			var newLeaf leafI
			newLeaf, added = leaf.put(k, v)
			newTable.replace(idx, newLeaf)
		} else {
			var tmpTable = nh.createTable(depth+1, leaf, newFlatLeaf(k, v))
			newTable.replace(idx, tmpTable)
			added = true
		}
	}

	if added {
		nh.nentries++
	}

	nh.persist(curTable, newTable, path)

	return nh, added
}

func (h *HamtFunctional) Del(k key.Key) (Hamt, interface{}, bool) {
	if h.IsEmpty() {
		return h, nil, false
	}

	var path, leaf, idx = h.find(k)

	var curTable = path.pop()

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

	var newTable tableI = curTable.copy()
	if newLeaf != nil { //leaf was a CollisionLeaf
		newTable.replace(idx, newLeaf)
	} else { //leaf was a FlatLeaf
		newTable.remove(idx)

		// Side-Effects of removing a KeyVal from the table
		switch {
		case newTable.nentries() == 0:
			newTable = nil
		case h.grade && newTable.nentries() == downgradeThreshold:
			newTable = downgradeToCompressedTable(
				newTable.Hash30(), newTable.entries())
		}
	}

	var nh = new(HamtFunctional)
	*nh = *h

	nh.nentries--

	nh.persist(curTable, newTable, path)

	return nh, val, deleted
}
