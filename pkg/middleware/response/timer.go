package response

/**
 * @author: Alex Kozadaev
 */

import (
	"net/http"
	"time"

	"github.com/snobb/susanin/pkg/logging"
	"github.com/snobb/susanin/pkg/middleware"
)

// NewTimer adds time counter middleware
func NewTimer(logger logging.Logger) middleware.Middleware {
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
