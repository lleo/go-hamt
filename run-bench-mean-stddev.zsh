#!/usr/bin/env zsh

time repeat 10 do
       go test -run=xxx -bench="Hamt(32|64)(Get|Put|Del)" -timeout=25m
done | ./hamt-mean-stddev.pl
