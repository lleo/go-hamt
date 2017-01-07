package hamt32

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/lleo/go-hamt/stringkey"
	"github.com/lleo/stringutil"
	"github.com/pkg/errors"
)

var numHugeKvs int = 5 * 1024 * 1024

//var numMaxKvs int = math.MaxInt32

var hugeKvs []keyVal
var maxKvs []keyVal

var LookupMap = make(map[string]int, numHugeKvs)
var DeleteMap = make(map[string]int, numHugeKvs)

var LookupHamt32 *Hamt
var DeleteHamt32 *Hamt

var StartTime = make(map[string]time.Time)
var RunTime = make(map[string]time.Duration)

var options int

func TestMain(m *testing.M) {
	var fullonly, componly, hybrid, all bool
	flag.BoolVar(&fullonly, "F", false, "Use fullTables only and exclude C and H options.")
	flag.BoolVar(&componly, "C", false, "Use compressedTables only and exclude F and H options.")
	flag.BoolVar(&hybrid, "H", false, "Use compressed tables initially and exclude F and C options.")
	flag.BoolVar(&all, "A", false, "Run all Tests w/ options set to FULLONLY, COMPONLY, and HYBRID; in that order.")

	flag.Parse()

	// If all flag set, ignore fullonly, componly, and hybrid.
	if !all {

		// only one flag may be set between fullonly, componly, and hybrid
		if (fullonly && (componly || hybrid)) ||
			(componly && (fullonly || hybrid)) ||
			(hybrid && (componly || fullonly)) {
			flag.PrintDefaults()
			os.Exit(1)
		}
	}

	// If no flags given, run all tests.
	if !(all || fullonly || componly || hybrid) {
		all = true
	}

	log.SetFlags(log.Lshortfile)

	var logfile, err = os.Create("test.log")
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to os.Create(\"test.log\")"))
	}
	defer logfile.Close()

	log.SetOutput(logfile)

	StartTime["TestMain"] = time.Now()
	log.Println("begin TestMain")

	log.Printf("\nall=%t\nfullonly=%t\ncomponly=%t\nhybrid=%t\n", all, fullonly, componly, hybrid)

	StartTime["TestMain: build Lookup/Delete Map/Hamt32"] = time.Now()

	hugeKvs = buildKeyVals(numHugeKvs)

	StartTime["TestMain: build Lookup/Delete Hamt32"] = time.Now()

	LookupHamt32 = New(options)
	DeleteHamt32 = New(options)

	for _, kv := range genRandomizedKvs(hugeKvs) {
		strkey := kv.key.(*stringkey.StringKey).Str()
		intval := kv.val.(int)
		LookupMap[strkey] = intval
		DeleteMap[strkey] = intval

		inserted := LookupHamt32.Put(kv.key, kv.val)
		if !inserted {
			log.Fatalf("failed to LookupHamt32.Put(%s, %v)", kv.key, kv.val)
		}

		inserted = DeleteHamt32.Put(kv.key, kv.val)
		if !inserted {
			log.Fatalf("failed to DeleteHamt32.Put(%s, %v)", kv.key, kv.val)
		}
	}

	RunTime["TestMain: build Lookup/Delete Hamt32"] = time.Since(StartTime["TestMain: build Lookup/Delete Hamt32"])

	var xit int

	if all {
		options = FULLONLY
		log.Printf("options == %s", OPTIONS[options])
		rebuildDeleteMap(hugeKvs)
		rebuildDeleteHamt32(hugeKvs)
		xit = m.Run()
		if xit != 0 {
			goto SKIPTESTS
		}

		options = COMPONLY
		log.Printf("options == %s", OPTIONS[options])
		rebuildDeleteMap(hugeKvs)
		rebuildDeleteHamt32(hugeKvs)
		xit = m.Run()
		if xit != 0 {
			goto SKIPTESTS
		}

		options = HYBRID
		log.Printf("options == %s", OPTIONS[options])
		rebuildDeleteMap(hugeKvs)
		rebuildDeleteHamt32(hugeKvs)
		xit = m.Run()

	SKIPTESTS:
	} else {
		if fullonly {
			options = FULLONLY
		}
		if componly {
			options = COMPONLY
		}
		if hybrid || (!fullonly && !componly) {
			options = HYBRID
		}
		log.Printf("\noptions=%d\n", options)

		log.Printf("options == %s", OPTIONS[options])

		xit = m.Run()
	}

	RunTime["TestMain"] = time.Since(StartTime["TestMain"])

	log.Println("\n", RunTimes())
	log.Println("end TestMain")

	os.Exit(xit)
}

func rebuildDeleteMap(kvs []keyVal) {
	for _, kv := range kvs {
		sk := kv.key.(*stringkey.StringKey)
		str := sk.Str()
		val := kv.val.(int)

		if _, ok := DeleteMap[str]; !ok {
			DeleteMap[str] = val
		} else {
			// we delete inorder so we can stop rebuilding when the entries start existing
			break
		}
	}
}

func rebuildDeleteHamt32(kvs []keyVal) {
	for _, kv := range kvs {
		inserted := DeleteHamt32.Put(kv.key, kv.val)
		if !inserted {
			//log.Printf("BenchmarkHamt32Del: inserted,%v := DeleteHamt32.Put(%s, %d)", inserted, kv.key, kv.val)

			// we delete inorder so we can stop rebuilding when the entries start existing
			break
		}
	}
}

func buildKeyVals(num int) []keyVal {
	var kvs = make([]keyVal, num, num)

	//var s = stringutil.Str("aaa")
	s := "aaa"

	for i := 0; i < num; i++ {
		kvs[i].key = stringkey.New(s)
		kvs[i].val = i

		//s = s.DigitalInc(1)
		s = stringutil.DigitalInc(s)
	}

	return kvs
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

func RunTimes() string {
	var s = ""

	s += "Key                                      Val\n"
	s += "========================================+==========\n"

	for key, val := range RunTime {
		s += fmt.Sprintf("%-40s %s\n", key, val)
	}
	return s
}
