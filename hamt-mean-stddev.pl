#!/usr/bin/env perl

use v5.12;
use Statistics::Basic qw(:all);
use Data::Dumper;
use Getopt::Long;
use File::Basename qw(basename dirname);
use File::Spec::Functions qw(catfile);

$| = 1; #turn on autoflush

sub UsageAndExit {
	my ($xit, $fh, @msgs) = @_;
	$xit = defined $xit ? $xit : 0;
	$fh = defined $fh ? $fh : $xit == 0 ? \*STDOUT : \*STDERR;
	$fh->print($_) for @msgs;
	$fh->print(<<EOU);
Usage $0 [-h] [-d|--data_fn <data filename>]
              [-i|--input_fn <input data filename>]
ex. $0 -i data.pcf
ex. $0 -i=data.pcf
EOU
	exit($xit)
}

my $data_fn = "data.pcf";
my $input_fn = "-";
my $help;
GetOptions("h|help" => \$help,
           "d|data_fn=s" => \$data_fn,
           "i|input_fn=s" => \$input_fn);

$help && UsageAndExit(0);

#say "perl: data_fn=>$data_fn<";
#say "perl: input_fn=>$input_fn<";

my $s;
if ($input_fn eq '-') {
	$s = {'32'=>{'transient'=>
	             {'fixed' =>{'get'=>{'data'=>[]},
	                         'put'=>{'data'=>[]},
	                         'del'=>{'data'=>[]}},
	              'sparse'=>{'get'=>{'data'=>[]},
	                         'put'=>{'data'=>[]},
	                         'del'=>{'data'=>[]}},
	              'hybrid'=>{'get'=>{'data'=>[]},
	                         'put'=>{'data'=>[]},
	                         'del'=>{'data'=>[]}}
	             },
	             'functional'=>
	             {'fixed' =>{'get'=>{'data'=>[]},
	                         'put'=>{'data'=>[]},
	                         'del'=>{'data'=>[]}},
	              'sparse'=>{'get'=>{'data'=>[]},
	                         'put'=>{'data'=>[]},
	                         'del'=>{'data'=>[]}},
	              'hybrid'=>{'get'=>{'data'=>[]},
	                         'put'=>{'data'=>[]},
	                         'del'=>{'data'=>[]}}
	             }
	            },
	      '64'=>{'transient'=>
	             {'fixed' =>{'get'=>{'data'=>[]},
	                         'put'=>{'data'=>[]},
	                         'del'=>{'data'=>[]}},
	              'sparse'=>{'get'=>{'data'=>[]},
	                         'put'=>{'data'=>[]},
	                         'del'=>{'data'=>[]}},
	              'hybrid'=>{'get'=>{'data'=>[]},
	                         'put'=>{'data'=>[]},
	                         'del'=>{'data'=>[]}}
	             },
	             'functional'=>
	             {'fixed' =>{'get'=>{'data'=>[]},
	                         'put'=>{'data'=>[]},
	                         'del'=>{'data'=>[]}},
	              'sparse'=>{'get'=>{'data'=>[]},
	                         'put'=>{'data'=>[]},
	                         'del'=>{'data'=>[]}},
	              'hybrid'=>{'get'=>{'data'=>[]},
	                         'put'=>{'data'=>[]},
	                         'del'=>{'data'=>[]}}
	             }
	            }
	     };

	#say "reading from stdin";

	my $mode;
	my $ttype;
	LINE: while (my $l = <>) {
		chomp $l;
		#say STDERR $l;
		# Set $mode
		$l =~ m/Functional=true/ && do {
			$mode = 'functional';
			next LINE;
		};
		$l =~ m/Functional=false/ && do {
			$mode = 'transient';
			next LINE;
		};
		# Set the current $ttype
		#$l =~ m/TableOption=(Fixed|Sparse|Hybrid)/ && do {
		$l =~ m/TableOption=(\p{Alpha}+)T/ && do {
			$ttype = lc($1); # 'fixed' | 'sparse' | 'hybrid'
			#if (length($ttype) < 6) {
			#	$ttype .= ' 'x(6-length($ttype))
			#}
			next LINE;
		};
		# record the benchmark result
		$l =~ m/^BenchmarkHamt(32|64)(Get|Put|Del)-/ && do {
			my $bit = $1; # '32' | '64'
			my $op = lc($2); # 'get' | 'put' | 'del'
			my @fields = split(' ', $l);
			push @{$s->{$bit}{$mode}{$ttype}{$op}{'data'}}, $fields[2];
			my $n = scalar(@{$s->{$bit}{$mode}{$ttype}{$op}{'data'}});
			say STDERR "#$n: ", join(" ", $bit, $mode, $ttype, $op, $fields[2]);
			#local $Data::Dumper::Indent = 0;
			#say STDERR Dumper( $s->{$bit}{$mode}{$ttype}{$op}{'data'} );
			next LINE;
		};
	}

	if (defined $data_fn) {
		open(my $fh, ">", $data_fn)
		  or die "Can't open > output.txt: $!";
		$fh->print(Dumper($s));
		$fh->close();
	}
} else {
	$s = do $input_fn;

	# fixup old $ttype names
	my ($need_write);
	for my $bit (keys %$s) {
		for my $mode (keys %{$s->{$bit}}) {
			if (ref($s->{$bit}{$mode}{'fixed '}) eq 'HASH') {
				$s->{$bit}{$mode}{'fixed'} = delete $s->{$bit}{$mode}{'fixed '};
				$need_write = "true";
			}
			if (ref($s->{$bit}{$mode}{'spars'}) eq 'HASH') {
				$s->{$bit}{$mode}{'sparse'} = delete $s->{$bit}{$mode}{'spars'};
				$need_write = "true";
			}
			if (ref($s->{$bit}{$mode}{'hybri'}) eq 'HASH') {
				$s->{$bit}{$mode}{'hybrid'} = delete $s->{$bit}{$mode}{'hybri'};
				$need_write = "true";
			}
		}
	}
	for my $bit (keys %$s) {
		for my $mode (keys %{$s->{$bit}}) {
			for my $ttype (keys %{$s->{$bit}{$mode}}) {
				for my $op (keys %{$s->{$bit}{$mode}{$ttype}}) {
					if (ref($s->{$bit}{$mode}{$ttype}{$op}{'mean'})) {
						delete $s->{$bit}{$mode}{$ttype}{$op}{'mean'};
						$need_write = "true"
					}
					if (ref $s->{$bit}{$mode}{$ttype}{$op}{'stdev'}) {
						delete $s->{$bit}{$mode}{$ttype}{$op}{'stdev'};
						$need_write = "true"
					}
				}
			}
		}
	}

	if ($need_write) {
		my $fixedup_input_fn = catfile(dirname($input_fn),
		                               "fixedup-" . basename($input_fn));
		open(my $fh, ">", $fixedup_input_fn)
		  or die "Can't open > $fixedup_input_fn: $!";
		$fh->print(Dumper($s));
		$fh->close();
		say "WROTE $fixedup_input_fn";
	}
}


# Calculate and Display mean & stddev for each data vector
say "### REPORT ###";
#for my $bit ('32', '64') {
for my $bit (sort keys %$s) {
	for my $mode (sort keys %{$s->{$bit}}) {
		#for my $ttype (sort keys %{$s->{$bit}{$mode}}) {
		for my $ttype ('fixed', 'sparse', 'hybrid') {
			#for my $op (sort keys %{$s->{$bit}{$mode}{$ttype}}) {
			for my $op ('get', 'put', 'del') {
				my $entry = $s->{$bit}{$mode}{$ttype}{$op};
				@{$entry->{'data'}} == 0 && next;
				my $v = vector(@{$entry->{'data'}});
				my $m = mean($v);
				my $d = stddev($v);
				my $pct = ($d / $m) * 100;
				my $pct_s = sprintf("%%%-4.1f", $pct);
				#say "$bit $mode $ttype $op => $m +/- $pct_s ($d) ns";
				printf("%s %- 10s %- 6s %s => % 8.1f +/- %%%5.2f (%8.2f) ns\n",
				       $bit, $mode, $ttype, $op, $m, $pct, $d);
			}
		}
	}
}
