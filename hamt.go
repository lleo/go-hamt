/*
Package hamt is the unifying package between the 32bit and 64bit implementations
of Hash Array Mapped Tries (HAMT). HAMT datastructure make an efficient hashed
map data structure. You can `import hamt "github.com/lleo/go-hamt"`
then instantiate either a hamt32 or hamt64 datastructure with the
`hamt.NewHamt32()` or `hamt.NewHamt64()` functions. Both datastructures have
the same exported API defined by the Hamt interface.

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
indexes into a 32 entry table nodes for 32 bit HAMTs and 6 bit index into 64
entry table nodes for 64 bit HAMTs; isn't that symmetrical :).

For a this HAMT implementation, when key/value pair must be created, deleted,
or changed the key is hashed into a 30 or 60 bit value (described above) and
that hash30 or hash60 value represents a path of 5 or 6 bit values to place a
leaf containing the key, value pair. For a Get() or Del() operation we lookup
the deepest node along that pate that is not-nil. For a Put() operation we
lookup the deepest location that is nil and not beyond the lenth of the path.

You may implement your own Key type by implementeding the Key interface
defined in "github.com/lleo/go-hamt/key" or you may used the example
StringKey interface described in "github.com/lleo/go-hamt/stringkey".
*/
package hamt

import (
	"github.com/lleo/go-hamt/hamt32"
	"github.com/lleo/go-hamt/hamt64"
	"github.com/lleo/go-hamt/key"
)

// Hamt interface defines all behavior for implementations of the
// Hash Array Mapped Trie datastructures in hammt32/ and hamt64/.
type Hamt interface {
	Get(key.Key) (interface{}, bool)
	Put(key.Key, interface{}) bool
	Del(key.Key) (interface{}, bool)
	IsEmpty() bool
	String() string
}

// NewHamt32 ...
func NewHamt32() Hamt {
	//return hamt32.NewHamt()
	return hamt32.New(hamt32.HybridTables)
}

// NewHamt64 ...
func NewHamt64() Hamt {
	return hamt64.New(hamt64.HybridTables)
}
