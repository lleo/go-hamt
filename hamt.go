/*
Package hamt is just a trivial front door to the hamt32 and hamt64 packages
which really contain the HAMT implementations. Those HAMT implementations are
identical in every way but the size of the computed hash, called Hashval. Those
are either uint32 or uint64 values for hamt32 and hamt64 respectively. To repeat
myself, the hamt32 and hamt64 HAMT implementations are almost completely
identical code.

This package merely implements New32() and New64() functions and the table
option constants FixedTables, SparseTables, HybridTables, and the map
TableOptionName (eg. hamt.TableOptionName[hamt.FixedTables] ==
"FixedTables").

Choices

There are several choices to make: Hashval hamt32 versus hamt64, FixedTables
versus SparseTables versus HybridTables, and Functional versus
Transient. Then there is a hidden choice; you can change the source code
constant, IndexBits, to a value other than the current setting of 5.


Hashval hamt64 versus hamt32

Just use hamt64. I implemnted both before I really understood HAMT. I was
conflating 32 bit hash values with 32 wide branching factor (that was just a
feature of the other implmentations I was looking at).

While 32bit FNV hash values are still pretty random I have seen plenty of
collisions in my benchmarks.

I have never seen 64bit FNV hash values collide and in the current state of
computing having 64bit CPUs as the norm. I recommend using hamt64. If you are
on 32bit CPUs then maybe you could choose hamt32.

FixedTables versus SparseTables versus HybridTables

This is the classic speed versus memory choice with a twist. The facts to
consider are: The tree is indexed by essentially random values (the parts of the
hash value of the key), so the tree is going to be "balanced" to a statistical
likelihood. The inner branching nodes will be very densely populated, and the
outer branching nodes will probably be very sparsely populated.

FixedTables are fastest to access and modify, because it is a simple mater of
getting and setting preallocated fixed sized arrays. However, they will be
wasting most of their allocated space most of the time.

SparseTables are slowest because we always have to calculate bit counts of bit
arrays for get or set type operations. Further, inserting or removing values
is a matter of manipulating slice values. On the other hand, given the way
slice memory allocation works, they will usually waste less than half their
allocated memory.

If you are going to have very large number of entries in your HAMT to the degree
that your program will start paging occationally, DEFINITELY switch away from
FixedTables to HybridTables or even SparseTables.

As a general rule if the number of entry is only in the 100,000s or less, choose
FixedTables

Transient versus Ffunctional

Good question!

IndexBits

Both hamt32 and hamt64 have a constant IndexBits which determines all the other
constants defining the HAMT structures. For both hamt32 and hamt64, the
IndexBits constant is set to 5. You can manually change the
source code to set IndexBits to some uint other than 5. IndexBits is set to 5
because that is how other people do it.

IndexBits determines the branching factor (IndexLimit) and the depth
(DepthLimit) of the HAMT data structure. Given IndexBits=5 IndexLimit=32, and
DepthLimit=6 for hamt32 and DepthLimit=10 for hamt64.

*/
package hamt

import (
	"github.com/lleo/go-hamt/hamt32"
	"github.com/lleo/go-hamt/hamt64"
)

// WARNING!!! Duplicated code with both hamt32 and hamt64. Must have
// test to guarantee they stay in lock step.
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

// New32() takes two arguments and producest a value that conforms to the
// hamt32.Hamt interface. The arguments are a bool and an int. The bool argument
// determines if a functional structure(true) or transient stucture(false) is
// produced. The int option is either 0, 1, or 2 conforming to the
// FixedTables, SparseTables, or HybridTables constants.
func New32(functional bool, opt int) hamt32.Hamt {
	return hamt32.New(functional, opt)
}

// New32() takes two arguments and producest a value that conforms to the
// hamt64.Hamt interface. The arguments are a bool and an int. The bool argument
// determines if a functional structure(true) or transient stucture(false) is
// produced. The int option is either 0, 1, or 2 conforming to the
// FixedTables, SparseTables, or HybridTables constants.
func New64(functional bool, opt int) hamt64.Hamt {
	return hamt64.New(functional, opt)
}
