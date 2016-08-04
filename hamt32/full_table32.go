package hamt32

import (
	"fmt"
	"strings"
)

type fullTable32 struct {
	hashPath uint32
	nodeMap  uint32
	nodes    [TABLE_CAPACITY32]node32I
}

func UpgradeToFullTable32(hashPath uint32, tabEnts []tableEntry32) table32I {
	var ft = new(fullTable32)
	ft.hashPath = hashPath
	//ft.nodeMap = 0 //unnecessary

	for _, ent := range tabEnts {
		var nodeBit = uint32(1 << ent.idx)
		ft.nodeMap |= nodeBit
		ft.nodes[ent.idx] = ent.node
	}

	return ft
}

func (t *fullTable32) hash30() uint32 {
	return t.hashPath
}

func (t *fullTable32) String() string {
	return fmt.Sprintf("fullTable32{hashPath=%s, nentries()=%d}", hash30String(t.hashPath), t.nentries())
}

func (t *fullTable32) LongString(indent string, depth uint) string {
	var strs = make([]string, 3+len(t.nodes))

	strs[1] = indent + "\tnodeMap=" + nodeMapString(t.nodeMap) + ","

	for i, n := range t.nodes {
		if t.nodes[i] == nil {
			strs[2+i] = indent + fmt.Sprintf("\tt.nodes[%d]: nil", i)
		} else {
			if t, isTable := t.nodes[i].(table32I); isTable {
				strs[2+i] = indent + fmt.Sprintf("\tt.nodes[%d]:\n%s", i, t.LongString(indent+"\t", depth+1))
			} else {
				strs[2+i] = indent + fmt.Sprintf("\tt.nodes[%d]: %s", i, n)
			}
		}
	}

	strs[len(strs)-1] = indent + "}"

	return strings.Join(strs, "\n")
}

func (t *fullTable32) nentries() uint {
	return BitCount32(t.nodeMap)
}

func (t *fullTable32) entries() []tableEntry32 {
	var n = t.nentries()
	var ents = make([]tableEntry32, n)
	for i, j := uint(0), 0; i < TABLE_CAPACITY32; i++ {
		var nodeBit = uint32(1 << i)
		if (t.nodeMap & nodeBit) > 0 {
			ents[j] = tableEntry32{i, t.nodes[i]}
			j++
		}
	}
	return ents
}

func (t *fullTable32) get(idx uint) node32I {
	var nodeBit = uint32(1 << idx)

	if (t.nodeMap & nodeBit) == 0 {
		return nil
	}

	return t.nodes[idx]
}

func (t *fullTable32) set(idx uint, nn node32I) {
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
