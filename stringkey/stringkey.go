/*
Package stringkey is an implementation of the "github.com/lleo/go-hamt/key"
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
	key.Base
	str string
}

// New allocates and initializes a StringKey data structure.
func New(str string) *StringKey {
	var k = new(StringKey)

	//Bad for straight assignment  all strings are pointers to structs.
	//But thats ok, cuz all strings are immutable.
	k.str = str

	k.Initialize([]byte(k.str))

	return k
}

// String return a string representation of StringKey data structure.
func (sk StringKey) String() string {
	return fmt.Sprintf("StringKey{%s, str:%q}", sk.Base.String(), sk.str)
}

func (sk StringKey) String30() string {
	return fmt.Sprintf("{%s, %q}", sk.Base.String30(), sk.str)
}

func (sk StringKey) String60() string {
	return fmt.Sprintf("{%s, %q}", sk.Base.String60(), sk.str)
}

// Equals returns true iff the StringKey exactly matches the key passed it. If
// The key.Key passed as an argument is not also a StringKey Equals()
// automatically returns false.
func (sk StringKey) Equals(key key.Key) bool {
	var k, isStrinKey = key.(*StringKey)
	if !isStrinKey {
		panic("type mismatch")
	}
	return sk.str == k.str
}

func (sk StringKey) toByteSlice() []byte {
	return []byte(sk.str)
}

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

// Str returns the internal string of StringKey. This allows for read-only
// access to the string field of StringKey.
func (sk StringKey) Str() string {
	return sk.str
}
