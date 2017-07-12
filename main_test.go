package hamt_test

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/lleo/go-hamt"
	"github.com/lleo/go-hamt/hamt32"
	"github.com/lleo/go-hamt/hamt64"
	"github.com/lleo/go-hamt/stringkey32"
	"github.com/lleo/go-hamt/stringkey64"
	"github.com/lleo/stringutil"
	"github.com/pkg/errors"
)

// 1 million & change
var InitHamtNumKvsForPut = 1024 * 1024
var InitHamtNumKvs = InitHamtNumKvsForPut + (2 * 1024 * 1024)
var numKvs = InitHamtNumKvs + (4 * 1024)
var TwoKK = 2 * 1024 * 1024
var KVS32 []hamt32.KeyVal
var KVS64 []hamt64.KeyVal

var Functional bool
var TableOption int

var Hamt32 hamt32.Hamt
var Hamt64 hamt64.Hamt

var Inc = stringutil.Lower.Inc

var StartTime = make(map[string]time.Time)
var RunTime = make(map[string]time.Duration)

func TestMain(m *testing.M) {
	var fixedonly, sparseonly, hybrid, all bool
	flag.BoolVar(&fixedonly, "F", false,
		"Use fixed tables only and exclude C and H Options.")
	flag.BoolVar(&sparseonly, "S", false,
		"Use sparse tables only and exclude F and H Options.")
	flag.BoolVar(&hybrid, "H", false,
		"Use sparse tables initially and exclude F and S Options.")
	flag.BoolVar(&all, "A", false,
		"Run all Tests w/ Options set to FixedTablesOnly, SparseTablesOnly, and HybridTables")

	var functional, transient, both bool
	flag.BoolVar(&functional, "f", false,
		"Run Tests against HamtFunctional struct; excludes transient option")
	flag.BoolVar(&transient, "t", false,
		"Run Tests against HamtFunctional struct; excludes functional option")
	flag.BoolVar(&both, "b", false,
		"Run Tests against both transient and functional Hamt types.")

	flag.Parse()

	// If all flag set, ignore fixedonly, sparseonly, and hybrid.
	if !all {

		// only one flag may be set between fixedonly, sparseonly, and hybrid
		if (fixedonly && (sparseonly || hybrid)) ||
			(sparseonly && (fixedonly || hybrid)) ||
			(hybrid && (sparseonly || fixedonly)) {
			flag.PrintDefaults()
			os.Exit(1)
		}
	}

	// If no flags given, run all tests.
	if !(all || fixedonly || sparseonly || hybrid) {
		all = true
	}

	if !both {
		if functional && transient {
			flag.PrintDefaults()
			os.Exit(1)
		}
	}

	if !(both || functional || transient) {
		both = true
	}

	log.SetFlags(log.Lshortfile)

	var logfile, err = os.Create("test.log")
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to os.Create(\"test.log\")"))
	}
	defer logfile.Close()

	log.SetOutput(logfile)

	log.Println("TestMain: and so it begins...")

	KVS32 = buildKeyVals32("TestMain", numKvs)
	KVS64 = buildKeyVals64("TestMain", numKvs)

	// execute
	var xit int
	if all {
		if both {
			Functional = false
			log.Printf("TestMain: Functional=%t;\n", Functional)
			fmt.Printf("TestMain: Functional=%t;\n", Functional)

			xit = executeAll(m)
			if xit != 0 {
				log.Printf("\n", RunTimes())
				os.Exit(xit)
			}

			Hamt32 = nil
			Hamt64 = nil

			Functional = true
			log.Printf("TestMain: Functional=%t;\n", Functional)
			fmt.Printf("TestMain: Functional=%t;\n", Functional)

			xit = executeAll(m)
		} else if functional {
			Functional = true
			log.Printf("TestMain: Functional=%t;\n", Functional)
			fmt.Printf("TestMain: Functional=%t;\n", Functional)

			xit = executeAll(m)
		} else if transient {
			Functional = false
			log.Printf("TestMain: Functional=%t;\n", Functional)
			fmt.Printf("TestMain: Functional=%t;\n", Functional)

			xit = executeAll(m)
		}
	} else {
		if hybrid {
			TableOption = hamt32.HybridTables
		} else if fixedonly {
			TableOption = hamt32.FixedTablesOnly
		} else /* if sparseonly */ {
			TableOption = hamt32.SparseTablesOnly
		}

		if both {
			Functional = false

			log.Printf("TestMain: Functional=%t;\n", Functional)
			fmt.Printf("TestMain: Functional=%t;\n", Functional)
			log.Printf("TestMain: TableOption=%s;\n",
				hamt32.TableOptionName[TableOption])
			fmt.Printf("TestMain: TableOption=%s;\n",
				hamt32.TableOptionName[TableOption])

			xit = m.Run()
			if xit != 0 {
				log.Printf("\n", RunTimes())
				os.Exit(xit)
			}

			Hamt32 = nil
			Hamt64 = nil
			Functional = true

			log.Printf("TestMain: Functional=%t;\n", Functional)
			fmt.Printf("TestMain: Functional=%t;\n", Functional)
			log.Printf("TestMain: TableOption=%s;\n",
				hamt32.TableOptionName[TableOption])
			fmt.Printf("TestMain: TableOption=%s;\n",
				hamt32.TableOptionName[TableOption])

			xit = m.Run()
		} else {
			if functional {
				Functional = true
			} else /* if transient */ {
				Functional = false
			}

			log.Printf("TestMain: Functional=%t;\n", Functional)
			fmt.Printf("TestMain: Functional=%t;\n", Functional)
			log.Printf("TestMain: TableOption=%s;\n",
				hamt32.TableOptionName[TableOption])
			fmt.Printf("TestMain: TableOption=%s;\n",
				hamt32.TableOptionName[TableOption])
			xit = m.Run()
		}
	}

	log.Println("\n", RunTimes())
	os.Exit(xit)
}

