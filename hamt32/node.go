package hamt32

import (
	"fmt"

	"github.com/lleo/go-hamt/key"
)

type nodeI interface {
	Hash30() key.HashVal30
	String() string
}

type leafI interface {
	nodeI

	get(key key.Key) (interface{}, bool)
	put(key key.Key, val interface{}) (leafI, bool)
	del(key key.Key) (leafI, interface{}, bool)
	keyVals() []key.KeyVal
}

type tableI interface {
	nodeI

	copy() tableI
	deepCopy() tableI

	LongString(indent string, depth uint) string

	nentries() uint
	entries() []tableEntry

	get(idx uint) nodeI

	//occupied(idx uint) bool

	//set(idx uint, entry nodeI)
	insert(idx uint, n nodeI)
	replace(idx uint, n nodeI)
	remove(idx uint)
}

type tableEntry struct {
	idx  uint
	node nodeI
}

func (ent tableEntry) String() string {
	return fmt.Sprintf("tableEntry{idx:%d, node:%s}", ent.idx, ent.node.String())
}
