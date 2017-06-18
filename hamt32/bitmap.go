package hamt32

import (
	"fmt"
	"strings"
)

// BitmapShift is 5 because we are using uint32 as the base Bitmap type.
const BitmapShift uint = 5

// BitmapSize is the number of uint32 needed to cover IndexLimit bits.
const BitmapSize uint = (IndexLimit + (1 << BitmapShift) - 1) / (1 << BitmapShift)

type Bitmap [BitmapSize]uint32

func (bm *Bitmap) String() string {
	if BitmapSize == 1 {
		//only show IndexLimit bits
		var fmtStr = fmt.Sprintf("%%0%db", IndexLimit)
		return fmt.Sprintf(fmtStr, bm[0])
	}

	// Show all bits in Bitmap because IndexLimit is a multiple of the
	// Bitmap base type.
	var strs = make([]string, BitmapSize)
	var fmtStr = fmt.Sprintf("%%0%db", 1<<BitmapShift)
	for i := uint(0); i < BitmapSize; i++ {
		strs[i] = fmt.Sprintf(fmtStr, bm[i])
	}

	return strings.Join(strs, " ")
}

// IsSet returns a bool indicating whether the i'th bit is 1(true) or 0(false).
func (bm *Bitmap) IsSet(i uint) bool {
	var idx = i >> BitmapShift
	var bit = i & ((1 << BitmapShift) - 1)

	return (bm[idx] & (1 << bit)) > 0
}

// Set marks the i'th bit 1.
func (bm *Bitmap) Set(i uint) {
	var idx = i >> BitmapShift
	var bit = i & ((1 << BitmapShift) - 1)

	bm[idx] |= (1 << bit)

	return
}

// Unset marks the i'th bit to 0.
func (bm *Bitmap) Unset(i uint) {
	var idx = i >> BitmapShift
	var bit = i & ((1 << BitmapShift) - 1)

	if bm[idx]&(1<<bit) > 0 {
		bm[idx] &^= 1 << bit
	}

	return
}

// Count returns the numbers of bits set below the i'th bit.
func (bm *Bitmap) Count(i uint) uint {
	var idx = i >> BitmapShift
	var bit = i & ((1 << BitmapShift) - 1)

	var count uint
	for idxN := uint(0); idxN < idx; idxN++ {
		count += bitCount32(bm[idxN])
	}

	count += bitCount32(bm[idx] & ((1 << bit) - 1))

	return count
}
