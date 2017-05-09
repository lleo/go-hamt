package hamt_test

import (
	"fmt"
	"log"
	"testing"
)

func BenchmarkMapGet(b *testing.B) {
	var name = fmt.Sprintf("BenchmarkMapGet#%d", b.N)
	var lookupMap = buildMap(name, b.N)

	var svs = make([]StrVal, b.N)
	var j int
	for k, v := range lookupMap {
		svs[j] = StrVal{k, v}
		j++
	}

	//svs = genRandomizedSvs(svs)

	b.ResetTimer()

	for _, sv := range svs {
		var str = sv.Str
		var val = sv.Val

		var v, ok = lookupMap[str]
		if !ok {
			b.Fatalf("LookupMap[%s] not ok", str)
		}
		if val != v {
			b.Fatalf("v,%v != val,%v", v, val)
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
