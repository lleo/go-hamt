package hamt32

import (
	"fmt"
)

type nodeI interface {
	Hash30() uint32
	String() string
}

type tableI interface {
	nodeI

	LongString(indent string, depth uint) string

	nentries() uint
	entries() []tableEntry

	get(idx uint) nodeI
	set(idx uint, entry nodeI)
}

type tableEntry struct {
	idx  uint
	node nodeI
}

func (ent tableEntry) String() string {
	return fmt.Sprintf("tableEntry{idx:%d, node:%s}", ent.idx, ent.node.String())
}

func (h *Hamt) newRootTable(depth uint, hashPath uint32, lf leafI) tableI {
	if h.fullinit {
		return newRootFullTable(depth, hashPath, lf)
	}
	return newRootCompressedTable(depth, hashPath, lf)
}

func (h *Hamt) newTable(depth uint, hashPath uint32, leaf1 leafI, leaf2 *flatLeaf) tableI {
	if h.fullinit {
		return newFullTable(depth, hashPath, leaf1, leaf2)
	}
	return newCompressedTable(depth, hashPath, leaf1, leaf2)
}
