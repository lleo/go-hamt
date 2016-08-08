package hamt32

import (
	"fmt"
	"strings"
)

type fullTable struct {
	hashPath uint32
	nodeMap  uint32
	nodes    [TABLE_CAPACITY]nodeI
}

func UpgradeToFullTable(hashPath uint32, tabEnts []tableEntry) tableI {
	var ft = new(fullTable)
	ft.hashPath = hashPath
	//ft.nodeMap = 0 //unnecessary

	for _, ent := range tabEnts {
		var nodeBit = uint32(1 << ent.idx)
		ft.nodeMap |= nodeBit
		ft.nodes[ent.idx] = ent.node
	}

	return ft
}

func (t *fullTable) hash30() uint32 {
	return t.hashPath
}

func (t *fullTable) String() string {
	return fmt.Sprintf("fullTable{hashPath=%s, nentries()=%d}", hash30String(t.hashPath), t.nentries())
}

func (t *fullTable) LongString(indent string, depth uint) string {
	var strs = make([]string, 3+len(t.nodes))

	strs[1] = indent + "\tnodeMap=" + nodeMapString(t.nodeMap) + ","

	for i, n := range t.nodes {
		if t.nodes[i] == nil {
			strs[2+i] = indent + fmt.Sprintf("\tt.nodes[%d]: nil", i)
		} else {
			if t, isTable := t.nodes[i].(tableI); isTable {
				strs[2+i] = indent + fmt.Sprintf("\tt.nodes[%d]:\n%s", i, t.LongString(indent+"\t", depth+1))
			} else {
				strs[2+i] = indent + fmt.Sprintf("\tt.nodes[%d]: %s", i, n)
			}
		}
	}

	strs[len(strs)-1] = indent + "}"

	return strings.Join(strs, "\n")
}

func (t *fullTable) nentries() uint {
	return BitCount32(t.nodeMap)
}

func (t *fullTable) entries() []tableEntry {
	var n = t.nentries()
	var ents = make([]tableEntry, n)
	for i, j := uint(0), 0; i < TABLE_CAPACITY; i++ {
		var nodeBit = uint32(1 << i)
		if (t.nodeMap & nodeBit) > 0 {
			ents[j] = tableEntry{i, t.nodes[i]}
			j++
		}
	}
	return ents
}

func (t *fullTable) get(idx uint) nodeI {
	var nodeBit = uint32(1 << idx)

	if (t.nodeMap & nodeBit) == 0 {
		return nil
	}

	return t.nodes[idx]
}

func (t *fullTable) set(idx uint, nn nodeI) {
	var nodeBit = uint32(1 << idx)

	if nn != nil {
		t.nodeMap |= nodeBit
		t.nodes[idx] = nn
	} else /* if nn == nil */ {
		t.nodeMap &^= nodeBit
		t.nodes[idx] = nn
	}

	return
}
