package hamt32_test

import (
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/lleo/go-hamt/hamt32"
	"github.com/lleo/go-hamt/stringkey"
	"github.com/lleo/stringutil"
)

func TestHamt32New(t *testing.T) {
	var h = hamt32.New(TableOption)
	if h == nil {
		t.Fatal("no new Hamt struct")
	}
}

func TestBuildHamt32(t *testing.T) {
	var h = hamt32.New(TableOption)

	for _, kv := range hugeKvs {
		inserted := h.Put(kv.Key, kv.Val)
		if !inserted {
			t.Fatalf("failed to h.Put(%s, %v)", kv.Key, kv.Val)
		}
	}
}

func TestDeleteHamt32Del(t *testing.T) {
	var h = hamt32.New(TableOption)

	// build one up
	for _, kv := range hugeKvs {
		inserted := h.Put(kv.Key, kv.Val)
		if !inserted {
			t.Fatalf("failed to h.Put(%s, %v)", kv.Key, kv.Val)
		}
	}

	// then tear it down.
	for _, kv := range hugeKvs {
		val, deleted := h.Del(kv.Key)
		if !deleted {
			t.Fatalf("failed to h.Del(%q)", kv.Key)
		}
		if val != kv.Val {
			t.Fatalf("bad result of h.Del(%q); %v != %v", kv.Key, kv.Val, val)
		}
	}

	// make sure it is empty
	if !h.IsEmpty() {
		t.Fatalf("failed to empty h")
	}
}

func BenchmarkHamt32Get(b *testing.B) {
	log.Printf("BenchmarkHamt32Get: b.N=%d", b.N)

	for i := 0; i < b.N; i++ {
		var j = int(rand.Int31()) % numHugeKvs
		var key = hugeKvs[j].Key
		var val0 = hugeKvs[j].Val
		//var j = int(rand.Int31()) % numMidKvs
		//var key = midKvs[j].Key
		//var val0 = midKvs[j].Val
		var val, found = LookupHamt32.Get(key)
		if !found {
			b.Fatalf("H.Get(%s) not found", key)
		}
		if val != val0 {
			b.Fatalf("val,%v != hugeKvs[%d].Val,%v", val, j, val0)
			//b.Fatalf("val,%v != midKvs[%d].Val,%v", val, j, val0)
		}
	}
}

func BenchmarkHamt32Put(b *testing.B) {
	log.Printf("BenchmarkHamt32Put: b.N=%d", b.N)

	var h = hamt32.New(TableOption)
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

	// We rebuild the DeleteHamt32 datastructure because this Benchmark will probably be
	// rereun with different b.N values to get a better/more-accurate benchmark.

	StartTime["BenchmarkHamt32Del:rebuildDeleteHamt32"] = time.Now()
	rebuildDeleteHamt32(hugeKvs)
	RunTime["BenchmarkHamt32Del:rebuildDeleteHamt"] = time.Since(StartTime["BenchmarkHamt32Del:rebuildDeleteHamt32"])

	b.ResetTimer()

	StartTime["run BenchmarkHamt32Del"] = time.Now()
	for i := 0; i < b.N; i++ {
		kv := hugeKvs[i]
		key := kv.Key
		val := kv.Val

		v, ok := DeleteHamt32.Del(key)
		if !ok {
			b.Fatalf("failed to v, ok := DeleteHamt32.del(%s)", key)
		}
		if v != val {
			b.Fatalf("DeleteHamt32.del(%s) v,%d != i,%d", key, v, val)
		}
	}

	if DeleteHamt32.IsEmpty() {
		b.Fatal("DeleteHamt32.IsEmpty() => true")
	}

	RunTime["run BenchmarkHamt32Del"] = time.Since(StartTime["run BenchmarkHamt32Del"])
}
