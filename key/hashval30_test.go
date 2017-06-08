package key_test

import (
	"testing"

	"github.com/lleo/go-hamt/key"
)

func TestBuildHashPath30(t *testing.T) {
	var h30 key.HashVal30
	var manual_h30 uint32

	for d, i := range []uint{3, 6, 7, 11, 13, 17} {
		manual_h30 += uint32(i << (uint(d) * key.BitsPerLevel30))
		h30 = h30.BuildHashPath(i, uint(d))

		if key.HashVal30(manual_h30) != h30 {
			t.Fatalf("depth %d: manual_30,%s != h30,%s",
				key.HashVal30(manual_h30), h30)
		}
	}
}

func TestHashPathString30(t *testing.T) {
	var h30 key.HashVal30
	// for HashVal60        []uint{3, 5, 7, 11, 13, 17, 19, 21, 23, 27}
	for depth, idx := range []uint{3, 5, 7, 11, 13, 17} {
		h30 |= key.HashVal30(idx << (uint(depth) * key.BitsPerLevel30))
	}

	var s = h30.HashPathString(6)
	if s != "/03/05/07/11/13/17" {
		t.Fatalf("%q != \"/03/05/07/11/13/17\"", s)
	}
}

func TestParseHashPath30(t *testing.T) {
	var h30 key.HashVal30

	for depth, idx := range []uint{3, 5, 7, 11, 13, 17} {
		h30 |= key.HashVal30(idx << (uint(depth) * key.BitsPerLevel30))
		var s = h30.HashPathString(uint(depth + 1))
		//log.Printf("TestParseHashPath30: h30.HashPathString(uint(%s) => %q\n",
		//	depth+1, s)
		var parsedH30 = key.ParseHashPath30(s)
		if parsedH30 != h30 {
			t.Fatal("parsedH30,%#08x != h30,%#08x", parsedH30, h30)
		}
	}
}
