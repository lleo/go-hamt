package hamt32

import (
	"context"
	"unsafe"
)

// hashSize is the size of hashVal in bits.
const hashSize uint = uint(unsafe.Sizeof(hashVal(0))) * 8

// NumIndexBits is the fundemental setting for the Hamt data structure. Given
// that we hash every key ([]byte slice) into a hashVal, that hashVal must be
// split into DepthLimit number of NumIndexBits wide parts. Each of those parts
// of the hashVal is used as the index into the given level of the Hamt tree.
// So NumIndexBits determines how wide and how deep the Hamt can be.
const NumIndexBits uint = 5

// DepthLimit is the maximum number of levels of the Hamt. It is calculated as
// DepthLimit = floor(hashSize / NumIndexBits) or a strict integer division.
const DepthLimit = hashSize / NumIndexBits
const remainder = hashSize - (DepthLimit * NumIndexBits)

// IndexLimit is the maximum number of entries in a Hamt interior node. In other
// words it is the width of the Hamt data structure.
const IndexLimit = 1 << NumIndexBits

// maxDepth is the maximum value of a depth variable. maxDepth = DepthLimit - 1
const maxDepth = DepthLimit - 1

// maxIndex is the maximum value of a index variable. maxIndex = IndexLimit - 1
const maxIndex = IndexLimit - 1

// DowngradeThreshold is the constant that sets the threshold for the size of a
// table, such that when a table decreases to the threshold size, the table is
// converted from a FixedTable to a SparseTable.
//
// This conversion only happens if the Hamt structure has be constructed with
// the HybridTables option.
const DowngradeThreshold uint = IndexLimit * 3 / 8 //12 for NumIndexBits=5

// UpgradeThreshold is the constant that sets the threshold for the size of a
// table, such that when a table increases to the threshold size, the table is
// converted from a SparseTable to a FixedTable.
//
// This conversion only happens if the Hamt structure has be constructed with
// the HybridTables option.
const UpgradeThreshold uint = IndexLimit / 2 //16 for NumIndexBits=5

// Configuration contants to be passed to `hamt32.New(int) *Hamt`.
const (
	// FixedTable indicates the structure should use fixedTables ONLY.
	// This was intended to be for speed, as compressed tables use a software
	// bitCount function to access individual cells.
	FixedTables = iota
	// SparseTables indicates the structure should use sparseTables ONLY.
	// This was intended just save space, but also seems to be faster; CPU cache
	// locality maybe?
	SparseTables
	// HybridTables indicates the structure should use sparseTable
	// initially, then upgrade to fixedTable when appropriate.
	HybridTables
)

// TableOptionName is a lookup table to map the integer value of
// FixedTables, SparseTables, and HybridTables to a string representing
// that option.
//     var option = hamt32.FixedTables
//     hamt32.TableOptionName[option] == "FixedTables"
var TableOptionName [3]string

// Could have used...
//var TableOptionName = [3]string{
//	"FixedTables",
//	"SparseTables",
//	"HybridTables",
//}

func init() {
	TableOptionName[FixedTables] = "FixedTables"
	TableOptionName[SparseTables] = "SparseTables"
	TableOptionName[HybridTables] = "HybridTables"
}

// Hamt defines the interface that both the HamtFunctional and HamtTransient
// data structures must (and do) implement.
type Hamt interface {
	IsEmpty() bool
	Nentries() uint
	ToFunctional() Hamt
	ToTransient() Hamt
	DeepCopy() Hamt
	Get([]byte) (interface{}, bool)
	Put([]byte, interface{}) (Hamt, bool)
	Del([]byte) (Hamt, interface{}, bool)
	String() string
	LongString(string) string
	visit(visitFn) uint
	Stats() *Stats
	Iter() IterFunc
	IterChan(chanBufLen int) <-chan KeyVal
	IterChanWithContext(chanBufLen int, ctx context.Context) <-chan KeyVal
}

// New constructs a datastucture that implements the Hamt interface. When the
// functional argument is true it implements a HamtFunctional data structure.
// When the functional argument is false it implements a HamtTransient
// data structure. In either case the opt argument is handed to the to the
// contructore for either NewFunctional(opt) or NewTransient(opt).
func New(functional bool, opt int) Hamt {
	if functional {
		return NewFunctional(opt)
	}
	return NewTransient(opt)
}

type Stats struct {
	// Depth of deepest table
	MaxDepth uint

	// TableCountsByNentries is a Hash table of the number of tables with each
	// given number of entries in the tatble. There are slots for
	// [0..IndexLimit] inclusive (so there are IndexLimit+1 slots). Technically,
	// there should never be a table with zero entries, but I allow counting
	// tables with zero entries just to catch those errors.
	TableCountsByNentries [IndexLimit + 1]uint // [0..IndexLimit] inclusive

	// TableCountsByDepth is a Hash table of the number of tables at a given
	// depth. There are slots for [0..DepthLimit).
	TableCountsByDepth [DepthLimit]uint // [0..DepthLimit)

	// Nils is the total count of allocated slots that are unused in the HAMT.
	Nils uint

	// Nodes is the total count of nodeI capable structs in the HAMT.
	Nodes uint

	// Tables is the total count of tableI capable structs in the HAMT.
	Tables uint

	// Leafs is the total count of leafI capable structs in the HAMT.
	Leafs uint

	// FixedTables is the total count of fixedTable structs in the HAMT.
	FixedTables uint

	// SparseTables is the total count of sparseTable structs in the HAMT.
	SparseTables uint

	// FlatLeafs is the total count of flatLeaf structs in the HAMT.
	FlatLeafs uint

	// CollisionLeafs is the total count of collisionLeaf structs in the HAMT.
	CollisionLeafs uint

	// KeyVals is the total number of Key,Val pairs int the HAMT.
	KeyVals uint
}
