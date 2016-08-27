package hamt64

import "github.com/lleo/go-hamt/key"

type leafI interface {
	nodeI
	// from nodeI
	// hash60() uint64
	// String() string

	get(key key.Key) (interface{}, bool)
	put(key key.Key, val interface{}) (leafI, bool)
	del(key key.Key) (interface{}, leafI, bool)
	keyVals() []keyVal
}
