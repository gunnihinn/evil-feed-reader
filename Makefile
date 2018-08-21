binary=evil-feed-reader
source = $(shell find -name "*.go")
package = github.com/gunnihinn/evil-feed-reader/cmd/evil-feed-reader
manpage=evil-feed-reader.1

$(binary): $(source)
	go build -o $(binary) $(package)

check: $(source)
	go test ./...

debug:
	go build -gcflags=all='-l -N' -o $(binary) $(package)
	dlv exec ./$(binary) -- -config test.yaml

clean:
	rm -f $(binary)

$(manpage): man.adoc
	asciidoctor -b manpage man.adoc -o $(manpage)

package: $(binary) $(manpage)

