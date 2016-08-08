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
