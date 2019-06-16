package response

import (
	"net/http"
	"time"

	"github.com/snobb/susanin/pkg/logging"
	"github.com/snobb/susanin/pkg/middleware"
)

// NewLogger middleware
func NewLogger(logger logging.Logger) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wbuf := NewBuffer(w)
			next.ServeHTTP(wbuf, r)

			body := wbuf.Body.Bytes()

			logger.Trace(
				"status", wbuf.Status,
				"type", "response",
				"headers", wbuf.Header(),
				"body", string(body),
				"elapsed", time.Since(start))

			wbuf.Flush()
		})
	}
}
