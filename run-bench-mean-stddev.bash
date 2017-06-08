#!/usr/bin/env bash

time for ((i=0; 10 - $i; i++)) do
    go test -run=xxx -bench="Hamt(32|64)(Get|Put|Del)" -timeout=25m
done | ./hamt-mean-stddev.pl
