package hamt32_test

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/lleo/go-hamt/hamt32"
	"github.com/lleo/go-hamt/key"
	"github.com/lleo/go-hamt/stringkey"
	"github.com/lleo/stringutil"
	"github.com/pkg/errors"
)

var numHugeKvs int = 5 * 1024 * 1024
var hugeKvs []key.KeyVal

var LookupHamt32 *hamt32.Hamt
var DeleteHamt32 *hamt32.Hamt

var StartTime = make(map[string]time.Time)
var RunTime = make(map[string]time.Duration)

var Inc = stringutil.Lower.Inc

var TableOption int

func TestMain(m *testing.M) {
	// flags
	var fullonly, componly, hybrid, all bool
	flag.BoolVar(&fullonly, "F", false, "Use full tables only and exclude C and H TableOption.")
	flag.BoolVar(&componly, "C", false, "Use compressed tables only and exclude F and H TableOption.")
	flag.BoolVar(&hybrid, "H", false, "Use compressed tables initially and exclude F and C TableOption.")
	flag.BoolVar(&all, "A", false, "Run all Tests w/ TableOption set to hamt32.FullTablesOnly, hamt32.CompTablesOnly, and hamt32.HybridTables; in that order.")

	flag.Parse()

	// If all flag not set, only one of -F, -C, or -H can be set.
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

	// log
	log.SetFlags(log.Lshortfile)

	var logfile, err = os.Create("test.log")
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to os.Create(\"test.log\")"))
	}
	defer logfile.Close()

	log.SetOutput(logfile)

	log.Println("TestMain: and so it begins...")

	hugeKvs = buildKeyVals(numHugeKvs)

	// execute
	var xit int
	if all {
		TableOption = hamt32.FullTablesOnly
		log.Printf("TestMain: TableOption == %s\n", hamt32.TableOptionName[TableOption])
		initialize(TableOption)
		xit = m.Run()
		if xit != 0 {
			os.Exit(1)
		}

		TableOption = hamt32.CompTablesOnly
		log.Printf("TestMain: TableOption == %s\n", hamt32.TableOptionName[TableOption])
		initialize(TableOption)
		xit = m.Run()
		if xit != 0 {
			os.Exit(1)
		}

		TableOption = hamt32.HybridTables
		log.Printf("TestMain: TableOption == %s\n", hamt32.TableOptionName[TableOption])
		initialize(TableOption)
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
		initialize(TableOption)
		xit = m.Run()
	}

	log.Println("\n", RunTimes())
	log.Println("TestMain: the end.")
	os.Exit(xit)
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

func initialize(tableOption int) {
	var funcName = fmt.Sprintf("hamt32: initialize(%s)", hamt32.TableOptionName[tableOption])

	var metricName = fmt.Sprintf("%s: build Lookup/Delete Hamt32", funcName)
	StartTime[metricName] = time.Now()

	LookupHamt32 = hamt32.New(tableOption)
	DeleteHamt32 = hamt32.New(tableOption)

	for _, kv := range genRandomizedKvs(hugeKvs) {
		inserted := LookupHamt32.Put(kv.Key, kv.Val)
		if !inserted {
			log.Fatalf("failed to LookupHamt32.Put(%s, %v)", kv.Key, kv.Val)
		}

		inserted = DeleteHamt32.Put(kv.Key, kv.Val)
		if !inserted {
			log.Fatalf("failed to DeleteHamt32.Put(%s, %v)", kv.Key, kv.Val)
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

func rebuildDeleteHamt32(kvs []key.KeyVal) {
	for _, kv := range kvs {
		inserted := DeleteHamt32.Put(kv.Key, kv.Val)
		if !inserted {
			//log.Printf("BenchmarkHamt32Del: inserted,%v := DeleteHamt32.Put(%s, %d)", inserted, kv.Key, kv.Val)

			// we delete inorder so we can stop rebuilding when the entries start existing
			break
		}
	}
}
