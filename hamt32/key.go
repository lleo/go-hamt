package hamt32

import (
	"fmt"
	"hash/fnv"
)

type Key struct {
	hash HashVal
	bs   []byte
}

func newKey(bs []byte) *Key {
	var k = new(Key)

	k.bs = make([]byte, len(bs))
	copy(k.bs, bs)

	k.hash = HashVal(fold(hash(bs), remainder))

	return k
}

// Hash return the HashVal of KeyBase
func (k *Key) Hash() HashVal {
	return k.hash
}

func (k *Key) Equals(k0 *Key) bool {
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
func (k *Key) String() string {
	return fmt.Sprintf("Key{%s}", k.hash)
}

func hash(bs []byte) uint32 {
	var h = fnv.New32()
	h.Write(bs)
	return h.Sum32()
}

func mask(size uint) uint32 {
	return uint32(1<<size) - 1
}

func fold(hash uint32, rem uint) uint32 {
	return (hash >> (HashSize - rem)) ^ (hash & mask(HashSize-rem))
}
