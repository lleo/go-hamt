package hamt32

import (
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/lleo/go-hamt/stringkey"
	"github.com/lleo/stringutil"
)

func TestNewHamt32(t *testing.T) {
	var h = New(options)
	if h == nil {
		t.Fatal("no new Hamt struct")
	}
}

func TestBuildHamt32(t *testing.T) {
	var h = New(options)

	for _, kv := range hugeKvs {
		inserted := h.Put(kv.key, kv.val)
		if !inserted {
			t.Fatalf("failed to h.Put(%s, %v)", kv.key, kv.val)
		}
	}
}

func TestDeleteHamt32(t *testing.T) {
	var h = New(options)

	for _, kv := range hugeKvs {
		inserted := h.Put(kv.key, kv.val)
		if !inserted {
			t.Fatalf("failed to h.Put(%s, %v)", kv.key, kv.val)
		}
	}

	//for _, kv := range genRandomizedKvs(hugeKvs) {
	for _, kv := range hugeKvs {
		val, deleted := h.Del(kv.key)
		if !deleted {
			t.Fatalf("failed to h.Del(%q)", kv.key)
		}
		if val != kv.val {
			t.Fatalf("bad result of h.Del(%q); %v != %v", kv.key, kv.val, val)
		}
	}

	if !h.IsEmpty() {
		t.Fatalf("failed to empty h")
	}
}

func BenchmarkHamt32Get(b *testing.B) {
	log.Printf("BenchmarkHamt32Get: b.N=%d", b.N)

	for i := 0; i < b.N; i++ {
		var j = int(rand.Int31()) % numHugeKvs
		var key = hugeKvs[j].key
		var val0 = hugeKvs[j].val
		//var j = int(rand.Int31()) % numMidKvs
		//var key = midKvs[j].key
		//var val0 = midKvs[j].val
		var val, found = LookupHamt32.Get(key)
		if !found {
			b.Fatalf("H.Get(%s) not found", key)
		}
		if val != val0 {
			b.Fatalf("val,%v != hugeKvs[%d].val,%v", val, j, val0)
			//b.Fatalf("val,%v != midKvs[%d].val,%v", val, j, val0)
		}
	}
}

func BenchmarkHamt32Put(b *testing.B) {
	log.Printf("BenchmarkHamt32Put: b.N=%d", b.N)

	var h = New(options)
	var s = "aaa"
	for i := 0; i < b.N; i++ {
		key := stringkey.New(s)
		val := i + 1
		h.Put(key, val)
		s = stringutil.DigitalInc(s)
	}
}

func BenchmarkHamt32Del(b *testing.B) {
	log.Printf("BenchmarkHamt32Del: b.N=%d", b.N)

	StartTime["build BenchmarkHamt32Del:DeleteHamt32"] = time.Now()

	//debug := false
	for _, kv := range hugeKvs {
		//if kv.key.(*stringkey.StringKey).Str() == "aaa" {
		//	debug = true
		//}
		inserted := DeleteHamt32.Put(kv.key, kv.val)
		if !inserted {
			log.Printf("BenchmarkHamt32Del: inserted,%v := DeleteHamt32.Put(%s, %d)", inserted, kv.key, kv.val)
			break
		}
	}

	RunTime["build BenchmarkHamt32Del:DeleteHamt"] = time.Since(StartTime["build BenchmarkHamt32Del:DeleteHamt32"])

	b.ResetTimer()

	//log.Printf("BenchmarkHamt32Del: build BenchmarkHamt32Del:DeleteHamt32: took %s", RunTime["build numHugeKvs Hamt"])

	StartTime["run BenchmarkHamt32Del"] = time.Now()
	//s := "aaa"
	for i := 0; i < b.N; i++ {
		kv := hugeKvs[i]
		key := kv.key
		val := kv.val

		v, ok := DeleteHamt32.Del(key)
		if !ok {
			b.Fatalf("failed to v, ok := DeleteHamt32.del(%s)", key)
		}
		if v != val {
			b.Fatalf("DeleteHamt32.del(%s) v,%d != i,%d", key, v, val)
		}
		//s = stringutil.DigitalInc(s)
	}

	if DeleteHamt32.IsEmpty() {
		b.Fatal("DeleteHamt32.IsEmpty() => true")
	}

	RunTime["run BenchmarkHamt32Del"] = time.Since(StartTime["run BenchmarkHamt32Del"])
}
