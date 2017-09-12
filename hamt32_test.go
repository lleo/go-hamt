package hamt_test

import (
	"bytes"
	"context"
	"log"
	"testing"
	"time"

	"github.com/lleo/go-hamt/hamt32"
)

func TestBuild32(t *testing.T) {
	var name = "TestBuild32"
	if Functional {
		name += ":functional:" + hamt32.TableOptionName[TableOption]
	} else {
		name += ":transient:" + hamt32.TableOptionName[TableOption]
	}

	var h = hamt32.New(Functional, TableOption)

	for _, kv := range KVS[:30] {
		var k = kv.Key
		var v = kv.Val

		var inserted bool
		h, inserted = h.Put(k, v)
		if !inserted {
			log.Printf("%s: failed to insert s=%q, v=%d", name, string(k), v)
			t.Fatalf("%s: failed to insert s=%q, v=%d", name, string(k), v)
		}

		//log.Print(h.LongString(""))
	}
}

func TestHamt32Put(t *testing.T) {
	runTestHamt32Put(t, KVS, Functional, TableOption)
}

func runTestHamt32Put(
	t *testing.T,
	kvs []KeyVal,
	functional bool,
	tblOpt int,
) {
	var name = "TestHamt32Put"
	if functional {
		name += ":functional:" + hamt32.TableOptionName[tblOpt]
	} else {
		name += ":transient:" + hamt32.TableOptionName[tblOpt]
	}

	StartTime[name] = time.Now()
	Hamt32 = hamt32.New(functional, tblOpt)
	for _, kv := range kvs {
		var k = kv.Key
		var v = kv.Val

		var inserted bool
		Hamt32, inserted = Hamt32.Put(k, v)
		if !inserted {
			log.Printf("%s: failed to Hamt32.Put(%q, %v)", name, string(k), v)
			t.Fatalf("%s: failed to Hamt32.Put(%q, %v)", name, string(k), v)
		}

		var val, found = Hamt32.Get(k)
		if !found {
			log.Printf("%s: failed to Hamt32.Get(%q)", name, string(k))
			//log.Print(Hamt32.LongString(""))
			t.Fatalf("%s: failed to Hamt32.Get(%q)", name, string(k))
		}
		if val != v {
			log.Printf("%s: returned val,%d != expected v,%d for s=%q", name, val, v, string(k))
			t.Fatalf("%s: returned val,%d != expected v,%d for s=%q", name, val, v, string(k))
		}
	}
	RunTime[name] = time.Since(StartTime[name])

	StartTime["Hamt32.Stats()"] = time.Now()
	var stats = Hamt32.Stats()
	RunTime["Hamt32.Stats()"] = time.Since(StartTime["Hamt32.Stats()"])
	log.Printf("%s: stats=%+v;\n", name, stats)
}

func TestHamt32IterFunc(t *testing.T) {
	runTestHamt32IterFunc(t, KVS, Functional, TableOption)
}

func runTestHamt32IterFunc(
	t *testing.T,
	kvs []KeyVal,
	functional bool,
	tblOpt int,
) {
	var name = "TestHamt32IterFunc"
	if functional {
		name += ":functional:" + hamt32.TableOptionName[tblOpt]
	} else {
		name += ":transient:" + hamt32.TableOptionName[tblOpt]
	}

	if Hamt32 == nil {
		var err error
		Hamt32, err = buildHamt32(name, kvs, functional, tblOpt)
		if err != nil {
			log.Printf("%s: failed buildHamt32(%q, kvs#%d, %t, %s) => %s", name,
				name, len(kvs), functional,
				hamt32.TableOptionName[tblOpt], err)
			t.Fatalf("%s: failed buildHamt32(%q, kvs#%d, %t, %s) => %s", name,
				name, len(kvs), functional,
				hamt32.TableOptionName[tblOpt], err)
		}

		StartTime["Hamt32.Stats()"] = time.Now()
		var stats = Hamt32.Stats()
		RunTime["Hamt32.Stats()"] = time.Since(StartTime["Hamt32.Stats()"])
		log.Printf("%s: stats=%+v;\n", name, stats)
	}

	StartTime[name] = time.Now()

	var i int
	var next = Hamt32.Iter()
	for kv, ok := next(); ok; kv, ok = next() {
		var val, ok = Hamt32.Get(kv.Key)
		if !ok {
			t.Fatalf("failed to lookup %s in Hamt32", kv.Key)
		}

		if val != kv.Val {
			t.Fatalf("val,%v != kv.Val,%v\n", val, kv.Val)
		}

		i++
	}

	if len(kvs) != i {
		t.Fatalf("Expected len(kvs),%d go i,%d; Hamt32.Nentries()=%d;",
			len(kvs), i, Hamt32.Nentries())
	}

	RunTime[name] = time.Since(StartTime[name])
}

