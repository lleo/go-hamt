package hamt32

import (
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
	//var fullonly, componly, hybrid bool
	//flag.BoolVar(&fullonly, "F", false, "Use fullTables only and exclude C and H options.")
	//flag.BoolVar(&componly, "C", false, "Use compressedTables only and exclude F and H options.")
	//flag.BoolVar(&hybrid, "H", false, "Use bot compressed tables initially and exclude F and C options.")
	//
	//flag.Parse()
	//
	//if (fullonly && (componly || hybrid)) ||
	//	(componly && (fullonly || hybrid)) ||
	//	(hybrid && (componly || fullonly)) {
	//	flag.PrintDefaults()
	//	os.Exit(1)
	//}
	//
	//if fullonly {
	//	options = FULLONLY
	//}
	//if componly {
	//	options = COMPONLY
	//}
	//if hybrid || (!fullonly && !componly) {
	//	options = HYBRID
	//}

	log.SetFlags(log.Lshortfile)

	var logfile, err = os.Create("test.log")
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to os.Create(\"test.log\")"))
	}
	defer logfile.Close()

	log.SetOutput(logfile)
	StartTime["TestMain"] = time.Now()
	log.Println("begin TestMain")

	//log.Printf("\nfullonly=%t\ncomponly=%t\nhybrid=%t\n", fullonly, componly, hybrid)
	//log.Printf("\noptions=%d\n", options)
	//
	//if options == FULLONLY {
	//	log.Println("options == FULLONLY")
	//}
	//if options == COMPONLY {
	//	log.Println("options == COMPONLY")
	//}
	//if options == HYBRID {
	//	log.Println("options == HYBRID")
	//}

	StartTime["TestMain: build Lookup/Delete Map/Hamt32"] = time.Now()

	hugeKvs = buildKeyVals(numHugeKvs)

	StartTime["TestMain: build Lookup/Delete Hamt32"] = time.Now()

	for _, kv := range genRandomizedKvs(hugeKvs) {
		strkey := kv.key.(*stringkey.StringKey).Str()
		intval := kv.val.(int)
		LookupMap[strkey] = intval
		DeleteMap[strkey] = intval
	}

	RunTime["TestMain: build Lookup/Delete Hamt32"] = time.Since(StartTime["TestMain: build Lookup/Delete Hamt32"])

	options = FULLONLY
	buildLookupDeleteHamt32(options, hugeKvs)
	xit := m.Run()
	if xit != 0 {
		os.Exit(xit)
	}
	log.Println("\n", RunTimes())

	options = COMPONLY
	buildLookupDeleteHamt32(options, hugeKvs)
	xit = m.Run()
	if xit != 0 {
		os.Exit(xit)
	}
	log.Println("\n", RunTimes())

	options = COMPONLY
	buildLookupDeleteHamt32(options, hugeKvs)
	xit = m.Run()
	if xit != 0 {
		os.Exit(xit)
	}
	log.Println("\n", RunTimes())

	RunTime["TestMain"] = time.Since(StartTime["TestMain"])

	log.Println("end TestMain")

	os.Exit(xit)
}

func buildLookupDeleteHamt32(options int, kvs []keyVal) {
	k := fmt.Sprintf("buildLookupDeleteHamt32: options=%d", options)

	log.Println(k)

	StartTime[k] = time.Now()

	LookupHamt32 = New(options)
	DeleteHamt32 = New(options)

	for _, kv := range genRandomizedKvs(kvs) {
		inserted := LookupHamt32.Put(kv.key, kv.val)
		if !inserted {
			log.Fatalf("failed to LookupHamt32.Put(%s, %v)", kv.key, kv.val)
		}

		inserted = DeleteHamt32.Put(kv.key, kv.val)
		if !inserted {
			log.Fatalf("failed to DeleteHamt32.Put(%s, %v)", kv.key, kv.val)
		}
	}

	RunTime[k] = time.Since(StartTime[k])
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
