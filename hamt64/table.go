package hamt64

type nodeI interface {
	hash60() uint64
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
	HEXI_FIVES  = uint64(0x5555555555555555)
	HEXI_THREES = uint64(0x3333333333333333)
	HEXI_ONES   = uint64(0x0101010101010101)
	HEXI_FS     = uint64(0x0f0f0f0f0f0f0f0f)
)

func BitCount64(n uint64) uint {
	n = n - ((n >> 1) & HEXI_FIVES)
	n = (n & HEXI_THREES) + ((n >> 2) & HEXI_THREES)
	return uint((((n + (n >> 4)) & HEXI_FS) * HEXI_ONES) >> 56)
}
