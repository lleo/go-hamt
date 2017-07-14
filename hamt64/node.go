package hamt64

import "fmt"

type nodeI interface {
	Hash() hashVal
	String() string
	visit(fn visitFn, arg interface{}, depth uint) uint
}

type leafI interface {
	nodeI

	get(key *iKey) (interface{}, bool)
	put(key *iKey, val interface{}) (leafI, bool)
	del(key *iKey) (leafI, interface{}, bool)
	keyVals() []KeyVal
}

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
}

type tableEntry struct {
	idx  uint
	node nodeI
}

func (ent tableEntry) String() string {
	return fmt.Sprintf("tableEntry{idx:%d, node:%s}", ent.idx, ent.node.String())
}
