package string_key

import "hash/fnv"

type StringKey string

func (sk StringKey) Equals(key hamt_key.Key) bool {
	var k, typeMatches = key.(StringKey)
	if !typeMatches {
		panic("type mismatch")
	}
	return string(sk) == string(k)
}

func (sk StringKey) Hash30() uint32 {
	return hash30(sk.hash32())
}

func (sk StringKey) Hash60() uint64 {
	return hash60(sk.hash64())
}

func (sk StringKey) hash32() uint32 {
	var h = fnv.New32()
	h.Write([]byte(sk))
	return h.Sum32()
}

func (sk StringKey) hash64() uint64 {
	var h = fnv.New64()
	h.Write([]byte(sk))
	return h.Sum64()
}

const mask30 = uint32(1<<30) - 1
const mask60 = uint64(1<<60) - 1

func hash30(h30 uint32) uint32 {
	return (h30 >> 30) ^ (h30 & mask30)
}

func hash60(h64 uint64) uint64 {
	return (h64 >> 60) ^ (h64 & mask60)
}

func (sk StringKey) String() string {
	return string(sk)
}
