package middleware

/**
 * @author: Alex Kozadaev
 */

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
)

type debugResponseWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

// DebugMiddleware adds logging of HTTP protocol information
func DebugMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("<<<< %s %s %s", r.Method, r.URL.Path, r.Proto)
		for name, values := range r.Header {
			log.Printf(" %s : %v", name, values)
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			next(w, r)
			return
		}

		// setting the body back for futher processing
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		if len(body) > 0 {
			log.Printf(" << Body [content-length: %d]", r.ContentLength)
			for _, line := range strings.Split(string(body), "\n") {
				if len(line) > 0 {
					log.Printf(" %s", line)
				}
			}
		}

		rec := httptest.NewRecorder()
		next(rec, r)

		body = rec.Body.Bytes()
		for name, values := range rec.Header() {
			w.Header()[name] = values
		}

		w.WriteHeader(rec.Code)
		w.Write(body)

		res := rec.Result()
		log.Printf(">>>> %s %s", res.Status, res.Proto)

		// write headers from recorder to response
		for name, values := range res.Header {
			log.Printf(" %s : %v", name, values)
		}

		if len(body) > 0 {
			log.Printf(" << Body [content-length: %d]", len(body))
			for _, line := range strings.Split(string(body), "\n") {
				if len(line) > 0 {
					log.Printf(" %s", line)
				}
			}
		}
	}
}
