package main

import (
	"flag"
	"log/slog" // TODO: look into https://pkg.go.dev/go.uber.org/zap
	"net/http"
	"os"
)

// application struct to hold application-wide dependencies
type application struct {
	logger *slog.Logger
}

func main() {
	// new command-line flag 'addr', default value ":4000", and short description
	// value of flag is stored in addr variable at runtime
	addr := flag.String("addr", ":4000", "HTTP network address")
	// read command-line flage and assign it to addr variable
	// call before use addr variable otherwise it will always contian the default value
	flag.Parse()

	// initialize new structured logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// init a new instance of application stuct containing dependencies
	app := &application{
		logger: logger,
	}

	mux := http.NewServeMux()

	// serves files out of specified directory
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// register file server to handler all URL paths starting with "/static/"
	// for matching paths, strip the "/static" prefix before request reaches file server
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home) // Restrict this route to exact matches on / only.
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	logger.Info("starting server", "addr", *addr)

	err := http.ListenAndServe(*addr, mux)
	logger.Error(err.Error())
	os.Exit(1)
}
