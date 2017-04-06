package hamt_test

import (
	"testing"
)

func BenchmarkInc(b *testing.B) {
	s := "aaa"
	for i := 0; i < b.N; i++ {
		s = Inc(s)
	}
}