func TestHamt32IterChan(t *testing.T) {
	runTestHamt32IterChan(t, KVS, Functional, TableOption)
}

func runTestHamt32IterChan(
	t *testing.T,
	kvs []KeyVal,
	functional bool,
	tblOpt int,
) {
	var name = "TestHamt32IterChan"
	if functional {
		name += ":functional:" + hamt32.TableOptionName[tblOpt]
	} else {
		name += ":transient:" + hamt32.TableOptionName[tblOpt]
	}

	if Hamt32 == nil {
		var err error
		Hamt32, err = buildHamt32(name, kvs, functional, tblOpt)
		if err != nil {
			log.Printf("%s: failed buildHamt32(%q, kvs#%d, %t, %s) => %s", name,
				name, len(kvs), functional,
				hamt32.TableOptionName[tblOpt], err)
			t.Fatalf("%s: failed buildHamt32(%q, kvs#%d, %t, %s) => %s", name,
				name, len(kvs), functional,
				hamt32.TableOptionName[tblOpt], err)
		}

		StartTime["Hamt32.Stats()"] = time.Now()
		var stats = Hamt32.Stats()
		RunTime["Hamt32.Stats()"] = time.Since(StartTime["Hamt32.Stats()"])
		log.Printf("%s: stats=%+v;\n", name, stats)
	}

	StartTime[name] = time.Now()

	var i int
	for kv := range Hamt32.IterChan(0, nil) {
		var val, ok = Hamt32.Get(kv.Key)
		if !ok {
			t.Fatalf("failed to lookup %s in Hamt32", kv.Key)
		}

		if val != kv.Val {
			t.Fatalf("val,%v != kv.Val,%v\n", val, kv.Val)
		}

		i++
	}

	if len(kvs) != i {
		t.Fatalf("Expected len(kvs),%d go i,%d", len(kvs), i)
	}

	RunTime[name] = time.Since(StartTime[name])
}

func TestHamt32IterChanContext(t *testing.T) {
	runTestHamt32IterChanContext(t, KVS, Functional, TableOption)
}

func runTestHamt32IterChanContext(
	t *testing.T,
	kvs []KeyVal,
	functional bool,
	tblOpt int,
) {
	var name = "TestHamt32IterChanContext"
	if functional {
		name += ":functional:" + hamt32.TableOptionName[tblOpt]
	} else {
		name += ":transient:" + hamt32.TableOptionName[tblOpt]
	}

	if Hamt32 == nil {
		var err error
		Hamt32, err = buildHamt32(name, kvs, functional, tblOpt)
		if err != nil {
			log.Printf("%s: failed buildHamt32(%q, kvs#%d, %t, %s) => %s", name,
				name, len(kvs), functional,
				hamt32.TableOptionName[tblOpt], err)
			t.Fatalf("%s: failed buildHamt32(%q, kvs#%d, %t, %s) => %s", name,
				name, len(kvs), functional,
				hamt32.TableOptionName[tblOpt], err)
		}

		StartTime["Hamt32.Stats()"] = time.Now()
		var stats = Hamt32.Stats()
		RunTime["Hamt32.Stats()"] = time.Since(StartTime["Hamt32.Stats()"])
		log.Printf("%s: stats=%+v;\n", name, stats)
	}

	StartTime[name] = time.Now()

	var i int
	var stopKey = kvs[0].Key // "aaa" but key from iter are random
	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	var iterChan = Hamt32.IterChan(0, ctx)
	for kv := range iterChan {
		var val, ok = Hamt32.Get(kv.Key)
		if !ok {
			t.Fatalf("failed to lookup %s in Hamt32", kv.Key)
		}

		if val != kv.Val {
			t.Fatalf("val,%v != kv.Val,%v\n", val, kv.Val)
		}

		i++

		if bytes.Equal(kv.Key, stopKey) {
			break
		}
	}

	log.Printf("%s: stopped after %d iterations", name, i)

	RunTime[name] = time.Since(StartTime[name])
}

func TestHamt32Get(t *testing.T) {
	runTestHamt32Get(t, KVS, Functional, TableOption)
}

