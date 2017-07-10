package hamt64

import (
	"fmt"
	"log"
)

// This is here as the Hamt base data struture.
type Common struct {
	root     tableI
	nentries uint
	grade    bool
	compinit bool
}

// IsEmpty simply returns if the HamtFunctional datastucture has no entries.
func (h *Common) IsEmpty() bool {
	return h.root == nil
	//return h.nentries == 0
}

// Nentries return the number of (key,value) pairs are stored in the
// HamtFunctional datastructure.
func (h *Common) Nentries() uint {
	return h.nentries
}

// DeepCopy() copies the HamtFunctional datastructure and every table it
// contains recursively. This is expensive, but usefull, if you want to use
// ToTransient() and ToFunctional().
func (h *Common) DeepCopy() Hamt {
	var nh = new(HamtFunctional)
	nh.root = h.root.deepCopy()
	nh.nentries = h.nentries
	nh.grade = h.grade
	nh.compinit = h.compinit
	return nh
}

func (h *Common) find(k Key) (tableStack, leafI, uint) {
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

// This is slower due to extraneous code and allocations in find().
//func (h *Common) Get(k Key) (interface{}, bool) {
//	var _, leaf, _ = h.find(k)
//
//	if leaf == nil {
//		return nil, false
//	}
//
//	return leaf.get(k)
//}

// Get retrieves the value related to the key in the HamtFunctional
// datastructure. It also return a bool to indicate the value was found. This
// allows you to store nil values in the HamtFunctional datastructure.
func (h *Common) Get(k Key) (interface{}, bool) {
	if h.IsEmpty() {
		return nil, false
	}

	var hv = k.Hash()

	var curTable = h.root //ISA tableI

	for depth := uint(0); depth <= MaxDepth; depth++ {
		var idx = hv.Index(depth)
		var curNode = curTable.get(idx) //nodeI

		if curNode == nil {
			return nil, false
		}

		if leaf, isLeaf := curNode.(leafI); isLeaf {
			return leaf.get(k)
		}

		if depth == MaxDepth {
			panic("SHOULD NOT HAPPEN")
		}
		curTable = curNode.(tableI)
	}

	panic("SHOULD NEVER BE REACHED")
}

func (h *Common) createRootTable(leaf leafI) tableI {
	if h.compinit {
		return createRootSparseTable(leaf)
	}
	return createRootFixedTable(leaf)
}

func (h *Common) createTable(depth uint, leaf1 leafI, leaf2 *flatLeaf) tableI {
	if h.compinit {
		return createSparseTable(depth, leaf1, leaf2)
	}
	return createFixedTable(depth, leaf1, leaf2)
}

// String returns a string representation of the Common stastructure.
// Secifically it returns a representation of the datastructure with the
// nentries value of Nentries() and a representation of the root table.
func (h *Common) String() string {
	return fmt.Sprintf(
		"Common{ nentries: %d, root: %s }",
		h.nentries,
		h.root.LongString("", 0),
	)
}

// LongString returns a complete listing of the entire Hamt data structure
// recursively indented..
func (h *Common) LongString(indent string) string {
	var str string
	if h.root != nil {
		str = indent +
			fmt.Sprintf("Common{ nentries: %d, root:\n", h.nentries)
		str += indent + h.root.LongString(indent, 0)
		str += indent + "} //Common"
	} else {
		str = indent +
			fmt.Sprintf("Common{ nentries: %d, root: nil }", h.nentries)
	}
	return str
}
