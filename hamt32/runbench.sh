#!/usr/bin/env bash

#go test -F -run=xxx -bench=Map | tee map.b
go test -F -run=xxx -bench=Hamt32 | tee fullonly.b
go test -C -run=xxx -bench=Hamt32 | tee componly.b
go test -H -run=xxx -bench=Hamt32 | tee hybrid.b

#perl -pi -e 's/Map//' map.b
perl -pi -e 's/Hamt32//' fullonly.b
perl -pi -e 's/Hamt32//' componly.b
perl -pi -e 's/Hamt32//' hybrid.b

./summary.sh
