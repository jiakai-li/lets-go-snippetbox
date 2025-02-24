package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// Create a file server which serves files out of the "./ui/static" directory
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})

	// Use mux.Handle() function to register the file server as the handler
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// Register other handlers
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	// recoverPanic -> logRequest -> commonHeaders -> servemux -> application handler
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
