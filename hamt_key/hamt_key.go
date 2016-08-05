package hamt_key

type Key interface {
	Equals(Key) bool
	Hash30() uint32
	Hash60() uint64
	String() string
}
