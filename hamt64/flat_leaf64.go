package hamt64

import (
	"fmt"

	"github.com/lleo/go-hamt/hamt_key"
)

type flatLeaf64 struct {
	_hash60 uint64
	key     hamt_key.Key
	val     interface{}
}

func NewFlatLeaf64(h60 uint64, key hamt_key.Key, val interface{}) *flatLeaf64 {
	var fl = new(flatLeaf64)
	fl._hash60 = h60
	fl.key = key
	fl.val = val
	return fl
}

func (l flatLeaf64) hash60() uint64 {
	return l._hash60
}

func (l flatLeaf64) String() string {
	return fmt.Sprintf("flatLeaf64{hash60: %s, key: %s, val: %v}", hash60String(l._hash60), l.key, l.val)
}

func (l flatLeaf64) get(key hamt_key.Key) (interface{}, bool) {
	if l.key.Equals(key) {
		return l.val, true
	}
	return nil, false
}

func (l flatLeaf64) put(key hamt_key.Key, val interface{}) (leaf64I, bool) {
	if l.key.Equals(key) {
		l.val = val
		return l, false
	}
	var newLeaf = newCollisionLeaf64(l.hash60(), []keyVal{keyVal{l.key, l.val}, keyVal{key, val}})
	return newLeaf, true // key,val was added
}

func (l flatLeaf64) del(key hamt_key.Key) (interface{}, leaf64I, bool) {
	if l.key.Equals(key) {
		return l.val, nil, true
	}
	return nil, l, false
}

func (l flatLeaf64) keyVals() []keyVal {
	return []keyVal{keyVal{l.key, l.val}}
}
