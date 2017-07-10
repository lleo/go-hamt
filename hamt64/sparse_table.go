package hamt64

import (
	"fmt"
	"strings"
)

// sparseTableInitCap constant sets the default capacity of a new
// sparseTable.
const sparseTableInitCap int = 2

// New sparseTable layout size == 48
type sparseTable struct {
	nodes    []nodeI // 24
	depth    uint    // 8; amd64 cpu
	nodeMap  Bitmap  // 8; uint32 x 2
	hashPath HashVal // 8
}

func (t *sparseTable) copy() tableI {
	var nt = new(sparseTable)
	nt.hashPath = t.hashPath
	nt.depth = t.depth
	nt.nodeMap = t.nodeMap

	nt.nodes = make([]nodeI, len(t.nodes), cap(t.nodes))
	copy(nt.nodes, t.nodes)

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
			//leafI's are functional, so no need to copy them.
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
	ct.nodeMap.Set(idx)
	ct.nodes = make([]nodeI, 1, sparseTableInitCap)
	ct.nodes[0] = lf

	return ct
}

func createSparseTable(depth uint, leaf1 leafI, leaf2 *flatLeaf) tableI {
	_ = AssertOn && assert(depth > 0, "createSparseTable(): depth < 1")

	_ = AssertOn && assertf(
		leaf1.Hash().HashPath(depth) == leaf2.Hash().HashPath(depth),
		"createSparseTable(): hp1,%s != hp2,%s",
		leaf1.Hash().HashPath(depth),
		leaf2.Hash().HashPath(depth))

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
		nt.nodeMap.Set(ent.idx)
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

	strs[1] = indent + "\tnodeMap=" + t.nodeMap.String() + ","

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
	//return t.nodeMap.Count(IndexLimit)
}

func (t *sparseTable) entries() []tableEntry {
	var n = t.nentries()
	var ents = make([]tableEntry, n)

	for idx, j := uint(0), uint(0); idx < IndexLimit; idx++ {
		if t.nodeMap.IsSet(idx) {
			ents[j] = tableEntry{idx, t.nodes[j]}
			j++
		}
	}

	return ents
}

func (t *sparseTable) get(idx uint) nodeI {
	if !t.nodeMap.IsSet(idx) {
		return nil
	}

	var j = t.nodeMap.Count(idx)

	return t.nodes[j]
}

func (t *sparseTable) insert(idx uint, n nodeI) {
	_ = AssertOn && assert(!t.nodeMap.IsSet(idx),
		"t.insert(idx, n) where idx slot is NOT empty; this should be a replace")

	var j = t.nodeMap.Count(idx)
	if j == uint(len(t.nodes)) {
		t.nodes = append(t.nodes, n)
	} else {
		t.nodes = append(t.nodes[:j], append([]nodeI{n}, t.nodes[j:]...)...)
	}

	t.nodeMap.Set(idx)
}

func (t *sparseTable) replace(idx uint, n nodeI) {
	_ = AssertOn && assert(t.nodeMap.IsSet(idx),
		"t.replace(idx, n) where idx slot is empty; this should be an insert")

	var j = t.nodeMap.Count(idx)
	t.nodes[j] = n
}

func (t *sparseTable) remove(idx uint) {
	_ = AssertOn && assert(t.nodeMap.IsSet(idx),
		"t.remove(idx) where idx slot is already empty")

	var j = t.nodeMap.Count(idx)
	if int(j) == len(t.nodes)-1 {
		t.nodes = t.nodes[:j]
	} else {
		t.nodes = append(t.nodes[:j], t.nodes[j+1:]...)
	}

	t.nodeMap.Unset(idx)
}
