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
	"github.com/lleo/go-hamt/key"
	"github.com/lleo/go-hamt/stringkey"
	"github.com/lleo/stringutil"
	"github.com/pkg/errors"
)

// 1 million & change
var InitHamtNumKvsForPut = 1024 * 1024 //1000000

// 4 million & change
var InitHamtNumKvs = InitHamtNumKvsForPut + (3 * 1024 * 1024) //+ 3000000

//var numKvs = (4 * 1024 * 1024) + (4 * 1024)
var KVS []key.KeyVal

var Functional bool
var TableOption int

var Hamt32 hamt32.Hamt
var Hamt64 hamt64.Hamt

var Inc = stringutil.Lower.Inc

var StartTime = make(map[string]time.Time)
var RunTime = make(map[string]time.Duration)

func TestMain(m *testing.M) {
	var fullonly, componly, hybrid, all bool
	flag.BoolVar(&fullonly, "F", false,
		"Use full tables only and exclude C and H Options.")
	flag.BoolVar(&componly, "C", false,
		"Use compressed tables only and exclude F and H Options.")
	flag.BoolVar(&hybrid, "H", false,
		"Use compressed tables initially and exclude F and C Options.")
	flag.BoolVar(&all, "A", false,
		"Run all Tests w/ Options set to FullTablesOnly, CompTablesOnly, and HybridTables")

	var functional, transient, both bool
	flag.BoolVar(&functional, "f", false,
		"Run Tests against HamtFunctional struct; excludes transient option")
	flag.BoolVar(&transient, "t", false,
		"Run Tests against HamtFunctional struct; excludes functional option")
	flag.BoolVar(&both, "b", false,
		"Run Tests against both transient and functional Hamt types.")

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

	KVS = buildKeyVals("TestMain", InitHamtNumKvs)

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
		} else if fullonly {
			TableOption = hamt32.FullTablesOnly
		} else /* if componly */ {
			TableOption = hamt32.CompTablesOnly
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
	TableOption = hamt32.FullTablesOnly

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
	TableOption = hamt32.CompTablesOnly

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

	for _, k := range ks {
		v := RunTime[k]
		s += fmt.Sprintf("%-50s %s\n", k, v)
	}
	return s
}

func TestConstantsInSync(t *testing.T) {
	if hamt.FullTablesOnly != hamt32.FullTablesOnly {
		t.Fatal("hamt.FullTablesOnly != hamt32.FullTablesOnly")
	}
	if hamt.CompTablesOnly != hamt32.CompTablesOnly {
		t.Fatal("hamt.CompTablesOnly != hamt32.CompTablesOnly")
	}
	if hamt.HybridTables != hamt32.HybridTables {
		t.Fatal("hamt.HybridTables != hamt32.HybridTables")
	}
	if hamt.TableOptionName != hamt32.TableOptionName {
		t.Fatal("TableOptionName != hamt32.TableOptionName")
	}

	if hamt.FullTablesOnly != hamt64.FullTablesOnly {
		t.Fatal("hamt.FullTablesOnly != hamt64.FullTablesOnly")
	}
	if hamt.CompTablesOnly != hamt64.CompTablesOnly {
		t.Fatal("hamt.CompTablesOnly != hamt64.CompTablesOnly")
	}
	if hamt.HybridTables != hamt64.HybridTables {
		t.Fatal("hamt.HybridTables != hamt64.HybridTables")
	}
	if hamt.TableOptionName != hamt64.TableOptionName {
		t.Fatal("TableOptionName != hamt64.TableOptionName")
	}

	// Well... the communative property makes these true BUT...
	// aah dont truck wit nun of dat fancy mathamagical bullshit! Vote Trump!
	if hamt32.FullTablesOnly != hamt64.FullTablesOnly {
		t.Fatal("hamt32.FullTablesOnly != hamt64.FullTablesOnly")
	}
	if hamt32.CompTablesOnly != hamt64.CompTablesOnly {
		t.Fatal("hamt32.CompTablesOnly != hamt64.CompTablesOnly")
	}
	if hamt32.HybridTables != hamt64.HybridTables {
		t.Fatal("hamt32.HybridTables != hamt64.HybridTables")
	}
	if hamt32.TableOptionName != hamt64.TableOptionName {
		t.Fatal("hamt32.TableOptionName != hamt64.TableOptionName")
	}
}
