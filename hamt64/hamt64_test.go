package hamt64_test

import (
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/lleo/go-hamt/hamt64"
	"github.com/lleo/go-hamt/stringkey"
	"github.com/lleo/stringutil"
)

func TestNewHamt64(t *testing.T) {
	var h = hamt64.New(TableOption)
	if h == nil {
		t.Fatal("no new Hamt struct")
	}
}

func TestBuildHamt64(t *testing.T) {
	var h = hamt64.New(TableOption)

	for _, kv := range hugeKvs {
		inserted := h.Put(kv.Key, kv.Val)
		if !inserted {
			t.Fatalf("failed to h.Put(%s, %v)", kv.Key, kv.Val)
		}
	}
}

func TestDeleteHamt64(t *testing.T) {
	var h = hamt64.New(TableOption)

	for _, kv := range hugeKvs {
		inserted := h.Put(kv.Key, kv.Val)
		if !inserted {
			t.Fatalf("failed to h.Put(%s, %v)", kv.Key, kv.Val)
		}
	}

	//for _, kv := range genRandomizedKvs(hugeKvs) {
	for _, kv := range hugeKvs {
		val, deleted := h.Del(kv.Key)
		if !deleted {
			t.Fatalf("failed to h.Del(%q)", kv.Key)
		}
		if val != kv.Val {
			t.Fatalf("bad result of h.Del(%q); %v != %v", kv.Key, kv.Val, val)
		}
	}

	if !h.IsEmpty() {
		t.Fatalf("failed to empty h")
	}
}

func BenchmarkHamt64Get(b *testing.B) {
	log.Printf("BenchmarkHamt64Get: b.N=%d", b.N)

	for i := 0; i < b.N; i++ {
		var j = int(rand.Int31()) % numHugeKvs
		var key = hugeKvs[j].Key
		var val0 = hugeKvs[j].Val
		//var j = int(rand.Int31()) % numMidKvs
		//var key = midKvs[j].Key
		//var val0 = midKvs[j].val
		var val, found = LookupHamt64.Get(key)
		if !found {
			b.Fatalf("H.Get(%s) not found", key)
		}
		if val != val0 {
			b.Fatalf("val,%v != hugeKvs[%d].val,%v", val, j, val0)
			//b.Fatalf("val,%v != midKvs[%d].val,%v", val, j, val0)
		}
	}
}

func BenchmarkHamt64Put(b *testing.B) {
	log.Printf("BenchmarkHamt64Put: b.N=%d", b.N)

	var h = hamt64.New(TableOption)
	var s = "aaa"
	for i := 0; i < b.N; i++ {
		key := stringkey.New(s)
		val := i + 1
		h.Put(key, val)
		s = stringutil.DigitalInc(s)
	}
}

func BenchmarkHamt64Del(b *testing.B) {
	log.Printf("BenchmarkHamt64Del: b.N=%d", b.N)

	// We rebuild the DeleteHamt64 datastructure because this Benchmark will probably be
	// rereun with different b.N values to get a better/more-accurate benchmark.

	StartTime["BenchmarkHamt64Del:rebuildDeleteHamt64"] = time.Now()
	rebuildDeleteHamt64(hugeKvs)
	RunTime["BenchmarkHamt64Del:rebuildDeleteHamt"] = time.Since(StartTime["BenchmarkHamt64Del:rebuildDeleteHamt64"])

	b.ResetTimer()

	StartTime["run BenchmarkHamt64Del"] = time.Now()
	for i := 0; i < b.N; i++ {
		kv := hugeKvs[i]
		key := kv.Key
		val := kv.Val

		v, ok := DeleteHamt64.Del(key)
		if !ok {
			b.Fatalf("failed to v, ok := DeleteHamt64.del(%s)", key)
		}
		if v != val {
			b.Fatalf("DeleteHamt64.del(%s) v,%d != i,%d", key, v, val)
		}
	}

	if DeleteHamt64.IsEmpty() {
		b.Fatal("DeleteHamt64.IsEmpty() => true")
	}

	RunTime["run BenchmarkHamt64Del"] = time.Since(StartTime["run BenchmarkHamt64Del"])
}
