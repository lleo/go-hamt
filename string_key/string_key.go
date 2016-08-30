/*
Package string_key is an implementation of the
"github.com/lleo/go-hamt/key.Key" interface.

Both of these packages are are in separate package namespaces to prevent circular
references between "github.com/lleo/go-hamt" and "github.com/lleo/go-hamt/hamt32"
and/or "github.com/lleo/go-hamt/hamt64".

Additionally, "github.com/lleo/go-hamt/string_key" is used by the functional
variation of the Hash Array Mapped Trie in "github.com/lleo/go-hamt-fuctional".
*/
package string_key

import "github.com/lleo/go-hamt/key"

// StringKey is a simple string implementing the key.Key interface.
type StringKey string

// Equals returns true iff the StringKey exactly matches the key passed it. The
// key.Key passed as an argument MUST BE a StringKey or the method Equals()
// automatically returns false.
func (sk StringKey) Equals(key key.Key) bool {
	var k, typeMatches = key.(StringKey)
	if !typeMatches {
		panic("type mismatch")
	}
	return string(sk) == string(k)
}

// Convert the string in StringKey to a []byte.
func (sk StringKey) ToByteSlice() []byte {
	return []byte(sk)
}

// Return the string in StringKey.
func (sk StringKey) String() string {
	return string(sk)
}
