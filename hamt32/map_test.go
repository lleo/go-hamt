package hamt32

import (
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/lleo/go-hamt/stringkey"
	"github.com/lleo/stringutil"
)

func BenchmarkMapGet(b *testing.B) {
	log.Printf("BenchmarkMapGet: b.N=%d", b.N)

	for i := 0; i < b.N; i++ {
		var j = int(rand.Int31()) % numHugeKvs
		var s = hugeKvs[j].key.(*stringkey.StringKey).Str()
		var val, ok = LookupMap[s]
		if !ok {
			b.Fatalf("LookupMap[%s] not ok", string(s))
		}
		if val != hugeKvs[j].val {
			b.Fatalf("val,%v != hugeKvs[%d].val,%v", val, j, hugeKvs[j].val)
		}
	}
}

func BenchmarkMapPut(b *testing.B) {
	log.Printf("BenchmarkMapPut: b.N=%d", b.N)

	var m = make(map[string]int)
	var s = "aaa"
	for i := 0; i < b.N; i++ {
		m[string(s)] = i + 1
		s = stringutil.DigitalInc(s)
	}
}

func BenchmarkMapDel(b *testing.B) {
	log.Printf("BenchmarkMapDel: b.N=%d", b.N)

	StartTime["BenchmarkMapDel:DeleteMap"] = time.Now()
	for _, kv := range hugeKvs {
		sk := kv.key.(*stringkey.StringKey)
		str := sk.Str()
		val := kv.val.(int)

		if _, ok := DeleteMap[str]; !ok {
			DeleteMap[str] = val
		}
	}
	RunTime["build BenchmarkMapDel:DeleteMap"] = time.Since(StartTime["BenchmarkMapDel:DeleteMap"])

	b.ResetTimer()

	s := "aaa"
	for i := 0; i < b.N; i++ {
		delete(DeleteMap, s)
		s = stringutil.DigitalInc(s)
	}

	if len(DeleteMap) == 0 {
		b.Fatal("len(DeleteMap) == 0")
	}
}
