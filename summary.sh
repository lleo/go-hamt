#!/usr/bin/env bash

function msgExit() {
	xit=$1
	msg=$2

	echo $2
	exit $xit
}

for f in map.b fullonly-hamt32.b componly-hamt32.b hybrid-hamt32.b \
			   fullonly-hamt64.b componly-hamt64.b hybrid-hamt64.b ;
do
	[ -f $f ] || msgExit 1 "I can't find a $f file."
done

perl -pi -e 's/Map//' map.b
perl -pi -e 's/Hamt32//' fullonly-hamt32.b
perl -pi -e 's/Hamt32//' componly-hamt32.b
perl -pi -e 's/Hamt32//' hybrid-hamt32.b
perl -pi -e 's/Hamt64//' fullonly-hamt64.b
perl -pi -e 's/Hamt64//' componly-hamt64.b
perl -pi -e 's/Hamt64//' hybrid-hamt64.b

echo "Go's Map VS Hamt32 FullTablesOnly"
benchcmp map.b fullonly-hamt32.b
echo

echo "Go's Map VS Hamt64 FullTablesOnly"
benchcmp map.b fullonly-hamt64.b
echo

echo "Go's Map VS Hamt32 CompTablesOnly"
benchcmp map.b componly-hamt32.b
echo

echo "Go's Map VS Hamt64 CompTablesOnly"
benchcmp map.b componly-hamt64.b
echo

echo "Go's Map VS Hamt32 HybridTables"
benchcmp map.b hybrid-hamt32.b
echo

echo "Go's Map VS Hamt64 HybridTables"
benchcmp map.b hybrid-hamt64.b
echo
