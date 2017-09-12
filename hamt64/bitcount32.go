// +build go1.9

package hamt64

import (
	"math/bits"
)

func bitCount32(n uint32) uint {
	// Using go1.9+ implementation of popcount as it uses ASM when available.
	return uint(bits.OnesCount32(n))
}
