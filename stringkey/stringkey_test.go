package stringkey_test

import (
	"testing"

	"github.com/lleo/go-hamt/stringkey"
)

func TestNew(t *testing.T) {
	var k = stringkey.New("test")

	var hash30 uint32 = 0xbc2c0be9 // []byte of []("test")
	//var hash30 uint32 = 0xca8c8619 // binary.Write uint32(0) & []byte(sk.s)

	if k.Hash30() != hash30 {
		t.Errorf("k.Hash30() != %#v", hash30)
	}
}
