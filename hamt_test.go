package hamt

import (
	"log"
	"math/rand"
	"os"
	"testing"

	"github.com/lleo/go-hamt/hamt32"
	"github.com/lleo/go-hamt/string_key"

	"github.com/lleo/stringutil"
)

var numMidKvs int
var numHugeKvs int
var midKvs []keyVal
var hugeKvs []keyVal

var M map[string]int
var H Hamt

func TestMain(m *testing.M) {
	//SETUP
	//var fh, err = os.OpenFile("test.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	//if err != nil {
	//	os.Exit(1)
	//}
	//defer fh.Close()

	log.SetFlags(log.Lshortfile)
	//log.SetOutput(fh)

	midKvs = make([]keyVal, 0, 32)
	var s0 = stringutil.Str("aaa")
	//numMidKvs := 10000 //ten thousand
	numMidKvs = 1000 // 10 million

	for i := 0; i < numMidKvs; i++ {
		var key = string_key.StringKey(s0)
		var val = i

		//log.Printf("numHugeKvs[%d] val=%d; key=%s", i, i, s1)
		midKvs = append(midKvs, keyVal{key, val})
		s0 = s0.DigitalInc(1) //get off "" first
	}

	hugeKvs = make([]keyVal, 0, 32)
	var s1 = stringutil.Str("aaa")
	//numHugeKvs = 8 * 1024
	numHugeKvs = 1 * 1024 * 1024 // one mega-entries
	//numHugeKvs = 256 * 1024 * 1024 //256 MB
	for i := 0; i < numHugeKvs; i++ {
		var key = string_key.StringKey(s1)
		var val = i

		//log.Printf("numHugeKvs[%d] val=%d; key=%s", i, i, s1)
		hugeKvs = append(hugeKvs, keyVal{key, val})
		s1 = s1.DigitalInc(1)
	}

	// Build map & hamt, for h.Get() and h.Del() benchmarks
	M = make(map[string]int)
	H = hamt32.NewHamt32()
	var s = stringutil.Str("aaa")
	for i := 0; i < numHugeKvs; i++ {
		var key = string_key.StringKey(s)
		var val = i

		M[string(s)] = val
		H.Put(key, val)
		s = s.DigitalInc(1)
	}

	//RUN
	var xit = m.Run()

	//TEARDOWN
	os.Exit(xit)
}

func genRandomizedKvs(kvs []keyVal) []keyVal {
	randKvs := make([]keyVal, len(kvs))
	copy(randKvs, kvs)

	//From: https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle#The_modern_algorithm
	for i := len(randKvs) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		randKvs[i], randKvs[j] = randKvs[j], randKvs[i]
	}

	return randKvs
}

func TestNewHamt32(t *testing.T) {
	//log.Println("=== TestNewHamt32 ===")
	var h = hamt32.NewHamt32()
	if !h.IsEmpty() {
		t.Fatal("!?!? a brand new Hamt !IsEmpty()")
	}
	//log.Println("TestNewHamt32 ok")
}

func TestPutGetOne(t *testing.T) {
	//log.Println("=== TestPutGetOne ===")
	var h = hamt32.NewHamt32()

	var s = stringutil.Str("aaa")
	var k = string_key.StringKey(s)
	var v int = 1

	var inserted = h.Put(k, v)
	if !inserted {
		t.Fatalf("h.Put(%s, %v) returned false", k, v)
	}

	//log.Println(h.LongString(""))

	var vv, found = h.Get(k)
	if !found {
		t.Fatalf("h.Get(%s) returned !found", k)
	}
	var val = vv.(int)
	if val != v {
		t.Fatalf("h.Get(%s) val,%d != v,%d", k, val, v)
	}

}

func TestPutDelOne(t *testing.T) {
	//log.Println("=== TestPutDelOne ===")
	var h = hamt32.NewHamt32()

	var s = stringutil.Str("aaa")
	var k = string_key.StringKey(s)
	var v int = 1

	var inserted = h.Put(k, v)
	if !inserted {
		t.Fatalf("h.Put(%s, %v) returned false", k, v)
	}

	//log.Println(h.LongString(""))

	var vv, deleted = h.Del(k)
	if !deleted {
		t.Fatalf("h.Del(%s) returned !deleted", k)
	}
	var val = vv.(int)
	if val != v {
		t.Fatalf("h.Del(%s) val,%d != v,%d", k, val, v)
	}

	//log.Println("h = ", h.LongString(""))

	if !h.IsEmpty() {
		t.Fatalf("h is not empty h=\n%s", h.LongString(""))
	}
}

