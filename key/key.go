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
	//toByteSlice() []byte
}

// Base struct of all structs that satisfy Key interface
type KeyBase struct {
	hash30 uint32
}

const mask30 = uint32(1<<30) - 1

// NBITS constant is the number of bits(5) a 30bit hash value is split into
// to provied the index of a HAMT.
const NBITS uint = 5

func indexMask(depth uint) uint32 {
	return uint32(uint8(1<<NBITS)-1) << (depth * NBITS)
}

func index(h30 uint32, depth uint) uint {
	var idxMask = indexMask(depth)
	var idx = uint((h30 & idxMask) >> (depth * NBITS))
	return idx
}

func hashPathString(hashPath uint32, depth uint) string {
	if depth == 0 {
		return "/"
	}
	var strs = make([]string, depth)

	for d := depth; d > 0; d-- {
		var idx = index(hashPath, d-1)
		strs[d-1] = fmt.Sprintf("%02d", idx)
	}

	return "/" + strings.Join(strs, "/")
}

func hash30String(h30 uint32) string {
	return hashPathString(h30, 6)
}

func hash30(bs []byte) uint32 {
	var h = fnv.New32()
	h.Write(bs)
	return h.Sum32()
}

func (kb *KeyBase) String() string {
	//return fmt.Sprintf("KeyBase{hash30:%#v}", kb.hash30)
	return fmt.Sprintf("KeyBase{hash30:(%d)%s}", kb.hash30, hash30String(kb.hash30))
}

// Initialize the KeyBase part of any struct that has the KeyBase struct embeded.
// will be called by the New() function of any decendent class.
//     k.Initialize(bs)
// where k is interpreted as a *KeyBase and bs is a unique []byte to calculate
// the 30bit hash from.
func (kb *KeyBase) Initialize(bs []byte) {
	kb.hash30 = hash30(bs)
}

func (kb KeyBase) Hash30() uint32 {
	return kb.hash30
}
