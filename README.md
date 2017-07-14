# Evil RSS reader

I'm frustrated with the RSS readers available, and decided to write my own that
does what I want, silly name and all.

Note that what I personally want may not at all be what you want.

## To run

You need a Go compiler. I'm using Go 1.8.1. Our dependencies are all in the
standard library so far, so you can run the program with:

    $ git clone https://github.com/gunnihinn/evil-rss-reader.git
    $ cd rdr
    $ go run main.go

This will launch an HTTP server on `localhost:8080` that renders the latest
items from the feeds I'm interested in. (At the moment this is one feed that's
hardcoded in the program source code.)

## Todo

These are the initial release goals.

- [x] Fetch feed from hardcoded URL
- [x] Fetch more than one feed from hardcoded URLs
- [x] Render more than one feed
- Guard against feeds including N-thousand word posts in the description.
- Async fetching of feeds
- Don't hardcode feed URLs
- Make the rendering pretty
- Render the N newest items across all feeds separately

Nice-to-haves (?):

- Refactor project structure
- Allow easy swapping of feed rendering styles
- Allow for refreshing of individual feeds
- Keep some history of feed updates

## Evil

See [evilwm](http://www.6809.org.uk/evilwm/).
