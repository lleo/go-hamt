package hamt64

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

func (l *collisionLeaf) copy() *collisionLeaf {
	var nl = new(collisionLeaf)
	nl.kvs = append(nl.kvs, l.kvs...)
	return nl
}

func newCollisionLeaf(kvs []key.KeyVal) *collisionLeaf {
	var leaf = new(collisionLeaf)
	leaf.kvs = append(leaf.kvs, kvs...)

	return leaf
}

func (l collisionLeaf) Hash60() key.HashVal60 {
	return l.kvs[0].Key.Hash60()
}

func (l collisionLeaf) String() string {
	var kvstrs = make([]string, len(l.kvs))
	for i := 0; i < len(l.kvs); i++ {
		kvstrs[i] = l.kvs[i].String()
	}
	var jkvstr = strings.Join(kvstrs, ",")

	return fmt.Sprintf("collisionLeaf{hash60:%s, kvs:[]key.KeyVal{%s}}",
		l.Hash60(), jkvstr)
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

	log.Printf("%s : %d\n", l.Hash60(), len(l.kvs))

	return l, true // key_,val was added
}

func (l collisionLeaf) del(key key.Key) (leafI, interface{}, bool) {
	var nl = l.copy()
	for i, kv := range l.kvs {
		if kv.Key.Equals(key) {
			nl.kvs = append(nl.kvs[:i], nl.kvs[i+1:]...)
			if len(nl.kvs) == 1 {
				var fl = newFlatLeaf(nl.kvs[0].Key, nl.kvs[0].Val)
				return fl, kv.Val, true
			}
			return nl, kv.Val, true
		}
	}
	return l, nil, false
}

func (l collisionLeaf) keyVals() []key.KeyVal {
	var r = make([]key.KeyVal, 0, len(l.kvs))
	r = append(r, l.kvs...)
	return r
	//return l.kvs
}
