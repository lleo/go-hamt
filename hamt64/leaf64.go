package hamt64

import "github.com/lleo/go-hamt/hamt_key"

type leafI interface {
	nodeI
	// from nodeI
	// hash60() uint64
	// String() string

	get(key hamt_key.Key) (interface{}, bool)
	put(key hamt_key.Key, val interface{}) (leafI, bool)
	del(key hamt_key.Key) (interface{}, leafI, bool)
	keyVals() []keyVal
}
