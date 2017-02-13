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
	String() string // ISA/statisfies fmt.Stringer interface
	String30() string
	String60() string
	Equals(Key) bool
	Hash30() uint32
	Hash60() uint64
}

// Base struct is intended to be the base struct of all structs that satisfy
// Key interface. It caches the calculated hash30 and hash60 generated when
// the method Initialize([]byte) is called.
type Base struct {
	hash30 uint32
	hash60 uint64
}

// INTERESTING IDEA JUST NOT YET...
//
// type Hash30 uint32
//
// func (h30 Hash30) String() {
// 	return hash30String(uint32(h30))
// }
//
// func (h30 Hash30) index(depth uint) {
// 	return index(uint32(h30), depth)
// }
//
// type Hash60 uint64
//
// func (h Hash60) String() {
// 	return hash60String(uint64(h))
// }
//
// func (h60 Hash60) index(depth uint) {
// 	return index(uint64(h60), depth)
// }

//
// 30bit hash value and string generator functions
// Hash values are generated when Initialize([]byte) is called.
// Hash strings are generated as needed.
//

const nBits30 uint = 5

func indexMask30(depth uint) uint32 {
	return uint32(uint8(1<<nBits30)-1) << (depth * nBits30)
}

func index30(h30 uint32, depth uint) uint {
	var idxMask = indexMask30(depth)
	var idx = uint((h30 & idxMask) >> (depth * nBits30))
	return idx
}

func hashPathString30(hashPath uint32, depth uint) string {
	if depth == 0 {
		return "/"
	}
	var strs = make([]string, depth+1)

	for d := uint(0); d <= depth; d++ {
		var idx = index30(hashPath, d)
		strs[d] = fmt.Sprintf("%02d", idx)
	}

	return "/" + strings.Join(strs, "/")
}

func hash30String(h30 uint32) string {
	return hashPathString30(h30, 5)
}

//
// 60bit hash value and string generator functions.
// Hash values are generated when Initialize([]byte) is called.
// Hash strings are generated as needed.
//

const nBits60 uint = 6

func indexMask60(depth uint) uint64 {
	return uint64(uint64(1<<nBits60)-1) << (depth * nBits60)
}

func index60(h60 uint64, depth uint) uint {
	var idxMask = indexMask60(depth)
	var idx = uint((h60 & idxMask) >> (depth * nBits60))
	return idx
}

func hashPathString60(hashPath uint64, depth uint) string {
	if depth == 0 {
		return "/"
	}
	var strs = make([]string, depth+1)

	for d := uint(0); d <= depth; d++ {
		var idx = index60(hashPath, d)
		strs[d] = fmt.Sprintf("%02d", idx)
	}

	return "/" + strings.Join(strs, "/")
}

func hash60String(h60 uint64) string {
	return hashPathString60(h60, 9)
}

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

// String provides a human readable form of the Base's hash30 and hash60
// values. Both in raw integer and slash '/' separated integers of the 5 and
// 6 bit path values; for hash30 and hash60 respectively.
func (kb Base) String() string {
	return fmt.Sprintf("Base{%s, %s}", hash30String(kb.hash30), hash60String(kb.hash60))
}

func (kb Base) String30() string {
	return hash30String(kb.hash30)
}

func (kb Base) String60() string {
	return hash60String(kb.hash60)
}

// Initialize the Base part of any struct that has the Base struct embedded.
// will be called by the New() function of any descendant class.
//     k.Initialize(bs)
// where k is interpreted as a *Base and bs is a unique []byte to calculate
// the 30bit hash from.
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
