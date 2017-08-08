package hamt32

import (
	"fmt"
)

// This is here as the Hamt base data struture.
type hamtBase struct {
	root       fixedTable
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
	//return h.root == nil
	return h.nentries == 0
}

// Nentries return the number of (key,value) pairs are stored in the
// HamtFunctional data structure.
func (h *hamtBase) Nentries() uint {
	return h.nentries
}

// DeepCopy copies the HamtFunctional data structure and every table it
// contains recursively. This is expensive, but usefull, if you want to use
// ToTransient and ToFunctional.
func (h *hamtBase) DeepCopy() Hamt {
	var nh = new(HamtFunctional)
	nh.root = *h.root.deepCopy().(*fixedTable)
	nh.nentries = h.nentries
	nh.grade = h.grade
	nh.startFixed = h.startFixed
	return nh
}

// copyKey is meant to guard against the data of the slice being modified
// during two periods it may be modified outside the call to Get, Put, and/or
// Del. First the lookup from the call site to the match for the op. Second,
// during the storage as the key in the leaf which is a much longer time.
// The First applies to Get, Put, and Del, the second applies only to Put.
// We hope this function is inlined.
func copyKey(key []byte) []byte {
	var k = make([]byte, len(key))
	copy(k, key)
	return k
}

func (h *hamtBase) find(hv hashVal) (tableStack, leafI, uint) {
	var curTable tableI = &h.root

	var path = newTableSlice() //conforms to tableStack interface
	var leaf leafI
	var idx uint

DepthIter:
	for depth := uint(0); depth <= maxDepth; depth++ {
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
			_ = assertOn && assert(depth != maxDepth,
				"Invalid Hamt: TableI found at maxDepth.")
			curTable = n
		default:
			// Can not happend. Just here becuase I am paranoid.
			_ = assertOn && assert(false,
				"Invalid Hamt: curNode != nil || LeafI || TableI")
		}
	}

	return path, leaf, idx
}

// This is slower due to extraneous code and allocations in find().
//func (h *hamtBase) Get(key []byte) (interface{}, bool) {
//	key = copyKey(key)
//	var hv = calcHashVal(key)
//	var _, leaf, _ = h.find(hv)
//
//	if leaf == nil {
//		return nil, false
//	}
//
//	return leaf.get(key)
//}

// Get retrieves the value related to the key in the HamtFunctional
// data structure. It also return a bool to indicate the value was found. This
// allows you to store nil values in the HamtFunctional data structure.
func (h *hamtBase) Get(key []byte) (interface{}, bool) {
	if h.IsEmpty() {
		return nil, false
	}

	//key = copyKey(key)

	var hv = calcHashVal(key)
	var curTable tableI = &h.root

	var val interface{}
	var found bool

DepthIter:
	for depth := uint(0); depth <= maxDepth; depth++ {
		var idx = hv.Index(depth)
		var curNode = curTable.get(idx) //nodeI

		switch n := curNode.(type) {
		case nil:
			val, found = nil, false
			break DepthIter
		case leafI:
			val, found = n.get(key)
			break DepthIter
		case tableI:
			_ = assertOn && assert(depth != maxDepth,
				"Invalid Hamt: TableI found at maxDepth.")
			curTable = n
		default:
			// Can not happend. Just here becuase I am paranoid.
			_ = assertOn && assert(false,
				"Invalid Hamt: curNode != nil || LeafI || TableI")
		}
	}

	return val, found
}

func (h *hamtBase) createTable(depth uint, l1 leafI, l2 *flatLeaf) tableI {
	if h.startFixed {
		return createFixedTable(depth, l1, l2)
	}
	return createSparseTable(depth, l1, l2)
}

// String returns a string representation of the hamtBase stastructure.
// Secifically it returns a representation of the data structure with the
// nentries value of Nentries() and a representation of the root table.
func (h *hamtBase) String() string {
	return fmt.Sprintf(
		"hamtBase{ nentries: %d, root: %s }",
		h.nentries,
		h.root.String(),
	)
}

// LongString returns a complete recusive listing of the entire hamtBase
// data structure.
func (h *hamtBase) LongString(indent string) string {
	var str string

	str = indent +
		fmt.Sprintf("hamtBase{ nentries: %d, root:\n", h.nentries)
	str += indent + h.root.LongString(indent, 0)
	str += indent + "} //hamtBase"

	return str
}

type visitFn func(nodeI)

func (h *hamtBase) visit(fn visitFn) uint {
	return h.root.visit(fn, 0)
}

// Stats returns various measures of the Hamt; for example counts of the numbers
// of various struct types in the HAMT.
func (h *hamtBase) Stats() *Stats {
	var stats = new(Stats)

	// statFn closes over the stats variable
	var statFn = func(n nodeI) {
		switch x := n.(type) {
		case nil:
			stats.Nils++
		case *fixedTable:
			stats.Nodes++
			stats.Tables++
			stats.FixedTables++
			stats.TableCountsByNentries[x.nentries()]++
			stats.TableCountsByDepth[x.depth]++
		case *sparseTable:
			stats.Nodes++
			stats.Tables++
			stats.SparseTables++
			stats.TableCountsByNentries[x.nentries()]++
			stats.TableCountsByDepth[x.depth]++
		case *flatLeaf:
			stats.Nodes++
			stats.Leafs++
			stats.FlatLeafs++
			stats.KeyVals += 1
		case *collisionLeaf:
			stats.Nodes++
			stats.Leafs++
			stats.CollisionLeafs++
			stats.KeyVals += uint(len(x.kvs))
		}
	}

	stats.MaxDepth = h.visit(statFn)
	return stats
}

