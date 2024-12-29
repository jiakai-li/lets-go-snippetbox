package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// Parse command line argument
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// Configure logger
	// Custom loggers created by slog.New() are concurrency-safe
	// You can share a single logger and use it across multiple goroutines
	// and in your HTTP handlers without needing to worry about race conditions

	// But if there are multiple structured loggers writing to the same destination
	// then you need to be careful and ensure that the destination's underlying `write()`
	// method is also safe for concurrent use
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	mux := http.NewServeMux()

	// Create a file server which serves files out of the "./ui/static" directory
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})

	// Use mux.Handle() function to register the file server as the handler
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// Register other handlers
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	logger.Info("starting server", slog.String("addr", *addr))

	err := http.ListenAndServe(*addr, mux)
	logger.Error(err.Error())
	os.Exit(1)
}

// Custom FileSystem for disabling directory listings
type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		indexFile := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(indexFile); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
