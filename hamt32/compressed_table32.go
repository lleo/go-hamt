package hamt32

import (
	"fmt"
	"log"
	"strings"
)

type compressedTable32 struct {
	hashPath uint32
	nodeMap  uint32
	nodes    []node32I
}

func newCompressedTable32(depth uint, hashPath uint32, lf leaf32I) table32I {
	var idx = index(hashPath, depth)

	var ct = new(compressedTable32)
	ct.hashPath = hashPath & hashPathMask(depth)
	ct.nodeMap = uint32(1 << idx)
	ct.nodes = make([]node32I, 1)
	ct.nodes[0] = lf

	return ct
}

func newCompressedTable32_2(depth uint, hashPath uint32, leaf1 leaf32I, leaf2 *flatLeaf32) table32I {
	var retTable = new(compressedTable32)
	retTable.hashPath = hashPath & hashPathMask(depth)

	var curTable = retTable
	var d uint
	for d = depth; d < DEPTHLIMIT32; d++ {
		var idx1 = index(leaf1.hash30(), d)
		var idx2 = index(leaf2.hash30(), d)

		if idx1 != idx2 {
			curTable.nodes = make([]node32I, 2)

			curTable.nodeMap |= 1 << idx1
			curTable.nodeMap |= 1 << idx2
			if idx1 < idx2 {
				curTable.nodes[0] = leaf1
				curTable.nodes[1] = leaf2
			} else {
				curTable.nodes[0] = leaf2
				curTable.nodes[1] = leaf1
			}

			break //leaving the for-loop
		}
		// idx1 == idx2 && continue

		curTable.nodes = make([]node32I, 1)

		var newTable = new(compressedTable32)

		hashPath = buildHashPath(hashPath, idx1, d)
		newTable.hashPath = hashPath

		curTable.nodeMap = uint32(1 << idx1) //Set the idx1'th bit
		curTable.nodes[0] = newTable

		curTable = newTable
	}
	// We either BREAK out of the loop,
	// OR we hit d = DEPTHLIMIT32.
	if d == DEPTHLIMIT32 {
		// leaf1.hashcode() == leaf2.hashcode()
		var idx = index(leaf1.hash30(), d)
		hashPath = buildHashPath(hashPath, idx, d)
		var kvs = append(leaf1.keyVals(), leaf2.keyVals()...)
		var leaf = newCollisionLeaf32(hashPath, kvs)
		curTable.set(idx, leaf)
	}

	return retTable
}

// DowngradeToCompressedTable32() converts fullTable32 structs that have less than
// TABLE_CAPACITY32/2 tableEntry32's. One important thing we know is that none of
// the entries will collide with another.
//
// The ents []tableEntry32 slice is guaranteed to be in order from lowest idx to
// highest. table32I.entries() also adhears to this contract.
func DowngradeToCompressedTable32(hashPath uint32, ents []tableEntry32) *compressedTable32 {
	var nt = new(compressedTable32)
	nt.hashPath = hashPath
	//nt.nodeMap = 0
	nt.nodes = make([]node32I, len(ents))

	for i := 0; i < len(ents); i++ {
		var ent = ents[i]
		var nodeBit = uint32(1 << ent.idx)
		nt.nodeMap |= nodeBit
		nt.nodes[i] = ent.node
	}

	return nt
}

func (t *compressedTable32) hash30() uint32 {
	return t.hashPath
}

func (t *compressedTable32) String() string {
	return fmt.Sprintf("compressedTable32{hashPath:%s, nentries()=%d}",
		hash30String(t.hashPath), t.nentries())
}

func (t *compressedTable32) LongString(indent string, depth uint) string {
	var strs = make([]string, 3+len(t.nodes))

	strs[0] = indent + fmt.Sprintf("compressedTable32{hashPath=%s, nentries()=%d,", hashPathString(t.hashPath, depth), t.nentries())

	strs[1] = indent + "\tnodeMap=" + nodeMapString(t.nodeMap) + ","

	for i, n := range t.nodes {
		if t, isTable := n.(table32I); isTable {
			strs[2+i] = indent + fmt.Sprintf("\tt.nodes[%d]:\n%s", i, t.LongString(indent+"\t", depth+1))
		} else {
			strs[2+i] = indent + fmt.Sprintf("\tt.nodes[%d]: %s", i, n)
		}
	}

	strs[len(strs)-1] = indent + "}"

	return strings.Join(strs, "\n")
}

func (t *compressedTable32) nentries() uint {
	return BitCount32(t.nodeMap)
}

func (t *compressedTable32) entries() []tableEntry32 {
	var n = t.nentries()
	var ents = make([]tableEntry32, n)

	for i, j := uint(0), uint(0); i < TABLE_CAPACITY32; i++ {
		var nodeBit = uint32(1 << i)

		if (t.nodeMap & nodeBit) > 0 {
			ents[j] = tableEntry32{i, t.nodes[j]}
			j++
		}
	}

	return ents
}

func (t *compressedTable32) get(idx uint) node32I {
	var nodeBit = uint32(1 << idx)

	if (t.nodeMap & nodeBit) == 0 {
		return nil
	}

	var m = uint32(1<<idx) - 1
	var i = BitCount32(t.nodeMap & m)

	return t.nodes[i]
}

func (t *compressedTable32) set(idx uint, nn node32I) {
	var nodeBit = uint32(1 << idx)
	var bitMask = nodeBit - 1
	var i = BitCount32(t.nodeMap & bitMask)

	if nn != nil {
		if (t.nodeMap & nodeBit) == 0 {
			t.nodeMap |= nodeBit
			t.nodes = append(t.nodes[:i], append([]node32I{nn}, t.nodes[i:]...)...)
		} else {
			t.nodes[i] = nn
		}
	} else /* if nn == nil */ {
		if (t.nodeMap & nodeBit) > 0 {
			t.nodeMap &^= nodeBit
			t.nodes = append(t.nodes[:i], t.nodes[i+1:]...)
		} else if (t.nodeMap & nodeBit) == 0 {
			log.Panicf("compressedTable32.set(%02d, nil): when no node was set here in the first place", idx)
			// do nothing
		}
	}
	return
}
