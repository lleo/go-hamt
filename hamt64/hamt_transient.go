package hamt64

import (
	"fmt"
	"log"
)

// HamtTransient is the datastructure which the Transient Hamt methods are
// called upon. In fact it is identical to the HamtFunctional datastructure and
// all the table and leaf datastructures it uses are the same ones used by the
// HamtTransient implementation. It is its own type so that the methods it calls
// are the transient version of the hamt64.Hamt interface.
//
// The Transient version of the Hamt datastructure, does all modifications
// in-place. So sharing this datastruture between threads is NOT safe unless
// you were to implement a locking stategy CORRECTLY.
type HamtTransient struct {
	root     tableI
	nentries uint
	grade    bool
	compinit bool
}

// NewTransient constructs a new HamtTransient datastructure based on the opt
// argument.
func NewTransient(opt int) *HamtTransient {
	var h = new(HamtTransient)

	switch opt {
	case HybridTables:
		h.grade = true
		h.compinit = true
	case SparseTablesOnly:
		h.grade = false
		h.compinit = true
	case FixedTablesOnly:
		fallthrough
	default:
		h.grade = false
		h.compinit = false
	}

	return h
}

// IsEmpty simply returns if the HamtTransient datastucture has no entries.
func (h *HamtTransient) IsEmpty() bool {
	return h.root == nil
}

// Nentries return the number of (key,value) pairs are stored in the
// HamtTransient datastructure.
func (h *HamtTransient) Nentries() uint {
	return h.nentries
}

// ToFunctional creates a new HamtFunctional datastructure and simply copies the
// values stored in the HamtTransient datastructure over to the HamtFunctional
// datastructure, then it returns a pointer to the HamtFunctional datastructure
// as a hamt64.Hamt interface.
//
// WARNING: given that ToFunctional() just copies pointers to a new
// HamtTransient datastruture, ANY modification to the original HamtTransient
// datastructure will change the new HamtFunctional datastructure (as they
// share the exact same tables & leafs).
//
// Modifications to the new HamtFunctional datastructure WILL NOT affect the
// original HamtTransient datastructure because it does all its modifiation in
// a copy-on-write manner.
//
// The only way to convert a HamtTransient to a HamtFunctional and keep the
// functionality of both is to first perfom a DeepCopy().
func (h *HamtTransient) ToFunctional() Hamt {
	return &HamtFunctional{
		root:     h.root,
		nentries: h.nentries,
		grade:    h.grade,
		compinit: h.compinit,
	}
}

// ToTransient does nothing to a HamtTransient datastructure. This method only
// returns the HamtTransient datastructure pointer as a hamt64.Hamt interface.
func (h *HamtTransient) ToTransient() Hamt {
	return h
}

// DeepCopy() copies the HamtTransient datastructure and every table it contains
// recursively. This is expensive, but usefull, if you want to use ToTransient()
// and ToFunctional().
func (h *HamtTransient) DeepCopy() Hamt {
	var nh = new(HamtTransient)
	nh.root = h.root.deepCopy()
	nh.nentries = h.nentries
	nh.grade = h.grade
	nh.compinit = h.compinit
	return nh
}

