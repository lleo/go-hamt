#!/usr/bin/env bash

function msgExit() {
	xit=$1
	msg=$2

	echo $2
	exit $xit
}

for f in map.b hamt32-hybrid.b hamt32-comp.b hamt32-full.b \
			   hamt64-hybrid.b hamt64-comp.b hamt64-full.b ;
do
	[ -f $f ] || msgExit 1 "I can't find a $f file."
done

# dup but it doesn't matter
perl -pi -e 's/Map//' map.b
perl -pi -e 's/Hamt32//' hamt32-hybrid.b
perl -pi -e 's/Hamt32//' hamt32-comp.b
perl -pi -e 's/Hamt32//' hamt32-full.b
perl -pi -e 's/Hamt64//' hamt64-hybrid.b
perl -pi -e 's/Hamt64//' hamt64-comp.b
perl -pi -e 's/Hamt64//' hamt64-full.b

#echo
echo "Go's Map" vs "Hamt32 FullTablesOnly"
benchcmp map.b hamt32-full.b

echo
echo "Go's Map" vs "Hamt64 FullTablesOnly"
benchcmp map.b hamt64-full.b

echo
echo "Go's Map" vs "Hamt32 CompTablesOnly"
benchcmp map.b hamt32-comp.b

echo
echo "Go's Map" vs "Hamt64 CompTablesOnly"
benchcmp map.b hamt64-comp.b

echo
echo "Go's Map" vs "Hamt32 HybridTables"
benchcmp map.b hamt32-hybrid.b

echo
echo "Go's Map" vs "Hamt64 HybridTables"
benchcmp map.b hamt64-hybrid.b
