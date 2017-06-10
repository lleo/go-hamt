/*
Package hamt is the unifying package between the 32bit and 64bit implementations
of Hash Array Mapped Tries (HAMT). HAMT datastructure make an efficient hashed
map data structure. You can `import "github.com/lleo/go-hamt"` then instantiate
either a hamt32 or hamt64 datastructure with the `hamt.New32()` or
`hamt.New64()` functions. Both datastructures have the same exported API defined
by the hamt32.Hamt and hamt64.Hamt interfaces.

Given how wide a HAMT node is (either 32 or 64 nodes wide) HAMT datastructures
not very deep; either 6, for 32bit, or 10, for 64bit implementations, nodes
deep. This neans HAMTs are effectively O(1) for Search, Insertions, and
Deletions.

Both 32 and 64 bit implementations of HAMTs are of fixed depth is because they
are [Tries](https://en.wikipedia.org/wiki/Trie). The key of a Trie is split
into n-number smaller indecies and each node from the root uses each successive
index.

In the case of a this HAMT implementation the key is hashed into a 30 or 60 bit
number. In the case of the stringkey we take the []byte slice of the string
and feed it to hash.fnv.New32() or New64() hash generator. Since these
generate 32 and 64 bit hash values respectively and we need 30 and 60 bit
values, we use the [xor-fold technique](http://www.isthe.com/chongo/tech/comp/fnv/index.html#xor-fold)
to "fold" the high 2 or 4 bits of the 32 and 64 bit hash values into 30 and
60 bit values for our needs.

We want 30 and 60 bit values because they split nicely into six 5bit and ten
6bit values respectively. Each of these 5 and 6 bit values become the indexies
of our Trie nodes with a maximum depth of 6 or 10 respectively. Further 5 bits
indexe into a 32 entry table nodes for 32 bit HAMTs and 6 bit index into 64
entry table nodes for 64 bit HAMTs; isn't that symmetrical :).

For a this HAMT implementation, when key/value pair must be created, deleted,
or changed the key is hashed into a 30 or 60 bit value (described above) and
that hash30 or hash60 value represents a path of 5 or 6 bit values to place a
leaf containing the key, value pair. For a Get() or Del() operation we lookup
the deepest node along that path that is not-nil. For a Put() operation we
lookup the deepest location that is a leaf or nil and not beyond the lenth of
the path.

You may implement your own Key type by implementeding the Key interface
defined in "github.com/lleo/go-hamt/key" or you may used the example
StringKey interface described in "github.com/lleo/go-hamt/stringkey".
*/
package hamt

import (
	"github.com/lleo/go-hamt/hamt32"
	"github.com/lleo/go-hamt/hamt64"
)

// WARNING!!! Duplicated code with both hamt32 and hamt64. Must have
// test to guarantee they stay in lock step.
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
	TableOptionName[FullTablesOnly] = "FullTablesOnly"
	TableOptionName[CompTablesOnly] = "CompTablesOnly"
	TableOptionName[HybridTables] = "HybridTables"
}

// New32() takes two arguments and producest a value that conforms to the
// hamt32.Hamt interface. The arguments are a bool and an int. The bool argument
// determines if a functional structure(true) or transient stucture(false) is
// produced. The int option is either 0, 1, or 2 conforming to the
// FullTablesOnly, CompTablesOnly, or HybridTables constants.
func New32(functional bool, opt int) hamt32.Hamt {
	return hamt32.New(functional, opt)
}

// New32() takes two arguments and producest a value that conforms to the
// hamt64.Hamt interface. The arguments are a bool and an int. The bool argument
// determines if a functional structure(true) or transient stucture(false) is
// produced. The int option is either 0, 1, or 2 conforming to the
// FullTablesOnly, CompTablesOnly, or HybridTables constants.
func New64(functional bool, opt int) hamt64.Hamt {
	return hamt64.New(functional, opt)
}
