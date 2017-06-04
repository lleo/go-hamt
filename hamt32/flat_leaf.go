package hamt32

import (
	"fmt"

	"github.com/lleo/go-hamt-key"
)

type flatLeaf struct {
	key key.Key
	val interface{}
}

func newFlatLeaf(key key.Key, val interface{}) *flatLeaf {
	var fl = new(flatLeaf)
	fl.key = key
	fl.val = val
	return fl
}

func (l flatLeaf) Hash30() key.HashVal30 {
	return l.key.Hash30()
}

func (l flatLeaf) String() string {
	return fmt.Sprintf("flatLeaf{key: %s, val: %v}", l.key, l.val)
}

func (l flatLeaf) get(key key.Key) (interface{}, bool) {
	if l.key.Equals(key) {
		return l.val, true
	}
	return nil, false
}

func (l flatLeaf) put(k key.Key, v interface{}) (leafI, bool) {
	if l.key.Equals(k) {
		l.val = v
		return l, false
	}
	var newLeaf = newCollisionLeaf([]key.KeyVal{{l.key, l.val}, {k, v}})
	return newLeaf, true // key,val was added
}

func (l flatLeaf) del(key key.Key) (leafI, interface{}, bool) {
	if l.key.Equals(key) {
		return nil, l.val, true
	}
	return l, nil, false
}

func (l flatLeaf) keyVals() []key.KeyVal {
	return []key.KeyVal{{l.key, l.val}}
}
