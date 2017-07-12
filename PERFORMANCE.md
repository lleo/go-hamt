# Benchmarking & Performance

From go-hamt/ the basic benchmark command is:

    go test -run=xxx -bench=Hamt(32|64)(Get|Put|Del) -timeout=20m

The -timeout=20m option is required because go test, by default, times out at
ten minutes of total run time. Where running the six Hamt(32|64)(Get|Put|Del)
benchmarks for the six configuation options functional/transient and
FixedTablesOnly/SparseTablesOnly/HybridOnly require a little over 10 minutes.
So I chose to -timeout=20m.

## Performance of []byte instead of Key for v2 version of library

Given the fastest configuation and Benchmark is transient FixedTablesOnly
running the BenchmarkHamt32Get I used the following test:

    go test -F -t -run=xxx -bench=Hamt32Get

Hamt32Get took ~400 ns/op with prepared Stringkey32 keys.
Hamt32Get took ~600 ns/op with prefared []byte slices.

Given that the new []byte API is doing the same work you would have to do by
allocating stringkey32.New(string) Keys, I don't think this is much of a loss.
However, If your app would somehow be able to cache the string->stringkey32
mapping you would be loosing performance with each subsequent call to the
new []byte Public API for Get/Put/Del ops. I think it is worth it for the
simpler API and no dependency on github.com/lleo/go-hamt/string32 library.
Further, Users of this library would not need to implement their own Key type
for other keys that would have had to be converted to []byte slices anyways.
Again simpler API for the Win.

## Performance Notes on Specific Subsystems

### Assert

TLDR; assert() calls with the short-circuit AsserOn constant does not make
any performance difference. I just wish I could have gotten that performance
with an environment variable setting to turn the assert() code on or off. The
performance difference with a var AssertOn rather than a constant is not so
bad, so I could have implemented the environment variable solution, but given
these assert() calls were ment for development ONLY, it does not matter.

First I tested hamt32 isolating just the hamt.Get() path benchmarks for the
fastest type of hamt (that is FixedTablesOnly and Transient). In this path there
is only one '_ = AssertOn && assert()' statement in the 'hv.Index(depth)' call
(hamt32/hamt_transient.go line 168).

I test performance with (in zsh)

    repeat 10 do
    	go test -F -t -run=xxx -bench=Hamt32Get
    done | ../hamt-mean-stddev.pl

If I comment out the assert statement on line 19 of hamt32/hashval.go, meaning
that there is no AssertOn short-circuit and assert() call at all, we get the
following performance result (stored in
hamt32/data/data-Ft-Hamt32Get-nothing.pcf):

    32/transient/fixed/get => 459.8 +/- %1.2 (5.49) ns

If I uncomment out the short-circuit and assert() on line 19 of
hamt32/hashval.go, but set AssertOn to false (line 10 of hamt32/assert.go) we
get the following results (stored in
hamt32/data/data-Ft-Hamt32Get-OFF-assert.pcf):


    32/transient/fixed/get => 467.6 +/- %1.0 (4.57) ns

This benchmark shows that the short-circuit code with AsserOn set to false is
roughly equivelent to no test at all (ie. the short-circuit code becomes a
noop).

If I set AssertOn to true we get (stored in
hamt32/data/data-Ft-Hamt32Get-ON-assert.pcf):

    32/transient/fixed/get => 502.9 +/- %2.2 (11.23) ns

If I put the original non-assert based check in (which is a simple if-test-panic
statement) we get (stored in hamt32/data/data-Ft-if-plain-panic.pcf):

    32/transient/fixed/get => 487.4 +/- %0.8 (3.93) ns

These last two benchmark results show that short-circuit assert() with AssertOn
set to true is roughly equivelent (very roughly) to a straight if-test-panic
statement.

The most significant benchmark result is when AssertOn is set to true and we do
an assertf() (ie. a vararg version of assert()) we get a big difference (stored
in hamt32/data/data-Ft-Hamt32Get-ON-assertf.pcf):

    32/transient/fixed/get => 801.3 +/- %0.6 (5.08) ns

Lastly, setting AssertOn to false still transforms the short-circuit assertf()
to a noop (stored in hamt32/data/data-Ft-Hamt32Get-OFF-assertf.pcf):

    32/transient/fixed/get => 470.7 +/- %1.0 (4.61) ns
