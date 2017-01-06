package hamt32

import (
	"fmt"
	"log"
	"strings"
)

// COMPRESSED_TABLE_INIT_CAP constant sets the default capacity of a new compressedTable.
const COMPRESSED_TABLE_INIT_CAP int = 8

type compressedTable struct {
	hashPath uint32
	nodeMap  uint32
	nodes    []nodeI
}

func newRootCompressedTable(depth uint, hashPath uint32, lf leafI) tableI {
	var idx = index(lf.Hash30(), depth)

	var ct = new(compressedTable)
	ct.hashPath = hashPath & hashPathMask(depth)
	ct.nodeMap = uint32(1 << idx)
	ct.nodes = make([]nodeI, 1, COMPRESSED_TABLE_INIT_CAP)
	ct.nodes[0] = lf

	return ct
}

func newCompressedTable(depth uint, hashPath uint32, leaf1 leafI, leaf2 *flatLeaf) tableI {
	var retTable = new(compressedTable)
	retTable.hashPath = hashPath & hashPathMask(depth)

	var curTable = retTable
	var d uint
	for d = depth; d <= MAXDEPTH; d++ {
		var idx1 = index(leaf1.Hash30(), d)
		var idx2 = index(leaf2.Hash30(), d)

		if idx1 != idx2 {
			//curTable.nodes = make([]nodeI, 0, COMPRESSED_TABLE_INIT_CAP)
			//curTable.set(idx1, leaf1)
			//curTable.set(idx2, leaf2)

			// This is faster
			curTable.nodes = make([]nodeI, 2, COMPRESSED_TABLE_INIT_CAP)
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

		curTable.nodes = make([]nodeI, 0, COMPRESSED_TABLE_INIT_CAP)

		var newTable = new(compressedTable)

		hashPath = buildHashPath(hashPath, idx1, d)
		newTable.hashPath = hashPath

		curTable.set(idx1, newTable)

		curTable = newTable
	}
	// We either BREAK out of the loop,
	// OR we hit d > MAXDEPTH.
	if d > MAXDEPTH {
		// leaf1.Hash30() == leaf2.Hash30()
		log.Printf("newCompressedTable: d > MAXDEPTH branch taken.")
		var idx = index(leaf1.Hash30(), MAXDEPTH)
		var kvs = append(leaf1.keyVals(), leaf2.keyVals()...)
		var leaf = newCollisionLeaf(kvs)
		curTable.set(idx, leaf)
	}

	return retTable
}

func nodeMapString(nodeMap uint32) string {
	var strs = make([]string, 4)

	var top2 = nodeMap >> 30
	strs[0] = fmt.Sprintf("%02b", top2)

	const tenBitMask uint32 = 1<<10 - 1
	for i := uint(0); i < 3; i++ {
		var tenBitVal = (nodeMap & (tenBitMask << (i * 10))) >> (i * 10)
		strs[3-i] = fmt.Sprintf("%010b", tenBitVal)
	}

	return strings.Join(strs, " ")
}

// downgradeToCompressedTable() converts fullTable structs that have less than
// TABLE_CAPACITY/2 tableEntry's. One important thing we know is that none of
// the entries will collide with another.
//
// The ents []tableEntry slice is guaranteed to be in order from lowest idx to
// highest. tableI.entries() also adhears to this contract.
func downgradeToCompressedTable(hashPath uint32, ents []tableEntry) *compressedTable {
	var nt = new(compressedTable)
	nt.hashPath = hashPath
	//nt.nodeMap = 0
	nt.nodes = make([]nodeI, len(ents))

	for i := 0; i < len(ents); i++ {
		var ent = ents[i]
		var nodeBit = uint32(1 << ent.idx)
		nt.nodeMap |= nodeBit
		nt.nodes[i] = ent.node
	}

	return nt
}

func (t *compressedTable) Hash30() uint32 {
	return t.hashPath
}

func (t *compressedTable) String() string {
	return fmt.Sprintf("compressedTable{hashPath:%s, nentries()=%d}",
		hash30String(t.hashPath), t.nentries())
}

func (t *compressedTable) LongString(indent string, depth uint) string {
	var strs = make([]string, 3+len(t.nodes))

	strs[0] = indent + fmt.Sprintf("compressedTable{hashPath=%s, nentries()=%d,", hashPathString(t.hashPath, depth), t.nentries())

	strs[1] = indent + "\tnodeMap=" + nodeMapString(t.nodeMap) + ","

	for i, n := range t.nodes {
		if t, isTable := n.(tableI); isTable {
			strs[2+i] = indent + fmt.Sprintf("\tt.nodes[%d]:\n%s", i, t.LongString(indent+"\t", depth+1))
		} else {
			strs[2+i] = indent + fmt.Sprintf("\tt.nodes[%d]: %s", i, n)
		}
	}

	strs[len(strs)-1] = indent + "}"

	return strings.Join(strs, "\n")
}

func (t *compressedTable) nentries() uint {
	return bitCount32(t.nodeMap)
}

func (t *compressedTable) entries() []tableEntry {
	var n = t.nentries()
	var ents = make([]tableEntry, n)

	for i, j := uint(0), uint(0); i < TABLE_CAPACITY; i++ {
		var nodeBit = uint32(1 << i)

		if (t.nodeMap & nodeBit) > 0 {
			ents[j] = tableEntry{i, t.nodes[j]}
			j++
		}
	}

	return ents
}

func (t *compressedTable) get(idx uint) nodeI {
	var nodeBit = uint32(1 << idx)

	if (t.nodeMap & nodeBit) == 0 {
		return nil
	}

	var m = uint32(1<<idx) - 1
	var i = bitCount32(t.nodeMap & m)

	return t.nodes[i]
}

func (t *compressedTable) set(idx uint, nn nodeI) {
	var nodeBit = uint32(1 << idx)
	var bitMask = nodeBit - 1
	var i = bitCount32(t.nodeMap & bitMask)

	if nn != nil {
		if (t.nodeMap & nodeBit) == 0 {
			t.nodeMap |= nodeBit
			t.nodes = append(t.nodes[:i], append([]nodeI{nn}, t.nodes[i:]...)...)
		} else {
			t.nodes[i] = nn
		}
	} else /* if nn == nil */ {
		if (t.nodeMap & nodeBit) > 0 {
			t.nodeMap &^= nodeBit
			t.nodes = append(t.nodes[:i], t.nodes[i+1:]...)
		} /* else {
			// do nothing
		} */
	}
	return
}

//POPCNT Implementation
// copied from https://github.com/jddixon/xlUtil_go/blob/master/popCount.go
//  was MIT License

const (
	octo_fives  = uint32(0x55555555)
	octo_threes = uint32(0x33333333)
	octo_ones   = uint32(0x01010101)
	octo_fs     = uint32(0x0f0f0f0f)
)

func bitCount32(n uint32) uint {
	n = n - ((n >> 1) & octo_fives)
	n = (n & octo_threes) + ((n >> 2) & octo_threes)
	return uint((((n + (n >> 4)) & octo_fs) * octo_ones) >> 24)
}
