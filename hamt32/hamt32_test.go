package hamt32_test

import (
	"bytes"
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

	for _, sv := range BVS[:30] {
		var bs = sv.Bsl
		var v = sv.Val

		var inserted bool
		h, inserted = h.Put(bs, v)
		if !inserted {
			log.Printf("%s: failed to insert s=%q, v=%d", name, string(bs), v)
			t.Fatalf("%s: failed to insert s=%q, v=%d", name, string(bs), v)
		}

		//log.Print(h.LongString(""))
	}
}

func TestHamt32Put(t *testing.T) {
	var name = "TestHamt32Put"
	if Functional {
		name += ":functional:" + hamt32.TableOptionName[TableOption]
	} else {
		name += ":transient:" + hamt32.TableOptionName[TableOption]
	}

	StartTime[name] = time.Now()
	Hamt32 = hamt32.New(Functional, TableOption)
	for _, sv := range BVS {
		var bs = sv.Bsl
		var v = sv.Val

		var inserted bool
		Hamt32, inserted = Hamt32.Put(bs, v)
		if !inserted {
			log.Printf("%s: failed to Hamt32.Put(%q, %v)", name, string(bs), v)
			t.Fatalf("%s: failed to Hamt32.Put(%q, %v)", name, string(bs), v)
		}

		var val, found = Hamt32.Get(bs)
		if !found {
			log.Printf("%s: failed to Hamt32.Get(%q)", name, string(bs))
			//log.Print(Hamt32.LongString(""))
			t.Fatalf("%s: failed to Hamt32.Get(%q)", name, string(bs))
		}
		if val != v {
			log.Printf("%s: returned val,%d != expected v,%d for s=%q", name, val, v, string(bs))
			t.Fatalf("%s: returned val,%d != expected v,%d for s=%q", name, val, v, string(bs))
		}
	}
	RunTime[name] = time.Since(StartTime[name])

	StartTime["Hamt32.Stats()"] = time.Now()
	var stats = Hamt32.Stats()
	RunTime["Hamt32.Stats()"] = time.Since(StartTime["Hamt32.Stats()"])
	log.Printf("%s: stats=%+v;\n", name, stats)
}

func TestHamt32IterFunc(t *testing.T) {
	var name = "TestHamt32IterFunc"
	if Functional {
		name += ":functional:" + hamt32.TableOptionName[TableOption]
	} else {
		name += ":transient:" + hamt32.TableOptionName[TableOption]
	}

	if Hamt32 == nil {
		var err error
		Hamt32, err = buildHamt32(name, BVS, Functional, TableOption)
		if err != nil {
			log.Printf("%s: failed buildHamt32(%q, BVS#%d, %t, %s) => %s", name,
				name, len(BVS), Functional,
				hamt32.TableOptionName[TableOption], err)
			t.Fatalf("%s: failed buildHamt32(%q, BVS#%d, %t, %s) => %s", name,
				name, len(BVS), Functional,
				hamt32.TableOptionName[TableOption], err)
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

	if len(BVS) != i {
		t.Fatalf("Expected len(BVS),%d got i,%d; Hamt32.Nentries()=%d;",
			len(BVS), i, Hamt32.Nentries())
	}

	RunTime[name] = time.Since(StartTime[name])
}

func TestHamt32IterChan(t *testing.T) {
	var name = "TestHamt32IterChan"
	if Functional {
		name += ":functional:" + hamt32.TableOptionName[TableOption]
	} else {
		name += ":transient:" + hamt32.TableOptionName[TableOption]
	}

	if Hamt32 == nil {
		var err error
		Hamt32, err = buildHamt32(name, BVS, Functional, TableOption)
		if err != nil {
			log.Printf("%s: failed buildHamt32(%q, BVS#%d, %t, %s) => %s", name,
				name, len(BVS), Functional,
				hamt32.TableOptionName[TableOption], err)
			t.Fatalf("%s: failed buildHamt32(%q, BVS#%d, %t, %s) => %s", name,
				name, len(BVS), Functional,
				hamt32.TableOptionName[TableOption], err)
		}

		StartTime["Hamt32.Stats()"] = time.Now()
		var stats = Hamt32.Stats()
		RunTime["Hamt32.Stats()"] = time.Since(StartTime["Hamt32.Stats()"])
		log.Printf("%s: stats=%+v;\n", name, stats)
	}

	StartTime[name] = time.Now()

	var i int
	for kv := range Hamt32.IterChan(0) {
		var val, ok = Hamt32.Get(kv.Key)
		if !ok {
			t.Fatalf("failed to lookup %s in Hamt32", kv.Key)
		}

		if val != kv.Val {
			t.Fatalf("val,%v != kv.Val,%v\n", val, kv.Val)
		}

		i++
	}

	if len(BVS) != i {
		t.Fatalf("Expected len(BVS),%d go i,%d", len(BVS), i)
	}

	RunTime[name] = time.Since(StartTime[name])
}

func TestHamt32IterChanCancel(t *testing.T) {
	var name = "TestHamt32IterChanCancel"
	if Functional {
		name += ":functional:" + hamt32.TableOptionName[TableOption]
	} else {
		name += ":transient:" + hamt32.TableOptionName[TableOption]
	}

	if Hamt32 == nil {
		var err error
		Hamt32, err = buildHamt32(name, BVS, Functional, TableOption)
		if err != nil {
			log.Printf("%s: failed buildHamt32(%q, BVS#%d, %t, %s) => %s", name,
				name, len(BVS), Functional,
				hamt32.TableOptionName[TableOption], err)
			t.Fatalf("%s: failed buildHamt32(%q, BVS#%d, %t, %s) => %s", name,
				name, len(BVS), Functional,
				hamt32.TableOptionName[TableOption], err)
		}

		StartTime["Hamt32.Stats()"] = time.Now()
		var stats = Hamt32.Stats()
		RunTime["Hamt32.Stats()"] = time.Since(StartTime["Hamt32.Stats()"])
		log.Printf("%s: stats=%+v;\n", name, stats)
	}

	StartTime[name] = time.Now()

	var i int
	var stopKey = BVS[0].Bsl // "aaa" but key from iter are random
	var iterChan, iterChanCancel = Hamt32.IterChanWithCancel(0)
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
			iterChanCancel()
			break
		}
	}

	log.Printf("%s: stopped after %d iterations", name, i)

	RunTime[name] = time.Since(StartTime[name])
}

