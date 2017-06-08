package key_test

import (
	"testing"

	"github.com/lleo/go-hamt/key"
)

func TestBuildHashPath60(t *testing.T) {
	var h60 key.HashVal60
	var manual_h60 uint64

	for d, i := range []uint{3, 6, 7, 11, 13, 17, 19, 21, 23, 27} {
		manual_h60 += uint64(i << (uint(d) * key.BitsPerLevel60))
		h60 = h60.BuildHashPath(i, uint(d))

		if key.HashVal60(manual_h60) != h60 {
			t.Fatalf("depth %d: manual_60,%s != h60,%s",
				key.HashVal60(manual_h60), h60)
		}
	}
}

func TestHashPathString60(t *testing.T) {
	var h60 key.HashVal60

	for depth, idx := range []uint{3, 5, 7, 11, 13, 17, 19, 21, 23, 27} {
		h60 |= key.HashVal60(idx << (uint(depth) * key.BitsPerLevel60))
	}

	var s = h60.HashPathString(10)
	if s != "/03/05/07/11/13/17/19/21/23/27" {
		t.Fatalf("%q != \"/03/05/07/11/13/17/19/21/23/27\"", s)
	}
}

func TestParseHashPath60(t *testing.T) {
	var h60 key.HashVal60

	for depth, idx := range []uint{3, 6, 7, 11, 13, 17, 19, 21, 23, 27} {
		h60 |= key.HashVal60(idx << (uint(depth) * key.BitsPerLevel60))
		var s = h60.HashPathString(uint(depth + 1))
		//log.Printf("TestParseHashPath60: h60.HashPathString(uint(%s) => %q\n",
		//	depth+1, s)
		var parsedH60 = key.ParseHashPath60(s)
		if parsedH60 != h60 {
			t.Fatal("parsedH60,%#016x != h60,%#016x", parsedH60, h60)
		}
	}
}
