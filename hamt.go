package hamt

import (
	"fmt"

	"github.com/lleo/go-hamt/hamt32"
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

//type Key interface {
//	Equals(Key) bool
//	Hash30() uint32
//	Hash60() uint64
//	String() string
//}

type keyVal struct {
	key hamt_key.Key
	val interface{}
}

func (kv keyVal) String() string {
	return fmt.Sprintf("keyVal{%s, %v}", kv.key, kv.val)
}

func NewHamt32() Hamt {
	return hamt32.NewHamt32()
}
