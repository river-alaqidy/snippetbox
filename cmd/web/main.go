package main

import (
	"log" // TODO: look into https://pkg.go.dev/go.uber.org/zap
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// serves files out of specified directory
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// register file server to handler all URL paths starting with "/static/"
	// for matching paths, strip the "/static" prefix before request reaches file server
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home) // Restrict this route to exact matches on / only.
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	log.Print("starting server on :4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
