package hamt32

import (
	"fmt"

	"github.com/lleo/go-hamt/key"
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

func (l flatLeaf) Hash30() uint32 {
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

func (l flatLeaf) put(key key.Key, val interface{}) (leafI, bool) {
	if l.key.Equals(key) {
		l.val = val
		return l, false
	}
	var newLeaf = newCollisionLeaf([]keyVal{keyVal{l.key, l.val}, keyVal{key, val}})
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
