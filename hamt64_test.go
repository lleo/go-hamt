package hamt_test

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
)

func TestHamt64Put(t *testing.T) {
	var name = "TestHamt64Put"

	var _, err = buildHamt64(name, KVS, TableOption)
	if err != nil {
		t.Fatalf("failed to build Hamt64: %s", err)
	}

	//var h = hamt64.New(TableOption)
	//for _, kv := range KVS {
	//	inserted := h.Put(kv.Key, kv.Val)
	//	if !inserted {
	//		t.Fatalf("failed to h.Put(%s, %v)", kv.Key, kv.Val)
	//	}
	//}
}

func TestHamt64Del(t *testing.T) {
	var name = "TestHamt64Del"

	// build one up
	var h, err = buildHamt64(name, KVS, TableOption)
	if err != nil {
		t.Fatalf("failed to build Hamt64: %s", err)
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

func BenchmarkHamt64Get(b *testing.B) {
	var name = fmt.Sprintf("BenchmarkHamt64Get:%d", b.N)
	log.Printf("BenchmarkHamt64Get: b.N=%d", b.N)

	var kvs = buildKeyVals(name, b.N)
	var h, err = buildHamt64(name, kvs, TableOption)
	if err != nil {
		b.Fatalf("Failed to buildHamt64: %s", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var j = int(rand.Int31()) % b.N
		var key = kvs[j].Key
		var val = kvs[j].Val

		var v, found = h.Get(key)
		if !found {
			b.Fatalf("H.Get(%s) not found", key)
		}
		if val != v {
			b.Fatalf("v,%v != kvs[%d].Val,%v", val, j, v)
		}
	}
}

func BenchmarkHamt64Put(b *testing.B) {
	var name = fmt.Sprintf("BenchmarkHamt64Put:%d", b.N)
	log.Printf("BenchmarkHamt64Put: b.N=%d", b.N)

	var kvs = buildKeyVals(name, b.N)

	b.ResetTimer()

	var _, err = buildHamt64(name, kvs, TableOption)
	if err != nil {
		b.Fatalf("failed to buildHamt64: %s", err)
	}

	//var h = hamt64.New(TableOption)
	//var s = "aaa"
	//for i := 0; i < b.N; i++ {
	//	key := stringkey.New(s)
	//	val := i + 1
	//	h.Put(key, val)
	//	s = Inc(s)
	//}
}

func BenchmarkHamt64Del(b *testing.B) {
	var name = fmt.Sprintf("BenchmarkHamt64Del:%d", b.N)
	log.Printf("BenchmarkHamt64Del: b.N=%d", b.N)

	var kvs = buildKeyVals(name, b.N)

	var h, err = buildHamt64(name, kvs, TableOption)
	if err != nil {
		b.Fatalf("failed to buildHamt64: %s", err)
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
}
