package hamt32

import (
	"fmt"
	"log"
	"strings"

	"github.com/lleo/go-hamt-key"
)

// implements nodeI
// implements leafI
type collisionLeaf struct {
	kvs []key.KeyVal
}

func newCollisionLeaf(kvs []key.KeyVal) *collisionLeaf {
	var leaf = new(collisionLeaf)
	leaf.kvs = append(leaf.kvs, kvs...)

	return leaf
}

func (l collisionLeaf) Hash30() key.HashVal30 {
	return l.kvs[0].Key.Hash30()
}

func (l collisionLeaf) String() string {
	var kvstrs = make([]string, len(l.kvs))
	for i := 0; i < len(l.kvs); i++ {
		kvstrs[i] = l.kvs[i].String()
	}
	var jkvstr = strings.Join(kvstrs, ",")

	return fmt.Sprintf("collisionLeaf{hash30:%s, kvs:[]kv{%s}}",
		l.Hash30(), jkvstr)
}

func (l collisionLeaf) get(key key.Key) (interface{}, bool) {
	for _, kv := range l.kvs {
		if kv.Key.Equals(key) {
			return kv.Val, true
		}
	}
	return nil, false
}

func (l collisionLeaf) put(k key.Key, v interface{}) (leafI, bool) {
	for _, kv := range l.kvs {
		if kv.Key.Equals(k) {
			kv.Val = v
			return l, false //key,val was not added, merely replaced
		}
	}
	l.kvs = append(l.kvs, key.KeyVal{k, v})

	log.Printf("%s : %d\n", l.Hash30(), len(l.kvs))

	return l, true // key_,val was added
}

func (l collisionLeaf) del(key key.Key) (interface{}, leafI, bool) {
	for i, kv := range l.kvs {
		if kv.Key.Equals(key) {
			l.kvs = append(l.kvs[:i], l.kvs[i+1:]...)
			if len(l.kvs) == 1 {
				var fl = newFlatLeaf(l.kvs[0].Key, l.kvs[0].Val)
				return kv.Val, fl, true
			}
			return kv.Val, l, true
		}
	}
	return nil, l, false
}

func (l collisionLeaf) keyVals() []key.KeyVal {
	var r = make([]key.KeyVal, len(l.kvs))
	r = append(r, l.kvs...)
	return r
	//return l.kvs
}
