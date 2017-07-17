package provider

/*
 * Package provider defines functions that fetch raw bytes from a input.
 */

import (
	"io/ioutil"
	"net/http"
	"os"
)

// Provider is the interface that any provider function has to satisfy.
type Provider interface {
	Get(string) ([]byte, error)
}

// provider implements the Provider interface.
type provider struct {
	fromFunc func(string) ([]byte, error)
}

func (p provider) Get(resource string) ([]byte, error) {
	return p.fromFunc(resource)
}

// HTTP defines an HTTP provider.
func HTTP() Provider {
	return provider{
		fromFunc: fromHTTP,
	}
}

// File defines a File provider.
func File() Provider {
	return provider{
		fromFunc: fromFile,
	}
}

// String defines a String provider.
func String() Provider {
	return provider{
		fromFunc: fromString,
	}
}

func fromHTTP(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	blob, err := ioutil.ReadAll(resp.Body)

	return blob, err
}

func fromFile(filename string) ([]byte, error) {
	fh, err := os.Open(filename)
	if err != nil {
		return []byte{}, err
	}
	defer fh.Close()

	blob, err := ioutil.ReadAll(fh)

	return blob, err
}

func fromString(content string) ([]byte, error) {
	return []byte(content), nil
}
