package hamt64

import "bytes"

type ByteSliceKey []byte

func (bsk ByteSliceKey) Hash() uint64 {
	return CalcHash(bsk)
}

func (bsk ByteSliceKey) Equals(K KeyI) bool {
	var k, ok = K.(ByteSliceKey)
	if !ok {
		return false
	}
	return bytes.Equal(bsk, k)
}

type StringKey string

func (sk StringKey) Hash() uint64 {
	return CalcHash([]byte(sk))
}

func (sk StringKey) Equals(K KeyI) bool {
	var k, ok = K.(StringKey)
	if !ok {
		return false
	}
	return sk == k
}

type Int32Key int32

func (ik Int32Key) Hash() uint64 {
	return CalcHash([]byte{
		byte(0xff000000 & uint32(ik) >> 3 * 8),
		byte(0x00ff0000 & uint32(ik) >> 2 * 8),
		byte(0x0000ff00 & uint32(ik) >> 1 * 8),
		byte(0x000000ff & uint32(ik)),
	})
}

func (ik Int32Key) Equals(K KeyI) bool {
	var k, ok = K.(Int32Key)
	if !ok {
		return false
	}
	return ik == k
}

type Int64Key int64

func (ik Int64Key) Hash() uint64 {
	return CalcHash([]byte{
		byte(0xff00000000000000 & uint64(ik) >> 7 * 8),
		byte(0x00ff000000000000 & uint64(ik) >> 6 * 8),
		byte(0x0000ff0000000000 & uint64(ik) >> 5 * 8),
		byte(0x000000ff00000000 & uint64(ik) >> 4 * 8),
		byte(0x00000000ff000000 & uint64(ik) >> 3 * 8),
		byte(0x0000000000ff0000 & uint64(ik) >> 2 * 8),
		byte(0x000000000000ff00 & uint64(ik) >> 1 * 8),
		byte(0x00000000000000ff & uint64(ik)),
	})
}

func (ik Int64Key) Equals(K KeyI) bool {
	var k, ok = K.(Int64Key)
	if !ok {
		return false
	}
	return ik == k
}

type Uint32Key int32

func (ik Uint32Key) Hash() uint64 {
	return CalcHash([]byte{
		byte(0xff000000 & uint32(ik) >> 3 * 8),
		byte(0x00ff0000 & uint32(ik) >> 2 * 8),
		byte(0x0000ff00 & uint32(ik) >> 1 * 8),
		byte(0x000000ff & uint32(ik)),
	})
}

func (ik Uint32Key) Equals(K KeyI) bool {
	var k, ok = K.(Uint32Key)
	if !ok {
		return false
	}
	return ik == k
}

type Uint64Key int64

func (ik Uint64Key) Hash() uint64 {
	return CalcHash([]byte{
		byte(0xff00000000000000 & uint64(ik) >> 7 * 8),
		byte(0x00ff000000000000 & uint64(ik) >> 6 * 8),
		byte(0x0000ff0000000000 & uint64(ik) >> 5 * 8),
		byte(0x000000ff00000000 & uint64(ik) >> 4 * 8),
		byte(0x00000000ff000000 & uint64(ik) >> 3 * 8),
		byte(0x0000000000ff0000 & uint64(ik) >> 2 * 8),
		byte(0x000000000000ff00 & uint64(ik) >> 1 * 8),
		byte(0x00000000000000ff & uint64(ik)),
	})
}

func (ik Uint64Key) Equals(K KeyI) bool {
	var k, ok = K.(Uint64Key)
	if !ok {
		return false
	}
	return ik == k
}