func runTestHamt32Get(
	t *testing.T,
	kvs []KeyVal,
	functional bool,
	tblOpt int,
) {
	var name = "TestHamt32Get"
	if functional {
		name += ":functional:" + hamt32.TableOptionName[tblOpt]
	} else {
		name += ":transient:" + hamt32.TableOptionName[tblOpt]
	}

	if Hamt32 == nil {
		var err error
		Hamt32, err = buildHamt32(name, kvs, functional, tblOpt)
		if err != nil {
			log.Printf("%s: failed buildHamt32(%q, kvs#%d, %t, %s) => %s", name,
				name, len(kvs), functional,
				hamt32.TableOptionName[tblOpt], err)
			t.Fatalf("%s: failed buildHamt32(%q, kvs#%d, %t, %s) => %s", name,
				name, len(kvs), functional,
				hamt32.TableOptionName[tblOpt], err)
		}

		StartTime["Hamt32.Stats()"] = time.Now()
		var stats = Hamt32.Stats()
		RunTime["Hamt32.Stats()"] = time.Since(StartTime["Hamt32.Stats()"])
		log.Printf("%s: stats=%+v;\n", name, stats)
	}

	StartTime[name] = time.Now()
	for _, kv := range kvs {
		var k = kv.Key
		var v = kv.Val

		var val, found = Hamt32.Get(k)
		if !found {
			log.Printf("%s: Failed to Hamt32.Get(%q)", name, string(k))
			log.Print(Hamt32.LongString(""))
			t.Fatalf("%s: Failed to Hamt32.Get(%q)", name, string(k))
		}
		if val != v {
			log.Printf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, string(k))
			t.Fatalf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, string(k))
		}
	}
	RunTime[name] = time.Since(StartTime[name])
}

func TestHamt32Del(t *testing.T) {
	runTestHamt32Del(t, KVS, Functional, TableOption)
}

func runTestHamt32Del(
	t *testing.T,
	kvs []KeyVal,
	functional bool,
	tblOpt int,
) {
	var name = "TestHamt32Del"
	if functional {
		name += ":functional:" + hamt32.TableOptionName[tblOpt]
	} else {
		name += ":transient:" + hamt32.TableOptionName[tblOpt]
	}

	if Hamt32 == nil {
		var err error
		Hamt32, err = buildHamt32(name, kvs, functional, tblOpt)
		if err != nil {
			log.Printf("%s: failed buildHamt32(%q, kvs#%d, %t, %s) => %s", name,
				name, len(kvs), functional,
				hamt32.TableOptionName[tblOpt], err)
			t.Fatalf("%s: failed buildHamt32(%q, kvs#%d, %t, %s) => %s", name,
				name, len(kvs), functional,
				hamt32.TableOptionName[tblOpt], err)
		}

		StartTime["Hamt32.Stats()"] = time.Now()
		var stats = Hamt32.Stats()
		RunTime["Hamt32.Stats()"] = time.Since(StartTime["Hamt32.Stats()"])
		log.Printf("%s: stats=%+v;\n", name, stats)
	}

	StartTime[name] = time.Now()
	for _, kv := range kvs {
		var k = kv.Key
		var v = kv.Val

		var val interface{}
		var deleted bool
		Hamt32, val, deleted = Hamt32.Del(k)
		if !deleted {
			log.Printf("%s: Failed to Hamt32.Del(%q)", name, string(k))
			log.Print(Hamt32.LongString(""))
			t.Fatalf("%s: Failed to Hamt32.Del(%q)", name, string(k))
		}
		if val != v {
			log.Printf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, string(k))
			t.Fatalf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, string(k))
		}
	}
	RunTime[name] = time.Since(StartTime[name])
}

func BenchmarkHamt32Get(b *testing.B) {
	runBenchmarkHamt32Get(b, KVS, Functional, TableOption)
}

func runBenchmarkHamt32Get(
	b *testing.B,
	kvs []KeyVal,
	functional bool,
	tblOpt int,
) {
	var name = "BenchmarkHamt32Get"
	if functional {
		name += ":functional:" + hamt32.TableOptionName[tblOpt]
	} else {
		name += ":transient:" + hamt32.TableOptionName[tblOpt]
	}

	var BenchHamt32Get, err = buildHamt32(name, kvs, functional, tblOpt)
	if err != nil {
		log.Printf("%s: failed buildHamt32(%q, kvs#%d, %t, %s) => %s", name,
			name, len(kvs), false, hamt32.TableOptionName[tblOpt], err)
		b.Fatalf("%s: failed buildHamt32(%q, kvs#%d, %t, %s) => %s", name,
			name, len(kvs), false, hamt32.TableOptionName[tblOpt], err)
	}

	log.Printf("%s: b.N=%d", name, b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var j = i % len(kvs)
		var k = kvs[j].Key
		var v = kvs[j].Val

		var val, found = BenchHamt32Get.Get(k)
		if !found {
			log.Printf("%s: Failed to h.Get(%q)", name, string(k))
			//log.Print(h.LongString(""))
			b.Fatalf("%s: Failed to h.Get(%q)", name, string(k))
		}
		if val != v {
			log.Printf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, string(k))
			b.Fatalf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, string(k))
		}
	}
}

