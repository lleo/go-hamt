/*
Package stringkey is an implementation of the "github.com/lleo/go-hamt/key"
interface.

Both of these packages are are in separate package namespaces to prevent
circular references between "github.com/lleo/go-hamt" and
"github.com/lleo/go-hamt/hamt32" and/or "github.com/lleo/go-hamt/hamt64"; as
 well as "github.com/lleo/go-hamt-functional",
"github.com/lleo/go-hamt-functional/hamt32", and
"github.com/lleo/go-hamt-functional/hamt64".
*/
package stringkey

import (
	"fmt"

	"github.com/lleo/go-hamt/key"
)

// StringKey is a simple string implementing the key.Key interface.
type StringKey struct {
	key.Base
	str string
}

// New allocates and initializes a StringKey data structure.
func New(str string) *StringKey {
	var k = new(StringKey)

	k.str = str

	k.Initialize([]byte(k.str))

	return k
}

// String return a string representation of StringKey data structure.
func (sk StringKey) String() string {
	return fmt.Sprintf("StringKey{%s, str:%q}", sk.Base.String(), sk.str)
}

// Equals returns true iff the StringKey exactly matches the key passed it. If
// The key.Key passed as an argument is not also a StringKey Equals()
// automatically returns false.
func (sk StringKey) Equals(k1 key.Key) bool {
	var k, isStrinKey = k1.(*StringKey)
	if !isStrinKey {
		panic("type mismatch")
	}
	return sk.str == k.str
}

// Str() returns the internal string this key is based on.
func (sk StringKey) Str() string {
	return sk.str
}
