package hamt_test

import (
	"log"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkMapGet(b *testing.B) {
	log.Printf("BenchmarkMapGet: b.N=%d", b.N)

	for i := 0; i < b.N; i++ {
		var j = int(rand.Int31()) % numHugeKvs
		var s = keyStrings[j]
		var val, ok = LookupMap[s]
		if !ok {
			b.Fatalf("LookupMap[%s] not ok", s)
		}
		if val != j {
			b.Fatalf("val,%v != %v", val, j)
		}
	}
}

func BenchmarkMapPut(b *testing.B) {
	log.Printf("BenchmarkMapPut: b.N=%d", b.N)

	var m = make(map[string]int, b.N)
	var s = "aaa"
	for i := 0; i < b.N; i++ {
		m[s] = i
		s = Inc(s)
	}
}

func BenchmarkMapDel(b *testing.B) {
	log.Printf("BenchmarkMapDel: b.N=%d", b.N)

	StartTime["BenchmarkMapDel:rebuildDeleteMap"] = time.Now()

	rebuildDeleteMap(keyStrings)

	RunTime["build BenchmarkMapDel:rebuildDeleteMap"] = time.Since(StartTime["BenchmarkMapDel:rebuildDeleteMap"])

	b.ResetTimer()

	s := "aaa"
	for i := 0; i < b.N; i++ {
		delete(DeleteMap, s)
		s = Inc(s)
	}

	if len(DeleteMap) == 0 {
		b.Fatal("len(DeleteMap) == 0")
	}
}
