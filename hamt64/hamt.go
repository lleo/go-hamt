/*
Package hamt64 is the package that implements two Hamt structures for both
functional and transient implementations. The first structure is HamtFunctional,
and the second is HamtTransient. Each of these datastructures implemnents the
hamt64.Hamt interface.
*/
package hamt64

const tableCapacity uint = IndexLimit

// DowngradeThreshold is the constant that sets the threshold for the size of a
// table, that when a table decreases to the threshold size, the table is
// converted from a FullTable to a CompressedTable.
//
// This conversion only happens if the Hamt structure has be constructed with
// the HybridTables option.
const DowngradeThreshold uint = 10 // floor(tableCapacity / 3)

// UpgradeThreshold is the constant that sets the threshold for the size of a
// table, that when a table increases to the threshold size, the table is
// converted from a CompressedTable to a FullTable.
//
// This conversion only happens if the Hamt structure has be constructed with
// the HybridTables option.
const UpgradeThreshold uint = 21 // round(tableCapacity * 2 / 3)

// Configuration contants to be passed to `hamt64.New(int) *Hamt`.
const (
	// FullTableOnly indicates the structure should use fullTables ONLY.
	// This was intended to be for speed, as compressed tables use a software
	// bitCount function to access individual cells.
	FullTablesOnly = iota
	// CompTablesOnly indicates the structure should use compressedTables ONLY.
	// This was intended just save space, but also seems to be faster; CPU cache
	// locality maybe?
	CompTablesOnly
	// HybridTables indicates the structure should use compressedTable
	// initially, then upgrade to fullTable when appropriate.
	HybridTables
)

// TableOptionName is a lookup table to map the integer value of FullTablesOnly,
// CompTablesOnly, and HybridTables to a string representing that option.
//     var option = hamt64.FullTablesOnly
//     hamt64.TableOptionName[option] == "FullTablesOnly"
var TableOptionName [3]string

// Could have used...
//var TableOptionName = [3]string{
//	"FullTablesOnly",
//	"CompTablesOnly",
//	"HybridTables",
//}

func init() {
	TableOptionName[FullTablesOnly] = "FullTablesOnly"
	TableOptionName[CompTablesOnly] = "CompTablesOnly"
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
