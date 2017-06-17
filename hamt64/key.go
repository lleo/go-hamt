package hamt64

import (
	"fmt"
	"hash/fnv"
)

// HashSize is the size of the basic hash function output.
const HashSize uint = 64

// BitsPerLevel is the fundemental setting along with HashSize for the Key
// constants. 2..HashSize/2 step 1
const BitsPerLevel uint = 6

// DepthLimit is the maximum number of levels of the Hamt. It is calculated as
// DepthLimit = floor(HashSize / BitsPerLevel) or a strict integer division.
const DepthLimit = HashSize / BitsPerLevel
const remainder = HashSize - (DepthLimit * BitsPerLevel)

// IndexLimit is the maximum number of entries in the Hamt interior nodes.
// IndexLimit = 1 << BitsPerLevel
const IndexLimit = 1 << BitsPerLevel

// MaxDepth is the maximum value of a depth variable. MaxDepth = DepthLimit - 1
const MaxDepth = DepthLimit - 1

// MaxIndex is the maximum value of a index variable. MaxIndex = IndexLimie - 1
const MaxIndex = IndexLimit - 1

// Key interface descibes the methods any struct needs to implement to be used
// as a Key in github.com/lleo/go-hamt/hamt64
type Key interface {
	Hash() HashVal
	Equals(Key) bool
	String() string
}

// KeyBase is the fundemental struct of any derived key structure.
type KeyBase struct {
	hash HashVal
}

// Hash return the HashVal of KeyBase
func (kb *KeyBase) Hash() HashVal {
	return kb.hash
}

// String return a human readable representation of KeyBase
func (kb *KeyBase) String() string {
	return fmt.Sprintf("KeyBase{%s}", kb.hash)
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
	return (hash >> (HashSize - rem)) ^ (hash & mask(HashSize-rem))
}

// Initialize MUST be call when a derived Key structure is created as it sets
// the KeyBase.hash value.
func (b *KeyBase) Initialize(basis []byte) {
	b.hash = HashVal(fold(hash(basis), remainder))
}
