package hamt32

import (
	"fmt"
	"strings"

	"github.com/lleo/go-hamt-key"
)

type fullTable struct {
	hashPath key.HashVal30
	nents    uint
	nodes    [tableCapacity]nodeI
}

func newRootFullTable(depth uint, hashPath key.HashVal30, lf leafI) tableI {
	var idx = lf.Hash30().Index(depth)

	var ft = new(fullTable)
	ft.hashPath = hashPath & key.HashPathMask30(depth)
	//ft.nents = 0
	ft.set(idx, lf)

	return ft
}

func newFullTable(depth uint, hashPath key.HashVal30, leaf1 leafI, leaf2 *flatLeaf) tableI {
	var retTable = new(fullTable)
	retTable.hashPath = hashPath & key.HashPathMask30(depth)

	var curTable = retTable
	var d uint
	for d = depth; d <= maxDepth; d++ {
		var idx1 = leaf1.Hash30().Index(d)
		var idx2 = leaf2.Hash30().Index(d)

		if idx1 != idx2 {
			curTable.set(idx1, leaf1)
			curTable.set(idx2, leaf2)

			break
		}
		// idx1 == idx2 ...

		var newTable = new(fullTable)

		hashPath = hashPath.BuildHashPath(idx1, d)
		newTable.hashPath = hashPath

		curTable.set(idx1, newTable)

		curTable = newTable
	}
	// We either BREAK out of loops,
	// OR we hit d > maxDepth.
	if d > maxDepth {
		var idx = leaf1.Hash30().Index(maxDepth)
		var kvs = append(leaf1.keyVals(), leaf2.keyVals()...)
		var leaf = newCollisionLeaf(kvs)
		curTable.set(idx, leaf)
	}

	return retTable
}

func upgradeToFullTable(hashPath key.HashVal30, ents []tableEntry) *fullTable {
	var ft = new(fullTable)
	ft.hashPath = hashPath
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
	return fmt.Sprintf("fullTable{hashPath=%s, nentries()=%d}", t.hashPath, t.nentries())
}

func (t *fullTable) LongString(indent string, depth uint) string {
	var strs = make([]string, 3+len(t.nodes))

	strs[0] = indent + "fullTable{"
	strs[1] = indent + fmt.Sprintf("\tnents=%d,", t.nents)

	for i, n := range t.nodes {
		if t.nodes[i] == nil {
			strs[2+i] = indent + fmt.Sprintf("\tnodes[%d]: nil", i)
		} else {
			if t, isTable := t.nodes[i].(tableI); isTable {
				strs[2+i] = indent + fmt.Sprintf("\tnodes[%d]:\n%s", i, t.LongString(indent+"\t", depth+1))
			} else {
				strs[2+i] = indent + fmt.Sprintf("\tnodes[%d]: %s", i, n)
			}
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
