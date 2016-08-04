package hamt32

type leaf32I interface {
	node32I
	// from node32I
	// hash30() uint32
	// String() string

	get(key Key) (interface{}, bool)
	put(key Key, val interface{}) (leaf32I, bool)
	del(key Key) (interface{}, leaf32I, bool)
	keyVals() []keyVal
}
