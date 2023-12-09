package middlewares

import (
	"goblog/pkg/session"
	"net/http"
)

// StartSession initialize a session instance, must be called before using Session
func StartSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Start a session
		session.StartSession(w, r)
		// Process the next middleware or the main handler
		next.ServeHTTP(w, r)
	})
}
