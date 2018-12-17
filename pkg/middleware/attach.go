package middleware

import (
	"net/http"
)

// Middleware is a type for Middleware function
type Middleware func(http.Handler) http.Handler

// Attach middleware to a handler
func Attach(handler http.Handler, mws ...Middleware) http.HandlerFunc {
	for _, mw := range mws {
		handler = mw(handler)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})
}
