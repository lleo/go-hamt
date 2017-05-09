package hamt_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/lleo/go-hamt/hamt32"
)

func TestHamt32Get(t *testing.T) {
	var name = "TestHamt32Get"

	var h, err = buildHamt32(name, KVS, TableOption)
	if err != nil {
		t.Fatalf("failed to build Hamt32: %s", err)
	}

	for _, kv := range KVS {
		var key = kv.Key
		var val = kv.Val

		var v, found = h.Get(key)
		if !found {
			t.Fatalf("failed to Get(%s)", key)
		}
		if val != v {
			t.Fatalf("val,%d != v,%d for key=%s", val, v, key)
		}
	}
}

func TestHamt32Put(t *testing.T) {
	//var name = "TestHamt32Put"

	//var _, err = buildHamt32(name, KVS, TableOption)
	//if err != nil {
	//	t.Fatalf("failed to build Hamt32: %s", err)
	//}

	var h = hamt32.New(TableOption)
	for _, kv := range KVS {
		var key = kv.Key
		var val = kv.Val

		inserted := h.Put(kv.Key, kv.Val)
		if !inserted {
			t.Fatalf("failed to h.Put(%s, %v)", kv.Key, kv.Val)
		}

		var v, found = h.Get(key)
		if !found {
			t.Fatalf("failed to Get(%s)", key)
		}
		if val != v {
			t.Fatalf("val,%d != v,%d for key=%s", val, v, key)
		}
	}
}

func TestHamt32Del(t *testing.T) {
	var name = "TestHamt32Del"

	// build one up
	var h, err = buildHamt32(name, KVS, TableOption)
	if err != nil {
		t.Fatalf("failed to build Hamt32: %s", err)
	}

	// then tear it down.
	for _, kv := range KVS {
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
	var name = fmt.Sprintf("BenchmarkHamt32Get:%d", b.N)
	log.Printf("BenchmarkHamt32Get: b.N=%d", b.N)

	var kvs = buildKeyVals(name, b.N)
	var h, err = buildHamt32(name, kvs, TableOption)
	if err != nil {
		b.Fatalf("Failed to buildHamt32: %s", err)
	}

	//kvs = genRandomizedKvs(kvs)
	b.ResetTimer()

	for _, kv := range kvs {
		var key = kv.Key
		var val = kv.Val

		var v, found = h.Get(key)
		if !found {
			b.Fatalf("H.Get(%s) not found", key)
		}
		if val != v {
			b.Fatalf("v,%v != val,%v for key=%s", v, val, key)
		}
	}
}

func BenchmarkHamt32Put(b *testing.B) {
	var name = fmt.Sprintf("BenchmarkHamt32Put:%d", b.N)
	log.Printf("BenchmarkHamt32Put: b.N=%d", b.N)

	var kvs = buildKeyVals(name, b.N)

	b.ResetTimer()

	var _, err = buildHamt32(name, kvs, TableOption)
	if err != nil {
		b.Fatalf("failed to buildHamt32: %s", err)
	}

	//var h = hamt32.New(TableOption)
	//var s = "aaa"
	//for i := 0; i < b.N; i++ {
	//	key := stringkey.New(s)
	//	val := i + 1
	//	h.Put(key, val)
	//	s = Inc(s)
	//}
}

func BenchmarkHamt32Del(b *testing.B) {
	var name = fmt.Sprintf("BenchmarkHamt32Del:%d", b.N)
	log.Printf("BenchmarkHamt32Del: b.N=%d", b.N)

	var kvs = buildKeyVals(name, b.N+1) //b.N+1 so it won't be empty in the end

	var h, err = buildHamt32(name, kvs, TableOption)
	if err != nil {
		b.Fatalf("failed to buildHamt32: %s", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		kv := kvs[i]
		key := kv.Key
		val := kv.Val

		var v, ok = h.Del(key)
		if !ok {
			b.Fatalf("failed to v, ok := h.del(%s)", key)
		}
		if v != val {
			b.Fatalf("h.del(%s) v,%d != kvs[%d].Val,%d", key, v, i, val)
		}
	}

	if h.IsEmpty() {
		b.Fatal("h.IsEmpty() => true")
	}
	if h.Nentries() != 1 {
		b.Fatal("h.Nentries() != 1")
	}
}
