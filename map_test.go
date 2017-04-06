package hamt_test

import (
	"fmt"
	"log"
	"testing"
)

func BenchmarkMapGet(b *testing.B) {
	var name = fmt.Sprintf("BenchmarkMapGet:%d", b.N)
	log.Printf("BenchmarkMapGet: b.N=%d", b.N)

	var lookupMap = buildMap(name, b.N)

	var keyStrings = make([]string, b.N)
	for k, v := range lookupMap {
		keyStrings[v] = k
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		//var j = int(rand.Int31()) % b.N
		//var s = keyStrings[j]
		var s = keyStrings[i]

		var val, ok = lookupMap[s]
		if !ok {
			b.Fatalf("LookupMap[%s] not ok", s)
		}
		//if val != j {
		//	b.Fatalf("val,%v != %v", val, j)
		if val != i {
			b.Fatalf("val,%v != %v", val, i)
		}
	}
}

func BenchmarkMapPut(b *testing.B) {
	log.Printf("BenchmarkMapPut: b.N=%d", b.N)

	var strings = make([]string, b.N)
	var s = "aaa"
	for i := 0; i < b.N; i++ {
		strings[i] = s
		s = Inc(s)
	}

	b.ResetTimer()

	var m = make(map[string]int, b.N)
	for i := 0; i < b.N; i++ {
		m[strings[i]] = i
	}
}

func BenchmarkMapDel(b *testing.B) {
	var name = fmt.Sprintf("BenchmarkMapDel:%d", b.N)
	log.Printf("BenchmarkMapDel: b.N=%d", b.N)

	var deleteMap = buildMap(name, b.N)

	var keyStrings = make([]string, b.N)
	var i int
	for k := range deleteMap {
		keyStrings[i] = k
		i++
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		delete(deleteMap, keyStrings[i])
	}

	if len(deleteMap) != 0 {
		b.Fatal("len(deleteMap) != 0")
	}
}
