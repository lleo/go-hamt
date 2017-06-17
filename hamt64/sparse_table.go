package hamt64

import (
	"fmt"
	"log"
	"strings"
)

// sparseTableInitCap constant sets the default capacity of a new
// sparseTable.
const sparseTableInitCap int = 8

type sparseTable struct {
	hashPath HashVal
	depth    uint
	nodeMap  uint64
	nodes    []nodeI
}

func (t *sparseTable) copy() tableI {
	var nt = new(sparseTable)
	nt.hashPath = t.hashPath
	nt.depth = t.depth
	nt.nodeMap = t.nodeMap
	nt.nodes = append(nt.nodes, t.nodes...)
	return nt
}

func (t *sparseTable) deepCopy() tableI {
	var nt = new(sparseTable)
	nt.hashPath = t.hashPath
	nt.depth = t.depth
	nt.nodeMap = t.nodeMap
	nt.nodes = make([]nodeI, len(t.nodes), cap(t.nodes))
	for i := 0; i < len(t.nodes); i++ {
		if table, isTable := t.nodes[i].(tableI); isTable {
			nt.nodes[i] = table.deepCopy()
		} else {
			//leafs are functional, so no need to copy
			//nils can be copied just fine; duh!
			nt.nodes[i] = t.nodes[i]
		}
	}
	return nt
}

func createRootSparseTable(lf leafI) tableI {
	var idx = lf.Hash().Index(0)

	var ct = new(sparseTable)
	//ct.hashPath = 0
	//ct.depth = 0
	ct.nodeMap = uint64(1 << idx)
	ct.nodes = make([]nodeI, 1, sparseTableInitCap)
	ct.nodes[0] = lf

	return ct
}

func createSparseTable(depth uint, leaf1 leafI, leaf2 *flatLeaf) tableI {
	if depth < 1 {
		log.Panic("createSparseTable(): depth < 1")
	}
	var hp1 = leaf1.Hash().HashPath(depth)
	var hp2 = leaf2.Hash().HashPath(depth)
	if hp1 != hp2 {
		log.Panic("createSparseTable(): hp1,%s != hp2,%s",
			hp1.HashPathString(depth), hp2.HashPathString(depth))
	}
	//for d := uint(0); d < depth; d++ {
	//	if leaf1.Hash().Index(d) != leaf2.Hash().Index(d) {
	//		log.Panicf("createSparseTable(): leaf1.Hash().Index(%d) != leaf2.Hash().Index(%d)", d, d)
	//	}
	//}

	var retTable = new(sparseTable)
	retTable.hashPath = leaf1.Hash().HashPath(depth)
	retTable.depth = depth
	//retTable.nodeMap = 0
	retTable.nodes = make([]nodeI, 0, sparseTableInitCap)

	var idx1 = leaf1.Hash().Index(depth)
	var idx2 = leaf2.Hash().Index(depth)
	if idx1 != idx2 {
		retTable.insert(idx1, leaf1)
		retTable.insert(idx2, leaf2)
	} else { //idx1 == idx2
		var node nodeI
		if depth == MaxDepth {
			node = newCollisionLeaf(append(leaf1.keyVals(), leaf2.keyVals()...))
		} else {
			node = createSparseTable(depth+1, leaf1, leaf2)
		}
		retTable.insert(idx1, node)
	}

	return retTable
}

func nodeMapString(nodeMap uint64) string {
	var strs = make([]string, 7)

	var top4 = nodeMap >> 60
	strs[0] = fmt.Sprintf("%04b", top4)

	const tenBitMask uint64 = 1<<10 - 1
	for i := uint(0); i < 6; i++ {
		var tenBitVal = (nodeMap & (tenBitMask << (i * 10))) >> (i * 10)
		strs[6-i] = fmt.Sprintf("%010b", tenBitVal) // strs[6..1]
	}

	return strings.Join(strs, " ")
}

// downgradeToSparseTable() converts fixedTable structs that have less than
// or equal to downgradeThreshold tableEntry's. One important thing we know is
// that none of the entries will collide with another.
//
// The ents []tableEntry slice is guaranteed to be in order from lowest idx to
// highest. tableI.entries() also adhears to this contract.
func downgradeToSparseTable(
	hashPath HashVal,
	depth uint,
	ents []tableEntry,
) *sparseTable {
	var nt = new(sparseTable)
	nt.hashPath = hashPath
	//nt.nodeMap = 0
	nt.nodes = make([]nodeI, len(ents), len(ents)+1)

	for i := 0; i < len(ents); i++ {
		var ent = ents[i]
		var nodeBit = uint64(1 << ent.idx)
		nt.nodeMap |= nodeBit
		nt.nodes[i] = ent.node
	}

	return nt
}

