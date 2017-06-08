package key

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// HashVal30 stores 30 bits of a hash value.
type HashVal30 uint32

// BitsPerLevel30 is the number of bits per depth level of the HashVal30.
const BitsPerLevel30 uint = 5

// MaxDepth30 represents the maximum depth of the HashVal30.
const MaxDepth30 uint = 5

func indexMask30(depth uint) HashVal30 {
	return HashVal30((1<<BitsPerLevel30)-1) << (depth * BitsPerLevel30)
}

// Index() will return a 5bit (aka BitsPerLevel30) value 'depth' number of 5bits
// from the beginning of the HashVal30 (aka uint32) h30 value.
func (h30 HashVal30) Index(depth uint) uint {
	var idxMask = indexMask30(depth)
	var idx = uint((h30 & idxMask) >> (depth * BitsPerLevel30))
	return idx
}

// HashPathMask30() returns the mask for a 30 bit HashPath value.
func HashPathMask30(depth uint) HashVal30 {
	//return HashVal30(1<<(depth*BitsPerLevel30)) - 1
	return HashVal30(1<<((depth+1)*BitsPerLevel30)) - 1
}

// HashPath() calculates the path required to read the given depth. In other words
// it returns a uint32 that preserves the first depth-1 5bit index values.
// For depth=0 it always returns no path (aka a 0 value).
// For depth=MaxDepth60 it returns all but the last set of 5bit index values.
func (h30 HashVal30) HashPath(depth uint) HashVal30 {
	if depth == 0 {
		return 0
	}
	if depth > MaxDepth30 {
		log.Panicf("HashPath(): depth,%d > MaxDepth30,%d", depth, MaxDepth30)
	}
	return h30 & HashPathMask30(depth-1)
}

// BuildHashPath() method adds a idx at depth level of the hashPath.
// Given a hashPath = "/11/07/13" and you call hashPath.BuildHashPath(23, 3)
// the method will return hashPath "/11/07/13/23". hashPath is shown here
// in the string representation, but the real value is HashVal30 (aka uint32).
func (hashPath HashVal30) BuildHashPath(idx, depth uint) HashVal30 {
	//var mask = HashPathMask30(depth-1)
	var mask HashVal30 = (1 << (depth * BitsPerLevel30)) - 1
	var hp = hashPath & mask

	return hp | HashVal30(idx<<(depth*BitsPerLevel30))
}

// HashPathString() returns a string representation of the index path of a
// HashVal30 30 bit value; that is depth number of zero padded numbers between
// "00" and "31" separated by '/' characters and a leading '/'. If the limit
// parameter is 0 then the method will simply return a solitary "/".
// Warning: It will panic() if limit > MaxDepth30+1.
// Example: "/00/24/46/17" for limit=4 of a hash30 value represented
//       by "/00/24/46/17/34/08".
func (h30 HashVal30) HashPathString(limit uint) string {
	if limit > MaxDepth30+1 {
		panic(fmt.Sprintf("HashPathString: limit,%d > MaxDepth30+1,%d\n", limit, MaxDepth30+1))
	}

	if limit == 0 {
		return "/"
	}

	var strs = make([]string, limit)

	for d := uint(0); d < limit; d++ {
		var idx = h30.Index(d)
		strs[d] = fmt.Sprintf("%02d", idx)
	}

	return "/" + strings.Join(strs, "/")
}

// Return the HashVal30 as a 30 bit bit string separated into groups of 5 bits
// (aka BitsPerLevel30).
func (h30 HashVal30) BitString() string {
	var s = make([]string, MaxDepth30+1)
	for i := uint(0); i <= MaxDepth30; i++ {
		s[MaxDepth30-i] += fmt.Sprintf("%05b", h30.Index(i))
	}
	return "00 " + strings.Join(s, " ")
}

// String() returns a string representation of the h30 HashVal30 value. This
// is MaxDepth30+1(6) two digit numbers (zero padded) between "00" and "31"
// seperated by '/' characters and given a leading '/'.
// Example: "/08/14/28/20/00/31"
func (h30 HashVal30) String() string {
	return h30.HashPathString(MaxDepth30 + 1)
}

// ParseHashPath30() parses a string with a leading '/' and MaxDepth30+1 number
// of two digit numbers zero padded between "00" and "31" joined by '/' characters.
// Example: var h30 key.HashVal30 = key.ParseHashVal30("/00/01/02/03/04/05")
func ParseHashPath30(s string) HashVal30 {
	if !strings.HasPrefix(s, "/") {
		panic(errors.New("does not start with '/'"))
	}
	var s0 = s[1:]
	var as = strings.Split(s0, "/")

	var h30 HashVal30 = 0
	for i, s1 := range as {
		var ui, err = strconv.ParseUint(s1, 10, int(BitsPerLevel30))
		if err != nil {
			panic(errors.Wrap(err, fmt.Sprintf("strconv.ParseUint(%q, %d, %d) failed", s1, 10, BitsPerLevel30)))
		}
		h30 |= HashVal30(ui << (uint(i) * BitsPerLevel30))
		//fmt.Printf("%d: h30 = %q %2d %#02x %05b\n", i, s1, ui, ui, ui)
	}

	return h30
}
