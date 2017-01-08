/*
Package stringKey is an implementation of the "github.com/lleo/go-hamt/key"
interface.

Both of these packages are are in separate package namespaces to prevent circular
references between "github.com/lleo/go-hamt" and "github.com/lleo/go-hamt/hamt32"
and/or "github.com/lleo/go-hamt/hamt64".

Additionally, "github.com/lleo/go-hamt/stringKey" is used by the functional
variation of the Hash Array Mapped Trie in "github.com/lleo/go-hamt-fuctional".
*/
package stringkey

import (
	"fmt"

	"github.com/lleo/go-hamt/key"
)

// StringKey is a simple string implementing the key.Key interface.
type StringKey struct {
	key.KeyBase
	str string
}

func New(str string) *StringKey {
	var k = new(StringKey)
	k.str = str

	k.Initialize([]byte(k.str))

	return k
}

// Return the string in StringKey.
func (sk *StringKey) String() string {
	//return fmt.Sprintf("StringKey{hash30:%#v, s:%#v}", sk.Hash30(), sk.str)
	//return fmt.Sprintf("StringKey{KeyBase{hash30:%#v}, s:%#v}", sk.Hash30(), sk.str)
	//return fmt.Sprintf("StringKey{%#v, s:%#v}", sk.KeyBase, sk.str)
	return fmt.Sprintf("StringKey{%s, str:%q}", sk.KeyBase.String(), sk.str)
	//return fmt.Sprintf("%#v", sk)
}

// Equals returns true iff the StringKey exactly matches the key passed it. The
// key.Key passed as an argument MUST BE a StringKey or the method Equals()
// automatically returns false.
func (sk *StringKey) Equals(key key.Key) bool {
	var k, typeMatches = key.(*StringKey)
	if !typeMatches {
		panic("type mismatch")
	}
	return sk.str == k.str
}

// Convert the string in StringKey to a []byte.
func toByteSlice(str string) []byte {
	return []byte(str)

	////BROKEN: This does not work because StringKey is not a fixed size Struct.
	//var bytebuf bytes.Buffer
	//err := binary.Write(&bytebuf, binary.BigEndian, *sk)
	//if err != nil {
	//	panic(err)
	//}
	//return bytebuf.Bytes()

	////Variation on above: works!
	//var bytebuf bytes.Buffer
	//err := binary.Write(&bytebuf, binary.BigEndian, uint32(0) /* Hash30() */)
	//if err != nil {
	//	panic(err)
	//}
	//err = binary.Write(&bytebuf, binary.BigEndian, []byte(str))
	////err = binary.Write(&bytebuf, binary.BigEndian, str)
	//if err != nil {
	//	panic(err)
	//}
	//return bytebuf.Bytes()

	////ANOTHER variation: does not work, binary.Write() does not like
	//// struct with []byte
	//var d = struct {
	//	hash30 uint32
	//	s      []byte
	//}{0 /*Hash30()*/, []byte(str)}
	//var bytebuf bytes.Buffer
	//err := binary.Write(&bytebuf, binary.BigEndian, d)
	//if err != nil {
	//	panic(err)
	//}
	//return bytebuf.Bytes()

	//ANOTHER IDEA: use encoding/gob; but I don't know that api well enough.

	//ANOThER IDEA: use gopkg.in/mgo.v2/bson; but it requires a third party lib.
}

// Str() returns the internal string of StringKey. This allows for read-only
// access to the string field of StringKey.
func (sk *StringKey) Str() string {
	return sk.str
}