// IterFunc is the function to call repeatedly to iterate over the Hamt.
type IterFunc func() (KeyVal, bool)

type iterState struct {
	locStack   iterLocStack
	curTable   tableI
	idx        uint
	colLeaf    *collisionLeaf
	colLeafIdx uint
}

func newIterState() *iterState {
	var state = new(iterState)
	state.locStack = make([]iterLocation, 0, DepthLimit)
	return state
}

// Iter returns an IterFunc to be called repeatedly to iterat over the Hamt.
// No modifications should happend during the lifetime of the iterator. For
// HamtFunctional this is not a problem, but for HamtTransient this constaint
// is up to the Library user.
//
//    var next = h.Iter()
//    for kv, ok:= next(); ok; kv, ok = next() {
//        doSomething(kv)
//    }
//
func (h *hamtBase) Iter() IterFunc {
	if h.IsEmpty() {
		return func() (KeyVal, bool) {
			return KeyVal{nil, nil}, false
		}
	}

	//Closure variable
	var state *iterState
	var leaf leafI //not a closure variable

	state, leaf = findFirstLeaf(&h.root)

	if cl, ok := leaf.(*collisionLeaf); ok {
		state.colLeaf = cl
		state.colLeafIdx = 0
	}

	state.curTable, state.idx = state.locStack.pop()
	//state.curTable != nil cuz !h.IsEmpty()

	return func() (KeyVal, bool) {
		if state.colLeaf != nil {
			var kv = state.colLeaf.kvs[state.colLeafIdx]
			state.colLeafIdx++

			if state.colLeafIdx >= uint(len(state.colLeaf.kvs)) {
				state.colLeaf = nil
				state.colLeafIdx = 0
			}

			return KeyVal{copyKey(kv.Key), kv.Val}, true
		}

		var leaf = findNextLeaf(state)

		var retKV KeyVal
		switch x := leaf.(type) {
		case nil:
			// Is this the END? AKA found no leaf
			// If so, I expect state.locStack.len()==0 && state.idx==IndexLimit
			_ = assertOn && assertf(state.locStack.len() == 0,
				"leaf=nil && state.locStack.len(),%d != 0", state.locStack.len())
			_ = assertOn && assertf(state.idx == IndexLimit,
				"leaf=nil && state.idx,%d != IndexLimit,%d",
				state.idx, IndexLimit)

			return KeyVal{nil, nil}, false
		case *flatLeaf:
			retKV = KeyVal{copyKey(x.key), x.val}
		case *collisionLeaf:
			// This is the first time I've visited this colLeaf, the rest of the
			// colLeaf.kvs will be dealt with at the beginning of this func.
			state.colLeaf = x
			state.colLeafIdx = 0

			var kv = state.colLeaf.kvs[state.colLeafIdx]
			retKV = KeyVal{copyKey(kv.Key), kv.Val}

			state.colLeafIdx++
		default:
			panic("WTF! leaf isa leafI but not *flatLeaf nor *collisionLeaf")
		}
		return retKV, true
	}
}

func findFirstLeaf(root tableI) (*iterState, leafI) {
	var state = newIterState()

	var curTable = root
	var idx uint
	var leaf leafI

DepthLoop:
	for uint(state.locStack.len()) < DepthLimit {
	IndexIter:
		for idx = 0; idx < IndexLimit; idx++ {
			var curNode = curTable.get(idx)
			switch x := curNode.(type) {
			case nil:
				// implicit break; my C-trained brain rebels from go-switch
			case leafI:
				state.locStack.push(curTable, idx)
				leaf = x
				break DepthLoop
			case tableI:
				_ = assertOn && assert(uint(state.locStack.len()) != maxDepth,
					"Invalid Hamt: TableI found at maxDepth.")

				state.locStack.push(curTable, idx)
				curTable = x
				break IndexIter
			default:
				// Can not happend. Just here becuase I am paranoid.
				_ = assertOn && assert(false,
					"Invalid Hamt: curNode != nil || LeafI || TableI")
			} //switch
		} //IndexIter
	} //DepthLoop

	return state, leaf
}

func findNextLeaf(state *iterState) leafI {
	var leaf leafI

DepthLoop:
	for uint(state.locStack.len()) < DepthLimit {
	IndexIter:
		for ; state.idx < IndexLimit; state.idx++ {
			var curNode = state.curTable.get(state.idx)
			switch x := curNode.(type) {
			case nil:
				// do nothing
			case leafI:
				leaf = x
				state.idx++
				break DepthLoop
			case tableI:
				_ = assertOn && assert(uint(state.locStack.len()) != maxDepth,
					"Invalid Hamt: TableI found at maxDepth.")

				state.locStack.push(state.curTable, state.idx)
				state.curTable = x
				state.idx = 0
				break IndexIter
			default:
				// Can not happend. Just here becuase I am paranoid.
				_ = assertOn && assert(false,
					"Invalid Hamt: curNode != nil || LeafI || TableI")
			} //switch
		} //IndexIter

		if state.locStack.len() == 0 && state.idx == IndexLimit {
			//THE END
			//leaf = nil
			break DepthLoop
		}

		if state.idx == IndexLimit { //implicit state.itstk.len() > 0
			state.curTable, state.idx = state.locStack.pop()
			state.idx++
		}
	} //DepthLoop

	return leaf
}
