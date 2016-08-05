package hamt64

import "github.com/lleo/go-hamt/hamt_key"

type leaf64I interface {
	node64I
	// from node64I
	// hash60() uint64
	// String() string

	get(key hamt_key.Key) (interface{}, bool)
	put(key hamt_key.Key, val interface{}) (leaf64I, bool)
	del(key hamt_key.Key) (interface{}, leaf64I, bool)
	keyVals() []keyVal
}
