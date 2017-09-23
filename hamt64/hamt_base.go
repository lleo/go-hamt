package hamt64

import (
	"fmt"
)

// This is here as the Hamt base data struture.
type hamtBase struct {
	root       fixedTable
	nentries   uint
	nograde    bool
	startFixed bool
}

func (h *hamtBase) init(tblOpt int) {
	// boolean zero value is false
	switch tblOpt {
	case HybridTables:
		h.nograde = false
		//h.startFixed = false
	case SparseTables:
		h.nograde = true
		//h.startFixed = false
	case FixedTables:
		h.nograde = true
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
	nh.nograde = h.nograde
	nh.startFixed = h.startFixed
	return nh
}

// // copyKey is meant to guard against the data of the slice being modified
// // during two periods it may be modified outside the call to Get, Put, and/or
// // Del. First the lookup from the call site to the match for the op. Second,
// // during the storage as the key in the leaf which is a much longer time.
// // The First applies to Get, Put, and Del, the second applies only to Put.
// // We hope this function is inlined.
// func copyKey(key []byte) []byte {
// 	var k = make([]byte, len(key))
// 	copy(k, key)
// 	return k
// }

func (h *hamtBase) find(hv HashVal) (tableStack, leafI, uint) {
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
			curTable = n
		}
	}

	return path, leaf, idx
}

// This is slower due to extraneous code and allocations in find().
//func (h *hamtBase) Get(key KeyI) (interface{}, bool) {
//	var hv = CalcHash(key)
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
func (h *hamtBase) Get(key KeyI) (interface{}, bool) {
	if h.IsEmpty() {
		return nil, false
	}

	var hv = key.Hash()
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
			curTable = n
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

// walk traverses the Trie in pre-order traversal. For a Trie this is also a
// in-order traversal of all leaf nodes.
//
// walk returns false if the traversal stopped early.
func (h *hamtBase) walk(fn visitFn) bool {
	return h.root.visit(fn)
}

// Range executes the given function for every KeyVal pair in the Hamt. KeyVal
// pairs are visited in a seeminly random order.
//
// Note: we say "seemingly random order", becuase there is a predictable order
// based on the hash value of the Keys and the insertion order of the KeyVal
// pairs, so you cannot reley on the "randomness" of the order of KeyVal pairs.
func (h *hamtBase) Range(fn func(KeyI, interface{}) bool) {
	var visitLeafs = func(n nodeI) bool {
		var keepOn = true

		switch x := n.(type) {
		case nil, tableI:
			//ignore
		case leafI:
			for _, kv := range x.keyVals() {
				if !fn(kv.Key, kv.Val) {
					keepOn = false
					break //for
				}
			}
		}

		return keepOn
	} //end: visitLeafsFn = func(nodeI)

	h.walk(visitLeafs)
}

// Stats walks the Hamt in a pre-order traversal and populates a Stats data
// struture which it returns.
func (h *hamtBase) Stats() *Stats {
	var stats = new(Stats)

	// statFn closes over the stats variable
	var statFn = func(n nodeI) bool {
		var keepOn = true
		switch x := n.(type) {
		case nil:
			stats.Nils++
			keepOn = false
		case *fixedTable:
			stats.Nodes++
			stats.Tables++
			stats.FixedTables++
			stats.TableCountsByNentries[x.nentries()]++
			stats.TableCountsByDepth[x.depth]++
			if x.depth > stats.MaxDepth {
				stats.MaxDepth = x.depth
			}
		case *sparseTable:
			stats.Nodes++
			stats.Tables++
			stats.SparseTables++
			stats.TableCountsByNentries[x.nentries()]++
			stats.TableCountsByDepth[x.depth]++
			if x.depth > stats.MaxDepth {
				stats.MaxDepth = x.depth
			}
		case *flatLeaf:
			stats.Nodes++
			stats.Leafs++
			stats.FlatLeafs++
			stats.KeyVals += 1
			keepOn = false
		case *collisionLeaf:
			stats.Nodes++
			stats.Leafs++
			stats.CollisionLeafs++
			stats.KeyVals += uint(len(x.kvs))
			keepOn = false
		}
		return keepOn
	}

	h.walk(statFn)
	return stats
}