func TestHamt32Get(t *testing.T) {
	var name = "TestHamt32Get"
	if Functional {
		name += ":functional:" + hamt32.TableOptionName[TableOption]
	} else {
		name += ":transient:" + hamt32.TableOptionName[TableOption]
	}

	if Hamt32 == nil {
		var err error
		Hamt32, err = buildHamt32(name, BVS, Functional, TableOption)
		if err != nil {
			log.Printf("%s: failed buildHamt32(%q, BVS#%d, %t, %s) => %s", name,
				name, len(BVS), Functional,
				hamt32.TableOptionName[TableOption], err)
			t.Fatalf("%s: failed buildHamt32(%q, BVS#%d, %t, %s) => %s", name,
				name, len(BVS), Functional,
				hamt32.TableOptionName[TableOption], err)
		}

		StartTime["Hamt32.Stats()"] = time.Now()
		var stats = Hamt32.Stats()
		RunTime["Hamt32.Stats()"] = time.Since(StartTime["Hamt32.Stats()"])
		log.Printf("%s: stats=%+v;\n", name, stats)
	}

	StartTime[name] = time.Now()
	for _, sv := range BVS {
		var bs = sv.Bsl
		var v = sv.Val

		var val, found = Hamt32.Get(bs)
		if !found {
			log.Printf("%s: Failed to Hamt32.Get(%q)", name, string(bs))
			log.Print(Hamt32.LongString(""))
			t.Fatalf("%s: Failed to Hamt32.Get(%q)", name, string(bs))
		}
		if val != v {
			log.Printf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, string(bs))
			t.Fatalf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, string(bs))
		}
	}
	RunTime[name] = time.Since(StartTime[name])
}