// Hash returns an incomplete Hash of this table. Any levels past it's current
// depth should be zero.
func (t *sparseTable) Hash() HashVal {
	return t.hashPath
}

// String return a string representation of this table including the hashPath,
// depth, and number of entries.
func (t *sparseTable) String() string {
	return fmt.Sprintf("sparseTable{hashPath:%s, depth=%d, nentries()=%d}",
		t.hashPath, t.depth, t.nentries())
}

// LongString returns a string representation of this table and all the tables
// contained herein recursively.
func (t *sparseTable) LongString(indent string, depth uint) string {
	var strs = make([]string, 3+len(t.nodes))

	strs[0] = indent +
		fmt.Sprintf("sparseTable{hashPath=%s, depth=%d, nentries()=%d,",
			t.hashPath.HashPathString(depth), t.depth, t.nentries())

	strs[1] = indent + "\tnodeMap=" + nodeMapString(t.nodeMap) + ","

	for i, n := range t.nodes {
		var idx = n.Hash().Index(depth)
		if t, isTable := n.(tableI); isTable {
			strs[2+i] = indent +
				fmt.Sprintf("\tt.nodes[%d]:\n%s",
					idx, t.LongString(indent+"\t", depth+1))
		} else {
			strs[2+i] = indent + fmt.Sprintf("\tt.nodes[%d]: %s", idx, n)
		}
	}

	strs[len(strs)-1] = indent + "}"

	return strings.Join(strs, "\n")
}

func (t *sparseTable) nentries() uint {
	return uint(len(t.nodes))
	//return bitCount64(t.nodeMap)
}

func (t *sparseTable) entries() []tableEntry {
	var n = t.nentries()
	var ents = make([]tableEntry, n)

	for i, j := uint(0), uint(0); i < tableCapacity; i++ {
		var nodeBit = uint64(1 << i)

		if (t.nodeMap & nodeBit) > 0 {
			ents[j] = tableEntry{i, t.nodes[j]}
			j++
		}
	}

	return ents
}

func (t *sparseTable) get(idx uint) nodeI {
	var nodeBit = uint64(1 << idx)

	if (t.nodeMap & nodeBit) == 0 {
		return nil
	}

	var bitMask = nodeBit - 1
	var i = bitCount64(t.nodeMap & bitMask)

	return t.nodes[i]
}

func (t *sparseTable) set(idx uint, nn nodeI) {
	var nodeBit = uint64(1 << idx)
	var bitMask = nodeBit - 1
	var i = bitCount64(t.nodeMap & bitMask)

	if nn != nil {
		if (t.nodeMap & nodeBit) == 0 {
			t.nodeMap |= nodeBit
			t.nodes = append(t.nodes[:i], append([]nodeI{nn}, t.nodes[i:]...)...)
		} else {
			t.nodes[i] = nn
		}
	} else /* if nn == nil */ {
		if (t.nodeMap & nodeBit) > 0 {
			t.nodeMap &^= nodeBit
			t.nodes = append(t.nodes[:i], t.nodes[i+1:]...)
		} /* else {
			// do nothing
		} */
	}
	return
}

func (t *sparseTable) insert(idx uint, n nodeI) {
	var nodeBit = uint64(1 << idx)

	if (t.nodeMap & nodeBit) > 0 {
		panic("t.insert(idx, n) where idx slot is NOT empty; this should be a replace")
	}

	var bitMask = nodeBit - 1
	var i = bitCount64(t.nodeMap & bitMask)
	if i == uint(len(t.nodes)) {
		t.nodes = append(t.nodes, n)
	} else {
		t.nodes = append(t.nodes[:i], append([]nodeI{n}, t.nodes[i:]...)...)
	}
	t.nodeMap |= nodeBit
}

func (t *sparseTable) replace(idx uint, n nodeI) {
	var nodeBit = uint64(1 << idx)

	if (t.nodeMap & nodeBit) == 0 {
		panic("t.replace(idx, n) where idx slot is empty; this should be an insert")
	}

	var bitMask = nodeBit - 1
	var i = bitCount64(t.nodeMap & bitMask)
	t.nodes[i] = n
}

func (t *sparseTable) remove(idx uint) {
	var nodeBit = uint64(1 << idx)

	if (t.nodeMap & nodeBit) == 0 {
		panic("t.remove(idx) where idx slot is already empty")
	}

	var bitMask = nodeBit - 1
	var i = bitCount64(t.nodeMap & bitMask)
	if int(i) == len(t.nodes)-1 {
		t.nodes = t.nodes[:i]
	} else {
		t.nodes = append(t.nodes[:i], t.nodes[i+1:]...)
	}
	t.nodeMap &^= nodeBit
}
