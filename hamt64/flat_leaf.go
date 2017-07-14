package hamt64

import (
	"fmt"
)

type flatLeaf struct {
	key *Key
	val interface{}
}

func newFlatLeaf(k *Key, val interface{}) *flatLeaf {
	var fl = new(flatLeaf)
	fl.key = k
	fl.val = val
	return fl
}

func (l *flatLeaf) Hash() hashVal {
	return l.key.Hash()
}

func (l *flatLeaf) String() string {
	return fmt.Sprintf("flatLeaf{key: %s, val: %v}", l.key, l.val)
}

func (l *flatLeaf) get(k *Key) (interface{}, bool) {
	if l.key.Equals(k) {
		return l.val, true
	}
	return nil, false
}

func (l *flatLeaf) put(k *Key, v interface{}) (leafI, bool) {
	if l.key.Equals(k) {
		l.val = v
		return l, false
	}
	var newLeaf = newCollisionLeaf([]KeyVal{{l.key, l.val}, {k, v}})
	return newLeaf, true // key,val was added
}

func (l *flatLeaf) del(k *Key) (leafI, interface{}, bool) {
	if l.key.Equals(k) {
		return nil, l.val, true
	}
	return l, nil, false
}

func (l *flatLeaf) keyVals() []KeyVal {
	return []KeyVal{{l.key, l.val}}
}

func (l *flatLeaf) visit(fn visitFn, arg interface{}, depth uint) uint {
	fn(l, arg)
	return depth - 1 //remove cuz this method is called with depth+1
}
