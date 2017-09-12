#!/usr/bin/env bash
echo "Copy & Convert hamt64 to hamt32"
set -x

pkg_files="assert.go bitmap.go collision_leaf.go fixed_table.go flat_leaf.go hamt.go hamt_base.go hamt_functional.go hamt_transient.go hashval.go keyval.go node.go sizeof.go sparse_table.go table_iter_stack.go table_stack.go"

bitcount_files="bitcount32.go bitcount32_pre19.go bitcount64.go bitcount64_pre19.go"

test_files="main_test.go hamt64_test.go"

#hamt64_test.go
#key.go

( cd hamt64; cp $pkg_files $bitcount_files $test_files ../hamt32/ )

mv hamt32/hamt64_test.go hamt32/hamt32_test.go
#rm hamt32/bitcount64.go

( cd hamt32;
  perl -pi -e 's/64/32/g' $pkg_files main_test.go hamt32_test.go
  perl -pi -e 's/hamt64/hamt32/' $bitcount_files
)

cp hamt64_test.go hamt32_test.go
perl -pi -e 's/64/32/g' hamt32_test.go
