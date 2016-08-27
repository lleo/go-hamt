package hamt64

import (
	"fmt"

	"github.com/lleo/go-hamt/key"
)

type flatLeaf struct {
	_hash60 uint64
	key     key.Key
	val     interface{}
}

func newFlatLeaf(h60 uint64, key key.Key, val interface{}) *flatLeaf {
	var fl = new(flatLeaf)
	fl._hash60 = h60
	fl.key = key
	fl.val = val
	return fl
}

func (l flatLeaf) hash60() uint64 {
	return l._hash60
}

func (l flatLeaf) String() string {
	return fmt.Sprintf("flatLeaf{hash60: %s, key: %s, val: %v}", hash60String(l._hash60), l.key, l.val)
}

func (l flatLeaf) get(key key.Key) (interface{}, bool) {
	if l.key.Equals(key) {
		return l.val, true
	}
	return nil, false
}

func (l flatLeaf) put(key key.Key, val interface{}) (leafI, bool) {
	if l.key.Equals(key) {
		l.val = val
		return l, false
	}
	var newLeaf = newCollisionLeaf(l.hash60(), []keyVal{keyVal{l.key, l.val}, keyVal{key, val}})
	return newLeaf, true // key,val was added
}

func (l flatLeaf) del(key key.Key) (interface{}, leafI, bool) {
	if l.key.Equals(key) {
		return l.val, nil, true
	}
	return nil, l, false
}

func (l flatLeaf) keyVals() []keyVal {
	return []keyVal{keyVal{l.key, l.val}}
}
