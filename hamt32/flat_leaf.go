package hamt32

import (
	"fmt"

	"github.com/lleo/go-hamt/hamt_key"
)

type flatLeaf struct {
	_hash30 uint32
	key     hamt_key.Key
	val     interface{}
}

func newFlatLeaf(h30 uint32, key hamt_key.Key, val interface{}) *flatLeaf {
	var fl = new(flatLeaf)
	fl._hash30 = h30
	fl.key = key
	fl.val = val
	return fl
}

func (l flatLeaf) hash30() uint32 {
	return l._hash30
}

func (l flatLeaf) String() string {
	return fmt.Sprintf("flatLeaf{hash30: %s, key: %s, val: %v}", hash30String(l._hash30), l.key, l.val)
}

func (l flatLeaf) get(key hamt_key.Key) (interface{}, bool) {
	if l.key.Equals(key) {
		return l.val, true
	}
	return nil, false
}

func (l flatLeaf) put(key hamt_key.Key, val interface{}) (leafI, bool) {
	if l.key.Equals(key) {
		l.val = val
		return l, false
	}
	var newLeaf = newCollisionLeaf(l.hash30(), []keyVal{keyVal{l.key, l.val}, keyVal{key, val}})
	return newLeaf, true // key,val was added
}

func (l flatLeaf) del(key hamt_key.Key) (interface{}, leafI, bool) {
	if l.key.Equals(key) {
		return l.val, nil, true
	}
	return nil, l, false
}

func (l flatLeaf) keyVals() []keyVal {
	return []keyVal{keyVal{l.key, l.val}}
}