func TestPutGetMid(t *testing.T) {
	//log.Println("=== TestPutGetMid ===")
	var h = hamt32.NewHamt32()

	for i := 0; i < numMidKvs; i++ {
		var key = midKvs[i].key
		var val = midKvs[i].val

		var inserted = h.Put(key, val)
		if !inserted {
			t.Fatalf("h.Put(%s, %v): for i=%d returned false", key, val, i)
		}
	}

	for i := 0; i < numMidKvs; i++ {
		var key = midKvs[i].key
		var val = midKvs[i].val

		var vv, found = h.Get(key)
		if !found {
			t.Fatalf("h.Get(%s): for i=%d returned !found", key, i)
		}
		//v := vv.(int)
		if vv != midKvs[i].val {
			t.Fatalf("h.Get(%s): returned vv,%v != midKvs[%d].val,%v", key, vv, i, val)
		}
	}
}

func TestPutDelMid(t *testing.T) {
	//log.Println("=== TestPutDelMid ===")
	var h = hamt32.NewHamt32()

	for i := 0; i < numMidKvs; i++ {
		var key = midKvs[i].key
		var val = midKvs[i].val

		var inserted = h.Put(key, val)
		if !inserted {
			t.Fatalf("h.Put(%s, %v): for i=%d returned false", key, val, i)
		}
	}

	//log.Println("h =", h.LongString(""))

	for i := 0; i < numMidKvs; i++ {
		var key = midKvs[i].key
		var val = midKvs[i].val

		var vv, deleted = h.Del(key)
		if !deleted {
			t.Fatalf("h.Del(%s): for i=%d return !deleted", key, i)
		}
		if vv != val {
			t.Fatalf("h.Del(%s): returned vv,%v != midKvs[%d].val,%v", key, vv, i, val)
		}
		//log.Println("h =", h.LongString(""))
	}
}

func TestPutGetHuge(t *testing.T) {
	//log.Println("=== TestPutGetHuge ===")
	var h = hamt32.NewHamt32()

	for i := 0; i < numHugeKvs; i++ {
		var key = hugeKvs[i].key
		var val = hugeKvs[i].val

		var inserted = h.Put(key, val)
		if !inserted {
			t.Fatalf("h.Put(%s, %v): for i=%d returned false", key, val, i)
		}
	}

	for i := 0; i < numHugeKvs; i++ {
		var key = hugeKvs[i].key
		var val = hugeKvs[i].val

		var vv, found = h.Get(key)
		if !found {
			t.Fatalf("h.Get(%s): for i=%d returned !found", key, i)
		}
		//v := vv.(int)
		if vv != val {
			t.Fatalf("h.Get(%s): returned vv,%v != hugeKvs[%d].val,%v", key, vv, i, val)
		}
	}
}

func TestPutDelHuge(t *testing.T) {
	//log.Println("=== TestPutDelHuge ===")
	var h = hamt32.NewHamt32()

	for i := 0; i < numHugeKvs; i++ {
		key := hugeKvs[i].key
		val := hugeKvs[i].val

		var inserted = h.Put(key, val)
		if !inserted {
			t.Fatalf("h.Put(%s, %v): for i=%d returned false", key, val, i)
		}
	}

	for i := 0; i < numHugeKvs; i++ {
		key := hugeKvs[i].key
		val := hugeKvs[i].val

		var vv, deleted = h.Del(key)
		if !deleted {
			t.Fatalf("h.Del(%s): for i=%d returned !deleted", key, i)
		}
		if vv != val {
			t.Fatalf("h.Del(%s): returned vv,%v != hugeKvs[%d].val,%v", key, vv, i, val)
		}
	}
}

func BenchmarkMapGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var j = int(rand.Int31()) % numHugeKvs
		var s = hugeKvs[j].key.String()
		var val, ok = M[s]
		if !ok {
			b.Fatalf("M[%s] not ok", string(s))
		}
		if val != hugeKvs[j].val {
			b.Fatalf("val,%v != hugeKvs[%d].val,%v", val, j, hugeKvs[j].val)
		}
	}
}

func BenchmarkHamtGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var j = int(rand.Int31()) % numHugeKvs
		var key = hugeKvs[j].key
		var val0 = hugeKvs[j].val
		var val, found = H.Get(key)
		if !found {
			b.Fatalf("H.Get(%s) not found", key)
		}
		if val != val0 {
			b.Fatalf("val,%v != hugeKvs[%d].val,%v", val, j, val0)
		}
	}
}

func BenchmarkMapPut(b *testing.B) {
	var m = make(map[string]int)
	var s = stringutil.Str("aaa")
	for i := 0; i < b.N; i++ {
		m[string(s)] = i + 1
		s = s.DigitalInc(1)
	}
}

func BenchmarkHamtPut(b *testing.B) {
	var h = hamt32.NewHamt32()
	var s = stringutil.Str("aaa")
	for i := 0; i < b.N; i++ {
		key := string_key.StringKey(s)
		val := i + 1
		h.Put(key, val)
		s = s.DigitalInc(1)
	}
}
