package hamt64

import (
	"fmt"
	"strings"

	"github.com/lleo/go-hamt/hamt_key"
)

type collisionLeaf64 struct {
	_hash60 uint64
	kvs     []keyVal
}

func newCollisionLeaf64(hash60 uint64, kvs []keyVal) *collisionLeaf64 {
	var leaf = new(collisionLeaf64)
	leaf._hash60 = hash60
	leaf.kvs = append(leaf.kvs, kvs...)
	return leaf
}

func (l collisionLeaf64) hash60() uint64 {
	return l._hash60
}

func (l collisionLeaf64) String() string {
	var kvstrs = make([]string, len(l.kvs))
	for i := 0; i < len(l.kvs); i++ {
		kvstrs[i] = l.kvs[i].String()
	}
	var jkvstr = strings.Join(kvstrs, ",")

	return fmt.Sprintf("collisionLeaf64{hash60:%s, kvs:[]kv{%s}}", hash60String(l._hash60), jkvstr)
}

func (l collisionLeaf64) get(key hamt_key.Key) (interface{}, bool) {
	for _, kv := range l.kvs {
		if kv.key.Equals(key) {
			return kv.val, true
		}
	}
	return nil, false
}

func (l collisionLeaf64) put(key hamt_key.Key, val interface{}) (leaf64I, bool) {
	for _, kv := range l.kvs {
		if kv.key.Equals(key) {
			kv.val = val
			return l, false //key,val was not added, merely replaced
		}
	}
	l.kvs = append(l.kvs, keyVal{key, val})
	return l, true // key,val was added
}

func (l collisionLeaf64) del(key hamt_key.Key) (interface{}, leaf64I, bool) {
	for i, kv := range l.kvs {
		if kv.key.Equals(key) {
			l.kvs = append(l.kvs[:i], l.kvs[i+1:]...)
			if len(l.kvs) == 1 {
				var fl = NewFlatLeaf64(l.hash60(), l.kvs[0].key, l.kvs[0].val)
				return kv.val, fl, true
			}
			return kv.val, l, true
		}
	}
	return nil, l, false
}

func (l collisionLeaf64) keyVals() []keyVal {
	//var r []keyVal
	//r = append(r, l.kvs...)
	//return r
	return l.kvs
}
