// +build !go1.9

package hamt32

//POPCNT Implementation
// copied from https://github.com/jddixon/xlUtil_go/blob/master/popCount.go
//  was MIT License
// I found it explained at:
// http://stackoverflow.com/questions/22081738/how-does-this-algorithm-to-count-the-number-of-set-bits-in-a-32-bit-integer-work

const (
	hexiFives  = uint64(0x5555555555555555)
	hexiThrees = uint64(0x3333333333333333)
	hexiOnes   = uint64(0x0101010101010101)
	hexiFs     = uint64(0x0f0f0f0f0f0f0f0f)
)

func bitCount64(n uint64) uint {
	n = n - ((n >> 1) & hexiFives)
	n = (n & hexiThrees) + ((n >> 2) & hexiThrees)
	return uint((((n + (n >> 4)) & hexiFs) * hexiOnes) >> 56)
}
