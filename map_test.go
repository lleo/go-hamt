package hamt_test

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

//var SVS []StrVal

type StrVal struct {
	Str string
	Val int
}

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

func buildMap(prefix string, num int) map[string]int {
	var name = fmt.Sprintf("%s-buildMap", prefix)
	StartTime[name] = time.Now()

	var m = make(map[string]int, num)
	var s = "aaa"

	for i := 0; i < num; i++ {
		m[s] = i

		s = Inc(s)
	}

	RunTime[name] = time.Since(StartTime[name])
	return m
}

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
