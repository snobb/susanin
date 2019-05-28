package middleware

/**
 * @author: Alex Kozadaev
 */

import (
	"net/http"
	"time"

	"github.com/snobb/susanin/pkg/logging"
)

// Timer adds time counting
func Timer(logger logging.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			logger.Trace("msg", "accepted connection", "uri", r.URL.Path)

			next.ServeHTTP(w, r)

			elapsed := time.Since(start)
			logger.Trace("elapsed", elapsed, "elapsed_str", elapsed.String())
		})
	}
}
