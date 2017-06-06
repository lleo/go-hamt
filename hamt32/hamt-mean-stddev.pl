#!/usr/bin/env perl

use v5.12;
use Statistics::Basic qw(:all);
use Data::Dumper;

$| = 1; #turn on autoflush
#my @VEC;
my $s = {'32'=>{'transient'=>
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
	say $l;
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
		say(join(" ", $bit, $mode, $ttype, $op, $fields[2]));
		local $Data::Dumper::Indent = 0;
		say Dumper( $s->{$bit}{$mode}{$ttype}{$op}{'data'} );
		next LINE;
	};
	$l =~ m/^BenchmarkHamt32Put-8/ && do {
		my $bit = '32';
		my $op = 'put';
		my @fields = split(' ', $l);
		push @{$s->{$bit}{$mode}{$ttype}{$op}{'data'}}, $fields[2];
		say(join(" ", $bit, $mode, $ttype, $op, $fields[2]));
		local $Data::Dumper::Indent = 0;
		say Dumper( $s->{$bit}{$mode}{$ttype}{$op}{'data'} );
		next LINE;
	};
	$l =~ m/^BenchmarkHamt32Del-8/ && do {
		my $bit = '32';
		my $op = 'del';
		my @fields = split(' ', $l);
		push @{$s->{$bit}{$mode}{$ttype}{$op}{'data'}}, $fields[2];
		say(join(" ", $bit, $mode, $ttype, $op, $fields[2]));
		local $Data::Dumper::Indent = 0;
		say Dumper( $s->{$bit}{$mode}{$ttype}{$op}{'data'} );
		next LINE;
	};
	$l =~ m/^BenchmarkHamt64Get-8/ && do {
		my $bit = '64';
		my $op = 'get';
		my @fields = split(' ', $l);
		push @{$s->{$bit}{$mode}{$ttype}{$op}{'data'}}, $fields[2];
		say(join(" ", $bit, $mode, $ttype, $op, $fields[2]));
		local $Data::Dumper::Indent = 0;
		say Dumper( $s->{$bit}{$mode}{$ttype}{$op}{'data'} );
		next LINE;
	};
	$l =~ m/^BenchmarkHamt64Put-8/ && do {
		my $bit = '64';
		my $op = 'put';
		my @fields = split(' ', $l);
		push @{$s->{$bit}{$mode}{$ttype}{$op}{'data'}}, $fields[2];
		say(join(" ", $bit, $mode, $ttype, $op, $fields[2]));
		local $Data::Dumper::Indent = 0;
		say Dumper( $s->{$bit}{$mode}{$ttype}{$op}{'data'} );
		next LINE;
	};
	$l =~ m/^BenchmarkHamt64Del-8/ && do {
		my $bit = '64';
		my $op = 'del';
		my @fields = split(' ', $l);
		push @{$s->{$bit}{$mode}{$ttype}{$op}{'data'}}, $fields[2];
		say(join(" ", $bit, $mode, $ttype, $op, $fields[2]));
		local $Data::Dumper::Indent = 0;
		say Dumper( $s->{$bit}{$mode}{$ttype}{$op}{'data'} );
		next LINE;
	};
}

#for my $bl (@)
#say "[", join(', ',@VEC), "]";
#my $v = vector(@VEC);
#my $m = mean($v);
#my $d = stddev($v);
#say "$m +/- $d";

# Calculate mean & dev for each data vector
say "### REPORT ###";
for my $bit (keys %$s) {
	for my $mode (keys %{$s->{$bit}}) {
		for my $ttype (keys %{$s->{$bit}{$mode}}) {
			for my $op (keys %{$s->{$bit}{$mode}{$ttype}}) {
				my $entry = $s->{$bit}{$mode}{$ttype}{$op};
				@{$entry->{'data'}} == 0 && next;
				my $v = vector(@{$entry->{'data'}});
				my $m = $entry->{'mean'} = mean($v);
				my $d = $entry->{'stddev'} = stddev($v);
				say "$bit/$mode/$ttype/$op => $m +/- $d ns";
			}
		}
	}
}
