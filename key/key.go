/*
Package key contains a single Key interface. The key package was created to
prevent cicular dependencies between "github.com/lleo/go-hamt" and either
 "github.com/lleo/go-hamt/hamt32", "github.com/lleo/go-hamt/hamt64",
"github.com/lleo/go-hamt-functional/hamt32", or
"github.com/lleo/go-hamt-functional/hamt64"

Additionally, the "github.com/lleo/go-hamt/key" provides a Base structure. The
Base structure if added to a derivative key type will provide the k.Hash30() and
k.Hash60() methods. Base needs to be populated by the derivative key constructor
calling the k.Initialize([]byte) method.

Any key created using the Key interface must be read-only after construction.

The Hash30() returns a special type HashVal30. Which really just an alias for
uint32 that stores a 30 bit hash value, but it provides methods for viewing the
30 bit hash30 value. For instance the hv.Index() method will pull out the 5 bit
integer values, as a uint, for each depth of the Hamt datastructure that is
really the index into a table. Also, there are methods to produce string
representations of the underlying 30 bit hash value.

The Hash60() returns a special type HashVal60. Which really just an alias for
uint64 that stores a 60 bit hash value, but it provides methods for viewing the
60 bit hash60 value. For instance the hv.Index() method will pull out the 5 bit
integer values, as a uint, for each depth of the Hamt datastructure that is
really the index into a table. Also, there are methods to produce string
representations of the underlying 60 bit hash value.

The key package is also used by the functional HAMT variation in
"github.com/lleo/go-hamt-functional".
*/
package key

import (
	"fmt"
	"hash/fnv"
)

// Key interface descibes the methods any struct needs to implement to be used
// in either the github.com/lleo/go-hamt or github.com/lleo/go-hamt-functional
// packages.
type Key interface {
	Hash30() HashVal30
	Hash60() HashVal60
	Equals(Key) bool
	String() string
}

// Base struct is intended to be the base struct of all structs that satisfy
// Key interface. It caches the calculated hash30 and hash60 generated when
// the method Initialize([]byte) is called.
type Base struct {
	hash30 HashVal30
	hash60 HashVal60
}

// Hash30 returns the HashVal30 that contains the 30bit hash value.
func (kb Base) Hash30() HashVal30 {
	return kb.hash30
}

// Hash60 returns the HashVal60 that contains the 60bit hash value.
func (kb Base) Hash60() HashVal60 {
	return kb.hash60
}

// String() provides a human readable form of the Base's hash30 and hash60
// values. Both in raw integer and slash '/' separated integers of the 5 and
// 6 bit path values; for hash30 and hash60 respectively.
func (kb Base) String() string {
	return fmt.Sprintf("Base{%s, %s}", kb.hash30.String(), kb.hash60.String())
}

//
// Calculate Hash of byte slice functions
//

const mask30 = uint32(1<<30) - 1

// fold30 folds the top 2 bits into the bottom 30 to get 30bits of hash
// from: http://www.isthe.com/chongo/tech/comp/fnv/index.html#xor-fold
func fold30(h32 uint32) uint32 {
	return (h32 >> 30) ^ (h32 & mask30)
}

// hash64 calculates the 64bit hash value of a byte string using hash/fnv.
func hash32(bs []byte) uint32 {
	var h = fnv.New32()
	h.Write(bs)
	return h.Sum32()
}

const mask60 = uint64(1<<60) - 1

// fold60 folds the top 4 bits into the bottom 60 to get 60bits of hash
// from: http://www.isthe.com/chongo/tech/comp/fnv/index.html#xor-fold
func fold60(h64 uint64) uint64 {
	return (h64 >> 60) ^ (h64 & mask60)
}

// hash64 calculates the 64bit hash value of a byte string using hash/fnv.
func hash64(bs []byte) uint64 {
	var h = fnv.New64()
	h.Write(bs)
	return h.Sum64()
}

// Initialize the calculates the Base part of any struct that has the Base
// struct embedded. It will be called by the New() function of any descendant
// class.
//     k.Initialize(bs)
// where k is interpreted as a *Base and bs is a unique []byte to calculate
// the 30bit and 60bit hash values from.
func (kb *Base) Initialize(bs []byte) {
	kb.hash30 = HashVal30(fold30(hash32(bs)))
	kb.hash60 = HashVal60(fold60(hash64(bs)))
}
