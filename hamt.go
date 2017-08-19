/*
Package hamt is just a trivial front door to the hamt32 and hamt64 packages
which really contain the HAMT implementations. Those HAMT implementations are
identical in every way but the size of the computed hash, called Hashval. Those
are either uint32 or uint64 values for hamt32 and hamt64 respectively.

This package merely implements New(), New32() and New64() functions and the
table option constants FixedTables, SparseTables, HybridTables, and the map
TableOptionName (eg. hamt.TableOptionName[hamt.FixedTables] ==
"FixedTables").

Choices

There are several choices to make: Hashval hamt32 versus hamt64, FixedTables
versus SparseTables versus HybridTables, and Functional versus
Transient. Then there is a hidden choice; you can change the source code
constant, NumIndexBits, to a value other than the current setting of 5.

The New() function makes all the recommended choices for you. That is it
uses the 64 bit hashVal (aka hamt64), functional behavior, and hybrid tables.

Hashval hamt64 versus hamt32

Just use hamt64. I implemnted both before I really understood HAMT. I was
conflating 32 bit hash values with 32 wide branching factor (that was just a
feature of the other implmentations I was looking at).

While 32bit FNV hash values are still pretty random I have seen plenty of
collisions in my tests.

I have never seen 64bit FNV hash values collide and in the current state of
computing having 64bit CPUs as the norm. I recommend using hamt64. If you are
on 32bit CPUs then maybe you could choose hamt32.

FixedTables versus SparseTables versus HybridTables

Tables is the name I use to for the interior (aka non-leaf) nodes of the hamt
data structure. Those tables being the indexed arrays from the Hash-indexed
Array Mapped Tries (HAMT) name.

This is the classic speed versus memory choice with a twist. The facts to
consider are: The tree is indexed by essentially random values (the parts of the
hash value of the key), so the tree is going to be "balanced" to a statistical
likelihood. The inner branching nodes will be very densely populated, and the
outer branching nodes will be very sparsely populated.

FixedTables are fastest to access and modify, because it is a simple matter of
getting and setting preallocated fixed sized arrays. However, they will be
wasting most of their allocated space most of the time.

SparseTables are slowest because we always have to calculate bit counts of bit
arrays for get or set type operations. Further, inserting or removing values
is a matter of manipulating slice values. On the other hand, given the way
slice memory allocation works, they will usually waste less than half their
allocated memory.

According to tests, HybridTables setting behaves precisely the way we want it
to behave. For a test set of data with 3,149,824 KeyVal pairs, he distribution
of tables comes in two groups: tables with 25-32 entries and tables with 1-11
entries. There are no tables outside those two groupings. The 25-32 entry
tables are all fixed tables and the 1-11 entry tables are all sparse tables.
Of the sparse tables %40.1 have 1 or 2 entries, %85.4 have 4 or less and
%99.7 have 8 or less entries. Given sparse tables start at capacity of 2 and
capacity grows by doubling, the sparse tables are efficiently packed. The
conclusion from this test data is that HybridTables setting is an excellent
fit for memory efficiency.

Benchmarks show that fixed table hamts are slower than hybrid table hamts for
Put operations and comparable for Get & Del operations. I conclude that is due
to the massive over use of memory in fixed tables. This conclusion is partly
due to the fact that the Put operation differntial between fixed and hybrid
table hamts is twice as large for functional versus transient hamt behavior.

Clearly HybridTables table option for my HAMT data structure is the best
choice.

Transient versus Functional

The bottom line is that writing to transient behavior in a multiple
threads almost guarantees problems unless you implement a locking solution
(and that can be hard to do in a performant manner).

On the other hand, given that HamtFunctional behavior returns a new
HamtFunctional data structure upon any modification, HamtFunctional data
structures are inherently thread safe.

On your third hand, the copy-on-write strategy of HamtFunctional is inherently
slower than modify-in-place strategy of HamtTransient. How much slower? For
large hamt data structures (~3 million key/value pairs) the transient Put
operation takes ~1000ns, where the functional Put op takes ~3200ns. Which
really isn't that bad because they are within the same order of magnitude and
it is already fast.

Using the transient behavior begs the question: why not use the Go builtin map?
Of course, the reason is obvious if your key must be a slice; Go map keys can
not be slices. The ability to implement a reasonably efficient functional
behavior for HAMTs is the point of this library. The hamt transient speed
is definitely is slower than Go's builtin map: 435 vs 130 ns/Get; 950 vs 235
ns/Put; 900 vs 175 ns/Del; 235 vs 23 ns/KeyVal iterate.

The point of this library is to implement the functional map-like behavior, so
that is what I assume you will use it for. The transient behavior is useful for
faster single threaded bulk operations then transform it back to the functional
behavior.

NumIndexBits

Both hamt32 and hamt64 have a constant NumIndexBits which determines all the
other constants defining the HAMT structures. For both hamt32 and hamt64, the
NumIndexBits constant is set to 5. You can manually change the
source code to set NumIndexBits to some uint other than 5. IndexBits is set to 5
because that is how other people do it.

NumIndexBits determines the branching factor (IndexLimit) and the depth
(DepthLimit) of the HAMT data structure. Given IndexBits=5 IndexLimit=32, and
DepthLimit=6 for hamt32 and DepthLimit=12 for hamt64.

*/
package hamt

import (
	"github.com/lleo/go-hamt/hamt32"
	"github.com/lleo/go-hamt/hamt64"
)

const (
	// FixedTable indicates the structure should use fixedTables ONLY.
	// This was intended to be for speed, as sparse tables use a software
	// bitCount function to access individual cells.
	FixedTables = iota
	// SparseTables indicates the structure should use sparseTable's ONLY.
	// This was intended just save space, but also seems to be faster; CPU cache
	// locality maybe?
	SparseTables
	// HybridTables indicates the structure should use sparseTable
	// initially, then upgrad to fixedTable when appropriate.
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

// New() makes all the configuration choices for you. Specifically, it chooses
// functional behavior, 64bit hashes, and Hybrid tables. These are the
// recommended settings. See hamt64.Hamt for the API.
func New() hamt64.Hamt {
	return hamt64.New(true, HybridTables)
}

// New32 constructs a datastucture that implements the hamt32.Hamt interface.
//
// When the functional argument is true it implements a hamt32.HamtFunctional
// data structure. When the functional argument is false it implements a
// HamtTransient data structure.
//
// The tblOpt argument is the table option defined by the constants
// HybridTables, SparseTables, xor FixedTables.
//
func New32(functional bool, opt int) hamt32.Hamt {
	return hamt32.New(functional, opt)
}

// New64 constructs a datastucture that implements the hamt64.Hamt interface.
//
// When the functional argument is true it implements a hamt64.HamtFunctional
// data structure. When the functional argument is false it implements a
// HamtTransient data structure.
//
// The tblOpt argument is the table option defined by the constants
// HybridTables, SparseTables, xor FixedTables.
//
func New64(functional bool, opt int) hamt64.Hamt {
	return hamt64.New(functional, opt)
}
