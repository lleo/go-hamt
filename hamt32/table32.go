package hamt32

type node32I interface {
	hash30() uint32
	String() string
}

type table32I interface {
	node32I

	LongString(indent string, depth uint) string

	nentries() uint
	entries() []tableEntry32

	get(idx uint) node32I
	set(idx uint, entry node32I)
}

type tableEntry32 struct {
	idx  uint
	node node32I
}

//POPCNT Implementation
// copied from https://github.com/jddixon/xlUtil_go/blob/master/popCount.go
//  was MIT License

const (
	OCTO_FIVES  = uint32(0x55555555)
	OCTO_THREES = uint32(0x33333333)
	OCTO_ONES   = uint32(0x01010101)
	OCTO_FS     = uint32(0x0f0f0f0f)
)

func BitCount32(n uint32) uint {
	n = n - ((n >> 1) & OCTO_FIVES)
	n = (n & OCTO_THREES) + ((n >> 2) & OCTO_THREES)
	return uint((((n + (n >> 4)) & OCTO_FS) * OCTO_ONES) >> 24)
}
