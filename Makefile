all: assets
	go build github.com/gunnihinn/evil-feed-reader/cmd/evil-feed-reader

assets: cmd/evil-feed-reader/data
	go-bindata -o cmd/evil-feed-reader/bindata.go cmd/evil-feed-reader/data

devel: devel-assets
	go build github.com/gunnihinn/evil-feed-reader/cmd/evil-feed-reader

devel-assets: cmd/evil-feed-reader/data
	go-bindata -debug -o cmd/evil-feed-reader/bindata.go cmd/evil-feed-reader/data

clean:
	rm -f cmd/evil-feed-reader/bindata.go
	rm -f evil-feed-reader
