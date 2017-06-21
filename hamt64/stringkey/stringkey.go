/*
Package stringkey is an implementation of the
"github.com/lleo/go-hamt/hamt64" Key interface.
*/
package stringkey

import (
	"fmt"

	"github.com/lleo/go-hamt/hamt64"
)

// StringKey is a simple string implementing the hamt64.Key interface.
type StringKey struct {
	hamt64.KeyBase
	str string
}

// New allocates and initializes a StringKey data structure.
func New(str string) *StringKey {
	var k = new(StringKey)
	k.str = str
	k.Initialize([]byte(str))
	return k
}

// String return a string representation of StringKey data structure.
func (sk *StringKey) String() string {
	return fmt.Sprintf("StringKey{%s, str:%q}", sk.Hash(), sk.str)
}

// Equals returns true iff the StringKey exactly matches the key passed it. If
// The hamt64.Key passed as an argument is not also a StringKey Equals()
// automatically returns false.
func (sk *StringKey) Equals(k1 hamt64.Key) bool {
	var k, isStringKey = k1.(*StringKey)
	if !isStringKey {
		return false
	}
	return sk.str == k.str
}

// Str() returns the internal string this key is based on.
func (sk *StringKey) Str() string {
	return sk.str
}
