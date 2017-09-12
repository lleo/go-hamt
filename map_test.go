package hamt_test

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

func genRandomizedSvs(svs []StrVal) []StrVal {
	var randSvs = make([]StrVal, len(svs))
	copy(randSvs, svs)

	//From: https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle#The_modern_algorithm
	var limit = len(randSvs) //n-1
	for i := 0; i < limit; /* aka i_max = n-2 */ i++ {
		j := rand.Intn(i+1) - 1 // i <= j < n; j_min=n-(n-2+1)-1=0; j_max=n-0-1=n-1
		randSvs[i], randSvs[j] = randSvs[j], randSvs[i]
	}

	return randSvs
}

func buildMap(prefix string, svs []StrVal) map[string]int {
	var name = fmt.Sprintf("%s-buildMap", prefix)
	StartTime[name] = time.Now()

	var m = make(map[string]int, len(svs))

	for _, sv := range svs {
		m[sv.Str] = sv.Val
	}

	RunTime[name] = time.Since(StartTime[name])
	return m
}

func BenchmarkMapGet(b *testing.B) {
	var name = fmt.Sprintf("BenchmarkMapGet#%d", b.N)

	var svs = buildStrVals(name, b.N)
	//svs = genRandomizedSvs(svs)

	var lookupMap = buildMap(name, svs)

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

func BenchmarkMap_GetN30(b *testing.B) {
	runBenchmarkMapGetN(b, 30)
}

func BenchmarkMap_GetN1000(b *testing.B) {
	runBenchmarkMapGetN(b, 1000)
}

func BenchmarkMapGetN10000(b *testing.B) {
	runBenchmarkMapGetN(b, 10000)
}

func runBenchmarkMapGetN(b *testing.B, N int) {
	var name = "runBenchmarkMapGetN"

	log.Println(name, b.N)

	var svs = buildStrVals(name, N)
	var m = buildMap(name, svs)

	//var keys = make([]string, 0, N)
	//for k, _ := range m {
	//	keys = append(keys, k)
	//}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		//var k = keys[i%len(keys)]
		var sv = svs[i%len(svs)]
		var _, ok = m[sv.Str]
		if !ok {
			//b.Fatalf("failed to lookup k,%q in m", k)
			b.Fatalf("failed to lookup sv.Str,%q in m", sv.Str)
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

	var svs = buildStrVals(name, b.N)
	var deleteMap = buildMap(name, svs)

	//var keyStrings = make([]string, b.N)
	//var i int
	//for k := range deleteMap {
	//	keyStrings[i] = k
	//	i++
	//}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		//delete(deleteMap, keyStrings[i])
		delete(deleteMap, svs[i].Str)
	}

	if len(deleteMap) != 0 {
		b.Fatal("len(deleteMap) != 0")
	}
}

func BenchmarkMapIter(b *testing.B) {
	var name = fmt.Sprintf("BenchmarkMapIter:%d", b.N)
	log.Printf("BenchmarkMapIter: b.N=%d", b.N)

	var svs = buildStrVals(name, b.N)
	var iterMap = buildMap(name, svs)

	//var keyStrings = make([]string, b.N)
	//var i int
	//for k := range iterMap {
	//	keyStrings[i] = k
	//	i++
	//}

	b.ResetTimer()

	var i = 0
	for k := range iterMap {
		if len(k) < 1 {
			b.Fatal("len(k) == 0")
		}
		i++
	}
}
