#!/usr/bin/env bash
echo "Copy & Convert hamt64 to hamt32"
set -x

grep '^package hamt$' hamt.go >/dev/null
if [ $? -ne 0 ]; then
    echo "$0 must be run from base directory of go-hamt package" 2>&1
    exit 1
fi

pkg_files="assert.go collision_leaf.go fixed_table.go flat_leaf.go hamt.go hamt_base.go hamt_functional.go hamt_transient.go hashval.go keyval.go node.go sizeof.go sparse_table.go table_iter_stack.go table_stack.go"

specific_files="bitmap.go key_types.go bitcount32.go bitcount32_pre19.go bitcount64.go bitcount64_pre19.go"

test_files="main_test.go hamt64_test.go"

#hamt64_test.go
#key.go

cd hamt64
cp $pkg_files $specific_files $test_files ../hamt32/
if [ ! -d ../hamt32/key ]; then
    mkdir ../hamt32/key
fi
cp key/*.go ../hamt32/key/
cd ..

mv hamt32/hamt64_test.go hamt32/hamt32_test.go
#rm hamt32/bitcount64.go

cd hamt32
perl -pi -e 's/hamt64/hamt32/g' $pkg_files main_test.go hamt32_test.go
perl -pi -d 's/uint64/uint32/g' $pkg_files main_test.go hamt32_test.go
perl -pi -e 's/hamt64/hamt32/' $specific_files
perl -pi -e 's/64/32/g' key/*.go
cd ..

cp hamt64_test.go hamt32_test.go
perl -pi -e 's/64/32/g' hamt32_test.go
