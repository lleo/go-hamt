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
	Equals(Key) bool
	ToByteSlice() []byte
	String() string
}

const mask30 = uint32(1<<30) - 1

func hash30(h30 uint32) uint32 {
	return (h30 >> 30) ^ (h30 & mask30)
}

func hash32(k Key) uint32 {
	var bs = k.ToByteSlice()
	var h = fnv.New32()
	h.Write(bs)
	return h.Sum32()
}

func Hash30(k Key) uint32 {
	return hash30(hash32(k))
}

const mask60 = uint64(1<<60) - 1

func hash60(h64 uint64) uint64 {
	return (h64 >> 60) ^ (h64 & mask60)
}

func hash64(k Key) uint64 {
	var bs = k.ToByteSlice()
	var h = fnv.New64()
	h.Write(bs)
	return h.Sum64()
}

func Hash60(k Key) uint64 {
	return hash60(hash64(k))
}