func BenchmarkHamt32Put(b *testing.B) {
	runBenchmarkHamt32Put(b, KVS, Functional, TableOption)
}

func runBenchmarkHamt32Put(
	b *testing.B,
	kvs []KeyVal,
	functional bool,
	tblOpt int,
) {
	var name = "BenchmarkHamt32Put"
	if functional {
		name += ":functional:" + hamt32.TableOptionName[tblOpt]
	} else {
		name += ":transient:" + hamt32.TableOptionName[tblOpt]
	}

	if b.N+InitHamtNumKvsForPut > len(kvs) {
		log.Printf("%s: Can't run: b.N+num > len(kvs)", name)
		b.Fatalf("%s: Can't run: b.N+num > len(kvs)", name)
	}

	var initKvs = kvs[:InitHamtNumKvsForPut]

	var h, err = buildHamt32(name, initKvs, functional, tblOpt)
	if err != nil {
		log.Printf("%s: failed buildHamt32(%q, kvs[:%d], %t, %s) => %s", name,
			name, InitHamtNumKvsForPut, functional,
			hamt32.TableOptionName[tblOpt], err)
		b.Fatalf("%s: failed buildHamt32(%q, kvs[:%d], %t, %s) => %s", name,
			name, InitHamtNumKvsForPut, functional,
			hamt32.TableOptionName[tblOpt], err)
	}

	log.Printf("%s: b.N=%d;", name, b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var k = kvs[InitHamtNumKvsForPut+i].Key
		var v = kvs[InitHamtNumKvsForPut+i].Val

		var added bool
		h, added = h.Put(k, v)
		if !added {
			log.Printf("%s: failed to h.Put(%q, %d)", name, string(k), v)
			//log.Print(h.LongString(""))
			b.Fatalf("%s: failed to h.Put(%q, %d)", name, string(k), v)
		}
	}
}

func BenchmarkHamt32Del(b *testing.B) {
	runBenchmarkHamt32Del(b, KVS, Functional, TableOption)
}

func runBenchmarkHamt32Del(
	b *testing.B,
	kvs []KeyVal,
	functional bool,
	tblOpt int,
) {
	var name = "BenchmarkHamt32Del"
	if functional {
		name += ":functional:" + hamt32.TableOptionName[tblOpt]
	} else {
		name += ":transient:" + hamt32.TableOptionName[tblOpt]
	}

	var h, err = buildHamt32(name, kvs, functional, tblOpt)
	if err != nil {
		log.Printf("%s: failed buildHamt32(%q, kvs:%d, %t, %s) => %s", name,
			name, len(kvs), functional,
			hamt32.TableOptionName[tblOpt], err)
		b.Fatalf("%s: failed buildHamt32(%q, kvs:%d, %t, %s) => %s", name,
			name, len(kvs), functional,
			hamt32.TableOptionName[tblOpt], err)
	}

	log.Printf("%s: b.N=%d;", name, b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var k = kvs[i].Key
		var v = kvs[i].Val

		var deleted bool
		var val interface{}
		h, val, deleted = h.Del(k)
		if !deleted {
			log.Printf("%s: failed to h.Del(%q)", name, string(k))
			//log.Print(h.LongString(""))
			b.Fatalf("%s: failed to h.Del(%q)", name, string(k))
		}
		if val != v {
			log.Printf("%s: failed val,%d != v,%d", name, val, v)
			b.Fatalf("%s: failed val,%d != v,%d", name, val, v)
		}
	}
}

func BenchmarkHamt32IterFunc(b *testing.B) {
	runBenchmarkHamt32IterFunc(b, KVS, Functional, TableOption)
}

