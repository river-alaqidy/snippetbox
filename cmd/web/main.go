package main

import (
	"flag"
	"log" // TODO: look into https://pkg.go.dev/go.uber.org/zap
	"net/http"
)

func main() {
	// new command-line flag 'addr', default value ":4000", and short description
	// value of flag is stored in addr variable at runtime
	addr := flag.String("addr", ":4000", "HTTP network address")
	// read command-line flage and assign it to addr variable
	// call before use addr variable otherwise it will always contian the default value
	flag.Parse()

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

	log.Printf("starting server on %s", *addr)

	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
