package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/snobb/susanin/pkg/logging"
)

// RequestLogger middleware
func RequestLogger(logger logging.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fields := map[string]interface{}{
				"type":    "request",
				"method":  r.Method,
				"uri":     r.URL.Path,
				"proto":   r.Proto,
				"headers": r.Header,
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

					fields["body"] = jsonBody

				} else {
					fields["body"] = string(body)
				}
			}

			logger.Info(fields)

			next.ServeHTTP(w, r)
		})
	}
}

// ResponseLogger middleware
func ResponseLogger(logger logging.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wbuf := newResponseBuffer(w)
			next.ServeHTTP(wbuf, r)

			body := wbuf.Body.Bytes()

			var normBody interface{}

			if hdr, ok := w.Header()["Content-Type"]; ok && hdr[0] == "application/json" {
				if err := json.Unmarshal(body, &normBody); err != nil {
					w.WriteHeader(400)
					w.Write([]byte("JSON Expected"))
					return
				}

			} else {
				normBody = string(body)
			}

			fields := map[string]interface{}{
				"status":  wbuf.Status,
				"time":    time.Now(),
				"type":    "response",
				"headers": wbuf.Header(),
				"body":    normBody,
				"elapsed": time.Since(start),
			}

			logger.Info(fields)
			wbuf.flush()
		})
	}
}