func TestHamt32Del(t *testing.T) {
	var name = "TestHamt32Del"
	if Functional {
		name += ":functional:" + hamt32.TableOptionName[TableOption]
	} else {
		name += ":transient:" + hamt32.TableOptionName[TableOption]
	}

	if Hamt32 == nil {
		var err error
		Hamt32, err = buildHamt32(name, BVS, Functional, TableOption)
		if err != nil {
			log.Printf("%s: failed buildHamt32(%q, BVS#%d, %t, %s) => %s", name,
				name, len(BVS), Functional,
				hamt32.TableOptionName[TableOption], err)
			t.Fatalf("%s: failed buildHamt32(%q, BVS#%d, %t, %s) => %s", name,
				name, len(BVS), Functional,
				hamt32.TableOptionName[TableOption], err)
		}

		StartTime["Hamt32.Stats()"] = time.Now()
		var stats = Hamt32.Stats()
		RunTime["Hamt32.Stats()"] = time.Since(StartTime["Hamt32.Stats()"])
		log.Printf("%s: stats=%+v;\n", name, stats)
	}

	StartTime[name] = time.Now()
	for _, sv := range BVS {
		var bs = sv.Bsl
		var v = sv.Val

		var val interface{}
		var deleted bool
		Hamt32, val, deleted = Hamt32.Del(bs)
		if !deleted {
			log.Printf("%s: Failed to Hamt32.Del(%q)", name, string(bs))
			log.Print(Hamt32.LongString(""))
			t.Fatalf("%s: Failed to Hamt32.Del(%q)", name, string(bs))
		}
		if val != v {
			log.Printf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, string(bs))
			t.Fatalf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, string(bs))
		}
	}
	RunTime[name] = time.Since(StartTime[name])
}

var BenchHamt32Get hamt32.Hamt
var BenchHamt32Get_Functional bool

func BenchmarkHamt32Get(b *testing.B) {
	var name = "BenchmarkHamt32Get"
	if Functional {
		name += ":functional:" + hamt32.TableOptionName[TableOption]
	} else {
		name += ":transient:" + hamt32.TableOptionName[TableOption]
	}

	if BenchHamt32Get == nil || BenchHamt32Get_Functional != Functional {
		BenchHamt32Get_Functional = Functional

		var err error
		BenchHamt32Get, err = buildHamt32(name, BVS, Functional, TableOption)
		if err != nil {
			log.Printf("%s: failed buildHamt32(%q, BVS#%d, %t, %s) => %s", name,
				name, len(BVS), false, hamt32.TableOptionName[TableOption], err)
			b.Fatalf("%s: failed buildHamt32(%q, BVS#%d, %t, %s) => %s", name,
				name, len(BVS), false, hamt32.TableOptionName[TableOption], err)
		}
	}

	log.Printf("%s: b.N=%d", name, b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var j = i % len(BVS)
		var bs = BVS[j].Bsl
		var v = BVS[j].Val

		var val, found = BenchHamt32Get.Get(bs)
		if !found {
			log.Printf("%s: Failed to h.Get(%q)", name, string(bs))
			//log.Print(h.LongString(""))
			b.Fatalf("%s: Failed to h.Get(%q)", name, string(bs))
		}
		if val != v {
			log.Printf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, string(bs))
			b.Fatalf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, string(bs))
		}
	}
}

var BenchHamt32_T2F hamt32.Hamt

func BenchmarkHamt32_T2F_Get(b *testing.B) {
	var name = "BenchmarkHamt32_T2F_Get"
	name += ":functional:" + hamt32.TableOptionName[TableOption]

	if BenchHamt32_T2F == nil {
		var err error
		BenchHamt32_T2F, err = buildHamt32(name, BVS, false, TableOption)
		if err != nil {
			log.Printf("%s: failed buildHamt32(%q, BVS#%d, %t, %s) => %s", name,
				name, len(BVS), false, hamt32.TableOptionName[TableOption], err)
			b.Fatalf("%s: failed buildHamt32(%q, BVS#%d, %t, %s) => %s", name,
				name, len(BVS), false, hamt32.TableOptionName[TableOption], err)
		}
		BenchHamt32_T2F = BenchHamt32_T2F.ToFunctional()
	}

	log.Printf("%s: Transient-to-Functional; b.N=%d;", name, b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var j = i % len(BVS)
		var bs = BVS[j].Bsl
		var v = BVS[j].Val

		var val, found = BenchHamt32_T2F.Get(bs)
		if !found {
			log.Printf("%s: Failed to BenchHamt32_T2F.Get(%q)", name, string(bs))
			//log.Print(h.LongString(""))
			b.Fatalf("%s: Failed to BenchHamt32_T2F.Get(%q)", name, string(bs))
		}
		if val != v {
			log.Printf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, string(bs))
			b.Fatalf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, string(bs))
		}
	}
}

var BenchHamt32_F2T hamt32.Hamt

