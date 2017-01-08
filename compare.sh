#!/usr/bin/env bash

echo Comparing FULLONLY hamt32.Hamt versus hamt64.Hamt
benchcmp hamt32/fullonly.b hamt64/fullonly.b
echo

echo Comparing COMPONLY hamt32.Hamt versus hamt64.Hamt
benchcmp hamt32/componly.b hamt64/componly.b
echo

echo Comparing HYBRID hamt32.Hamt versus hamt64.Hamt
benchcmp hamt32/hybrid.b hamt64/hybrid.b

echo FOR EDUCATIONAL USE! here is how the two identical map.b benchmarks compare.
benchcmp hamt32/map.b hamt64/map.b
