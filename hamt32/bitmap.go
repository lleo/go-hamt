package hamt32

import (
	"fmt"
	"strings"
)

// bitmapShift is 5 because we are using uint32 as the base bitmap type.
const bitmapShift uint = 5

// bitmapSize is the number of uint32 needed to cover IndexLimit bits.
const bitmapSize uint = (IndexLimit + (1 << bitmapShift) - 1) / (1 << bitmapShift)

type bitmap [bitmapSize]uint32

func (bm *bitmap) String() string {
	if bitmapSize == 1 {
		//only show IndexLimit bits
		var fmtStr = fmt.Sprintf("%%0%db", IndexLimit)
		return fmt.Sprintf(fmtStr, bm[0])
	}

	// Show all bits in bitmap because IndexLimit is a multiple of the
	// bitmap base type.
	var strs = make([]string, bitmapSize)
	var fmtStr = fmt.Sprintf("%%0%db", 1<<bitmapShift)
	for i := uint(0); i < bitmapSize; i++ {
		strs[i] = fmt.Sprintf(fmtStr, bm[i])
	}

	return strings.Join(strs, " ")
}

// IsSet returns a bool indicating whether the i'th bit is 1(true) or 0(false).
func (bm *bitmap) IsSet(idx uint) bool {
	var nth = idx >> bitmapShift
	var bit = idx & ((1 << bitmapShift) - 1)

	return (bm[nth] & (1 << bit)) > 0
}

// Set marks the i'th bit 1.
func (bm *bitmap) Set(idx uint) {
	var nth = idx >> bitmapShift
	var bit = idx & ((1 << bitmapShift) - 1)

	bm[nth] |= 1 << bit

	return
}

// Unset marks the i'th bit to 0.
func (bm *bitmap) Unset(idx uint) {
	var nth = idx >> bitmapShift
	var bit = idx & ((1 << bitmapShift) - 1)

	//if bm[nth]&(1<<bit) > 0 {
	//	bm[nth] &^= 1 << bit
	//}
	bm[nth] &^= 1 << bit

	return
}

// Count returns the numbers of bits set below the i'th bit.
func (bm *bitmap) Count(idx uint) uint {
	var nth = idx >> bitmapShift
	var bit = idx & ((1 << bitmapShift) - 1)

	var count uint
	for i := uint(0); i < nth; i++ {
		count += bitCount32(bm[i])
	}

	count += bitCount32(bm[nth] & ((1 << bit) - 1))

	return count
}