func BenchmarkHamt32_F2T_Get(b *testing.B) {
	var name = "BenchmarkHamt32_F2T_Get"
	name += ":transient:" + hamt32.TableOptionName[TableOption]

	if BenchHamt32_F2T == nil {
		var err error
		BenchHamt32_F2T, err = buildHamt32(name, BVS, true, TableOption)
		if err != nil {
			log.Printf("%s: failed buildHamt32(%q, BVS#%d, %t, %s) => %s", name,
				name, len(BVS), false, hamt32.TableOptionName[TableOption], err)
			b.Fatalf("%s: failed buildHamt32(%q, BVS#%d, %t, %s) => %s", name,
				name, len(BVS), false, hamt32.TableOptionName[TableOption], err)
		}
		BenchHamt32_F2T = BenchHamt32_F2T.ToTransient()
	}

	log.Printf("%s: Functional-to-Transient; b.N=%d;", name, b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var j = i % len(BVS)
		var bs = BVS[j].Bsl
		var v = BVS[j].Val

		var val, found = BenchHamt32_F2T.Get(bs)
		if !found {
			log.Printf("%s: Failed to BenchHamt32_F2T.Get(%q)", name, string(bs))
			//log.Print(h.LongString(""))
			b.Fatalf("%s: Failed to BenchHamt32_F2T.Get(%q)", name, string(bs))
		}
		if val != v {
			log.Printf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, string(bs))
			b.Fatalf("%s: retrieved val,%d != expected v,%d for s=%q", name, val, v, string(bs))
		}
	}
}

func BenchmarkHamt32Put(b *testing.B) {
	var name = "BenchmarkHamt32Put"
	if Functional {
		name += ":functional:" + hamt32.TableOptionName[TableOption]
	} else {
		name += ":transient:" + hamt32.TableOptionName[TableOption]
	}

	if b.N+InitHamtNumBvsForPut > len(BVS) {
		log.Printf("%s: Can't run: b.N+num > len(BVS)", name)
		b.Fatalf("%s: Can't run: b.N+num > len(BVS)", name)
	}

	var bvs = BVS[:InitHamtNumBvsForPut]

	var h, err = buildHamt32(name, bvs, Functional, TableOption)
	if err != nil {
		log.Printf("%s: failed buildHamt32(%q, BVS[:%d], %t, %s) => %s", name,
			name, InitHamtNumBvsForPut, Functional,
			hamt32.TableOptionName[TableOption], err)
		b.Fatalf("%s: failed buildHamt32(%q, BVS[:%d], %t, %s) => %s", name,
			name, InitHamtNumBvsForPut, Functional,
			hamt32.TableOptionName[TableOption], err)
	}

	log.Printf("%s: b.N=%d;", name, b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var bs = BVS[InitHamtNumBvsForPut+i].Bsl
		var v = BVS[InitHamtNumBvsForPut+i].Val

		var added bool
		h, added = h.Put(bs, v)
		if !added {
			log.Printf("%s: failed to h.Put(%q, %d)", name, string(bs), v)
			//log.Print(h.LongString(""))
			b.Fatalf("%s: failed to h.Put(%q, %d)", name, string(bs), v)
		}
	}
}

func BenchmarkHamt32_T2F_Put(b *testing.B) {
	var name = "BenchmarkHamt32Put_T2F"
	name += ":functional:" + hamt32.TableOptionName[TableOption]

	var InitHamtNumBvsForPut int //= 1000000 // 1 million; allows b.N=3,000,000
	if b.N+InitHamtNumBvsForPut > len(BVS) {
		log.Printf("%s: Can't run: b.N+num > len(BVS)", name)
		b.Fatalf("%s: Can't run: b.N+num > len(BVS)", name)
	}

	var bvs = BVS[:InitHamtNumBvsForPut]

	var h, err = buildHamt32(name, bvs, false, TableOption)
	if err != nil {
		log.Printf("%s: failed buildHamt32(%q, BVS[:%d], %t, %s) => %s", name,
			name, InitHamtNumBvsForPut, Functional,
			hamt32.TableOptionName[TableOption], err)
		b.Fatalf("%s: failed buildHamt32(%q, BVS[:%d], %t, %s) => %s", name,
			name, InitHamtNumBvsForPut, Functional,
			hamt32.TableOptionName[TableOption], err)
	}
	h = h.ToFunctional()

	log.Printf("%s: Transient-to-Functional; b.N=%d;", name, b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var bs = BVS[InitHamtNumBvsForPut+i].Bsl
		var v = BVS[InitHamtNumBvsForPut+i].Val

		var added bool
		h, added = h.Put(bs, v)
		if !added {
			log.Printf("%s: failed to h.Put(%q, %d)", name, string(bs), v)
			//log.Print(h.LongString(""))
			b.Fatalf("%s: failed to h.Put(%q, %d)", name, string(bs), v)
		}
	}
}

