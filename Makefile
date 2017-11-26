binary=evilfr

lib_src = $(wildcard */*.go)
bin_src = $(wildcard *.go)
project_package = github.com/gunnihinn/evil-feed-reader

datadir = static
data = $(wildcard static/*)
bindata = bindata.go
bindata_gen = bindata_assetfs.go

$(binary): $(bindata) $(lib_src) $(bin_src)
	go build -o $(binary) $(project_package)

$(bindata): $(data)
	go-bindata-assetfs $(datadir)/...
	mv $(bindata_gen) $(bindata)

clean:
	rm -f $(bindata)
	rm -f $(binary)

deploy: $(bindata) $(lib_src) $(bin_src)
	GOOS=linux GOARCH=amd64 go build -o $(binary) $(project_package)
	./deploy.sh > /dev/null
