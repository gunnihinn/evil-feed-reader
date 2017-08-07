# Evil feed reader

I'm frustrated with the RSS and atom readers available, and decided to
write my own that does what I want, silly name and all.

Note that what I personally want may not at all be what you want.

## Features

Very few!

You setup the URLs of feeds you want to watch in a configuration file. We ping
those every 15 minutes and check if they have been updated. If a feed has
entries you haven't seen yet, it will be colored blue in the main UI.

Click on a feed to show the last few entries from it. If some of them are too
long for us to display comfortably, we'll print a message saying so. You can
then click through to the source's web page to read the entry.

The only history we keep is whether you've yet seen a feed with new entries.

Management is open to feature requests.

## Coming

I'm going to setup a reader with my config on a URL I own so I can read my
stuff from wherever. If you want, I can setup a reader for you too there. The
offer is free (you know, for now, while there are few to no takers).

## To run

Nota bene: If you think this is too much trouble, tell me. I'll figure
something better out.

You need a Go compiler. I'm using Go 1.8.1. The reader has one dependency
outside of the standard library (and only on the binary from that project),
which you can install by following the instructions on:

    https://github.com/jteeuwen/go-bindata

Then compile and run the reader with:

    $ git clone https://github.com/gunnihinn/evil-feed-reader.git
    $ make      # needs GNU make
    $ ./evil-feed-reader

This will launch an HTTP server on `localhost:8080` that renders the latest
items from the feeds you are interested in.

You'll need to setup some feed URLs for the reader to do anything. The reader
will look for a file named `feeds.cfg` in the directory where it is run by
default. You can have it look elsewhere by passing a `-feeds` flag to it when
run. The config file should have the format of one URL per line. Lines
beginning with "#" will be ignored.

## Evil

See [evilwm](http://www.6809.org.uk/evilwm/).
