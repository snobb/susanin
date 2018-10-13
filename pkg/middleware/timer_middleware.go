package middleware

/**
 * @author: Alex Kozadaev
 */

import (
	"log"
	"net/http"
	"time"
)

// TimerMiddleware adds time counting
func TimerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("accepting a request [uri: %s]", r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("elapsed time %v\n", time.Since(start))
	})
}
