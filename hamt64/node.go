package hamt64

import "fmt"

// visitFn will be passed a value for every slot in the Hamt; this includes
// leafs, tables, and nil.
type visitFn func(nodeI) bool

type nodeI interface {
	Hash() hashVal
	String() string
	visit(fn visitFn) bool
}

type leafI interface {
	nodeI

	get(key KeyI) (interface{}, bool)
	put(key KeyI, val interface{}) (leafI, bool)
	del(key KeyI) (leafI, interface{}, bool)
	keyVals() []KeyVal
}

type tableIterFunc func() nodeI

type tableI interface {
	nodeI

	copy() tableI
	deepCopy() tableI

	LongString(indent string, depth uint) string

	nentries() uint
	entries() []tableEntry

	get(idx uint) nodeI

	insert(idx uint, n nodeI)
	replace(idx uint, n nodeI)
	remove(idx uint)

	iter() tableIterFunc
}

type tableEntry struct {
	idx  uint
	node nodeI
}

func (ent tableEntry) String() string {
	return fmt.Sprintf("tableEntry{idx:%d, node:%s}", ent.idx, ent.node.String())
}
