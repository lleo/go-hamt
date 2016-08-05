package hamt32

import (
	"fmt"
	"strings"

	"github.com/lleo/go-hamt/hamt_key"
)

type collisionLeaf32 struct {
	_hash30 uint32
	kvs     []keyVal
}

func newCollisionLeaf32(hash30 uint32, kvs []keyVal) *collisionLeaf32 {
	var leaf = new(collisionLeaf32)
	leaf._hash30 = hash30
	leaf.kvs = append(leaf.kvs, kvs...)
	return leaf
}

func (l collisionLeaf32) hash30() uint32 {
	return l._hash30
}

func (l collisionLeaf32) String() string {
	var kvstrs = make([]string, len(l.kvs))
	for i := 0; i < len(l.kvs); i++ {
		kvstrs[i] = l.kvs[i].String()
	}
	var jkvstr = strings.Join(kvstrs, ",")

	return fmt.Sprintf("collisionLeaf32{hash30:%s, kvs:[]kv{%s}}", hash30String(l._hash30), jkvstr)
}

func (l collisionLeaf32) get(key hamt_key.Key) (interface{}, bool) {
	for _, kv := range l.kvs {
		if kv.key.Equals(key) {
			return kv.val, true
		}
	}
	return nil, false
}

func (l collisionLeaf32) put(key hamt_key.Key, val interface{}) (leaf32I, bool) {
	for _, kv := range l.kvs {
		if kv.key.Equals(key) {
			kv.val = val
			return l, false //key,val was not added, merely replaced
		}
	}
	l.kvs = append(l.kvs, keyVal{key, val})
	return l, true // key,val was added
}

func (l collisionLeaf32) del(key hamt_key.Key) (interface{}, leaf32I, bool) {
	for i, kv := range l.kvs {
		if kv.key.Equals(key) {
			l.kvs = append(l.kvs[:i], l.kvs[i+1:]...)
			if len(l.kvs) == 1 {
				var fl = NewFlatLeaf32(l.hash30(), l.kvs[0].key, l.kvs[0].val)
				return kv.val, fl, true
			}
			return kv.val, l, true
		}
	}
	return nil, l, false
}

func (l collisionLeaf32) keyVals() []keyVal {
	//var r []keyVal
	//r = append(r, l.kvs...)
	//return r
	return l.kvs
}
