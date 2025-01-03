package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	// we need the driver’s init() function to run
	_ "github.com/go-sql-driver/mysql"
)

// Hold the application-wide dependencies
type application struct {
	logger *slog.Logger
}

func main() {
	// Parse command line argument
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
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

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	// Initialize a new instance of application struct
	app := &application{
		logger: logger,
	}

	logger.Info("starting server", slog.String("addr", *addr))

	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
