binary=evilfr

bin_src = $(shell find -name "*.go")
project_package = github.com/gunnihinn/evil-feed-reader

datadir = static
data = $(wildcard static/*)
bindata = bindata.go

$(binary): $(bindata) $(bin_src)
	go build -o $(binary) $(project_package)

$(bindata): $(data)
	go-bindata-assetfs $(datadir)/...

debug: $(binary)
	dlv exec ./$(binary)

clean:
	rm -f $(bindata)
	rm -f $(binary)

deploy: $(bindata) $(bin_src)
	GOOS=linux GOARCH=amd64 go build -o $(binary) $(project_package)
	./deploy.sh > /dev/null
