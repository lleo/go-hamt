#!/usr/bin/env bash

OPTIND=1 # Reset in case getopts has been used previously in the shell.

function usage() {
    cat <<- EOU
Usage: $0 [-d <data_fn>] [-b <bench string>] [-t <test_args>]
Example: $0 -d data.pcf -b "Hamt32(Get|Put|Del)" -t "-F -functional"

<data_fn> is the name where the hamt-mean-stddev.pl program will put its
perl datastructure. ex ./hamt-mean-stddev.pl -d <data_fn>

<bench_str> is the argument that will be passed to go test's -bench argument.
ex. go test -run=xxx -bench="$bench_str" -timeout=25m

<test_args> is any number of valid arguments to go test.
ex. go test $test_args -run=xxx -bench="$bench_str"
EOU
}


data_arg=
bench_str="Hamt(32|64)(Get|Put|Del)"
test_args=
while getopts ":d:b:t:h" opt; do
    case "$opt" in
	d)
	    data_arg="-d=$OPTARG"
	    ;;
	b)
	    bench_str="$OPTARG"
	    ;;
	t)
	    test_args="$OPTARG"
	    ;;
	h)
	    usage
	    exit 0
	    ;;
	\?)
	    echo "unknown option $OPTARG" 1>&2
	    usage 1>&2
	    exit 1
	    ;;
    esac
done

echo "test_args=$test_args"

time for ((i=0; 10 - $i; i++)) do
    go test $test_args -run=xxx -bench="$bench_str" -timeout=25m
done | ./hamt-mean-stddev.pl $data_arg
