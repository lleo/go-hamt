package hamt_test

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	//"github.com/lleo/go-hamt/hamt32"
	//"github.com/lleo/go-hamt/hamt64"

	"github.com/lleo/go-hamt"
	"github.com/lleo/go-hamt/hamt32"
	"github.com/lleo/go-hamt/hamt64"
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

var TableOption int

var LookupHamt32 *hamt32.Hamt
var DeleteHamt32 *hamt32.Hamt

var LookupHamt64 *hamt64.Hamt
var DeleteHamt64 *hamt64.Hamt

var Inc = stringutil.Lower.Inc

var StartTime = make(map[string]time.Time)
var RunTime = make(map[string]time.Duration)

func TestMain(m *testing.M) {
	var fullonly, componly, hybrid, all bool
	flag.BoolVar(&fullonly, "F", false, "Use full tables only and exclude C and H Options.")
	flag.BoolVar(&componly, "C", false, "Use compressed tables only and exclude F and H Options.")
	flag.BoolVar(&hybrid, "H", false, "Use compressed tables initially and exclude F and C Options.")
	flag.BoolVar(&all, "A", false, "Run all Tests w/ Options set to hamt32.FullTablesOnly, hamt32.CompTablesOnly, and hamt32.HybridTables; in that order.")

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

	log.Println("TestMain: and so it begins...")

	StartTime["TestMain: build Looup/Delete Map"] = time.Now()

	hugeKvs = buildKeyVals(numHugeKvs)

	RunTime["TestMain: build Looup/Delete Map"] = time.Since(StartTime["TestMain: build Looup/Delete Map"])

	// execute
	var xit int
	if all {
		TableOption = hamt32.FullTablesOnly
		log.Printf("TestMain: TableOption == %s\n", hamt32.TableOptionName[TableOption])
		initData(TableOption, hugeKvs)
		xit = m.Run()
		if xit != 0 {
			os.Exit(1)
		}

		TableOption = hamt32.CompTablesOnly
		log.Printf("TestMain: TableOption == %s\n", hamt32.TableOptionName[TableOption])
		initData(TableOption, hugeKvs)
		xit = m.Run()
		if xit != 0 {
			os.Exit(1)
		}

		TableOption = hamt32.HybridTables
		log.Printf("TestMain: TableOption == %s\n", hamt32.TableOptionName[TableOption])
		initData(TableOption, hugeKvs)
		xit = m.Run()
	} else {
		if hybrid {
			TableOption = hamt32.HybridTables
		} else if fullonly {
			TableOption = hamt32.FullTablesOnly
		} else /* if componly */ {
			TableOption = hamt32.CompTablesOnly
		}

		log.Printf("TestMain: TableOption == %s\n", hamt32.TableOptionName[TableOption])
		initData(TableOption, hugeKvs)
		xit = m.Run()
	}

	log.Println("\n", RunTimes())

	os.Exit(xit)
}

func initData(tableOption int, kvs []key.KeyVal) {
	var funcName = fmt.Sprintf("hamt32: initialize(%s)", hamt32.TableOptionName[tableOption])

	var metricName = fmt.Sprintf("%s: build Lookup/Delete Hamt32", funcName)

	StartTime[metricName] = time.Now()

	LookupMap = make(map[string]int, numHugeKvs)
	DeleteMap = make(map[string]int, numHugeKvs)
	keyStrings = make([]string, numHugeKvs)

	LookupHamt32 = hamt32.New(tableOption)
	DeleteHamt32 = hamt32.New(tableOption)

	LookupHamt64 = hamt64.New(tableOption)
	DeleteHamt64 = hamt64.New(tableOption)

	for _, kv := range genRandomizedKvs(hugeKvs) {
		var str = kv.Key.(*stringkey.StringKey).Str()
		var val = kv.Val.(int)
		LookupMap[str] = val
		DeleteMap[str] = val
		keyStrings[val] = str

		inserted := LookupHamt32.Put(kv.Key, val)
		if !inserted {
			log.Fatalf("failed to LookupHamt32.Put(%s, %v)", kv.Key, kv.Val)
		}

		inserted = DeleteHamt32.Put(kv.Key, val)
		if !inserted {
			log.Fatalf("failed to DeleteHamt32.Put(%s, %v)", kv.Key, kv.Val)
		}

		inserted = LookupHamt64.Put(kv.Key, val)
		if !inserted {
			log.Fatalf("failed to LookupHamt64.Put(%s, %v)", kv.Key, kv.Val)
		}

		inserted = DeleteHamt64.Put(kv.Key, val)
		if !inserted {
			log.Fatalf("failed to DeleteHamt64.Put(%s, %v)", kv.Key, kv.Val)
		}
	}

	RunTime[metricName] = time.Since(StartTime[metricName])
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

func genRandomizedKvs(kvs []key.KeyVal) []key.KeyVal {
	randKvs := make([]key.KeyVal, len(kvs))
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

// This Test is to guarantee that hamt.go and hamt32/hamt.go and hamt64/hamt.go
// definitions of TableOptionName and table options stay in lock step.
func TestTableOptions(t *testing.T) {
	if hamt.HybridTables != hamt32.HybridTables {
		t.Fatal("hamt.HybridTables != hamt32.HybridTables")
	}
	if hamt32.HybridTables != hamt64.HybridTables {
		t.Fatal("hamt32.HybridTables != hamt64.HybridTables")
	}

	if hamt.FullTablesOnly != hamt32.FullTablesOnly {
		t.Fatal("hamt.FullTablesOnly != hamt32.FullTablesOnly")
	}
	if hamt32.FullTablesOnly != hamt64.FullTablesOnly {
		t.Fatal("hamt32.FullTablesOnly != hamt64.FullTablesOnly")
	}

	if hamt.CompTablesOnly != hamt32.CompTablesOnly {
		t.Fatal("hamt.CompTablesOnly != hamt32.CompTablesOnly")
	}
	if hamt32.CompTablesOnly != hamt64.CompTablesOnly {
		t.Fatal("hamt32.CompTablesOnly != hamt64.CompTablesOnly")
	}

	if len(hamt.TableOptionName) != len(hamt32.TableOptionName) {
		t.Fatal("len(hamt.TableOptionName) != len(hamt32.TableOptionName)")
	}
	if len(hamt32.TableOptionName) != len(hamt64.TableOptionName) {
		t.Fatal("len(hamt32.TableOptionName) != len(hamt64.TableOptionName)")
	}
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
