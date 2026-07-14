package main

import (
	"fmt"
	"log" // TODO: look into https://pkg.go.dev/go.uber.org/zap
	"net/http"
	"strconv"
)

// home page
func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	w.Write([]byte("Hello from Snippetbox!"))
}

// Display a specific snippet
func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	msg := fmt.Sprintf("Display a specific snppet with ID %d...", id)

	w.Write([]byte(msg))
}

// Display a form for creating a new snippet
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated) // 201 successfully created, reminder since handler hasn't created anything yet
	w.Write([]byte("Save a new snippet..."))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", home) // Restrict this route to exact matches on / only.
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	log.Print("starting server on :4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
