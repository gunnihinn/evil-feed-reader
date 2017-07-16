all: assets
	go build github.com/gunnihinn/evil-rss-reader/cmd/evil-rss-reader

assets: cmd/evil-rss-reader/data
	go-bindata -o cmd/evil-rss-reader/bindata.go cmd/evil-rss-reader/data

clean:
	rm -f cmd/evil-rss-reader/bindata.go
	rm -f evil-rss-reader
