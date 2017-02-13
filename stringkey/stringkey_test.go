package stringkey_test

import (
	"testing"

	"github.com/lleo/go-hamt/stringkey"
)

func TestHash30New(t *testing.T) {
	var k = stringkey.New("test")

	t.Logf("TestHash30New: k.Hash30()=0x%x\n", k.Hash30())

	var hash30 uint32 = 0x3c2c0beb // []byte of []("test")
	//var hash30 uint32 = 0xca8c8619 // binary.Write uint32(0) & []byte(sk.s)

	if k.Hash30() != hash30 {
		t.Errorf("k.Hash30() != %#v", hash30)
	}
}

func TestHash60New(t *testing.T) {
	var k = stringkey.New("test")

	t.Logf("TestHash60New: k.Hash60()=0x%016x\n", k.Hash60())

	var hash60 uint64 = 0x0c093f7e9fccbf61 // []byte of []("test") | hash60InitMask
	//var hash60 uint64 = 0xca8c8619 // binary.Write uint64(0) & []byte(sk.s)

	if k.Hash60() != hash60 {
		t.Errorf("k.Hash60() != %#v", hash60)
	}
}
