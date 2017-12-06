package main

import (
	"net/http"
)

type HTTP struct {
	config   Config
	response *http.Response
	err      error
}
