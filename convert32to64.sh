#!/usr/bin/env bash

for f in "$@"; do
	echo $f

	echo "perl -pi -e 's/32/64/g' $f"
	perl -pi -e 's/32/64/g' $f

	echo "perl -pi -e 's/30/60/g' $f"
	perl -pi -e 's/30/60/g' $f

	echo "perl -pi -e 's/\bsix/ten/g'"
	perl -pi -e 's/\bsix/ten/g' $f

	echo "perl -pi -e 's/\b5bit/6bit/g'"
	perl -pi -e 's/\b5bit/6bit/g' $f

	echo "perl -pi -e 's/\b6\b/9/g'"
	perl -pi -e 's/\b6\b/9/g' $f

	echo "perl -pi -e 's/\b5\b/6/g'"
	perl -pi -e 's/\b5\b/6/g' $f

done
