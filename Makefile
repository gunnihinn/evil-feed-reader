binary=evilfr

lib_src = $(wildcard */*.go)
bin_src = $(wildcard cmd/evil-feed-reader/*.go)
project_package = github.com/gunnihinn/evil-feed-reader/cmd/evil-feed-reader

datadir = cmd/evil-feed-reader/static
data = $(wildcard cmd/evil-feed-reader/static/*)
bindata = cmd/evil-feed-reader/bindata.go
bindata_gen = bindata_assetfs.go

$(binary): $(bindata) $(lib_src) $(bin_src)
	go build -o $(binary) $(project_package)

$(bindata): $(data)
	go-bindata-assetfs -prefix "cmd/evil-feed-reader/" $(datadir)/...
	mv $(bindata_gen) $(bindata)

clean:
	rm -f $(bindata)
	rm -f $(binary)

deploy: $(bindata) $(lib_src) $(bin_src)
	GOOS=freebsd GOARCH=amd64 go build -o $(binary) $(project_package)
	./deploy.sh > /dev/null
