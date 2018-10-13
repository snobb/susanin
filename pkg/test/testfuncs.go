package test

/**
 * @author: Alex Kozadaev
 */

import (
	"net/http"
)

// HandlerFactory is a helper function useful for unit tests
func HandlerFactory(code int, message string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Write([]byte(message))
	}
}
