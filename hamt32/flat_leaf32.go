package hamt32

import (
	"fmt"

	"github.com/lleo/go-hamt/hamt_key"
)

type flatLeaf32 struct {
	_hash30 uint32
	key     hamt_key.Key
	val     interface{}
}

func NewFlatLeaf32(h30 uint32, key hamt_key.Key, val interface{}) *flatLeaf32 {
	var fl = new(flatLeaf32)
	fl._hash30 = h30
	fl.key = key
	fl.val = val
	return fl
}

func (l flatLeaf32) hash30() uint32 {
	return l._hash30
}

func (l flatLeaf32) String() string {
	return fmt.Sprintf("flatLeaf32{hash30: %s, key: %s, val: %v}", hash30String(l._hash30), l.key, l.val)
}

func (l flatLeaf32) get(key hamt_key.Key) (interface{}, bool) {
	if l.key.Equals(key) {
		return l.val, true
	}
	return nil, false
}

func (l flatLeaf32) put(key hamt_key.Key, val interface{}) (leaf32I, bool) {
	if l.key.Equals(key) {
		l.val = val
		return l, false
	}
	var newLeaf = newCollisionLeaf32(l.hash30(), []keyVal{keyVal{l.key, l.val}, keyVal{key, val}})
	return newLeaf, true // key,val was added
}

func (l flatLeaf32) del(key hamt_key.Key) (interface{}, leaf32I, bool) {
	if l.key.Equals(key) {
		return l.val, nil, true
	}
	return nil, l, false
}

func (l flatLeaf32) keyVals() []keyVal {
	return []keyVal{keyVal{l.key, l.val}}
}
