package hamt32

import "github.com/lleo/go-hamt/hamt_key"

type leaf32I interface {
	node32I
	// from node32I
	// hash30() uint32
	// String() string

	get(key hamt_key.Key) (interface{}, bool)
	put(key hamt_key.Key, val interface{}) (leaf32I, bool)
	del(key hamt_key.Key) (interface{}, leaf32I, bool)
	keyVals() []keyVal
}