func executeAll(m *testing.M) int {
	TableOption = hamt32.FixedTablesOnly

	log.Printf("TestMain: TableOption=%s;\n",
		hamt32.TableOptionName[TableOption])
	fmt.Printf("TestMain: TableOption=%s;\n",
		hamt32.TableOptionName[TableOption])

	var xit = m.Run()
	if xit != 0 {
		log.Println("\n", RunTimes())
		os.Exit(1)
	}

	Hamt32 = nil
	Hamt64 = nil
	TableOption = hamt32.SparseTablesOnly

	log.Printf("TestMain: TableOption=%s;\n",
		hamt32.TableOptionName[TableOption])
	fmt.Printf("TestMain: TableOption=%s;\n",
		hamt32.TableOptionName[TableOption])

	xit = m.Run()
	if xit != 0 {
		log.Println("\n", RunTimes())
		os.Exit(1)
	}

	Hamt32 = nil
	Hamt64 = nil
	TableOption = hamt32.HybridTables

	log.Printf("TestMain: TableOption=%s;\n",
		hamt32.TableOptionName[TableOption])
	fmt.Printf("TestMain: TableOption=%s;\n",
		hamt32.TableOptionName[TableOption])

	xit = m.Run()

	return xit
}

func buildKeyVals32(prefix string, num int) []hamt32.KeyVal {
	var name = fmt.Sprintf("%s-buildKeyVals64-%d", prefix, num)
	StartTime[name] = time.Now()

	var kvs = make([]hamt32.KeyVal, num)
	var s = "aaa"

	for i := 0; i < num; i++ {
		var k = stringkey32.New(s)

		kvs[i] = hamt32.KeyVal{k, i}
		s = Inc(s)
	}

	RunTime[name] = time.Since(StartTime[name])
	return kvs
}

func buildKeyVals64(prefix string, num int) []hamt64.KeyVal {
	var name = fmt.Sprintf("%s-buildKeyVals64-%d", prefix, num)
	StartTime[name] = time.Now()

	var kvs = make([]hamt64.KeyVal, num)
	var s = "aaa"

	for i := 0; i < num; i++ {
		var k = stringkey64.New(s)

		kvs[i] = hamt64.KeyVal{k, i}
		s = Inc(s)
	}

	RunTime[name] = time.Since(StartTime[name])
	return kvs
}

