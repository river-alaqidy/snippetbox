package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// routes method returns a servemux containing our application routes.
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	// chain standar middleware
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	// return app.recoverPanic(app.logRequest(commonHeaders(mux))) // old non-chained middleware setup
	return standard.Then(mux)
}
