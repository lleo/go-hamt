package hamt32

import "github.com/lleo/go-hamt-key"

const nBits uint = key.BitsPerLevel30

const maxDepth uint = key.MaxDepth30

const tableCapacity uint = uint(1 << nBits)

const downgradeThreshold uint = 10 // floor(tableCapacity / 3)

const upgradeThreshold uint = 21 // round(tableCapacity * 2 / 3)

// Configuration contants to be passed to `hamt32.New(int) *Hamt`.
const (
	// HybridTables indicates the structure should use compressedTable
	// initially, then upgrad to fullTable when appropriate.
	HybridTables = iota //0
	// CompTablesOnly indicates the structure should use compressedTables ONLY.
	// This was intended just save space, but also seems to be faster; CPU cache
	// locality maybe?
	CompTablesOnly //1
	// FullTableOnly indicates the structure should use fullTables ONLY.
	// This was intended to be for speed, as compressed tables use a software
	// bitCount function to access individual cells. Turns out, not so much.
	FullTablesOnly //2
)

// TableOptionName is a map of the table option value Hybrid, CompTablesOnly,
// or FullTableOnly to a string representing that option.
//      var options = hamt32.FullTablesOnly
//      hamt32.TableOptionName[hamt32.FullTablesOnly] == "FullTablesOnly"
var TableOptionName = make(map[int]string, 3)

func init() {
	TableOptionName[HybridTables] = "HybridTables"
	TableOptionName[CompTablesOnly] = "CompTablesOnly"
	TableOptionName[FullTablesOnly] = "FullTablesOnly"
}

type Hamt interface {
	IsEmpty() bool
	Nentries() uint
	Get(key.Key) (interface{}, bool)
	Put(key.Key, interface{}) (Hamt, bool)
	Del(key.Key) (Hamt, interface{}, bool)
	//String() string
	//LongString(string) string
}
