#!/usr/bin/env bash

#set -x

function msgExit() {
	xit=$1
	msg=$2

	echo $2
	exit $xit
}

#echo `pwd`
gohamt="$GOPATH/src/github.com/lleo/go-hamt"
cd $gohamt 2>/dev/null
#echo `pwd`

# Am I in the base directory of .../github.com/lleo/go-hamt
[ -f hamt.go ] || msgExit 1  "I don't see a hamt.go; am i in the base directory of $gohamt ?"
[ -d hamt32 ] || msgExit 1  "I don't see a hamt32/; am i in the base directory of $gohamt ?"
[ -d hamt64 ] || msgExit 1  "I don't see a hamt32/; am i in the base directory of $gohamt ?"

echo "Builtin Map implementation"
go test -F -run=xxx -bench=Map | tee map.b
perl -pi -e 's/Map//' map.b

echo "Running Hamt32 benches"
echo "Hamt32 Hybrid/Compressed/Full"
go test -H -run=xxx -bench=Hamt32 | tee hamt32-hybrid.b
go test -C -run=xxx -bench=Hamt32 | tee hamt32-comp.b
go test -F -run=xxx -bench=Hamt32 | tee hamt32-full.b

perl -pi -e 's/Hamt32//' hamt32-hybrid.b
perl -pi -e 's/Hamt32//' hamt32-comp.b
perl -pi -e 's/Hamt32//' hamt32-full.b

echo "Hamt64 Hybrid/Compressed/Full"
go test -H -run=xxx -bench=Hamt64 | tee hamt64-hybrid.b
go test -C -run=xxx -bench=Hamt64 | tee hamt64-comp.b
go test -F -run=xxx -bench=Hamt64 | tee hamt64-full.b

perl -pi -e 's/Hamt64//' hamt64-hybrid.b
perl -pi -e 's/Hamt64//' hamt64-comp.b
perl -pi -e 's/Hamt64//' hamt64-full.b

./summary.sh
