package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
)

func main() {
	var port = flag.Int("port", 8080, "HTTP port")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("index.html")
		if err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}
		t.Execute(w, nil)
	})
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
