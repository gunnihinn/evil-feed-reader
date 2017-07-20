# Evil feed reader

I'm frustrated with the RSS and atom readers available, and decided to
write my own that does what I want, silly name and all.

Note that what I personally want may not at all be what you want.

## To run

You need a Go compiler. I'm using Go 1.8.1. The reader has one dependency
outside of the standard library (and only on the binary from that project),
which you can install by following the instructions on:

    https://github.com/jteeuwen/go-bindata

Then compile and run the reader with:

    $ git clone https://github.com/gunnihinn/evil-feed-reader.git
    $ make
    $ ./evil-feed-reader

This will launch an HTTP server on `localhost:8080` that renders the latest
items from the feeds you are interested in.

You'll need to setup some feed URLs for the reader to do anything. The reader
will look for a file named `feeds.cfg` in the directory where it is run by
default. You can have it look elsewhere by passing a `-feeds` flag to it when
run. The config file should have the format of one URL per line. Lines
beginning with "#" will be ignored.

## Todo

These are the initial release goals.

- [x] Fetch feed from hardcoded URL
- [x] Fetch more than one feed from hardcoded URLs
- [x] Render more than one feed
- [x] Async fetching of feeds
- [x] Guard against feeds including N-thousand word posts in the description.
- [x] Dynamically generate handlers for each feed
- [x] Don't hardcode feed URLs
- [x] Parse atom feeds (e.g. Scheier on security)
- [x] Rename project to "Evil feed reader"
- [x] Make the rendering pretty
- Render the N newest items across all feeds separately

Nice-to-haves (?):

- [x] Refactor project structure
- Allow easy swapping of feed rendering styles
- Allow for refreshing of individual feeds
- Keep some history of feed updates

## Evil

See [evilwm](http://www.6809.org.uk/evilwm/).
