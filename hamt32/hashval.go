package hamt32

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type HashVal uint32

func indexMask(depth uint) HashVal {
	return HashVal((1<<BitsPerLevel)-1) << (depth * BitsPerLevel)
}

// Index returns the BitsPerLevel bit value of the HashVal at 'depth' number of
// BitsPerLevel number of bits into HashVal.
func (hv HashVal) Index(depth uint) uint {
	if depth > MaxDepth {
		log.Printf("Index: depth,%d > MaxIndex,%d", depth, MaxDepth)
	}
	var idxMask = indexMask(depth)
	return uint((hv & idxMask) >> (depth * BitsPerLevel))
}

func hashPathMask(depth uint) HashVal {
	return HashVal(1<<((depth)*BitsPerLevel)) - 1
}

// HashPath calculates the path required to read the given depth. In other words
// it returns a HashVal that preserves the first depth-1 BitsPerLevel index
// values. For depth=0 it always returns no path (aka a 0 value).
// For depth=MaxDepth it returns all but the last index value.
func (hv HashVal) HashPath(depth uint) HashVal {
	if depth == 0 {
		return 0
	}
	if depth > MaxDepth {
		log.Panicf("HashPath(): depth,%d > MaxDepth,%d", depth, MaxDepth)
	}
	return hv & hashPathMask(depth)

}

// BuildHashPath method adds a idx at depth level of the hashPath. Given a
// hash Path = "/11/07/13" and you call hashPath.BuildHashPath(23, 3) the method
// will return hashPath "/11/07/13/23". hashPath is shown here in the string
// representation, but the real value is HashVal (aka uint32).
func (hv HashVal) BuildHashPath(idx, depth uint) HashVal {
	if idx > MaxIndex {
		log.Panicf("BuildHashPath: idx,%d >= IndexLimit,%d", idx, IndexLimit)
	}
	hv &= hashPathMask(depth)
	return hv | HashVal(idx<<(depth*BitsPerLevel))
}

// HashPathString returns a string representation of the index path of a
// HashVal. It will be string of the form "/idx0/idx1/..." where each idxN value
// will be a zero padded number between 0 and MaxIndex. There will be limit
// number of such values where limit <= DepthLimit.
// If the limit parameter is 0 then the method will simply return "/".
// Warning: It will panic() if limit > DepthLimit.
// Example: "/00/24/46/17" for limit=4 of a BitsPerLevel=5 hash value
// represented by "/00/24/46/17/34/08".
func (hv HashVal) HashPathString(limit uint) string {
	if limit > DepthLimit {
		log.Panicf("HashPathString: limit,%d > DepthLimit,%d\n",
			limit, DepthLimit)
	}

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

// BitString returns a HashVal as a string of bits separated into groups of
// BitsPerLevel bits.
func (hv HashVal) BitString() string {
	var strs = make([]string, DepthLimit)

	for d := uint(0); d < DepthLimit; d++ {
		var fmtStr = fmt.Sprintf("%%0%dd", BitsPerLevel)
		strs[MaxDepth-d] = fmt.Sprintf(fmtStr, hv.Index(d))
	}

	var remStr string
	if remainder > 0 {
		remStr = strings.Repeat("0", int(remainder)) + " "
	}

	return remStr + strings.Join(strs, " ")
}

// String returns a string representation of a full HashVal. This is simply
// hv.HashPathString(DepthLimit).
func (hv HashVal) String() string {
	return hv.HashPathString(DepthLimit)
}

// ParseHashPath
func ParseHashPath(s string) HashVal {
	if !strings.HasPrefix(s, "/") {
		log.Panicf("ParseHashPath: input, %q, does not start with '/'", s)
	}
	if len(s) == 1 { // s="/"
		return 0
	}
	if strings.HasSuffix(s, "/") {
		log.Panicf("ParseHashPath: input, %q, ends with '/'", s)
	}
	var s0 = s[1:] //take the leading '/' off
	var idxStrs = strings.Split(s0, "/")

	var hv HashVal
	for i, idxStr := range idxStrs {
		var idx, err = strconv.ParseUint(idxStr, 10, int(BitsPerLevel))
		if err != nil {
			log.Panicf("ParseHashPath: the %d'th index string failed to parse. err=%s", i, err)
		}

		//hv |= HashVal(idx << (uint(i) * BitsPerLevel))
		hv = hv.BuildHashPath(uint(idx), uint(i))
	}

	return hv
}
