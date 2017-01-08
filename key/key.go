/*
Package key contains a single Key interface. The key package was
created to prevent cicular depedencies betwee "github.com/lleo/go-hamt" and
either "github.com/lleo/go-hamt/hamt32" or "github.com/lleo/go-hamt/hamt64".

However the key pacakge is also used by the functional HAMT variation in
"github.com/lleo/go-hamt-functional".
*/
package key

import (
	"fmt"
	"hash/fnv"
	"strings"
)

type Key interface {
	String() string // ISA/statisfies fmt.Stringer interface
	Equals(Key) bool
	Hash30() uint32
	Hash60() uint64
}

// Base struct of all structs that satisfy Key interface
type KeyBase struct {
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
// 30bit hash values & functions
//

const NBITS30 uint = 5

func indexMask30(depth uint) uint32 {
	return uint32(uint8(1<<NBITS30)-1) << (depth * NBITS30)
}

func index30(h30 uint32, depth uint) uint {
	var idxMask = indexMask30(depth)
	var idx = uint((h30 & idxMask) >> (depth * NBITS30))
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
// 60bit hash values & functions
//

const NBITS60 uint = 6

func indexMask60(depth uint) uint64 {
	return uint64(uint64(1<<NBITS60)-1) << (depth * NBITS60)
}

func index60(h60 uint64, depth uint) uint {
	var idxMask = indexMask60(depth)
	var idx = uint((h60 & idxMask) >> (depth * NBITS60))
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

// Fold the top two bits into the bottom 30 to get 30bits of hash
// from: http://www.isthe.com/chongo/tech/comp/fnv/index.html#xor-fold
func fold30(h32 uint32) uint32 {
	return h32>>30 ^ h32&mask30
}

func hash32(bs []byte) uint32 {
	var h = fnv.New32()
	h.Write(bs)
	return h.Sum32()
}

const mask60 = uint64(1<<60) - 1

// Fold the top four bits into the bottom 60 to get 60bits of hash
// from: http://www.isthe.com/chongo/tech/comp/fnv/index.html#xor-fold
func fold60(h64 uint64) uint64 {
	return h64>>60 ^ h64&mask60
}

func hash64(bs []byte) uint64 {
	var h = fnv.New64()
	h.Write(bs)
	return h.Sum64()
}

func (kb *KeyBase) String() string {
	return fmt.Sprintf("KeyBase{hash30:(%d)%s, hash60:(%d)%s}",
		kb.hash30, hash30String(kb.hash30),
		kb.hash60, hash60String(kb.hash60))
}

// Initialize the KeyBase part of any struct that has the KeyBase struct embeded.
// will be called by the New() function of any decendent class.
//     k.Initialize(bs)
// where k is interpreted as a *KeyBase and bs is a unique []byte to calculate
// the 30bit hash from.
func (kb *KeyBase) Initialize(bs []byte) {
	kb.hash30 = fold30(hash32(bs))
	kb.hash60 = fold60(hash64(bs))
}

func (kb *KeyBase) Hash30() uint32 {
	return kb.hash30
}

func (kb *KeyBase) Hash60() uint64 {
	return kb.hash60
}
