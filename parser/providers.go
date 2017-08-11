package parser

import (
	"crypto/tls"
	"io"
	"net/http"
	"os"
	"strings"
)

type Provider func(string) (io.ReadCloser, error)

var client = http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

// HTTP defines an HTTP provider.
func HTTP(url string) (io.ReadCloser, error) {
	response, err := client.Get(url)

	return response.Body, err
}

// File defines a File provider.
func File(filename string) (io.ReadCloser, error) {
	return os.Open(filename)
}

// String defines a String provider.
func String(str string) (io.ReadCloser, error) {
	return readerCloser{reader: strings.NewReader(str)}, nil
}

// readerCloser defines a trivial Close() method on a Reader
type readerCloser struct {
	reader io.Reader
}

func (r readerCloser) Read(bs []byte) (n int, err error) {
	return r.reader.Read(bs)
}

func (r readerCloser) Close() error { return nil }
