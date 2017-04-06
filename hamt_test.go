package hamt_test

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
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

// Found depth 5 collision between 3m+2k & 3m+4k at depth=5
// Works for hamt32 testing coverage not for hamt64 coverage
var numKvs int = (3 * 1024 * 1024) + (4 * 1024) // between 3m+2k & 3m+4k

//var SVS []StrVal
var KVS []key.KeyVal

var TableOption int

var LookupHamt32 *hamt32.Hamt

var LookupHamt64 *hamt64.Hamt

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

	KVS = buildKeyVals("global", numKvs)

	// execute
	var xit int
	if all {
		TableOption = hamt32.FullTablesOnly
		log.Printf("TestMain: TableOption == %s\n", hamt32.TableOptionName[TableOption])
		fmt.Printf("Running TableOption = %s\n", hamt32.TableOptionName[TableOption])

		xit = m.Run()
		if xit != 0 {
			os.Exit(1)
		}

		TableOption = hamt32.CompTablesOnly
		log.Printf("TestMain: TableOption == %s\n", hamt32.TableOptionName[TableOption])
		fmt.Printf("Running TableOption = %s\n", hamt32.TableOptionName[TableOption])

		xit = m.Run()
		if xit != 0 {
			os.Exit(1)
		}

		TableOption = hamt32.HybridTables
		log.Printf("TestMain: TableOption == %s\n", hamt32.TableOptionName[TableOption])
		fmt.Printf("Running TableOption = %s\n", hamt32.TableOptionName[TableOption])

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
		fmt.Printf("Running TableOption = %s\n", hamt32.TableOptionName[TableOption])

		xit = m.Run()
	}

	log.Println("\n", RunTimes())

	os.Exit(xit)
}

func buildKeyVals(prefix string, num int) []key.KeyVal {
	var name = fmt.Sprintf("%s-buildKeyVals-%d", prefix, num)
	StartTime[name] = time.Now()

	var kvs = make([]key.KeyVal, num)
	var s = "aaa"

	for i := 0; i < num; i++ {
		var k = stringkey.New(s)

		kvs[i] = key.KeyVal{k, i}
		s = Inc(s)
	}

	RunTime[name] = time.Since(StartTime[name])
	return kvs
}

func buildMap(prefix string, num int) map[string]int {
	var name = fmt.Sprintf("%s-buildMap-%d", prefix, num)
	StartTime[name] = time.Now()

	var m = make(map[string]int, num)
	var s = "aaa"

	for i := 0; i < num; i++ {
		m[s] = i

		s = Inc(s)
	}

	RunTime[name] = time.Since(StartTime[name])
	return m
}

func buildHamt32(prefix string, kvs []key.KeyVal, opt int) (*hamt32.Hamt, error) {
	var name = fmt.Sprintf("%s-buildHamt32-%d", prefix, len(kvs))
	StartTime[name] = time.Now()

	var h = hamt32.New(opt)

	for _, kv := range kvs {
		var k = kv.Key
		var v = kv.Val

		var inserted = h.Put(k, v)
		if !inserted {
			return nil, fmt.Errorf("failed to Put(%s, %v)", k, v)
		}
	}

	RunTime[name] = time.Since(StartTime[name])
	return h, nil
}

func buildHamt64(prefix string, kvs []key.KeyVal, opt int) (*hamt64.Hamt, error) {
	var name = fmt.Sprintf("%s-buildHamt64-%d", prefix, len(kvs))
	StartTime[name] = time.Now()

	var h = hamt64.New(opt)

	for _, kv := range kvs {
		var k = kv.Key
		var v = kv.Val

		var inserted = h.Put(k, v)
		if !inserted {
			return nil, fmt.Errorf("failed to Put(%s, %v)", k, v)
		}
	}

	RunTime[name] = time.Since(StartTime[name])
	return h, nil
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
	// Grab list of keys from RunTime map; MAJOR un-feature of Go!
	var ks = make([]string, len(RunTime))
	var i int = 0
	for k := range RunTime {
		ks[i] = k
		i++
	}
	sort.Strings(ks)

	var s = ""

	s += "Key                                      Val\n"
	s += "========================================+==========\n"

	for _, k := range ks {
		v := RunTime[k]
		s += fmt.Sprintf("%-40s %s\n", k, v)
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

//stupid on many levels
func TestNewHamt32(t *testing.T) {
	var h = hamt.NewHamt32()
	if !h.IsEmpty() {
		t.Fatal("!?!? a brand new Hamt !IsEmpty()")
	}
}

//stupid on many levels
func TestNewHamt64(t *testing.T) {
	var h = hamt.NewHamt32()
	if !h.IsEmpty() {
		t.Fatal("!?!? a brand new Hamt !IsEmpty()")
	}
}
