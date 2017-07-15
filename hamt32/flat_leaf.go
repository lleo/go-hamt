package hamt32

import (
	"fmt"
)

type flatLeaf struct {
	key *iKey
	val interface{}
}

func newFlatLeaf(key *iKey, val interface{}) *flatLeaf {
	var fl = new(flatLeaf)
	fl.key = key
	fl.val = val
	return fl
}

func (l *flatLeaf) Hash() hashVal {
	return l.key.Hash()
}

func (l *flatLeaf) String() string {
	return fmt.Sprintf("flatLeaf{key: %s, val: %v}", l.key, l.val)
}

func (l *flatLeaf) get(key *iKey) (interface{}, bool) {
	if l.key.Equals(key) {
		return l.val, true
	}
	return nil, false
}

func (l *flatLeaf) put(k *iKey, v interface{}) (leafI, bool) {
	if l.key.Equals(k) {
		l.val = v
		return l, false
	}
	var newLeaf = newCollisionLeaf([]iKeyVal{{l.key, l.val}, {k, v}})
	return newLeaf, true // key,val was added
}

func (l *flatLeaf) del(key *iKey) (leafI, interface{}, bool) {
	if l.key.Equals(key) {
		return nil, l.val, true
	}
	return l, nil, false
}

func (l *flatLeaf) keyVals() []iKeyVal {
	return []iKeyVal{{l.key, l.val}}
}

func (l *flatLeaf) visit(fn visitFn, depth uint) uint {
	fn(l)
	return depth
}
