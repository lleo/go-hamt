repeat 10 do
	go test -run=xxx -bench=Hamt32(Get|Put|Del) -timeout 20m
done | ./hamt-mean-stddev.pl -d data/All-IndexBits=5.pcf
