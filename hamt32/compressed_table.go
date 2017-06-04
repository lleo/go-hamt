package hamt32

import (
	"fmt"
	"log"
	"strings"

	"github.com/lleo/go-hamt-key"
)

// compressedTableInitCap constant sets the default capacity of a new
// compressedTable.
const compressedTableInitCap int = 8

type compressedTable struct {
	hashPath key.HashVal30
	nodeMap  uint32
	nodes    []nodeI
}

func (t *compressedTable) copy() tableI {
	var nt = new(compressedTable)
	nt.hashPath = t.hashPath
	nt.nodeMap = t.nodeMap
	nt.nodes = append(nt.nodes, t.nodes...)
	return nt
}

func createRootCompressedTable(lf leafI) tableI {
	var idx = lf.Hash30().Index(0)

	var ct = new(compressedTable)
	//ct.hashPath = 0
	ct.nodeMap = uint32(1 << idx)
	ct.nodes = make([]nodeI, 1, compressedTableInitCap)
	ct.nodes[0] = lf

	return ct
}

func createCompressedTable(depth uint, leaf1 leafI, leaf2 *flatLeaf) tableI {
	if depth < 1 {
		log.Panic("createCompressedTable(): depth < 1")
	}
	var hp1 = leaf1.Hash30() & key.HashPathMask30(depth-1)
	var hp2 = leaf2.Hash30() & key.HashPathMask30(depth-1)
	if hp1 != hp2 {
		log.Panic("createCompressedTable(): hp1,%s != hp2,%s",
			hp1.HashPathString(depth), hp2.HashPathString(depth))
	}
	//for d := uint(0); d < depth; d++ {
	//	if leaf1.Hash30().Index(d) != leaf2.Hash30().Index(d) {
	//		log.Panicf("createCompressedTable(): leaf1.Hash30().Index(%d) != leaf2.Hash30().Index(%d)", d, d)
	//	}
	//}

	var retTable = new(compressedTable)
	//retTable.hashPath = leaf1.Hash30() & key.HashPathMask30(depth-1)
	retTable.hashPath = leaf1.Hash30().HashPath(depth)
	//retTable.nodeMap = 0
	retTable.nodes = make([]nodeI, 0, compressedTableInitCap)

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
			node = createCompressedTable(depth+1, leaf1, leaf2)
		}
		retTable.insert(idx1, node)
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

// downgradeToCompressedTable() converts fullTable structs that have less than or equal
// to downgradeThreshold tableEntry's. One important thing we know is that none of
// the entries will collide with another.
//
// The ents []tableEntry slice is guaranteed to be in order from lowest idx to
// highest. tableI.entries() also adhears to this contract.
func downgradeToCompressedTable(hashPath key.HashVal30, ents []tableEntry) *compressedTable {
	var nt = new(compressedTable)
	nt.hashPath = hashPath
	//nt.nodeMap = 0
	nt.nodes = make([]nodeI, len(ents), compressedTableInitCap)

	for i := 0; i < len(ents); i++ {
		var ent = ents[i]
		var nodeBit = uint32(1 << ent.idx)
		nt.nodeMap |= nodeBit
		nt.nodes[i] = ent.node
	}

	return nt
}

func (t *compressedTable) Hash30() key.HashVal30 {
	return t.hashPath
}

func (t *compressedTable) String() string {
	return fmt.Sprintf("compressedTable{hashPath:%s, nentries()=%d}",
		t.hashPath, t.nentries())
}

func (t *compressedTable) LongString(indent string, depth uint) string {
	var strs = make([]string, 3+len(t.nodes))

	strs[0] = indent + fmt.Sprintf("compressedTable{hashPath=%s, nentries()=%d,", t.hashPath.HashPathString(depth+1), t.nentries())

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
	return uint(len(t.nodes))
	//return bitCount32(t.nodeMap)
}

func (t *compressedTable) entries() []tableEntry {
	var n = t.nentries()
	var ents = make([]tableEntry, n)

	for i, j := uint(0), uint(0); i < tableCapacity; i++ {
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

	var bitMask = nodeBit - 1
	var i = bitCount32(t.nodeMap & bitMask)

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

func (t *compressedTable) insert(idx uint, n nodeI) {
	var nodeBit = uint32(1 << idx)

	if (t.nodeMap & nodeBit) > 0 {
		panic("t.insert(idx, n) where idx slot is NOT empty; this should be a replace")
	}

	var bitMask = nodeBit - 1
	var i = bitCount32(t.nodeMap & bitMask)
	if i == uint(len(t.nodes)) {
		t.nodes = append(t.nodes, n)
	} else {
		t.nodes = append(t.nodes[:i], append([]nodeI{n}, t.nodes[i:]...)...)
	}
	t.nodeMap |= nodeBit
}

func (t *compressedTable) replace(idx uint, n nodeI) {
	var nodeBit = uint32(1 << idx)

	if (t.nodeMap & nodeBit) == 0 {
		panic("t.replace(idx, n) where idx slot is empty; this should be an insert")
	}

	var bitMask = nodeBit - 1
	var i = bitCount32(t.nodeMap & bitMask)
	t.nodes[i] = n
}

func (t *compressedTable) remove(idx uint) {
	var nodeBit = uint32(1 << idx)

	if (t.nodeMap & nodeBit) == 0 {
		panic("t.remove(idx) where idx slot is already empty")
	}

	var bitMask = nodeBit - 1
	var i = bitCount32(t.nodeMap & bitMask)
	if int(i) == len(t.nodes)-1 {
		t.nodes = t.nodes[:i]
	} else {
		t.nodes = append(t.nodes[:i], t.nodes[i+1:]...)
	}
	t.nodeMap &^= nodeBit
}