func runBenchmarkHamt32IterFunc(
	b *testing.B,
	kvs []KeyVal,
	functional bool,
	tblOpt int,
) {
	var name = "BenchmarkHamt32IterFunc"
	if Functional {
		name += ":functional:" + hamt32.TableOptionName[tblOpt]
	} else {
		name += ":transient:" + hamt32.TableOptionName[tblOpt]
	}

	var BenchHamt32Get, err = buildHamt32(name, kvs, functional, tblOpt)
	if err != nil {
		log.Printf("%s: failed buildHamt32(%q, kvs#%d, %t, %s) => %s", name,
			name, len(kvs), false, hamt32.TableOptionName[tblOpt], err)
		b.Fatalf("%s: failed buildHamt32(%q, kvs#%d, %t, %s) => %s", name,
			name, len(kvs), false, hamt32.TableOptionName[tblOpt], err)
	}

	log.Printf("%s: b.N=%d", name, b.N)
	b.ResetTimer()

	var next = BenchHamt32Get.Iter()
	var kv, ok = next()
	for i := 0; i < b.N; i++ {
		if len(kv.Key) < 0 {
			b.Fatal("stupid test to touch the kv")
		}

		if i >= b.N {
			break
		}

		kv, ok = next()
		if !ok {
			next = BenchHamt32Get.Iter()
			kv, ok = next()
		}
	}
}

func BenchmarkHamt32IterChan(b *testing.B) {
	runBenchmarkHamt32IterChan(b, KVS, Functional, TableOption)
}

func runBenchmarkHamt32IterChan(
	b *testing.B,
	kvs []KeyVal,
	functional bool,
	tblOpt int,
) {
	var name = "BenchmarkHamt32IterChan"
	if Functional {
		name += ":functional:" + hamt32.TableOptionName[tblOpt]
	} else {
		name += ":transient:" + hamt32.TableOptionName[tblOpt]
	}

	var BenchHamt32Get, err = buildHamt32(name, kvs, functional, tblOpt)
	if err != nil {
		log.Printf("%s: failed buildHamt32(%q, kvs#%d, %t, %s) => %s", name,
			name, len(kvs), false, hamt32.TableOptionName[tblOpt], err)
		b.Fatalf("%s: failed buildHamt32(%q, kvs#%d, %t, %s) => %s", name,
			name, len(kvs), false, hamt32.TableOptionName[tblOpt], err)
	}

	log.Printf("%s: b.N=%d", name, b.N)
	b.ResetTimer()

	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	var nents = BenchHamt32Get.Nentries()
	var iterChan = BenchHamt32Get.IterChan(20, ctx)
ForLoop:
	for i := 0; i < b.N; i++ {
		//// This code is ~100ns slower than the code after this select
		//select {
		//case <-ctx.Done():
		//	break ForLoop
		//case kv, ok := <-iterChan:
		//	if !ok {
		//		// iterChan closed
		//		break ForLoop
		//	}
		//
		//	if len(kv.Key) < 0 {
		//		b.Fatal("stupid test to touch the kv")
		//	}
		//
		//	if uint(i) == nents {
		//		cancel()
		//		ctx, cancel = context.WithCancel(context.Background())
		//		defer cancel()
		//		iterChan = BenchHamt32Get.IterChan(20, ctx)
		//	}
		//}
		var kv = <-iterChan

		if len(kv.Key) < 0 {
			b.Fatal("stupid test to touch the kv")
			break ForLoop //meaningless...just to use the ForLoop label
		}

		if uint(i) == nents {
			cancel()
			ctx, cancel = context.WithCancel(context.Background())
			defer cancel()
			iterChan = BenchHamt32Get.IterChan(20, ctx)
		}
	}
}

func BenchmarkHamt32Stats(b *testing.B) {
	runBenchmarkHamt32Stats(b, KVS, Functional, TableOption)
}

func runBenchmarkHamt32Stats(
	b *testing.B,
	kvs []KeyVal,
	functional bool,
	tblOpt int,
) {
	var name = "BenchmarkHamt32Stats"
	if Functional {
		name += ":functional:" + hamt32.TableOptionName[tblOpt]
	} else {
		name += ":transient:" + hamt32.TableOptionName[tblOpt]
	}

	var h, err = buildHamt32(name, kvs, functional, tblOpt)
	if err != nil {
		log.Printf("%s: failed buildHamt32(%q, kvs#%d, %t, %s) => %s", name,
			name, len(kvs), false, hamt32.TableOptionName[tblOpt], err)
		b.Fatalf("%s: failed buildHamt32(%q, kvs#%d, %t, %s) => %s", name,
			name, len(kvs), false, hamt32.TableOptionName[tblOpt], err)
	}

	log.Printf("%s: b.N=%d", name, b.N)
	b.ResetTimer()

	var stats *hamt32.Stats
	for i := 0; i < b.N; i++ {
		stats = h.Stats()
	}

	log.Printf("%s: stats=%+v;\n", name, stats)
}