func buildHamt32(
	prefix string,
	kvs []hamt32.KeyVal,
	functional bool,
	opt int,
) (hamt32.Hamt, error) {
	var name = fmt.Sprintf("%s-buildHamt32-%d", prefix, len(kvs))

	StartTime[name] = time.Now()
	var h = hamt32.New(functional, opt)
	for _, kv := range kvs {
		var k = kv.Key
		var v = kv.Val

		var inserted bool
		h, inserted = h.Put(k, v)
		if !inserted {
			return nil, fmt.Errorf("failed to Put(%s, %v)", k, v)
		}
	}
	RunTime[name] = time.Since(StartTime[name])

	return h, nil
}

func buildHamt64(
	prefix string,
	kvs []hamt64.KeyVal,
	functional bool,
	opt int,
) (hamt64.Hamt, error) {
	var name = fmt.Sprintf("%s-buildHamt64-%d", prefix, len(kvs))

	StartTime[name] = time.Now()
	var h = hamt64.New(functional, opt)
	for _, kv := range kvs {
		var k = kv.Key
		var v = kv.Val

		var inserted bool
		h, inserted = h.Put(k, v)
		if !inserted {
			return nil, fmt.Errorf("failed to Put(%s, %v)", k, v)
		}
	}
	RunTime[name] = time.Since(StartTime[name])

	return h, nil
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

	s += "Key                                                Val\n"
	s += "==================================================+==========\n"

	var tot time.Duration
	for _, k := range ks {
		v := RunTime[k]
		s += fmt.Sprintf("%-50s %s\n", k, v)
		tot += v
	}
	s += fmt.Sprintf("%50s %s\n", "TOTAL", tot)

	return s
}

func TestConstantsInSync(t *testing.T) {
	if hamt.FixedTablesOnly != hamt32.FixedTablesOnly {
		t.Fatal("hamt.FixedTablesOnly != hamt32.FixedTablesOnly")
	}
	if hamt.SparseTablesOnly != hamt32.SparseTablesOnly {
		t.Fatal("hamt.SparseTablesOnly != hamt32.SparseTablesOnly")
	}
	if hamt.HybridTables != hamt32.HybridTables {
		t.Fatal("hamt.HybridTables != hamt32.HybridTables")
	}
	if hamt.TableOptionName != hamt32.TableOptionName {
		t.Fatal("TableOptionName != hamt32.TableOptionName")
	}

	if hamt.FixedTablesOnly != hamt64.FixedTablesOnly {
		t.Fatal("hamt.FixedTablesOnly != hamt64.FixedTablesOnly")
	}
	if hamt.SparseTablesOnly != hamt64.SparseTablesOnly {
		t.Fatal("hamt.SparseTablesOnly != hamt64.SparseTablesOnly")
	}
	if hamt.HybridTables != hamt64.HybridTables {
		t.Fatal("hamt.HybridTables != hamt64.HybridTables")
	}
	if hamt.TableOptionName != hamt64.TableOptionName {
		t.Fatal("TableOptionName != hamt64.TableOptionName")
	}

	// Well... the communative property makes these true BUT...
	// aah dont truck wit nun of dat fancy mathamagical bullshit! Vote Trump!
	if hamt32.FixedTablesOnly != hamt64.FixedTablesOnly {
		t.Fatal("hamt32.FixedTablesOnly != hamt64.FixedTablesOnly")
	}
	if hamt32.SparseTablesOnly != hamt64.SparseTablesOnly {
		t.Fatal("hamt32.SparseTablesOnly != hamt64.SparseTablesOnly")
	}
	if hamt32.HybridTables != hamt64.HybridTables {
		t.Fatal("hamt32.HybridTables != hamt64.HybridTables")
	}
	if hamt32.TableOptionName != hamt64.TableOptionName {
		t.Fatal("hamt32.TableOptionName != hamt64.TableOptionName")
	}
}
