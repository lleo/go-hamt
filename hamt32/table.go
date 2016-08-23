package hamt32

type nodeI interface {
	hash30() uint32
	String() string
}

type tableI interface {
	nodeI

	LongString(indent string, depth uint) string

	nentries() uint
	entries() []tableEntry

	get(idx uint) nodeI
	set(idx uint, entry nodeI)
}

type tableEntry struct {
	idx  uint
	node nodeI
}

//POPCNT Implementation
// copied from https://github.com/jddixon/xlUtil_go/blob/master/popCount.go
//  was MIT License

const (
	octo_fives  = uint32(0x55555555)
	octo_threes = uint32(0x33333333)
	octo_ones   = uint32(0x01010101)
	octo_fs     = uint32(0x0f0f0f0f)
)

func bitCount32(n uint32) uint {
	n = n - ((n >> 1) & octo_fives)
	n = (n & octo_threes) + ((n >> 2) & octo_threes)
	return uint((((n + (n >> 4)) & octo_fs) * octo_ones) >> 24)
}
