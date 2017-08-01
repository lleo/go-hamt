package hamt64

import (
	"bytes"
	"fmt"
)

type flatLeaf struct {
	hash hashVal
	key  []byte
	val  interface{}
}

func newFlatLeaf(hv hashVal, key []byte, val interface{}) *flatLeaf {
	var fl = new(flatLeaf)
	fl.hash = hv
	fl.key = key
	fl.val = val
	return fl
}

func (l *flatLeaf) Hash() hashVal {
	return l.hash
}

func (l *flatLeaf) String() string {
	return fmt.Sprintf("flatLeaf{key: %s, val: %v}", l.key, l.val)
}

func (l *flatLeaf) get(key []byte) (interface{}, bool) {
	if bytes.Equal(l.key, key) {
		return l.val, true
	}
	return nil, false
}

// put maintains the functional behavior that any modification returns a new
// leaf and the original remains unaltered. It returns the new leafI and a bool
// indicating if the key,val was added ontop of the current leaf key,val or if
// the val mearly replaced the current key's val (either way a new leafI is
// allocated and returned).
func (l *flatLeaf) put(key []byte, val interface{}) (leafI, bool) {
	var nl leafI
	if bytes.Equal(l.key, key) {
		// maintain functional behavior of flatLeaf
		nl = newFlatLeaf(l.hash, l.key, val)
		return nl, false //replaced
	}
	nl = newCollisionLeaf(l.hash, []KeyVal{{l.key, l.val}, {key, val}})
	return nl, true // key,val was added
}

func (l *flatLeaf) del(key []byte) (leafI, interface{}, bool) {
	if bytes.Equal(l.key, key) {
		return nil, l.val, true //found
	}
	return l, nil, false //not found
}

func (l *flatLeaf) keyVals() []KeyVal {
	return []KeyVal{{l.key, l.val}}
}

func (l *flatLeaf) visit(fn visitFn, depth uint) uint {
	fn(l)
	return depth
}
