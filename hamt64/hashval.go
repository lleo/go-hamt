package hamt64

import (
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// hashVal sets the numberer of bits of the hash value by being an alias to
// uint64 and establishes a type we can hang methods, like Index(), off of.
type hashVal uint64

// CalcHash deterministically calculates a randomized uint64 of a given byte
// slice .
func CalcHash(bs []byte) uint64 {
	return fold(hash(bs), remainder)
}

func hash(bs []byte) uint64 {
	var h = fnv.New64()
	h.Write(bs)
	return h.Sum64()
}

func mask(size uint) uint64 {
	return uint64(1<<size) - 1
}

func fold(hash uint64, rem uint) uint64 {
	return (hash >> (hashSize - rem)) ^ (hash & mask(hashSize-rem))
}

func indexMask(depth uint) hashVal {
	return hashVal((1<<NumIndexBits)-1) << (depth * NumIndexBits)
}

// Index returns the NumIndexBits bit value of the hashVal at 'depth' number of
// NumIndexBits number of bits into hashVal.
func (hv hashVal) Index(depth uint) uint {
	_ = assertOn && assert(depth < DepthLimit, "Index: depth > maxDepth")

	var idxMask = indexMask(depth)
	return uint((hv & idxMask) >> (depth * NumIndexBits))
}

func hashPathMask(depth uint) hashVal {
	return hashVal(1<<((depth)*NumIndexBits)) - 1
}

// hashPath calculates the path required to read the given depth. In other words
// it returns a hashVal that preserves the first depth-1 NumIndexBits index
// values. For depth=0 it always returns no path (aka a 0 value).
// For depth=maxDepth it returns all but the last index value.
func (hv hashVal) hashPath(depth uint) hashVal {
	_ = assertOn && assert(depth < DepthLimit, "hashPath(): dept > maxDepth")

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
	_ = assertOn && assert(idx < DepthLimit, "buildHashPath: idx > maxIndex")

	hv &= hashPathMask(depth)
	return hv | hashVal(idx<<(depth*NumIndexBits))
}

// HashPathString returns a string representation of the index path of a
// hashVal. It will be string of the form "/idx0/idx1/..." where each idxN value
// will be a zero padded number between 0 and maxIndex. There will be limit
// number of such values where limit <= DepthLimit.
// If the limit parameter is 0 then the method will simply return "/".
// Example: "/00/24/46/17" for limit=4 of a NumIndexBits=5 hash value
// represented by "/00/24/46/17/34/08".
func (hv hashVal) HashPathString(limit uint) string {
	_ = assertOn && assertf(limit <= DepthLimit,
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
// NumIndexBits bits.
func (hv hashVal) bitString() string {
	var strs = make([]string, DepthLimit)

	var fmtStr = fmt.Sprintf("%%0%db", NumIndexBits)
	for d := uint(0); d < DepthLimit; d++ {
		strs[maxDepth-d] = fmt.Sprintf(fmtStr, hv.Index(d))
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

// parseHashPath
func parseHashPath(s string) (hashVal, error) {
	if !strings.HasPrefix(s, "/") {
		return 0, errors.Errorf(
			"parseHashPath: input, %q, does not start with '/'", s)
	}

	if len(s) == 1 { // s="/"
		return 0, nil
	}

	if strings.HasSuffix(s, "/") {
		return 0, errors.Errorf("parseHashPath: input, %q, ends with '/'", s)
	}
	var s0 = s[1:] //take the leading '/' off
	var idxStrs = strings.Split(s0, "/")

	var hv hashVal
	for i, idxStr := range idxStrs {
		var idx, err = strconv.ParseUint(idxStr, 10, int(NumIndexBits))
		if err != nil {
			return 0, errors.Wrapf(err,
				"parseHashPath: the %d'th index string failed to parse.", i)
		}

		//hv |= hashVal(idx << (uint(i) * NumIndexBits))
		hv = hv.buildHashPath(uint(idx), uint(i))
	}

	return hv, nil
}
