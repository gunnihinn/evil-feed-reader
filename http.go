package main

import (
	"log"
	"net/http"
	"os"
)

type Handler struct {
	logger *log.Logger
	mux    *http.ServeMux
}

func NewHandler() Handler {
	return Handler{
		logger: log.New(os.Stdout, "", 0),
		mux:    http.NewServeMux(),
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}
