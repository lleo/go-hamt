package hamt32

import (
	"fmt"
	"log"
)

// This is here as the Hamt base data struture.
type HamtBase struct {
	root       tableI
	nentries   uint
	grade      bool
	startFixed bool
}

func (h *HamtBase) init(opt int) {
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
func (h *HamtBase) IsEmpty() bool {
	return h.root == nil
	//return h.nentries == 0
}

// Nentries return the number of (key,value) pairs are stored in the
// HamtFunctional datastructure.
func (h *HamtBase) Nentries() uint {
	return h.nentries
}

// DeepCopy() copies the HamtFunctional datastructure and every table it
// contains recursively. This is expensive, but usefull, if you want to use
// ToTransient() and ToFunctional().
func (h *HamtBase) DeepCopy() Hamt {
	var nh = new(HamtFunctional)
	nh.root = h.root.deepCopy()
	nh.nentries = h.nentries
	nh.grade = h.grade
	nh.startFixed = h.startFixed
	return nh
}

func (h *HamtBase) find(k *iKey) (tableStack, leafI, uint) {
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
//func (h *HamtBase) Get(bs []byte) (interface{}, bool) {
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
func (h *HamtBase) Get(bs []byte) (interface{}, bool) {
	if h.IsEmpty() {
		return nil, false
	}

	var k = newKey(bs)
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

func (h *HamtBase) createRootTable(leaf leafI) tableI {
	if h.startFixed {
		return createRootFixedTable(leaf)
	}
	return createRootSparseTable(leaf)
}

func (h *HamtBase) createTable(depth uint, leaf1 leafI, leaf2 *flatLeaf) tableI {
	if h.startFixed {
		return createFixedTable(depth, leaf1, leaf2)
	}
	return createSparseTable(depth, leaf1, leaf2)
}

// String returns a simple string representation of the HamtTransient data
// structure.
func (h *HamtBase) String() string {
	return fmt.Sprintf(
		"HamtBase{ nentries: %d, root: %s }",
		h.nentries,
		h.root.String(),
	)
}

// LongString returns a complete recusive listing of the entire HamtBase
// data structure.
func (h *HamtBase) LongString(indent string) string {
	var str string
	if h.root != nil {
		str = indent +
			fmt.Sprintf("HamtBase{ nentries: %d, root:\n", h.nentries)
		str += indent + h.root.LongString(indent, 0)
		str += indent + "} //HamtBase"
	} else {
		str = indent +
			fmt.Sprintf("HamtBase{ nentries: %d, root: nil }", h.nentries)
	}
	return str
}

func (h *HamtBase) visit(fn visitFn, arg interface{}) uint {
	return h.root.visit(fn, arg, 0)
}

// Count returns a break down of the number of items in the HAMT.
func (h *HamtBase) Count() (maxDepth uint, counts *Counts) {
	counts = new(Counts)
	maxDepth = h.visit(count, counts)
	return maxDepth, counts
}

type visitFn func(nodeI, interface{})

//var deepest uint

func count(n nodeI, arg interface{}) {
	var counts = arg.(*Counts)

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
