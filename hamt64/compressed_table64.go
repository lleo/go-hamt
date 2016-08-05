package hamt64

import (
	"fmt"
	"log"
	"strings"
)

type compressedTable64 struct {
	hashPath uint64
	nodeMap  uint64
	nodes    []node64I
}

func newCompressedTable64(depth uint, hashPath uint64, lf leaf64I) table64I {
	var idx = index(hashPath, depth)

	var ct = new(compressedTable64)
	ct.hashPath = hashPath & hashPathMask(depth)
	ct.nodeMap = uint64(1 << idx)
	ct.nodes = make([]node64I, 1)
	ct.nodes[0] = lf

	return ct
}

func newCompressedTable64_2(depth uint, hashPath uint64, leaf1 leaf64I, leaf2 *flatLeaf64) table64I {
	var retTable = new(compressedTable64)
	retTable.hashPath = hashPath & hashPathMask(depth)

	var curTable = retTable
	var d uint
	for d = depth; d < DEPTHLIMIT64; d++ {
		var idx1 = index(leaf1.hash60(), d)
		var idx2 = index(leaf2.hash60(), d)

		if idx1 != idx2 {
			curTable.nodes = make([]node64I, 2)

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

		curTable.nodes = make([]node64I, 1)

		var newTable = new(compressedTable64)

		hashPath = buildHashPath(hashPath, idx1, d)
		newTable.hashPath = hashPath

		curTable.nodeMap = uint64(1 << idx1) //Set the idx1'th bit
		curTable.nodes[0] = newTable

		curTable = newTable
	}
	// We either BREAK out of the loop,
	// OR we hit d = DEPTHLIMIT64.
	if d == DEPTHLIMIT64 {
		// leaf1.hashcode() == leaf2.hashcode()
		var idx = index(leaf1.hash60(), d)
		hashPath = buildHashPath(hashPath, idx, d)
		var kvs = append(leaf1.keyVals(), leaf2.keyVals()...)
		var leaf = newCollisionLeaf64(hashPath, kvs)
		curTable.set(idx, leaf)
	}

	return retTable
}

// DowngradeToCompressedTable64() converts fullTable64 structs that have less than
// TABLE_CAPACITY64/2 tableEntry64's. One important thing we know is that none of
// the entries will collide with another.
//
// The ents []tableEntry64 slice is guaranteed to be in order from lowest idx to
// highest. table64I.entries() also adhears to this contract.
func DowngradeToCompressedTable64(hashPath uint64, ents []tableEntry64) *compressedTable64 {
	var nt = new(compressedTable64)
	nt.hashPath = hashPath
	//nt.nodeMap = 0
	nt.nodes = make([]node64I, len(ents))

	for i := 0; i < len(ents); i++ {
		var ent = ents[i]
		var nodeBit = uint64(1 << ent.idx)
		nt.nodeMap |= nodeBit
		nt.nodes[i] = ent.node
	}

	return nt
}

func (t *compressedTable64) hash60() uint64 {
	return t.hashPath
}

func (t *compressedTable64) String() string {
	return fmt.Sprintf("compressedTable64{hashPath:%s, nentries()=%d}",
		hash60String(t.hashPath), t.nentries())
}

func (t *compressedTable64) LongString(indent string, depth uint) string {
	var strs = make([]string, 3+len(t.nodes))

	strs[0] = indent + fmt.Sprintf("compressedTable64{hashPath=%s, nentries()=%d,", hashPathString(t.hashPath, depth), t.nentries())

	strs[1] = indent + "\tnodeMap=" + nodeMapString(t.nodeMap) + ","

	for i, n := range t.nodes {
		if t, isTable := n.(table64I); isTable {
			strs[2+i] = indent + fmt.Sprintf("\tt.nodes[%d]:\n%s", i, t.LongString(indent+"\t", depth+1))
		} else {
			strs[2+i] = indent + fmt.Sprintf("\tt.nodes[%d]: %s", i, n)
		}
	}

	strs[len(strs)-1] = indent + "}"

	return strings.Join(strs, "\n")
}

func (t *compressedTable64) nentries() uint {
	return BitCount64(t.nodeMap)
}

func (t *compressedTable64) entries() []tableEntry64 {
	var n = t.nentries()
	var ents = make([]tableEntry64, n)

	for i, j := uint(0), uint(0); i < TABLE_CAPACITY64; i++ {
		var nodeBit = uint64(1 << i)

		if (t.nodeMap & nodeBit) > 0 {
			ents[j] = tableEntry64{i, t.nodes[j]}
			j++
		}
	}

	return ents
}

func (t *compressedTable64) get(idx uint) node64I {
	var nodeBit = uint64(1 << idx)

	if (t.nodeMap & nodeBit) == 0 {
		return nil
	}

	var m = uint64(1<<idx) - 1
	var i = BitCount64(t.nodeMap & m)

	return t.nodes[i]
}

func (t *compressedTable64) set(idx uint, nn node64I) {
	var nodeBit = uint64(1 << idx)
	var bitMask = nodeBit - 1
	var i = BitCount64(t.nodeMap & bitMask)

	if nn != nil {
		if (t.nodeMap & nodeBit) == 0 {
			t.nodeMap |= nodeBit
			t.nodes = append(t.nodes[:i], append([]node64I{nn}, t.nodes[i:]...)...)
		} else {
			t.nodes[i] = nn
		}
	} else /* if nn == nil */ {
		if (t.nodeMap & nodeBit) > 0 {
			t.nodeMap &^= nodeBit
			t.nodes = append(t.nodes[:i], t.nodes[i+1:]...)
		} else if (t.nodeMap & nodeBit) == 0 {
			log.Panicf("compressedTable64.set(%02d, nil): when no node was set here in the first place", idx)
			// do nothing
		}
	}
	return
}