func BenchmarkHamt32Del(b *testing.B) {
	var name = "BenchmarkHamt32Del"
	if Functional {
		name += ":functional:" + hamt32.TableOptionName[TableOption]
	} else {
		name += ":transient:" + hamt32.TableOptionName[TableOption]
	}

	var h, err = buildHamt32(name, BVS[:TwoKK], Functional, TableOption)
	if err != nil {
		log.Printf("%s: failed buildHamt32(%q, BVS:%d, %t, %s) => %s", name,
			name, len(BVS), Functional,
			hamt32.TableOptionName[TableOption], err)
		b.Fatalf("%s: failed buildHamt32(%q, BVS:%d, %t, %s) => %s", name,
			name, len(BVS), Functional,
			hamt32.TableOptionName[TableOption], err)
	}

	log.Printf("%s: b.N=%d;", name, b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var bs = BVS[i].Bsl
		var v = BVS[i].Val

		var deleted bool
		var val interface{}
		h, val, deleted = h.Del(bs)
		if !deleted {
			log.Printf("%s: failed to h.Del(%q)", name, string(bs))
			//log.Print(h.LongString(""))
			b.Fatalf("%s: failed to h.Del(%q)", name, string(bs))
		}
		if val != v {
			log.Printf("%s: failed val,%d != v,%d", name, val, v)
			b.Fatalf("%s: failed val,%d != v,%d", name, val, v)
		}
	}
}

func BenchmarkHamt32IterFunc(b *testing.B) {
	var name = "BenchmarkHamt32Iter"
	if Functional {
		name += ":functional:" + hamt32.TableOptionName[TableOption]
	} else {
		name += ":transient:" + hamt32.TableOptionName[TableOption]
	}

	// We will reuse the BenchHamtGet tree for iterator benchmarking.
	if BenchHamt32Get == nil || BenchHamt32Get_Functional != Functional {
		BenchHamt32Get_Functional = Functional

		var err error
		BenchHamt32Get, err = buildHamt32(name, BVS, Functional, TableOption)
		if err != nil {
			log.Printf("%s: failed buildHamt32(%q, BVS#%d, %t, %s) => %s", name,
				name, len(BVS), false, hamt32.TableOptionName[TableOption], err)
			b.Fatalf("%s: failed buildHamt32(%q, BVS#%d, %t, %s) => %s", name,
				name, len(BVS), false, hamt32.TableOptionName[TableOption], err)
		}
	}

	log.Printf("%s: b.N=%d", name, b.N)
	b.ResetTimer()

	var i int
	var next = BenchHamt32Get.Iter()
	for kv, ok := next(); ok; kv, ok = next() {
		if len(kv.Key) < 0 {
			b.Fatal("stupid test to touch the kv")
		}

		if i >= b.N {
			break
		}
		i++
	}

	if i != b.N {
		b.Fatalf("Failed to run b.N,%d iterations; only ran %d.", b.N, i)
	}
}

func BenchmarkHamt32IterChan(b *testing.B) {
	var name = "BenchmarkHamt32Iter"
	if Functional {
		name += ":functional:" + hamt32.TableOptionName[TableOption]
	} else {
		name += ":transient:" + hamt32.TableOptionName[TableOption]
	}

	// We will reuse the BenchHamtGet tree for iterator benchmarking.
	if BenchHamt32Get == nil || BenchHamt32Get_Functional != Functional {
		BenchHamt32Get_Functional = Functional

		var err error
		BenchHamt32Get, err = buildHamt32(name, BVS, Functional, TableOption)
		if err != nil {
			log.Printf("%s: failed buildHamt32(%q, BVS#%d, %t, %s) => %s", name,
				name, len(BVS), false, hamt32.TableOptionName[TableOption], err)
			b.Fatalf("%s: failed buildHamt32(%q, BVS#%d, %t, %s) => %s", name,
				name, len(BVS), false, hamt32.TableOptionName[TableOption], err)
		}
	}

	log.Printf("%s: b.N=%d", name, b.N)
	b.ResetTimer()

	var i int
	var iterChan, iterChanCancel = BenchHamt32Get.IterChanWithCancel(0)
	for kv := range iterChan {
		if len(kv.Key) < 0 {
			b.Fatal("stupid test to touch the kv")
		}

		if i >= b.N {
			iterChanCancel()
			break
		}
		i++
	}

	if i != b.N {
		b.Fatalf("Failed to run b.N,%d iterations; only ran %d.", b.N, i)
	}
}
