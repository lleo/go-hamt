// +build !go1.9

package hamt32

//POPCNT Implementation
// copied from https://github.com/jddixon/xlUtil_go/blob/master/popCount.go
//  was MIT License
//
// I found it explained at:
// http://stackoverflow.com/questions/22081738/how-does-this-algorithm-to-count-the-number-of-set-bits-in-a-32-bit-integer-work

const (
	octoFives  = uint32(0x55555555)
	octoThrees = uint32(0x33333333)
	octoOnes   = uint32(0x01010101)
	octoFs     = uint32(0x0f0f0f0f)
)

func bitCount32(n uint32) uint {
	n = n - ((n >> 1) & octoFives)
	n = (n & octoThrees) + ((n >> 2) & octoThrees)
	return uint((((n + (n >> 4)) & octoFs) * octoOnes) >> 24)
}
