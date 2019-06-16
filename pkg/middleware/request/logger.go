package request

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/snobb/susanin/pkg/logging"
	"github.com/snobb/susanin/pkg/middleware"
)

// NewLogger create new request logger middleware
func NewLogger(logger logging.Logger) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fields := []interface{}{
				"type", "request",
				"method", r.Method,
				"uri", r.URL.Path,
				"proto", r.Proto,
				"headers", r.Header,
			}

			if r.Method == "POST" || r.Method == "PUT" {
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					panic(err)
				}

				// setting the body back for futher processing
				r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

				if hdr, ok := r.Header["Content-Type"]; ok && hdr[0] == "application/json" {
					var jsonBody map[string]interface{}

					if err := json.Unmarshal(body, &jsonBody); err != nil {
						w.WriteHeader(400)
						w.Write([]byte("JSON expected"))
						return
					}

					fields = append(fields, "body", jsonBody)

				} else {
					fields = append(fields, "body", string(body))
				}
			}

			logger.Trace(fields...)

			next.ServeHTTP(w, r)
		})
	}
}
