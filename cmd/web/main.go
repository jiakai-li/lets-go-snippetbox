package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

// Hold the application-wide dependencies
type application struct {
	logger *slog.Logger
}

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
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Initialize a new instance of application struct
	app := &application{
		logger: logger,
	}

	logger.Info("starting server", slog.String("addr", *addr))

	err := http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}
