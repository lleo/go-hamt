/*
Package stringKey is an implementation of the
"github.com/lleo/go-hamt/key.Key" interface.

Both of these packages are are in separate package namespaces to prevent circular
references between "github.com/lleo/go-hamt" and "github.com/lleo/go-hamt/hamt32"
and/or "github.com/lleo/go-hamt/hamt64".

Additionally, "github.com/lleo/go-hamt/stringKey" is used by the functional
variation of the Hash Array Mapped Trie in "github.com/lleo/go-hamt-fuctional".
*/
package stringkey

import (
	"bytes"
	"encoding/binary"

	"github.com/lleo/go-hamt/key"
)

// StringKey is a simple string implementing the key.Key interface.
type StringKey struct {
	key.KeyBase
	s string
}

func New(str string) *StringKey {
	var k = new(StringKey)
	k.s = str
	k.Initialize(k.ToByteSlice())
	return k
}

// Return the string in StringKey.
func (sk *StringKey) String() string {
	return string(sk.s)
}

// Equals returns true iff the StringKey exactly matches the key passed it. The
// key.Key passed as an argument MUST BE a StringKey or the method Equals()
// automatically returns false.
func (sk *StringKey) Equals(key key.Key) bool {
	var k, typeMatches = key.(*StringKey)
	if !typeMatches {
		panic("type mismatch")
	}
	return string(sk.s) == string(k.s)
}

// Convert the string in StringKey to a []byte.
func (sk *StringKey) ToByteSlice() []byte {
	//return []byte(sk.s)

	var bytebuf bytes.Buffer
	binary.Write(&bytebuf, binary.BigEndian, sk)
	return bytebuf.Bytes()
}
