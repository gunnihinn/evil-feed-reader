binary=evil-feed-reader

lib_src = $(wildcard */*.go)
bin_src = $(wildcard cmd/evil-feed-reader/*.go)
project_package = github.com/gunnihinn/evil-feed-reader/cmd/evil-feed-reader

data = $(wildcard cmd/evil-feed-reader/data/*)
bindata = cmd/evil-feed-reader/bindata.go

$(binary): $(bindata) $(lib_src) $(bin_src)
	go build -o $(binary) $(project_package)

$(bindata): $(data)
	go-bindata -o $(bindata) $(data)

clean:
	rm -f $(bindata)
	rm -f $(binary)

deploy: $(bindata) $(lib_src) $(bin_src)
	GOOS=freebsd GOARCH=amd64 go build -o $(binary) $(project_package)
	./deploy.sh
