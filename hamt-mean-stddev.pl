#!/usr/bin/env perl

use v5.12;
use Statistics::Basic qw(:all);
use Data::Dumper;
use Getopt::Long;

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

say "perl: data_fn=>$data_fn<";
say "perl: input_fn=>$input_fn<";

my $s;
if ($input_fn == '-') {
	$s = {'32'=>{'transient'=>
	             {'full'=>{'get'=>{'data'=>[],'mean'=>0,'stddev'=>0},
	                       'put'=>{'data'=>[],'mean'=>0,'stddev'=>0},
	                       'del'=>{'data'=>[],'mean'=>0,'stddev'=>0}},
	              'comp'=>{'get'=>{'data'=>[],'mean'=>0,'stddev'=>0},
	                       'put'=>{'data'=>[],'mean'=>0,'stddev'=>0},
	                       'del'=>{'data'=>[],'mean'=>0,'stddev'=>0}},
	              'hybr'=>{'get'=>{'data'=>[],'mean'=>0,'stddev'=>0},
	                       'put'=>{'data'=>[],'mean'=>0,'stddev'=>0},
	                       'del'=>{'data'=>[],'mean'=>0,'stddev'=>0}}
	             },
	             'functional'=>
	             {'full'=>{'get'=>{'data'=>[],'mean'=>0,'stddev'=>0},
	                       'put'=>{'data'=>[],'mean'=>0,'stddev'=>0},
	                       'del'=>{'data'=>[],'mean'=>0,'stddev'=>0}},
	              'comp'=>{'get'=>{'data'=>[],'mean'=>0,'stddev'=>0},
	                       'put'=>{'data'=>[],'mean'=>0,'stddev'=>0},
	                       'del'=>{'data'=>[],'mean'=>0,'stddev'=>0}},
	              'hybr'=>{'get'=>{'data'=>[],'mean'=>0,'stddev'=>0},
	                       'put'=>{'data'=>[],'mean'=>0,'stddev'=>0},
	                       'del'=>{'data'=>[],'mean'=>0,'stddev'=>0}}
	             }
	            },
	      '64'=>{'transient'=>
	             {'full'=>{'get'=>{'data'=>[],'mean'=>0,'stddev'=>0},
	                       'put'=>{'data'=>[],'mean'=>0,'stddev'=>0},
	                       'del'=>{'data'=>[],'mean'=>0,'stddev'=>0}},
	              'comp'=>{'get'=>{'data'=>[],'mean'=>0,'stddev'=>0},
	                       'put'=>{'data'=>[],'mean'=>0,'stddev'=>0},
	                       'del'=>{'data'=>[],'mean'=>0,'stddev'=>0}},
	              'hybr'=>{'get'=>{'data'=>[],'mean'=>0,'stddev'=>0},
	                       'put'=>{'data'=>[],'mean'=>0,'stddev'=>0},
	                       'del'=>{'data'=>[],'mean'=>0,'stddev'=>0}}
	             },
	             'functional'=>
	             {'full'=>{'get'=>{'data'=>[],'mean'=>0,'stddev'=>0},
	                       'put'=>{'data'=>[],'mean'=>0,'stddev'=>0},
	                       'del'=>{'data'=>[],'mean'=>0,'stddev'=>0}},
	              'comp'=>{'get'=>{'data'=>[],'mean'=>0,'stddev'=>0},
	                       'put'=>{'data'=>[],'mean'=>0,'stddev'=>0},
	                       'del'=>{'data'=>[],'mean'=>0,'stddev'=>0}},
	              'hybr'=>{'get'=>{'data'=>[],'mean'=>0,'stddev'=>0},
	                       'put'=>{'data'=>[],'mean'=>0,'stddev'=>0},
	                       'del'=>{'data'=>[],'mean'=>0,'stddev'=>0}}
	             }
	            }
	     };

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
		$l =~ m/TableOption=Full/ && do {
			$ttype = 'full';
			next LINE;
		};
		$l =~ m/TableOption=Comp/ && do {
			$ttype = 'comp';
			next LINE;
		};
		$l =~ m/TableOption=Hybr/ && do {
			$ttype = 'hybr';
			next LINE;
		};
		# record the benchmark result
		$l =~ m/^BenchmarkHamt32Get-8/ && do {
			my $bit = '32';
			my $op = 'get';
			my @fields = split(' ', $l);
			push @{$s->{$bit}{$mode}{$ttype}{$op}{'data'}}, $fields[2];
			my $n = scalar(@{$s->{$bit}{$mode}{$ttype}{$op}{'data'}});
			say STDERR "#$n: ", join(" ", $bit, $mode, $ttype, $op, $fields[2]);
			#local $Data::Dumper::Indent = 0;
			#say STDERR Dumper( $s->{$bit}{$mode}{$ttype}{$op}{'data'} );
			next LINE;
		};
		$l =~ m/^BenchmarkHamt32Put-8/ && do {
			my $bit = '32';
			my $op = 'put';
			my @fields = split(' ', $l);
			push @{$s->{$bit}{$mode}{$ttype}{$op}{'data'}}, $fields[2];
			my $n = scalar(@{$s->{$bit}{$mode}{$ttype}{$op}{'data'}});
			say STDERR "#$n: ", join(" ", $bit, $mode, $ttype, $op, $fields[2]);
			#local $Data::Dumper::Indent = 0;
			#say STDERR Dumper( $s->{$bit}{$mode}{$ttype}{$op}{'data'} );
			next LINE;
		};
		$l =~ m/^BenchmarkHamt32Del-8/ && do {
			my $bit = '32';
			my $op = 'del';
			my @fields = split(' ', $l);
			push @{$s->{$bit}{$mode}{$ttype}{$op}{'data'}}, $fields[2];
			my $n = scalar(@{$s->{$bit}{$mode}{$ttype}{$op}{'data'}});
			say STDERR "#$n: ", join(" ", $bit, $mode, $ttype, $op, $fields[2]);
			#local $Data::Dumper::Indent = 0;
			#say STDERR Dumper( $s->{$bit}{$mode}{$ttype}{$op}{'data'} );
			next LINE;
		};
		$l =~ m/^BenchmarkHamt64Get-8/ && do {
			my $bit = '64';
			my $op = 'get';
			my @fields = split(' ', $l);
			push @{$s->{$bit}{$mode}{$ttype}{$op}{'data'}}, $fields[2];
			my $n = scalar(@{$s->{$bit}{$mode}{$ttype}{$op}{'data'}});
			say STDERR "#$n: ", join(" ", $bit, $mode, $ttype, $op, $fields[2]);
			#local $Data::Dumper::Indent = 0;
			#say STDERR Dumper( $s->{$bit}{$mode}{$ttype}{$op}{'data'} );
			next LINE;
		};
		$l =~ m/^BenchmarkHamt64Put-8/ && do {
			my $bit = '64';
			my $op = 'put';
			my @fields = split(' ', $l);
			push @{$s->{$bit}{$mode}{$ttype}{$op}{'data'}}, $fields[2];
			my $n = scalar(@{$s->{$bit}{$mode}{$ttype}{$op}{'data'}});
			say STDERR "#$n: ", join(" ", $bit, $mode, $ttype, $op, $fields[2]);
			#local $Data::Dumper::Indent = 0;
			#say STDERR Dumper( $s->{$bit}{$mode}{$ttype}{$op}{'data'} );
			next LINE;
		};
		$l =~ m/^BenchmarkHamt64Del-8/ && do {
			my $bit = '64';
			my $op = 'del';
			my @fields = split(' ', $l);
			push @{$s->{$bit}{$mode}{$ttype}{$op}{'data'}}, $fields[2];
			my $n = scalar(@{$s->{$bit}{$mode}{$ttype}{$op}{'data'}});
			say STDERR "#$n: ", join(" ", $bit, $mode, $ttype, $op, $fields[2]);
			#local $Data::Dumper::Indent = 0;
			#say STDERR Dumper( $s->{$bit}{$mode}{$ttype}{$op}{'data'} );
			next LINE;
		};
	}
} else {
	$s = do $input_fn
}


if (defined $data_fn) {
	open(my $fh, ">", $data_fn)
	  or die "Can't open > output.txt: $!";
	$fh->print(Dumper($s));
	$fh->close();
}

# Calculate mean & dev for each data vector
say "### REPORT ###";
#for my $bit (keys %$s) {
for my $bit ('32', '64') {
	#for my $mode (keys %{$s->{$bit}}) {
	for my $mode ('functional', 'transient') {
		#for my $ttype (keys %{$s->{$bit}{$mode}}) {
		for my $ttype ('full', 'comp', 'hybr') {
			#for my $op (keys %{$s->{$bit}{$mode}{$ttype}}) {
			for my $op ('get', 'put', 'del') {
				my $entry = $s->{$bit}{$mode}{$ttype}{$op};
				@{$entry->{'data'}} == 0 && next;
				my $v = vector(@{$entry->{'data'}});
				my $m = $entry->{'mean'} = mean($v);
				my $d = $entry->{'stddev'} = stddev($v);
				my $pct = $d / $m;
				my $pct_s = sprintf("%%%.1f", $pct);
				say "$bit/$mode/$ttype/$op => $m +/- $pct_s ($d) ns";
			}
		}
	}
}
