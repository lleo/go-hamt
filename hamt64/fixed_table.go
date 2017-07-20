package hamt64

import (
	"fmt"
	"strings"
)

type fixedTable struct {
	nodes    [IndexLimit]nodeI // 512; 32*16
	depth    uint              // 8; amd64
	nents    uint              // 8; amd64
	hashPath hashVal           // 8
}

func (t *fixedTable) copy() tableI {
	var nt = new(fixedTable)
	nt = t
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
	//for i, n := range t.nodes {
	//	switch x := n.(type) {
	//	case tableI:
	//		nt.nodes[i] = x.deepCopy()
	//	default:
	//		nt.nodes[i] = x
	//	}
	//}
	return nt
}

//func createRootFixedTable(lf leafI) tableI {
//	var idx = lf.Hash().Index(0)
//
//	var ft = new(fixedTable)
//	//ft.hashPath = 0
//	//ft.depth = 0
//	//ft.nents = 0
//	ft.insert(idx, lf)
//
//	return ft
//}

func createFixedTable(depth uint, leaf1 leafI, leaf2 *flatLeaf) tableI {
	_ = assertOn && assertf(depth > 0, "createFixedTable(): depth,%d < 1", depth)

	_ = assertOn && assertf(
		leaf1.Hash().hashPath(depth) == leaf2.Hash().hashPath(depth),
		"createFixedTable(): hp1,%s != hp2,%s",
		leaf1.Hash().hashPath(depth),
		leaf2.Hash().hashPath(depth))

	var retTable = new(fixedTable)
	retTable.hashPath = leaf1.Hash().hashPath(depth)
	retTable.depth = depth

	var idx1 = leaf1.Hash().Index(depth)
	var idx2 = leaf2.Hash().Index(depth)
	if idx1 != idx2 {
		retTable.insert(idx1, leaf1)
		retTable.insert(idx2, leaf2)
	} else { //idx1 == idx2
		var node nodeI
		if depth == maxDepth {
			node = newCollisionLeaf(append(leaf1.keyVals(), leaf2.keyVals()...))
		} else {
			//log.Printf("createFixedTable(depth=%d, ...) recursing\n", depth+1)
			node = createFixedTable(depth+1, leaf1, leaf2)
		}
		retTable.insert(idx1, node)
	}

	return retTable
}

func upgradeToFixedTable(
	hashPath hashVal,
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
func (t *fixedTable) Hash() hashVal {
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

func (t *fixedTable) insert(idx uint, n nodeI) {
	_ = assertOn && assert(t.nodes[idx] == nil,
		"t.insert(idx, n) where idx slot is NOT empty; this should be a replace")

	t.nodes[idx] = n
	t.nents++
}

func (t *fixedTable) replace(idx uint, n nodeI) {
	_ = assertOn && assert(t.nodes[idx] != nil,
		"t.replace(idx, n) where idx slot is empty; this should be an insert")

	t.nodes[idx] = n
}

func (t *fixedTable) remove(idx uint) {
	_ = assertOn && assert(t.nodes[idx] != nil,
		"t.remove(idx) where idx slot is already empty")

	t.nodes[idx] = nil
	t.nents--
}

func (t *fixedTable) visit(fn visitFn, depth uint) uint {
	fn(t)

	var maxDepth = depth + 1
	for _, n := range t.nodes {
		if n == nil {
			fn(n)
		} else {
			var md = n.visit(fn, depth+1)
			if md > maxDepth {
				maxDepth = md
			}
		}
	}

	return maxDepth
}
