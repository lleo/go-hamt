package hamt32

import (
	"fmt"
	"log"
	"strings"

	"github.com/lleo/go-hamt/key"
)

type fullTable struct {
	hashPath key.HashVal30
	depth    uint
	nents    uint
	nodes    [tableCapacity]nodeI
}

func (t *fullTable) copy() tableI {
	var nt = new(fullTable)
	nt.hashPath = t.hashPath
	nt.depth = t.depth
	nt.nents = t.nents
	nt.nodes = t.nodes
	return nt
}

func (t *fullTable) deepCopy() tableI {
	var nt = new(fullTable)
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

func createRootFullTable(lf leafI) tableI {
	var idx = lf.Hash30().Index(0)

	var ft = new(fullTable)
	//ft.hashPath = 0
	//ft.depth = 0
	//ft.nents = 0
	ft.set(idx, lf)

	return ft
}

func createFullTable(depth uint, leaf1 leafI, leaf2 *flatLeaf) tableI {
	if depth < 1 {
		log.Panic("createFullTable(): depth < 1")
	}
	var hp1 = leaf1.Hash30() & key.HashPathMask30(depth-1)
	var hp2 = leaf2.Hash30() & key.HashPathMask30(depth-1)
	if hp1 != hp2 {
		log.Panic("newCompressedTable(): hp1,%s != hp2,%s",
			hp1.HashPathString(depth), hp2.HashPathString(depth))
	}
	//for d := uint(0); d < depth; d++ {
	//	if leaf1.Hash30().Index(d) != leaf2.Hash30().Index(d) {
	//		log.Panicf("createFullTable(): leaf1.Hash30().Index(%d) != leaf2.Hash30().Index(%d)", d, d)
	//	}
	//}

	var retTable = new(fullTable)
	//retTable.hashPath = leaf1.Hash30() & key.HashPathMask30(depth-1)
	retTable.hashPath = leaf1.Hash30().HashPath(depth)
	retTable.depth = depth

	var idx1 = leaf1.Hash30().Index(depth)
	var idx2 = leaf2.Hash30().Index(depth)
	if idx1 != idx2 {
		retTable.insert(idx1, leaf1)
		retTable.insert(idx2, leaf2)
	} else { //idx1 == idx2
		var node nodeI
		if depth == maxDepth {
			node = newCollisionLeaf(append(leaf1.keyVals(), leaf2.keyVals()...))
		} else {
			node = createFullTable(depth+1, leaf1, leaf2)
		}
		retTable.insert(idx1, node)
	}

	return retTable
}

func upgradeToFullTable(
	hashPath key.HashVal30,
	depth uint,
	ents []tableEntry,
) *fullTable {
	var ft = new(fullTable)
	ft.hashPath = hashPath
	ft.depth = depth
	ft.nents = uint(len(ents))

	for _, ent := range ents {
		ft.nodes[ent.idx] = ent.node
	}

	return ft
}

func (t *fullTable) Hash30() key.HashVal30 {
	return t.hashPath
}

func (t *fullTable) String() string {
	return fmt.Sprintf("fullTable{hashPath=%s, depth=%d, nentries()=%d}",
		t.hashPath, t.depth, t.nentries())
}

func (t *fullTable) LongString(indent string, depth uint) string {
	var strs = make([]string, 3+t.nentries())

	strs[0] = indent + "fullTable{"
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

func (t *fullTable) nentries() uint {
	return t.nents
}

func (t *fullTable) entries() []tableEntry {
	var n = t.nentries()
	var ents = make([]tableEntry, n)
	var i, j uint
	for i, j = 0, 0; j < n && i < tableCapacity; i++ {
		if t.nodes[i] != nil {
			ents[j] = tableEntry{i, t.nodes[i]}
			j++
		}
	}
	return ents
}

func (t *fullTable) get(idx uint) nodeI {
	return t.nodes[idx]
}

func (t *fullTable) set(idx uint, nn nodeI) {
	if nn != nil && t.nodes[idx] == nil {
		t.nents++
	} else if nn == nil && t.nodes[idx] != nil {
		t.nents--
	}
	t.nodes[idx] = nn

	return
}

func (t *fullTable) insert(idx uint, n nodeI) {
	if t.nodes[idx] != nil {
		panic("t.insert(idx, n) where idx slot is NOT empty; this should be a replace")
	}
	t.nodes[idx] = n
	t.nents++
}

func (t *fullTable) replace(idx uint, n nodeI) {
	if t.nodes[idx] == nil {
		panic("t.replace(idx, n) where idx slot is empty; this should be an insert")
	}
	t.nodes[idx] = n
}

func (t *fullTable) remove(idx uint) {
	if t.nodes[idx] == nil {
		panic("t.remove(idx) where idx slot is already empty")
	}
	t.nodes[idx] = nil
	t.nents--
}
