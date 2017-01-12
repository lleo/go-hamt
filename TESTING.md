TESTING
=======

For now you test each implementation seperately. There are two forms of testing.
Testing all three table configurations; HybridTables, FullTablesOnly, and
CompTablesOnly (for compressed tables only). Or Testing each table config
individually.

To test all three configurations try:

    go-hamt/ $ cd hamt32
    go-hamt/hamt32 $ go test
	
or

    go-hamt/hamt32 $ go test -A

To test individual configuration types try:

    go-hamt/hamt32 $ go test -H  #for HybridTables
    go-hamt/hamt32 $ go test -C  #for CompTablesOnly
    go-hamt/hamt32 $ go test -H  #for FullTablesOnly
	
   
You can run benchmarks on the individual strategies like:

   go-hamt/hamt32 $ go test -C -run=xx -bench=.
   go-hamt/hamt32 $ go test -F -run=xx -bench=.
   go-hamt/hamt32 $ go test -H -run=xx -bench=.

To get the fullonly.b, componly.b, and hybid.b and a benchmark comparison
use the given script `runbench.sh`.

You can see the summary of the benchmark comparison of the `*.b` files
with `summary.sh`.
