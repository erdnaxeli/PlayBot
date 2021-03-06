#!/usr/bin/perl -w

use strict;

use MIME::Parser;
use URI::Encode;
use JSON;
use File::Temp qw/ tempfile /;

sub get_txt_path {
    my $entity = shift;
    my $num_parts = $entity->parts;

    if ($num_parts) {
        foreach (1..$num_parts) {
            return get_txt_path($entity->parts($_ - 1));
        }
    } else {
        # we only want text
        if ($entity->effective_type =~ /^text\/(?!(html|enriched))/) {
            my $path = $entity->bodyhandle->path;
            return $path;
        }
    }
}

sub get_content {
    my $file = shift;
    my %content;
    my $uri = URI::Encode->new( { encode_reserved => 0 } );
    open (my $fh, "<", $file);

    $content{'msg'} = '';

    while (<$fh>) {
        if (/^(.*) posted in NightIIEs/) {
            $content{'author'} = $1;
        } elsif (/^https:\/\/l.facebook.com\/l\/.*\/(.*)/) {
            $content{'msg'} .= ' '.$uri->decode($1);
        } elsif (/#/) {
            chomp;
            $content{'msg'} .= ' '.$_;
        }
    }

    return \%content;
}

my ($fh, $filename) = tempfile('playbot' . time . '_XXXX', UNLINK => 0, DIR => '/tmp');
while (<STDIN>) {
    print $fh $_;
}

my $parser = new MIME::Parser;
$parser->output_under('/tmp');
my $entity = $parser->parse_open($filename);

my $txt_path = get_txt_path($entity);
my $content = get_content($txt_path);

#remove_files($entity);

my $json = encode_json $content;

pipe(READ,WRITE);
select((select(READ), $| = 1)[0]);

if (my $pid = fork) {
    # parent
    # map READ to STDIN
    open(STDIN, "<&READ");
    exec('ssh -i /usr/local/lib/playbot/key morignot2011@perso.iiens.net ./PlayBot/tools/readjson.pl');
} else {
    # child
    print WRITE "$json\n";
}
