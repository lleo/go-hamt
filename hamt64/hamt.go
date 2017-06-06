package hamt64

import "github.com/lleo/go-hamt-key"

const nBits uint = key.BitsPerLevel60

const maxDepth uint = key.MaxDepth60

const tableCapacity uint = uint(1 << nBits)

const downgradeThreshold uint = 21 // floor(tableCapacity / 3)

const upgradeThreshold uint = 43 // round(tableCapacity * 2 / 3)

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
	// initially, then upgrad to fullTable when appropriate.
	HybridTables
)

// TableOptionName is a lookup table to map the integer value of FullTablesOnly,
// CompTablesOnly, and HybridTables to a string representing that option.
//     var option = hamt32.FullTablesOnly
//     hamt32.TableOptionName[option] == "FullTablesOnly"
var TableOptionName [3]string

// Could have used...
//var TableOptionName = [3]string{
//	"FullTablesOnly",
//	"CompTablesOnly",
//	"HybridTables",
//}

func init() {
	TableOptionName[HybridTables] = "HybridTables"
	TableOptionName[CompTablesOnly] = "CompTablesOnly"
	TableOptionName[FullTablesOnly] = "FullTablesOnly"
}

type Hamt interface {
	IsEmpty() bool
	Nentries() uint
	ToFunctional() Hamt
	ToTransient() Hamt
	Get(key.Key) (interface{}, bool)
	Put(key.Key, interface{}) (Hamt, bool)
	Del(key.Key) (Hamt, interface{}, bool)
	String() string
	LongString(string) string
}

func New(functional bool, opt int) Hamt {
	if functional {
		return NewFunctional(opt)
	}
	return NewTransient(opt)
}
