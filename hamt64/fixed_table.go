package hamt64

import (
	"fmt"
	"log"
	"strings"
)

type fixedTable struct {
	nodes    [IndexLimit]nodeI // 1024; 64*16
	depth    uint              // 8; amd64
	nents    uint              // 8; amd64
	hashPath HashVal           // 8
}

func (t *fixedTable) copy() tableI {
	var nt = new(fixedTable)
	nt.hashPath = t.hashPath
	nt.depth = t.depth
	nt.nents = t.nents
	nt.nodes = t.nodes
	return nt
}

func (t *fixedTable) deepCopy() tableI {
	var nt = new(fixedTable)
	nt.hashPath = t.hashPath
	nt.depth = t.depth
	nt.nents = t.nents
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

func createRootFixedTable(lf leafI) tableI {
	var idx = lf.Hash().Index(0)

	var ft = new(fixedTable)
	//ft.hashPath = 0
	//ft.depth = 0
	//ft.nents = 0
	ft.set(idx, lf)

	return ft
}

func createFixedTable(depth uint, leaf1 leafI, leaf2 *flatLeaf) tableI {
	if depth < 1 {
		log.Panic("createFixedTable(): depth < 1")
	}
	var hp1 = leaf1.Hash().HashPath(depth)
	var hp2 = leaf2.Hash().HashPath(depth)
	if hp1 != hp2 {
		log.Panicf("createFixedTable(): hp1,%s != hp2,%s",
			hp1.HashPathString(depth), hp2.HashPathString(depth))
	}
	//for d := uint(0); d < depth; d++ {
	//	if leaf1.Hash().Index(d) != leaf2.Hash().Index(d) {
	//		log.Panicf("createFixedTable(): leaf1.Hash().Index(%d) != leaf2.Hash().Index(%d)", d, d)
	//	}
	//}

	var retTable = new(fixedTable)
	retTable.hashPath = leaf1.Hash().HashPath(depth)
	retTable.depth = depth

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
			node = createFixedTable(depth+1, leaf1, leaf2)
		}
		retTable.insert(idx1, node)
	}

	return retTable
}

func upgradeToFixedTable(
	hashPath HashVal,
	depth uint,
	ents []tableEntry,
) *fixedTable {
	var ft = new(fixedTable)
	ft.hashPath = hashPath
	ft.depth = depth
	ft.nents = uint(len(ents))

	for _, ent := range ents {
		ft.nodes[ent.idx] = ent.node
	}

	return ft
}

// Hash returns an incomplete Hash of this table. Any levels past it's current
// depth should be zero.
func (t *fixedTable) Hash() HashVal {
	return t.hashPath
}

// String return a string representation of this table including the hashPath,
// depth, and number of entries.
func (t *fixedTable) String() string {
	return fmt.Sprintf("fixedTable{hashPath=%s, depth=%d, nentries()=%d}",
		t.hashPath, t.depth, t.nentries())
}

// LongString returns a string representation of this table and all the tables
// contained herein recursively.
func (t *fixedTable) LongString(indent string, depth uint) string {
	var strs = make([]string, 3+t.nentries())

	strs[0] = indent + "fixedTable{"
	strs[1] = indent + fmt.Sprintf("\thashPath=%s, depth=%d, nents=%d,",
		t.hashPath.HashPathString(depth+1), t.depth, t.nents)

	var j = 0
	for i, n := range t.nodes {
		if t.nodes[i] != nil {
			if t, isTable := t.nodes[i].(tableI); isTable {
				strs[2+j] = indent + fmt.Sprintf("\tnodes[%d]:\n", i) +
					t.LongString(indent+"\t", depth+1)
			} else {
				strs[2+j] = indent + fmt.Sprintf("\tnodes[%d]: %s", i, n)
			}
			j++
		}
	}

	strs[len(strs)-1] = indent + "}"

	return strings.Join(strs, "\n")
}

func (t *fixedTable) nentries() uint {
	return t.nents
}

func (t *fixedTable) entries() []tableEntry {
	var n = t.nentries()
	var ents = make([]tableEntry, n)
	var i, j uint
	for i, j = 0, 0; j < n && i < IndexLimit; i++ {
		if t.nodes[i] != nil {
			ents[j] = tableEntry{i, t.nodes[i]}
			j++
		}
	}
	return ents
}

func (t *fixedTable) get(idx uint) nodeI {
	return t.nodes[idx]
}

func (t *fixedTable) set(idx uint, nn nodeI) {
	if nn != nil && t.nodes[idx] == nil {
		t.nents++
	} else if nn == nil && t.nodes[idx] != nil {
		t.nents--
	}
	t.nodes[idx] = nn

	return
}

func (t *fixedTable) insert(idx uint, n nodeI) {
	if t.nodes[idx] != nil {
		panic("t.insert(idx, n) where idx slot is NOT empty; this should be a replace")
	}
	t.nodes[idx] = n
	t.nents++
}

func (t *fixedTable) replace(idx uint, n nodeI) {
	if t.nodes[idx] == nil {
		panic("t.replace(idx, n) where idx slot is empty; this should be an insert")
	}
	t.nodes[idx] = n
}

func (t *fixedTable) remove(idx uint) {
	if t.nodes[idx] == nil {
		panic("t.remove(idx) where idx slot is already empty")
	}
	t.nodes[idx] = nil
	t.nents--
}
