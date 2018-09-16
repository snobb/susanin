package susanin

import (
	"log"
	"net/http"
	"time"
)

// TimerMiddleware adds time counting
func TimerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("accepting a request [uri: %s]", r.URL.Path)
		next(w, r)
		log.Printf("elapsed time %v\n", time.Since(start))
	}
}
