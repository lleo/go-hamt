#!/usr/bin/env bash

#echo "Map verus Hamt32 w/fullTables only"
#benchcmp map.b fullonly.b
#echo
echo "fullTables only versus compressedTables only"
benchcmp fullonly.b componly.b
echo
echo "fullTables only versus hybrid strategy (default)"
benchcmp fullonly.b hybrid.b
echo
echo "compressedTables only versus hybrid strategy (default)"
benchcmp componly.b hybrid.b
