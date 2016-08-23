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
	hexi_fives  = uint64(0x5555555555555555)
	hexi_threes = uint64(0x3333333333333333)
	hexi_ones   = uint64(0x0101010101010101)
	hexi_fs     = uint64(0x0f0f0f0f0f0f0f0f)
)

func bitCount64(n uint64) uint {
	n = n - ((n >> 1) & hexi_fives)
	n = (n & hexi_threes) + ((n >> 2) & hexi_threes)
	return uint((((n + (n >> 4)) & hexi_fs) * hexi_ones) >> 56)
}