func (h *HamtTransient) find(k Key) (tableStack, leafI, uint) {
	if h.IsEmpty() {
		return nil, nil, 0
	}

	var hv = k.Hash()
	var curTable = h.root

	var path = newTableStack()
	var leaf leafI
	var idx uint

	var depth uint
DepthIter:
	for depth = 0; depth <= MaxDepth; depth++ {
		path.push(curTable)
		idx = hv.Index(depth)

		var curNode = curTable.get(idx)
		switch n := curNode.(type) {
		case nil:
			leaf = nil
			break DepthIter
		case leafI:
			leaf = n
			break DepthIter
		case tableI:
			if depth == MaxDepth {
				log.Panicf("SHOULD NOT BE REACHED; depth,%d == MaxDepth,%d & tableI entry found; %s", depth, MaxDepth, n)
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
// BenchHamt64:
//func (h *HamtTransient) Get(k Key) (interface{}, bool) {
//	var _, leaf, _ = h.find(k)
//
//	if leaf == nil {
//		return nil, false
//	}
//
//	return leaf.get(k)
//}

// Get retrieves the value related to the key in the HamtTransient
// datastructure. It also return a bool to indicate the value was found. This
// allows you to store nil values in the HamtTransient datastructure.
func (h *HamtTransient) Get(k Key) (interface{}, bool) {
	if h.IsEmpty() {
		return nil, false
	}

	var val interface{}
	var found bool

	var hv = k.Hash()

	var curTable = h.root //ISA tableI

	for depth := uint(0); depth <= MaxDepth; depth++ {
		var idx = hv.Index(depth)
		var curNode = curTable.get(idx) //nodeI

		if curNode == nil {
			return nil, false
		}

		if leaf, isLeaf := curNode.(leafI); isLeaf {
			val, found = leaf.get(k)
			return val, found
		}

		if depth == MaxDepth {
			panic("SHOULD NOT HAPPEN")
		}
		curTable = curNode.(tableI)
	}

	panic("SHOULD NEVER BE REACHED")
}

func (h *HamtTransient) createRootTable(leaf leafI) tableI {
	if h.compinit {
		return createRootSparseTable(leaf)
	}
	return createRootFixedTable(leaf)
}

func (h *HamtTransient) createTable(depth uint, leaf1 leafI, leaf2 *flatLeaf) tableI {
	if h.compinit {
		return createSparseTable(depth, leaf1, leaf2)
	}
	return createFixedTable(depth, leaf1, leaf2)
}

// Put stores a new (key,value) pair in the HamtTransient datastructure. It
// returns a bool indicating if a new pair were added or if the value replaced
// the value in a previously stored (key,value) pair. Either way it returns and
// new HamtTransient datastructure containing the modification.
func (h *HamtTransient) Put(k Key, v interface{}) (Hamt, bool) {
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
		if h.grade && (curTable.nentries()+1) == UpgradeThreshold {
			var newTable tableI
			newTable = upgradeToFixedTable(
				curTable.Hash(), depth, curTable.entries())
			if curTable == h.root {
				h.root = newTable
			} else {
				var parentTable = path.peek()
				var parentIdx = k.Hash().Index(depth - 1)
				parentTable.replace(parentIdx, newTable)
			}
			curTable = newTable
		}
		curTable.insert(idx, newFlatLeaf(k, v))
		added = true
	} else {
		// This is the condition that allows collision leafs to exist at a level
		// less than MaxDepth. I don't know if I want to allow this...
		if leaf.Hash() == k.Hash() {
			var newLeaf leafI
			// There are four possibilities here:
			// if leaf isa collision leaf
			//   k is identical to one of the kv pairs in collision leaf; hence
			//     we replace that ones value and added = false
			//   k is unique in the collision leaf and the kv pair is added;
			//     this is very rare; the underlying key basis is different but
			//     the Hash is identical.
			// if leaf isa flat leaf
			//   k is identical to the flat leaf's key; hence the value is
			//     replaced and added == false
			//   k is not identical to the flat leaf's key; and a collision leaf
			//     is created and added == true; again this is very rare; the
			//     underlying key basis is different but the Hash is identical
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

// Del searches the HamtTransient for the key argument and returns three
// values: a Hamt datastuture, a value, and a bool. If the key was found then
// the bool returned is true and the value is the value related to that key and
// the returned Hamt is a new HamtTransient datastructure without. If the
// (key, value) pair. If key was not found, then the bool is false, the value is
// nil, and the Hamt value is the original HamtTransient datastructure.
func (h *HamtTransient) Del(k Key) (Hamt, interface{}, bool) {
	if h.IsEmpty() {
		return h, nil, false
	}

	var path, leaf, idx = h.find(k)

	var curTable = path.pop()
	var depth = uint(path.len())

	if leaf == nil {
		return h, nil, false
	}

	var newLeaf, val, deleted = leaf.del(k)

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
		case curTable != h.root && curTable.nentries() == 1:
			var lastNode = curTable.entries()[0].node
			if _, isLeaf := lastNode.(leafI); isLeaf {
				var parentTable = path.peek()
				var parentIdx = k.Hash().Index(depth - 1)
				parentTable.replace(parentIdx, lastNode)
			}

			// else check if downgrade allowed and required
		case h.grade && curTable.nentries() == DowngradeThreshold:
			//when nentries is decr'd it will be <DowngradeThreshold
			var newTable = downgradeToSparseTable(
				curTable.Hash(), depth, curTable.entries())
			if curTable == h.root { //aka path.len() == 0 or path.peek() == nil
				h.root = newTable
			} else {
				var parentTable = path.peek()
				var parentIdx = k.Hash().Index(depth - 1)
				parentTable.replace(parentIdx, newTable)
			}
		}
	}

	return h, val, deleted
}

// String returns a string representation of the Hamt string.
func (h *HamtTransient) String() string {
	return fmt.Sprintf(
		"HamtTransient{ nentries: %d, root: %s }",
		h.nentries,
		h.root.LongString("", 0),
	)
}

// LongString returns a complete listing of the entire Hamt data structure.
func (h *HamtTransient) LongString(indent string) string {
	var str string
	if h.root != nil {
		str = indent +
			fmt.Sprintf("HamtTransient{ nentries: %d, root:\n", h.nentries)
		str += indent + h.root.LongString(indent, 0)
		str += indent + "} //HamtTransient"
	} else {
		str = indent +
			fmt.Sprintf("HamtTransient{ nentries: %d, root: nil }", h.nentries)
	}
	return str
}
