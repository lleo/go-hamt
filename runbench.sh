#!/usr/bin/env bash

set -x

function msgExit() {
	xit=$1
	msg=$2

	echo $2
	exit $xit
}

echo `pwd`
gohamt="$GOPATH/src/github.com/lleo/go-hamt"
cd $gohamt 2>/dev/null
echo `pwd`

# Am I in the base directory of .../github.com/lleo/go-hamt
[ -f hamt.go ] || msgExit 1  "I don't see a hamt.go; am i in the base directory of $gohamt ?"
[ -d hamt32 ] || msgExit 1  "I don't see a hamt32/; am i in the base directory of $gohamt ?"
[ -d hamt64 ] || msgExit 1  "I don't see a hamt32/; am i in the base directory of $gohamt ?"

echo "Running map benches"
go test -run=xxx -bench=Map | tee map.b

cd hamt32
echo "Running Hamt32 benches"
go test -F -run=xxx -bench=Hamt32 | tee ../fullonly-hamt32.b
go test -C -run=xxx -bench=Hamt32 | tee ../componly-hamt32.b
go test -H -run=xxx -bench=Hamt32 | tee ../hybrid-hamt32.b

cd ../hamt64

echo "Running Hamt64 benches"
go test -F -run=xxx -bench=Hamt64 | tee ../fullonly-hamt64.b
go test -C -run=xxx -bench=Hamt64 | tee ../componly-hamt64.b
go test -H -run=xxx -bench=Hamt64 | tee ../hybrid-hamt64.b

cd ..

./summary.sh
