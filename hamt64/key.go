package hamt64

import (
	"fmt"
	"hash/fnv"
)

// iKey is the internal key structure build from the []byte slice passed to the
// Get/Put/Del hamt methods.
type iKey struct {
	hash hashVal
	bs   []byte
}

// newKey builds a iKey data structure from a []byte slice.
func newKey(bs []byte) *iKey {
	var k = new(iKey)

	k.bs = make([]byte, len(bs))
	copy(k.bs, bs)

	k.hash = hashVal(fold(hash(bs), remainder))

	return k
}

// Hash return the hashVal of KeyBase
func (k *iKey) Hash() hashVal {
	return k.hash
}

func (k *iKey) Equals(k0 *iKey) bool {
	//return string(k.bs) == string(k0)
	if len(k.bs) == len(k0.bs) {
		for i, ke := range k.bs {
			if ke != k0.bs[i] {
				return false
			}
			return true
		}
	}
	return false
}

// String return a human readable representation of KeyBase
func (k *iKey) String() string {
	return fmt.Sprintf("KeyBase{%s}", k.hash)
}

func hash(bs []byte) uint64 {
	var h = fnv.New64()
	h.Write(bs)
	return h.Sum64()
}

func mask(size uint) uint64 {
	return uint64(1<<size) - 1
}

func fold(hash uint64, rem uint) uint64 {
	return (hash >> (hashSize - rem)) ^ (hash & mask(hashSize-rem))
}
