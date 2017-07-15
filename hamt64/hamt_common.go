package hamt64

import (
	"fmt"
	"log"
)

// This is here as the Hamt base data struture.
type hamtBase struct {
	root       tableI
	nentries   uint
	grade      bool
	startFixed bool
}

func (h *hamtBase) init(opt int) {
	// boolean zero value is false
	switch opt {
	case HybridTables:
		h.grade = true
		//h.startFixed = false
	case SparseTables:
		//h.grade = false
		//h.startFixed = false
	case FixedTables:
		//h.grade = false
		h.startFixed = true
	}
}

// IsEmpty simply returns if the HamtFunctional datastucture has no entries.
func (h *hamtBase) IsEmpty() bool {
	return h.root == nil
	//return h.nentries == 0
}

// Nentries return the number of (key,value) pairs are stored in the
// HamtFunctional datastructure.
func (h *hamtBase) Nentries() uint {
	return h.nentries
}

// DeepCopy() copies the HamtFunctional datastructure and every table it
// contains recursively. This is expensive, but usefull, if you want to use
// ToTransient() and ToFunctional().
func (h *hamtBase) DeepCopy() Hamt {
	var nh = new(HamtFunctional)
	nh.root = h.root.deepCopy()
	nh.nentries = h.nentries
	nh.grade = h.grade
	nh.startFixed = h.startFixed
	return nh
}

func (h *hamtBase) find(k *iKey) (tableStack, leafI, uint) {
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
//func (h *hamtBase) Get(bs []byte) (interface{}, bool) {
//	var k = newKey(bs)
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
func (h *hamtBase) Get(key []byte) (interface{}, bool) {
	if h.IsEmpty() {
		return nil, false
	}

	var k = newKey(key)
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

func (h *hamtBase) createRootTable(leaf leafI) tableI {
	if h.startFixed {
		return createRootFixedTable(leaf)
	}
	return createRootSparseTable(leaf)
}

func (h *hamtBase) createTable(depth uint, leaf1 leafI, leaf2 *flatLeaf) tableI {
	if h.startFixed {
		return createFixedTable(depth, leaf1, leaf2)
	}
	return createSparseTable(depth, leaf1, leaf2)
}

// String returns a string representation of the hamtBase stastructure.
// Secifically it returns a representation of the datastructure with the
// nentries value of Nentries() and a representation of the root table.
func (h *hamtBase) String() string {
	return fmt.Sprintf(
		"hamtBase{ nentries: %d, root: %s }",
		h.nentries,
		h.root.String(),
	)
}

// LongString returns a complete listing of the entire Hamt data structure
// recursively indented.
func (h *hamtBase) LongString(indent string) string {
	var str string
	if h.root != nil {
		str = indent +
			fmt.Sprintf("hamtBase{ nentries: %d, root:\n", h.nentries)
		str += indent + h.root.LongString(indent, 0)
		str += indent + "} //hamtBase"
	} else {
		str = indent +
			fmt.Sprintf("hamtBase{ nentries: %d, root: nil }", h.nentries)
	}
	return str
}

type visitFn func(nodeI)

func (h *hamtBase) visit(fn visitFn) uint {
	return h.root.visit(fn, 0)
}

// Count returns a break down of the number of items in the HAMT.
func (h *hamtBase) Count() (maxDepth uint, counts *Counts) {
	counts = new(Counts)

	// countFn closes over the counts variable
	var countFn = func(n nodeI) {
		switch x := n.(type) {
		case nil:
			counts.Nils++
		case *fixedTable:
			counts.Nodes++
			counts.Tables++
			counts.FixedTables++
			counts.TableCountsByNentries[x.nentries()]++
			counts.TableCountsByDepth[x.depth]++
		case *sparseTable:
			counts.Nodes++
			counts.Tables++
			counts.SparseTables++
			counts.TableCountsByNentries[x.nentries()]++
			counts.TableCountsByDepth[x.depth]++
		case *flatLeaf:
			counts.Nodes++
			counts.Leafs++
			counts.FlatLeafs++
			counts.KeyVals += 1
		case *collisionLeaf:
			counts.Nodes++
			counts.Leafs++
			counts.CollisionLeafs++
			counts.KeyVals += uint(len(x.kvs))
		}
	}

	maxDepth = h.visit(countFn)
	return maxDepth, counts
}
