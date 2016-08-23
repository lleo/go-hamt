/*
The "github.com/lleo/go-hamt/string_key package" is an implementation of the
"github.com/lleo/go-hamt/hamt_key.Key" interface.

Both of these packages are are in separate package namespaces to prevent circular
references between "github.com/lleo/go-hamt" and "github.com/lleo/go-hamt/hamt32"
and/or "github.com/lleo/go-hamt/hamt64".

Additionally, "github.com/lleo/go-hamt/string_key" is used by the functional
variation of the Hash Array Mapped Trie in "github.com/lleo/go-hamt-fuctional".
*/
package string_key

import (
	"hash/fnv"

	"github.com/lleo/go-hamt/hamt_key"
)

type StringKey string

// Equals returns true iff the StringKey exactly matches the key passed it. The
// hamt_key.Key passed as an argument MUST BE a StringKey or the method Equals()
// automatically returns false.
func (sk StringKey) Equals(key hamt_key.Key) bool {
	var k, typeMatches = key.(StringKey)
	if !typeMatches {
		panic("type mismatch")
	}
	return string(sk) == string(k)
}

// The Hash30() calculates the fnv1 32bit hash of the string (redered to its bytes).
// Then we use the xor-fold technique described <a href="http://www.isthe.com/chongo/tech/comp/fnv/index.html#xor-fold">here</a>
// to fold the top 2bits into the lower 30bits of the 32bit hash value.
func (sk StringKey) Hash30() uint32 {
	return hash30(sk.hash32())
}

// The Hash60() calculates the fnv1 64bit hash of the string (redered to its bytes).
// Then we use the xor-fold technique described <a href="http://www.isthe.com/chongo/tech/comp/fnv/index.html#xor-fold">here</a>
// to fold the top 4bits into the lower 60bits of the 64bit hash value.
func (sk StringKey) Hash60() uint64 {
	return hash60(sk.hash64())
}

func (sk StringKey) hash32() uint32 {
	var h = fnv.New32()
	h.Write([]byte(sk))
	return h.Sum32()
}

func (sk StringKey) hash64() uint64 {
	var h = fnv.New64()
	h.Write([]byte(sk))
	return h.Sum64()
}

const mask30 = uint32(1<<30) - 1
const mask60 = uint64(1<<60) - 1

func hash30(h30 uint32) uint32 {
	return (h30 >> 30) ^ (h30 & mask30)
}

func hash60(h64 uint64) uint64 {
	return (h64 >> 60) ^ (h64 & mask60)
}

// String() returns the StringKey as a basic string type.
func (sk StringKey) String() string {
	return string(sk)
}
