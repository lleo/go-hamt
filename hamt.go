package hamt

import (
	"fmt"

	"github.com/lleo/go-hamt/hamt32"
)

type Hamt interface {
	Get(Key) (interface{}, bool)
	Put(Key, interface{}) bool
	Del(Key) (interface{}, bool)
	String() string
	LongString(indent string) string
}

type Key interface {
	Equals(Key) bool
	Hash30() uint32
	Hash60() uint64
	String() string
}

type keyVal struct {
	key Key
	val interface{}
}

func (kv keyVal) String() string {
	return fmt.Sprintf("keyVal{%s, %v}", kv.key, kv.val)
}

func NewHamt32() hamt32.Hamt {
	return hamt32.NewHamt()
}
