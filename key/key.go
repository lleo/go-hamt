/*
Package key contains a single Key interface. The key package was
created to prevent cicular depedencies betwee "github.com/lleo/go-hamt" and
either "github.com/lleo/go-hamt/hamt32" or "github.com/lleo/go-hamt/hamt64".

However the key pacakge is also used by the functional HAMT variation in
"github.com/lleo/go-hamt-functional".
*/
package key

import "hash/fnv"

type Key interface {
	String() string // ISA/statisfies fmt.Stringer interface
	Equals(Key) bool
	Hash30() uint32
	ToByteSlice() []byte
}

// Base struct of all structs that satisfy Key interface
type KeyBase struct {
	hash30 uint32
}

const mask30 = uint32(1<<30) - 1

func hash30(bs []byte) uint32 {
	var h = fnv.New32()
	h.Write(bs)
	return h.Sum32()
}

// Initialize the KeyBase part of any struct that has the KeyBase struct embeded.
// will be called by the New() function of any super class.
//     k.Initialize(k)
// where the first k is interpreted as a *KeyBase and the second k is interpreted
// as a Key interface type.
func (kb *KeyBase) Initialize(bs []byte) {
	kb.hash30 = hash30(bs)
}

func (kb KeyBase) Hash30() uint32 {
	return kb.hash30
}
