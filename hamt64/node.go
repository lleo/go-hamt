package hamt64

import (
	"fmt"

	"github.com/lleo/go-hamt-key"
)

type nodeI interface {
	Hash60() key.HashVal60
	String() string
}

type leafI interface {
	nodeI
	// from nodeI in table.go
	// Hash60() uint64
	// String() string

	get(key key.Key) (interface{}, bool)
	put(key key.Key, val interface{}) (leafI, bool)
	del(key key.Key) (interface{}, leafI, bool)
	keyVals() []key.KeyVal
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

func (h *Hamt) newRootTable(depth uint, hashPath key.HashVal60, lf leafI) tableI {
	if h.fullinit {
		return newRootFullTable(depth, hashPath, lf)
	}
	return newRootCompressedTable(depth, hashPath, lf)
}

func (h *Hamt) newTable(depth uint, hashPath key.HashVal60, leaf1 leafI, leaf2 *flatLeaf) tableI {
	if h.fullinit {
		return newFullTable(depth, hashPath, leaf1, leaf2)
	}
	return newCompressedTable(depth, hashPath, leaf1, leaf2)
}
