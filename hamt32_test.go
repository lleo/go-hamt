package hamt_test

import (
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
		var k = hamt32.StringKey(kv.Key)
		var v = kv.Val

		var inserted bool
		h, inserted = h.Put(k, v)
		if !inserted {
			log.Printf("%s: failed to insert s=%q, v=%d", name, k, v)
			t.Fatalf("%s: failed to insert s=%q, v=%d", name, k, v)
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
		var k = hamt32.StringKey(kv.Key)
		var v = kv.Val

		var inserted bool
		Hamt32, inserted = Hamt32.Put(k, v)
		if !inserted {
			log.Printf("%s: failed to Hamt32.Put(%q, %v)", name, k, v)
			t.Fatalf("%s: failed to Hamt32.Put(%q, %v)", name, k, v)
		}

		var val, found = Hamt32.Get(k)
		if !found {
			log.Printf("%s: failed to Hamt32.Get(%q)", name, k)
			//log.Print(Hamt32.LongString(""))
			t.Fatalf("%s: failed to Hamt32.Get(%q)", name, k)
		}
		if val != v {
			log.Printf("%s: returned val,%d != expected v,%d for s=%q", name, val, v, k)
			t.Fatalf("%s: returned val,%d != expected v,%d for s=%q", name, val, v, k)
		}
	}
	RunTime[name] = time.Since(StartTime[name])

	StartTime["Hamt32.Stats()"] = time.Now()
	var stats = Hamt32.Stats()
	RunTime["Hamt32.Stats()"] = time.Since(StartTime["Hamt32.Stats()"])
	log.Printf("%s: stats=%+v;\n", name, stats)
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
		var k = hamt32.StringKey(kv.Key)
		var v = kv.Val

		var val, found = Hamt32.Get(k)
		if !found {
			log.Printf("%s: Failed to Hamt32.Get(%q)", name, k)
			log.Print(Hamt32.LongString(""))
			t.Fatalf("%s: Failed to Hamt32.Get(%q)", name, k)
		}
		if val != v {
			log.Printf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, k)
			t.Fatalf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, k)
		}
	}
	RunTime[name] = time.Since(StartTime[name])
}

func TestHamt32Range(t *testing.T) {
	runTestHamt32Range(t, KVS, Functional, TableOption)
}

func runTestHamt32Range(
	t *testing.T,
	kvs []KeyVal,
	functional bool,
	tblOpt int,
) {
	var name = "TestHamt32Range"
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

	var kvMap = make(map[string]int, len(KVS))
	for _, kv := range KVS {
		kvMap[kv.Key] = kv.Val
	}

	// Reconstruct KVS as tmpKVS
	var totalKvs int
	var visitKeyVal = func(k hamt32.KeyI, v interface{}) bool {
		var s = string(k.(hamt32.StringKey))
		var i = v.(int)
		var expected_i, found = kvMap[s]

		if !found {
			t.Fatalf("%s: Range(visitKeyVal) KeyI.(StringKey),%q not in kvMap",
				name, s)
		}

		if expected_i != i {
			t.Fatalf("%s: Range(visitKeyVal) for KeyI.(StringKey),%q found i,%d != expected_i,%d", name, s, i, expected_i)
		}

		totalKvs++
		return true
	}
	Hamt32.Range(visitKeyVal)

	if totalKvs != len(KVS) {
		t.Fatalf("%s: Range(visitKeyVal) found totalKvs,%d != len(KVS),%d",
			name, totalKvs, len(KVS))
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
		var k = hamt32.StringKey(kv.Key)
		var v = kv.Val

		var val interface{}
		var deleted bool
		Hamt32, val, deleted = Hamt32.Del(k)
		if !deleted {
			log.Printf("%s: Failed to Hamt32.Del(%q)", name, k)
			log.Print(Hamt32.LongString(""))
			t.Fatalf("%s: Failed to Hamt32.Del(%q)", name, k)
		}
		if val != v {
			log.Printf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, k)
			t.Fatalf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, k)
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
		var k = hamt32.StringKey(kvs[j].Key)
		var v = kvs[j].Val

		var val, found = BenchHamt32Get.Get(k)
		if !found {
			log.Printf("%s: Failed to h.Get(%q)", name, k)
			//log.Print(h.LongString(""))
			b.Fatalf("%s: Failed to h.Get(%q)", name, k)
		}
		if val != v {
			log.Printf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, k)
			b.Fatalf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, k)
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
		var k = hamt32.StringKey(kvs[InitHamtNumKvsForPut+i].Key)
		var v = kvs[InitHamtNumKvsForPut+i].Val

		var added bool
		h, added = h.Put(k, v)
		if !added {
			log.Printf("%s: failed to h.Put(%q, %d)", name, k, v)
			//log.Print(h.LongString(""))
			b.Fatalf("%s: failed to h.Put(%q, %d)", name, k, v)
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
		var k = hamt32.StringKey(kvs[i].Key)
		var v = kvs[i].Val

		var deleted bool
		var val interface{}
		h, val, deleted = h.Del(k)
		if !deleted {
			log.Printf("%s: failed to h.Del(%q)", name, k)
			//log.Print(h.LongString(""))
			b.Fatalf("%s: failed to h.Del(%q)", name, k)
		}
		if val != v {
			log.Printf("%s: failed val,%d != v,%d", name, val, v)
			b.Fatalf("%s: failed val,%d != v,%d", name, val, v)
		}
	}
}

