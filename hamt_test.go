package hamt

import (
	"log"
	"math/rand"
	"os"
	"testing"

	hamt32 "github.com/lleo/go-hamt/hamt32"

	"github.com/lleo/stringutil"
)

var numMidKvs int
var numHugeKvs int
var midKvs []keyVal
var hugeKvs []keyVal

var M map[string]int
var H *hamt32.Hamt

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
		var key = StringKey(s0)
		var val = i
		midKvs = append(midKvs, keyVal{key, val})
		s0 = s0.DigitalInc(1) //get off "" first
	}

	hugeKvs = make([]keyVal, 0, 32)
	var s1 = stringutil.Str("aaa")
	//numHugeKvs = 8 * 1024
	numHugeKvs = 1 * 1024 * 1024 // one mega-entries
	//numHugeKvs = 256 * 1024 * 1024 //256 MB
	for i := 0; i < numHugeKvs; i++ {
		var key = StringKey(s1)
		var val = i
		//log.Printf("numHugeKvs[%d] val=%d; key=%s", i, i, s1)
		hugeKvs = append(hugeKvs, keyVal{key, val})
		s1 = s1.DigitalInc(1)
	}

	// Build map & hamt, for h.Get() and h.Del() benchmarks
	M = make(map[string]int)
	H = hamt32.NewHamt()
	var s = stringutil.Str("aaa")
	for i := 0; i < numHugeKvs; i++ {
		var val = i
		M[string(s)] = val
		H.Put(StringKey(s), val)
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
	var k = StringKey(s)
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
	var k = StringKey(s)
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
		var inserted = h.Put(midKvs[i].key, midKvs[i].val)
		if !inserted {
			t.Fatalf("h.Put(%s, %v): for i=%d returned false", midKvs[i].key, midKvs[i].val, i)
		}
	}

	for i := 0; i < numMidKvs; i++ {
		var vv, found = h.Get(midKvs[i].key)
		if !found {
			t.Fatalf("h.Get(%s): for i=%d returned !found", midKvs[i].key, i)
		}
		//v := vv.(int)
		if vv != midKvs[i].val {
			t.Fatalf("h.Get(%s): returned vv,%v != midKvs[%d].val,%v", midKvs[i].key, vv, i, midKvs[i].val)
		}
	}
}

func TestPutDelMid(t *testing.T) {
	//log.Println("=== TestPutDelMid ===")
	var h = hamt32.NewHamt32()

	for i := 0; i < numMidKvs; i++ {
		var inserted = h.Put(midKvs[i].key, midKvs[i].val)
		if !inserted {
			t.Fatalf("h.Put(%s, %v): for i=%d returned false", midKvs[i].key, midKvs[i].val, i)
		}
	}

	//log.Println("h =", h.LongString(""))

	for i := 0; i < numMidKvs; i++ {

		//if midKvs[i].key.Equals(StringKey("aba")) {
		//	log.Println("before h.Del(%s) h =\n%s", midKvs[i].key, h.LongString(""))
		//}

		var vv, deleted = h.Del(midKvs[i].key)
		if !deleted {
			//log.Printf("h.Del(%s): failed h =\n%s", midKvs[i].key, h.LongString(""))
			t.Fatalf("h.Del(%s): for i=%d return !deleted", midKvs[i].key, i)
		}
		if vv != midKvs[i].val {
			t.Fatalf("h.Del(%s): returned vv,%v != midKvs[%d].val,%v", midKvs[i].key, vv, i, midKvs[i].val)
		}
		//log.Println("h =", h.LongString(""))
	}
}

func TestPutGetHuge(t *testing.T) {
	//log.Println("=== TestPutGetHuge ===")
	var h = hamt32.NewHamt32()

	for i := 0; i < numHugeKvs; i++ {
		var inserted = h.Put(hugeKvs[i].key, hugeKvs[i].val)
		if !inserted {
			t.Fatalf("h.Put(%s, %v): for i=%d returned false", hugeKvs[i].key, hugeKvs[i].val, i)
		}
	}

	for i := 0; i < numHugeKvs; i++ {
		var vv, found = h.Get(hugeKvs[i].key)
		if !found {
			t.Fatalf("h.Get(%s): for i=%d returned !found", hugeKvs[i].key, i)
		}
		//v := vv.(int)
		if vv != hugeKvs[i].val {
			t.Fatalf("h.Get(%s): returned vv,%v != hugeKvs[%d].val,%v", hugeKvs[i].key, vv, i, hugeKvs[i].val)
		}
	}
}

func TestPutDelHuge(t *testing.T) {
	//log.Println("=== TestPutDelHuge ===")
	var h = hamt32.NewHamt32()

	for i := 0; i < numHugeKvs; i++ {
		var inserted = h.Put(hugeKvs[i].key, hugeKvs[i].val)
		if !inserted {
			t.Fatalf("h.Put(%s, %v): for i=%d returned false", hugeKvs[i].key, hugeKvs[i].val, i)
		}
	}

	for i := 0; i < numHugeKvs; i++ {
		var vv, deleted = h.Del(hugeKvs[i].key)
		if !deleted {
			t.Fatalf("h.Del(%s): for i=%d returned !deleted", hugeKvs[i].key, i)
		}
		if vv != hugeKvs[i].val {
			t.Fatalf("h.Del(%s): returned vv,%v != hugeKvs[%d].val,%v", hugeKvs[i].key, vv, i, hugeKvs[i].val)
		}

		//if hugeKvs[i].key.Equals(StringKey("aaaabccdefghijklmnopqrstuvwxy")) {
		//	log.Println("h =", h.LongString(""))
		//}
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
		var val, found = H.Get(key)
		if !found {
			b.Fatalf("H.Get(%s) not found", key)
		}
		if val != hugeKvs[j].val {
			b.Fatalf("val,%v != hugeKvs[%d].val,%v", val, j, hugeKvs[j].val)
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
		h.Put(StringKey(s), i+1)
		s = s.DigitalInc(1)
	}
}
