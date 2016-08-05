package hamt64

type node64I interface {
	hash60() uint64
	String() string
}

type table64I interface {
	node64I

	LongString(indent string, depth uint) string

	nentries() uint
	entries() []tableEntry64

	get(idx uint) node64I
	set(idx uint, entry node64I)
}

type tableEntry64 struct {
	idx  uint
	node node64I
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
