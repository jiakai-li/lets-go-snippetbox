package main

import (
	"fmt"
	"net/http"
)

func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(
			"Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com",
		)

		w.Header().Set(
			"Referer-Policy",
			"origin-when-cross-origin",
		)

		w.Header().Set(
			"X-Content-Type-Options",
			"nosniff",
		)

		w.Header().Set(
			"X-Frame-Options",
			"deny",
		)

		w.Header().Set(
			"X-XSS-Protection",
			"0",
		)

		w.Header().Set(
			"Server",
			"Go",
		)

		// Any code before will execute on the way down the chain
		next.ServeHTTP(w, r)
		// Any code after will execute on the way back up the chain
	})
}

// The middleware has access to handler dependencies including the structured logger
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)
		app.logger.Info("received request", "ip", ip, "proto", proto, "method", method, "uri", uri)
		next.ServeHTTP(w, r)
	})
}

// Will only recover panics that happen in the same goroutine that executed the recoverPanic middleware
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Creates a deferred function that will always be run in the event of a panic
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				// Return a 500 Internal Server Error response
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
