binary=evil-feed-reader
source = $(shell find -name "*.go")
package = github.com/gunnihinn/evil-feed-reader/cmd/evil-feed-reader

$(binary): $(source)
	go build -o $(binary) $(package)

debug: $(binary)
	dlv exec ./$(binary)

clean:
	rm -f $(binary)

deploy: $(source)
	GOOS=linux GOARCH=amd64 go build -o $(binary) $(package)
	./deploy.sh > /dev/null
