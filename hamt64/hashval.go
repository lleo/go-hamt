package hamt64

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

// hashVal sets the numberer of bits of the hash value by being an alias to
// uint64 and establishes a type we can hang methods, like Index(), off of.
type hashVal uint64

func indexMask(depth uint) hashVal {
	return hashVal((1<<IndexBits)-1) << (depth * IndexBits)
}

// Index returns the IndexBits bit value of the hashVal at 'depth' number of
// IndexBits number of bits into hashVal.
func (hv hashVal) Index(depth uint) uint {
	_ = AssertOn && assert(depth < DepthLimit, "Index: depth > MaxDepth")

	var idxMask = indexMask(depth)
	return uint((hv & idxMask) >> (depth * IndexBits))
}

func hashPathMask(depth uint) hashVal {
	return hashVal(1<<((depth)*IndexBits)) - 1
}

// hashPath calculates the path required to read the given depth. In other words
// it returns a hashVal that preserves the first depth-1 IndexBits index
// values. For depth=0 it always returns no path (aka a 0 value).
// For depth=MaxDepth it returns all but the last index value.
func (hv hashVal) hashPath(depth uint) hashVal {
	_ = AssertOn && assert(depth < DepthLimit, "hashPath(): dept > MaxDepth")

	if depth == 0 {
		return 0
	}

	return hv & hashPathMask(depth)

}

// buildHashPath method adds a idx at depth level of the hashPath. Given a
// hash Path = "/11/07/13" and you call hashPath.buildHashPath(23, 3) the method
// will return hashPath "/11/07/13/23". hashPath is shown here in the string
// representation, but the real value is hashVal (aka uint64).
func (hv hashVal) buildHashPath(idx, depth uint) hashVal {
	_ = AssertOn && assert(idx < DepthLimit, "buildHashPath: idx > MaxIndex")

	hv &= hashPathMask(depth)
	return hv | hashVal(idx<<(depth*IndexBits))
}

// HashPathString returns a string representation of the index path of a
// hashVal. It will be string of the form "/idx0/idx1/..." where each idxN value
// will be a zero padded number between 0 and MaxIndex. There will be limit
// number of such values where limit <= DepthLimit.
// If the limit parameter is 0 then the method will simply return "/".
// Warning: It will panic() if limit > DepthLimit.
// Example: "/00/24/46/17" for limit=4 of a IndexBits=5 hash value
// represented by "/00/24/46/17/34/08".
func (hv hashVal) HashPathString(limit uint) string {
	_ = AssertOn && assertf(limit <= DepthLimit,
		"HashPathString: limit,%d > DepthLimit,%d\n", limit, DepthLimit)

	if limit == 0 {
		return "/"
	}

	var strs = make([]string, limit)

	for d := uint(0); d < limit; d++ {
		var idx = hv.Index(d)
		strs[d] = fmt.Sprintf("%02d", idx)
	}

	return "/" + strings.Join(strs, "/")
}

// bitString returns a hashVal as a string of bits separated into groups of
// IndexBits bits.
func (hv hashVal) bitString() string {
	var strs = make([]string, DepthLimit)

	for d := uint(0); d < DepthLimit; d++ {
		var fmtStr = fmt.Sprintf("%%0%dd", IndexBits)
		strs[MaxDepth-d] = fmt.Sprintf(fmtStr, hv.Index(d))
	}

	var remStr string
	if remainder > 0 {
		remStr = strings.Repeat("0", int(remainder)) + " "
	}

	return remStr + strings.Join(strs, " ")
}

// String returns a string representation of a full hashVal. This is simply
// hv.HashPathString(DepthLimit).
func (hv hashVal) String() string {
	return hv.HashPathString(DepthLimit)
}

// ParseHashPath
func ParseHashPath(s string) hashVal {
	_ = AssertOn && assertf(strings.HasPrefix(s, "/"),
		"ParseHashPath: input, %q, does not start with '/'", s)

	if len(s) == 1 { // s="/"
		return 0
	}

	_ = AssertOn && assertf(!strings.HasSuffix(s, "/"),
		"ParseHashPath: input, %q, ends with '/'", s)

	var s0 = s[1:] //take the leading '/' off
	var idxStrs = strings.Split(s0, "/")

	var hv hashVal
	for i, idxStr := range idxStrs {
		var idx, err = strconv.ParseUint(idxStr, 10, int(IndexBits))
		if err != nil {
			log.Panicf("ParseHashPath: the %d'th index string failed to parse. err=%s", i, err)
		}

		//hv |= hashVal(idx << (uint(i) * IndexBits))
		hv = hv.buildHashPath(uint(idx), uint(i))
	}

	return hv
}
