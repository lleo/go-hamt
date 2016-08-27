package hamt32

import "github.com/lleo/go-hamt/key"

type leafI interface {
	nodeI
	// from nodeI
	// hash30() uint32
	// String() string

	get(key key.Key) (interface{}, bool)
	put(key key.Key, val interface{}) (leafI, bool)
	del(key key.Key) (interface{}, leafI, bool)
	keyVals() []keyVal
}