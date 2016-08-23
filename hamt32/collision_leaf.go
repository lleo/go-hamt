package hamt32

import (
	"fmt"
	"strings"

	"github.com/lleo/go-hamt/hamt_key"
)

type collisionLeaf struct {
	_hash30 uint32
	kvs     []keyVal
}

func newCollisionLeaf(hash30 uint32, kvs []keyVal) *collisionLeaf {
	var leaf = new(collisionLeaf)
	leaf._hash30 = hash30
	leaf.kvs = append(leaf.kvs, kvs...)
	return leaf
}

func (l collisionLeaf) hash30() uint32 {
	return l._hash30
}

func (l collisionLeaf) String() string {
	var kvstrs = make([]string, len(l.kvs))
	for i := 0; i < len(l.kvs); i++ {
		kvstrs[i] = l.kvs[i].String()
	}
	var jkvstr = strings.Join(kvstrs, ",")

	return fmt.Sprintf("collisionLeaf{hash30:%s, kvs:[]kv{%s}}", hash30String(l._hash30), jkvstr)
}

func (l collisionLeaf) get(key hamt_key.Key) (interface{}, bool) {
	for _, kv := range l.kvs {
		if kv.key.Equals(key) {
			return kv.val, true
		}
	}
	return nil, false
}

func (l collisionLeaf) put(key hamt_key.Key, val interface{}) (leafI, bool) {
	for _, kv := range l.kvs {
		if kv.key.Equals(key) {
			kv.val = val
			return l, false //key,val was not added, merely replaced
		}
	}
	l.kvs = append(l.kvs, keyVal{key, val})
	return l, true // key,val was added
}

func (l collisionLeaf) del(key hamt_key.Key) (interface{}, leafI, bool) {
	for i, kv := range l.kvs {
		if kv.key.Equals(key) {
			l.kvs = append(l.kvs[:i], l.kvs[i+1:]...)
			if len(l.kvs) == 1 {
				var fl = newFlatLeaf(l.hash30(), l.kvs[0].key, l.kvs[0].val)
				return kv.val, fl, true
			}
			return kv.val, l, true
		}
	}
	return nil, l, false
}

func (l collisionLeaf) keyVals() []keyVal {
	//var r []keyVal
	//r = append(r, l.kvs...)
	//return r
	return l.kvs
}
