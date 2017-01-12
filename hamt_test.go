package hamt_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	//"github.com/lleo/go-hamt/hamt32"
	//"github.com/lleo/go-hamt/hamt64"

	"github.com/lleo/go-hamt"
	"github.com/lleo/go-hamt/key"
	"github.com/lleo/go-hamt/stringkey"
	"github.com/lleo/stringutil"
	"github.com/pkg/errors"
)

// This number was found by experimenting with TestBenchmarkDel() and the
// testing harness kept calling the test with greater and greater b.N values.
// It toped out at 3,000,000 .
var numHugeKvs = 5 * 1024 * 1024 // five mega-entries
var hugeKvs []key.KeyVal

var LookupMap map[string]int
var DeleteMap map[string]int
var keyStrings []string

var Inc = stringutil.Lower.Inc

var StartTime = make(map[string]time.Time)
var RunTime = make(map[string]time.Duration)

func TestMain(m *testing.M) {
	//	var fullonly, componly, hybrid, all bool
	//	flag.BoolVar(&fullonly, "F", false, "Use full tables only and exclude C and H Options.")
	//	flag.BoolVar(&componly, "C", false, "Use compressed tables only and exclude F and H Options.")
	//	flag.BoolVar(&hybrid, "H", false, "Use compressed tables initially and exclude F and C Options.")
	//	flag.BoolVar(&all, "A", false, "Run all Tests w/ Options set to hamt32.FullTablesOnly, hamt32.CompTablesOnly, and hamt32.HybridTables; in that order.")
	//
	//	flag.Parse()
	//
	//	// If all flag set, ignore fullonly, componly, and hybrid.
	//	if !all {
	//
	//		// only one flag may be set between fullonly, componly, and hybrid
	//		if (fullonly && (componly || hybrid)) ||
	//			(componly && (fullonly || hybrid)) ||
	//			(hybrid && (componly || fullonly)) {
	//			flag.PrintDefaults()
	//			os.Exit(1)
	//		}
	//	}
	//
	//	// If no flags given, run all tests.
	//	if !(all || fullonly || componly || hybrid) {
	//		all = true
	//	}

	log.SetFlags(log.Lshortfile)

	var logfile, err = os.Create("test.log")
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to os.Create(\"test.log\")"))
	}
	defer logfile.Close()

	log.SetOutput(logfile)

	StartTime["TestMain: build Looup/Delete Map"] = time.Now()

	hugeKvs = buildKeyVals(numHugeKvs)

	LookupMap = make(map[string]int, numHugeKvs)
	DeleteMap = make(map[string]int, numHugeKvs)
	keyStrings = make([]string, numHugeKvs)

	for i, kv := range hugeKvs {
		var str = kv.Key.(*stringkey.StringKey).Str()
		var val = i

		LookupMap[str] = val
		DeleteMap[str] = val
		keyStrings[i] = str
	}

	RunTime["TestMain: build Looup/Delete Map"] = time.Since(StartTime["TestMain: build Looup/Delete Map"])

	log.Println("TestMain: before Running Tests")

	//RUN
	var xit = m.Run()

	log.Println("TestMain: after Running Tests")
	log.Println("TestMain:\n", RunTimes())

	//TEARDOWN
	os.Exit(xit)
}

func buildKeyVals(num int) []key.KeyVal {
	var kvs = make([]key.KeyVal, num, num)

	s := "aaa"
	for i := 0; i < num; i++ {
		kvs[i].Key = stringkey.New(s)
		kvs[i].Val = i

		s = Inc(s)
	}

	return kvs
}

//func genRandomizedKvs(kvs []key.KeyVal) []key.KeyVal {
//	randKvs := make([]key.KeyVal, len(kvs))
//	copy(randKvs, kvs)
//
//	//From: https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle#The_modern_algorithm
//	for i := len(randKvs) - 1; i > 0; i-- {
//		j := rand.Intn(i + 1)
//		randKvs[i], randKvs[j] = randKvs[j], randKvs[i]
//	}
//
//	return randKvs
//}

func rebuildDeleteMap(strs []string) {
	for val, str := range strs {
		_, exists := DeleteMap[str]
		if exists {
			break
		}

		DeleteMap[str] = val
	}
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

func TestNewHamt32(t *testing.T) {
	log.Println("TestNewHamt32:")
	var h = hamt.NewHamt32()
	if !h.IsEmpty() {
		t.Fatal("!?!? a brand new Hamt !IsEmpty()")
	}
}

func TestNewHamt64(t *testing.T) {
	log.Println("TestNewHamt64:")
	var h = hamt.NewHamt32()
	if !h.IsEmpty() {
		t.Fatal("!?!? a brand new Hamt !IsEmpty()")
	}
}
