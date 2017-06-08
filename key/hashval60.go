package key

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// HashVal60 stores 60 bits of a hash value.
type HashVal60 uint64

// BitsPerLevel60 is the number of bits per depth level of the HashVal60.
const BitsPerLevel60 uint = 6

// MaxDepth60 represents the maximum depth of the HashVal60.
const MaxDepth60 uint = 9

func indexMask60(depth uint) HashVal60 {
	return HashVal60((1<<BitsPerLevel60)-1) << (depth * BitsPerLevel60)
}

// Index() will return a 6bit (aka BitsPerLevel60) value 'depth' number
// of 6bits from the beginning of the HashVal60 (aka uint64) h60 value.
func (h60 HashVal60) Index(depth uint) uint {
	var idxMask = indexMask60(depth)
	var idx = uint((h60 & idxMask) >> (depth * BitsPerLevel60))
	return idx
}

// HashPathMask60() returns the mask for a 60 bit HashPath value.
func HashPathMask60(depth uint) HashVal60 {
	//return HashVal60(1<<(depth*BitsPerLevel60)) - 1
	return HashVal60(1<<((depth+1)*BitsPerLevel60)) - 1
}

// HashPath() calculates the path required to read the given depth. In other
// words it returns a uint64 that preserves the first depth-1 6bit index values.
// For depth=0 it always returns no path (aka a 0 value).
// For depth=MaxDepth60 it returns all but the last set of 6bit index values.
func (h60 HashVal60) HashPath(depth uint) HashVal60 {
	if depth == 0 {
		return 0
	}
	if depth > MaxDepth60 {
		log.Panicf("HashPath(): depth,%d > MaxDepth60,%d", depth, MaxDepth60)
	}
	return h60 & HashPathMask60(depth-1)
}

// BuildHashPath() method adds a idx at depth level of the hashPath.
// Given a hashPath = "/11/07/13" and you call hashPath.BuildHashPath(23, 3)
// the method will return hashPath "/11/07/13/23". hashPath is shown here
// in the string representation, but the real value is HashVal60 (aka uint64).
func (hashPath HashVal60) BuildHashPath(idx, depth uint) HashVal60 {
	//var mask = HashPathMask60(depth-1)
	var mask HashVal60 = (1 << (depth * BitsPerLevel60)) - 1
	var hp = hashPath & mask

	return hp | HashVal60(idx<<(depth*BitsPerLevel60))
}

// HashPathString() returns a string representation of the index path of a
// HashVal60 60 bit value; that is depth number of zero padded numbers between
// "00" and "63" separated by '/' characters and a leading '/'. If the limit
// parameter is 0 then the method will simply return a solitary "/".
// Warning: It will panic() if limit > MaxDepth60+1.
// Example: "/00/24/46/17" for limit=4 of a hash60 value represented
//       by "/00/24/46/17/34/08".
func (h60 HashVal60) HashPathString(limit uint) string {
	if limit > MaxDepth60+1 {
		panic(fmt.Sprintf("HashPathString: limit,%d > MaxDepth60+1,%d\n", limit, MaxDepth60+1))
	}

	if limit == 0 {
		return "/"
	}

	var strs = make([]string, limit)

	for d := uint(0); d < limit; d++ {
		var idx = h60.Index(d)
		strs[d] = fmt.Sprintf("%02d", idx)
	}

	return "/" + strings.Join(strs, "/")
}

// Return the HashVal60 as a 60 bit bit string separated into groups of 6 bits
// (aka BitsPerLevel60).
func (h60 HashVal60) BitString() string {
	var s = make([]string, MaxDepth60+1)
	for i := uint(0); i <= MaxDepth60; i++ {
		s[MaxDepth60-i] += fmt.Sprintf("%05b", h60.Index(i))
	}
	return "00 " + strings.Join(s, " ")
}

// String() returns a string representation of the h60 HashVal60 value. This
// is MaxDepth60+1(10) two digit numbers (zero padded) between "00" and "63"
// seperated by '/' characters and given a leading '/'.
// Example: "/08/14/28/20/00/63"
func (h60 HashVal60) String() string {
	return h60.HashPathString(MaxDepth60 + 1)
}

// ParseHashPath60() parses a string with a leading '/' and MaxDepth60+1 number
// of two digit numbers zero padded between "00" and "63" joined by '/' characters.
// Example: var h60 key.HashVal60 = key.ParseHashVal60("/00/01/02/03/04/05")
func ParseHashPath60(s string) HashVal60 {
	if !strings.HasPrefix(s, "/") {
		panic(errors.New("does not start with '/'"))
	}
	var s0 = s[1:]
	var as = strings.Split(s0, "/")

	var h60 HashVal60 = 0
	for i, s1 := range as {
		var ui, err = strconv.ParseUint(s1, 10, int(BitsPerLevel60))
		if err != nil {
			panic(errors.Wrap(err, fmt.Sprintf("strconv.ParseUint(%q, %d, %d) failed", s1, 10, BitsPerLevel60)))
		}
		h60 |= HashVal60(ui << (uint(i) * BitsPerLevel60))
		//fmt.Printf("%d: h60 = %q %2d %#02x %05b\n", i, s1, ui, ui, ui)
	}

	return h60
}
