/*
Package key contains a single Key interface. The key package was
created to prevent cicular dependencies between "github.com/lleo/go-hamt" and
either "github.com/lleo/go-hamt/hamt32" or "github.com/lleo/go-hamt/hamt64".

Additionally, the github.com/lleo/go-hamt/key provides a Base structure.
The Base structure if added to a derivative key type will provide the
Hamt30() and Hamt60() methods. Base needs to be populated by the derivative
key constructor calling the Initialize([]byte) method.

Any key created using the Key interface must be read only after construction.

The key package is used by the functional HAMT variation in
"github.com/lleo/go-hamt-functional".
*/
package key

import (
	"fmt"
	"hash/fnv"
	"strings"
)

// KeyVal is a simple struct used to transfer lists ([]KeyVal) from one
// function to another.
type KeyVal struct {
	Key Key
	Val interface{}
}

func (kv KeyVal) String() string {
	return fmt.Sprintf("KeyVal{%s, %v}", kv.Key, kv.Val)
}

// Key interface descibes the methods any struct needs to implement to be used
// in either the github.com/lleo/go-hamt or github.com/lleo/go-hamt-functional
// packages.
type Key interface {
	Hash30() uint32
	Hash60() uint64
	Equals(Key) bool
	String() string
	String30() string
	String60() string
	Index30(depth uint) uint
	Index60(depth uint) uint
	HashPathString30(depth uint) string
	HashPathString60(depth uint) string
}

type HashVal30 uint32
type HashVal60 uint64

// Base struct is intended to be the base struct of all structs that satisfy
// Key interface. It caches the calculated hash30 and hash60 generated when
// the method Initialize([]byte) is called.
type Base struct {
	hash30 uint32
	hash60 uint64
}

const BitsPerLevel30 uint = 5
const BitsPerLevel60 uint = 6

const MaxDepth30 uint = 5
const MaxDepth60 uint = 9

// Index30() will return a 5bit (aka BitsPerLevel30) value 'depth' number
// of 5bits from the beginning of the uint32 kb.base30 value.
func (kb Base) Index30(depth uint) uint {
	var idxMask = indexMask30(depth)
	var idx = uint(kb.hash30&idxMask) >> (depth * BitsPerLevel30)
	return idx
}

func indexMask30(depth uint) uint32 {
	return uint32((1<<BitsPerLevel30)-1) << (depth * BitsPerLevel30)
}

// Index60() will return a 6bit (aka BitsPerLevel60) value 'depth' number
// of 6bits from the beginning of the uint64 h60 value.
func (kb Base) Index60(depth uint) uint {
	var idxMask = indexMask60(depth)
	var idx = uint(kb.hash60&idxMask) >> (depth * BitsPerLevel60)
	return idx
}

func indexMask60(depth uint) uint64 {
	return uint64((1<<BitsPerLevel60)-1) << (depth * BitsPerLevel60)
}

// HashPathString30() returns a string representation of the index path of
// a hash30 30 bit value; that is depth number of zero padded numbers between
// "00" and "63" separated by "/" characters.
// Warning: It will panic() if depth > MaxDepth30.
// Example: "/00/24/46/17" for depth=4 of a hash30 value represented
//       by "/00/24/46/17/34/08".
func (kb Base) HashPathString30(depth uint) string {
	if depth > MaxDepth30 {
		panic(fmt.Sprintf("HashPathString30: depth,%d > MaxDepth30,%d\n", depth, MaxDepth30))
	}

	if depth == 0 {
		return "/"
	}

	var strs = make([]string, depth)

	for d := uint(0); d < depth; d++ {
		var idx = kb.Index30(d)
		strs[d] = fmt.Sprintf("%02d", idx)
	}

	return "/" + strings.Join(strs, "/")
}

// HashPathString60() returns a string representation of the index path of
// a hash60 60 bit value; that is depth number of zero padded numbers between
// "00" and "63" separated by "/" characters.
// Warning: It will panic() if depth > MaxDepth60.
// Example: "/00/24/46/17/34/08/54" for depth=7 of a hash60 value represented
//       by "/00/24/46/17/34/08/54/28/59/51".
func (kb Base) HashPathString60(depth uint) string {
	if depth > MaxDepth60 {
		panic(fmt.Sprintf("PathString60: depth,%d > MaxDepth60,%d\n", depth, MaxDepth60))
	}

	if depth == 0 {
		return "/"
	}

	var strs = make([]string, depth)

	for d := uint(0); d < depth; d++ {
		var idx = kb.Index60(d)
		strs[d] = fmt.Sprintf("%02d", idx)
	}

	return "/" + strings.Join(strs, "/")
}

// String() provides a human readable form of the Base's hash30 and hash60
// values. Both in raw integer and slash '/' separated integers of the 5 and
// 6 bit path values; for hash30 and hash60 respectively.
func (kb Base) String() string {
	return fmt.Sprintf("Base{%s, %s}", kb.String30(), kb.String60())
}

// String30() returns a string representation of the kb.hash30 value. This
// is MaxDepth30+1(6) two digit numbers (zero padded) between "00" and "31"
// seperated by '/' characters.
// Example: "/08/14/28/20/00/31"
func (kb Base) String30() string {
	return kb.HashPathString30(MaxDepth30)
}

// String60() returns a string representation of the kb.hash60 value. This
// is MaxDepth60+1(10) two digit numbers (zero padded) between "00" and "63"
// seperated by '/' characters.
// Example: "/08/14/28/20/00/31/56/01/24/63"
func (kb Base) String60() string {
	return kb.HashPathString60(MaxDepth60)
}

//
// Calculate Hash of byte slice functions
//

const mask30 = uint32(1<<30) - 1

// fold30 folds the top 2 bits into the bottom 30 to get 30bits of hash
// from: http://www.isthe.com/chongo/tech/comp/fnv/index.html#xor-fold
func fold30(h32 uint32) uint32 {
	return h32>>30 ^ h32&mask30
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
	return h64>>60 ^ h64&mask60
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
	kb.hash30 = fold30(hash32(bs))
	kb.hash60 = fold60(hash64(bs))
}

// Hash30 returns the uint32 that contains the 30bit hash value.
func (kb Base) Hash30() uint32 {
	return kb.hash30
}

// Hash60 returns the uint32 that contains the 60bit hash value.
func (kb Base) Hash60() uint64 {
	return kb.hash60
}