func BenchmarkHamt32Range(b *testing.B) {
	runBenchmarkHamt32Range(b, KVS, Functional, TableOption)
}

func runBenchmarkHamt32Range(
	b *testing.B,
	kvs []KeyVal,
	functional bool,
	tblOpt int,
) {
	var name = "BenchmarkHamt32Range"
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

	var kvMap = make(map[string]int, len(KVS))
	for _, kv := range KVS {
		kvMap[kv.Key] = kv.Val
	}

	log.Printf("%s: b.N=%d", name, b.N)
	b.ResetTimer()

	var i int
	h.Range(func(k hamt32.KeyI, v interface{}) bool {
		var sk = string(k.(hamt32.StringKey))
		var iv = v.(int)

		if kvMap[sk] != iv {
			b.Fatalf("%s: for kvMap[%q],%d != i,%d", name, sk, kvMap[sk], iv)
		}

		i++
		if i >= b.N {
			return false //stop Range()
		}

		return true
	})
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

func BenchmarkHamt32_GetN30(b *testing.B) {
	runBenchmarkHamt32GetN(b, KVS[:30], Functional, TableOption)
}

func BenchmarkHamt32_GetN1000(b *testing.B) {
	runBenchmarkHamt32GetN(b, KVS[:1000], Functional, TableOption)
}

func BenchmarkHamt32_GetN10000(b *testing.B) {
	runBenchmarkHamt32GetN(b, KVS[:10000], Functional, TableOption)
}

func runBenchmarkHamt32GetN(
	b *testing.B,
	kvs []KeyVal,
	functional bool,
	tblOpt int,
) {
	var name = "runBenchmarkHamt32GetN"
	if Functional {
		name += ":functional:" + hamt32.TableOptionName[tblOpt]
	} else {
		name += ":transient:" + hamt32.TableOptionName[tblOpt]
	}

	log.Println(name, b.N)

	var h, err = buildHamt32(name, kvs, functional, tblOpt)
	if err != nil {
		log.Printf("%s: failed buildHamt32(%q, kvs#%d, %t, %s) => %s", name,
			name, len(kvs), false, hamt32.TableOptionName[tblOpt], err)
		b.Fatalf("%s: failed buildHamt32(%q, kvs#%d, %t, %s) => %s", name,
			name, len(kvs), false, hamt32.TableOptionName[tblOpt], err)
	}

	log.Printf("%s: b.N=%d", name, b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var kv = kvs[i%len(kvs)]
		var val, found = h.Get(hamt32.StringKey(kv.Key))
		if !found {
			b.Fatalf("Failed to find h.Get(%q)", kv.Key)
		}

		if val != kv.Val {
			b.Fatalf("Retrieved val,%d != kv.Val,%d", val, kv.Val)
		}
	}
}
