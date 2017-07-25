bin=evil-feed-reader
bindevel=$(bin)-devel

project_package=github.com/gunnihinn/evil-feed-reader/cmd/evil-feed-reader

data=cmd/evil-feed-reader/data
bindata=cmd/evil-feed-reader/bindata.go

service=evil-feed-reader

$(bin): assets
	go build -o $(bin) $(project_package)

assets: $(data)
	go-bindata -o $(bindata) $(data)

devel: devel-assets
	go build -o $(bindevel) $(project_package)

devel-assets: $(data)
	go-bindata -debug -o $(bindata) $(data)

clean:
	rm -f $(bindata)
	rm -f $(bin)
	rm -f $(bindevel)

deploy: feeds
	GOOS=freebsd GOARCH=amd64 go build -o remote/$(bin) $(project_package)
	./remote/deploy-reader.sh

feeds:
	cp feeds.cfg remote
	./remote/deploy-feeds.sh
