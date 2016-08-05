package hamt64

import (
	"fmt"
	"strings"
)

type fullTable64 struct {
	hashPath uint64
	nodeMap  uint64
	nodes    [TABLE_CAPACITY64]node64I
}

func UpgradeToFullTable64(hashPath uint64, tabEnts []tableEntry64) table64I {
	var ft = new(fullTable64)
	ft.hashPath = hashPath
	//ft.nodeMap = 0 //unnecessary

	for _, ent := range tabEnts {
		var nodeBit = uint64(1 << ent.idx)
		ft.nodeMap |= nodeBit
		ft.nodes[ent.idx] = ent.node
	}

	return ft
}

func (t *fullTable64) hash60() uint64 {
	return t.hashPath
}

func (t *fullTable64) String() string {
	return fmt.Sprintf("fullTable64{hashPath=%s, nentries()=%d}", hash60String(t.hashPath), t.nentries())
}

func (t *fullTable64) LongString(indent string, depth uint) string {
	var strs = make([]string, 3+len(t.nodes))

	strs[1] = indent + "\tnodeMap=" + nodeMapString(t.nodeMap) + ","

	for i, n := range t.nodes {
		if t.nodes[i] == nil {
			strs[2+i] = indent + fmt.Sprintf("\tt.nodes[%d]: nil", i)
		} else {
			if t, isTable := t.nodes[i].(table64I); isTable {
				strs[2+i] = indent + fmt.Sprintf("\tt.nodes[%d]:\n%s", i, t.LongString(indent+"\t", depth+1))
			} else {
				strs[2+i] = indent + fmt.Sprintf("\tt.nodes[%d]: %s", i, n)
			}
		}
	}

	strs[len(strs)-1] = indent + "}"

	return strings.Join(strs, "\n")
}

func (t *fullTable64) nentries() uint {
	return BitCount64(t.nodeMap)
}

func (t *fullTable64) entries() []tableEntry64 {
	var n = t.nentries()
	var ents = make([]tableEntry64, n)
	for i, j := uint(0), 0; i < TABLE_CAPACITY64; i++ {
		var nodeBit = uint64(1 << i)
		if (t.nodeMap & nodeBit) > 0 {
			ents[j] = tableEntry64{i, t.nodes[i]}
			j++
		}
	}
	return ents
}

func (t *fullTable64) get(idx uint) node64I {
	var nodeBit = uint64(1 << idx)

	if (t.nodeMap & nodeBit) == 0 {
		return nil
	}

	return t.nodes[idx]
}

func (t *fullTable64) set(idx uint, nn node64I) {
	var nodeBit = uint64(1 << idx)

	if nn != nil {
		t.nodeMap |= nodeBit
		t.nodes[idx] = nn
	} else /* if nn == nil */ {
		t.nodeMap &^= nodeBit
		t.nodes[idx] = nn
	}

	return
}
