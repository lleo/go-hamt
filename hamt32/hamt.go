/*
Package hamt32 is the package that implements two Hamt structures for both
functional and transient implementations. The first structure is HamtFunctional,
and the second is HamtTransient. Each of these datastructures implemnents the
hamt32.Hamt interface.
*/
package hamt32

const tableCapacity uint = IndexLimit

// DowngradeThreshold is the constant that sets the threshold for the size of a
// table, that when a table decreases to the threshold size, the table is
// converted from a FixedTable to a SparseTable.
//
// This conversion only happens if the Hamt structure has be constructed with
// the HybridTables option.
const DowngradeThreshold uint = 10 // floor(tableCapacity / 3)

// UpgradeThreshold is the constant that sets the threshold for the size of a
// table, that when a table increases to the threshold size, the table is
// converted from a SparseTable to a FixedTable.
//
// This conversion only happens if the Hamt structure has be constructed with
// the HybridTables option.
const UpgradeThreshold uint = 21 // round(tableCapacity * 2 / 3)

// Configuration contants to be passed to `hamt32.New(int) *Hamt`.
const (
	// FixedTableOnly indicates the structure should use fixedTables ONLY.
	// This was intended to be for speed, as compressed tables use a software
	// bitCount function to access individual cells.
	FixedTablesOnly = iota
	// SparseTablesOnly indicates the structure should use sparseTables ONLY.
	// This was intended just save space, but also seems to be faster; CPU cache
	// locality maybe?
	SparseTablesOnly
	// HybridTables indicates the structure should use sparseTable
	// initially, then upgrade to fixedTable when appropriate.
	HybridTables
)

// TableOptionName is a lookup table to map the integer value of FixedTablesOnly,
// SparseTablesOnly, and HybridTables to a string representing that option.
//     var option = hamt32.FixedTablesOnly
//     hamt32.TableOptionName[option] == "FixedTablesOnly"
var TableOptionName [3]string

// Could have used...
//var TableOptionName = [3]string{
//	"FixedTablesOnly",
//	"SparseTablesOnly",
//	"HybridTables",
//}

func init() {
	TableOptionName[FixedTablesOnly] = "FixedTablesOnly"
	TableOptionName[SparseTablesOnly] = "SparseTablesOnly"
	TableOptionName[HybridTables] = "HybridTables"
}

// Hamt defines the interface that both the HamtFunctional and HamtTransient
// datastructures must (and do) implement.
type Hamt interface {
	IsEmpty() bool
	Nentries() uint
	ToFunctional() Hamt
	ToTransient() Hamt
	DeepCopy() Hamt
	Get(Key) (interface{}, bool)
	Put(Key, interface{}) (Hamt, bool)
	Del(Key) (Hamt, interface{}, bool)
	String() string
	LongString(string) string
}

// New() constructs a datastucture that implements the Hamt interface. When the
// functional argument is true it implements a HamtFunctional datastructure.
// When the functional argument is false it implements a HamtTransient
// datastructure. In either case the opt argument is handed to the to the
// contructore for either NewFunctional(opt) or NewTransient(opt).
func New(functional bool, opt int) Hamt {
	if functional {
		return NewFunctional(opt)
	}
	return NewTransient(opt)
}
