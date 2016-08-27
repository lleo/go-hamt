/*
This is the unifying package between "github.com/lleo/go-hamt/hamt32" and
"github.com/lleo/go-hamt/hamt64". You can `import hamt "github.com/lleo/go-hamt"'
then instantiate either a hamt32 or hamt64 datastructure with the
hamt.NewHamt32() or hamt.NewHamt64() functions. Both datastructures have the
same exported API defined by the Hamt interface.

You may implement your own Key type by implementeding the Key interface
defined in "github.com/lleo/go-hamt/hamt_key" or you may used the example
StringKey interface described in "github.com/lleo/go-hamt/string_key".

A HAMT is a Hashed Array Mapped Trie datastructure. FIXME: explain HAMT
*/
package hamt

import (
	"fmt"

	"github.com/lleo/go-hamt/hamt32"
	"github.com/lleo/go-hamt/hamt64"
	"github.com/lleo/go-hamt/hamt_key"
)

type Hamt interface {
	Get(hamt_key.Key) (interface{}, bool)
	Put(hamt_key.Key, interface{}) bool
	Del(hamt_key.Key) (interface{}, bool)
	IsEmpty() bool
	String() string
	LongString(indent string) string
}

type keyVal struct {
	key hamt_key.Key
	val interface{}
}

func (kv keyVal) String() string {
	return fmt.Sprintf("keyVal{%s, %v}", kv.key, kv.val)
}

func NewHamt32() Hamt {
	return hamt32.NewHamt()
}

func NewHamt64() Hamt {
	return hamt64.NewHamt()
}
