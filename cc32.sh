#!/usr/bin/env bash
echo "Copy & Convert hamt64 to hamt32"
set -x

cp hamt64/*.go hamt32/
rm hamt32/bitcount64.go
perl -pi -e 's/64/32/g' hamt32/*.go
