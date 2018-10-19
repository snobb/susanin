package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// RequestLogger middleware
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fields := map[string]interface{}{
			"time":    time.Now(),
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

			var jsonBody map[string]interface{}
			if err := json.Unmarshal(body, &jsonBody); err != nil {
				panic(err)
			}

			// setting the body back for futher processing
			r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			fields["body"] = jsonBody
		}

		out, err := json.Marshal(fields)
		if err != nil {
			panic(err)
		}

		log.Printf("%v", string(out))

		next.ServeHTTP(w, r)
	})
}

// ResponseLogger middleware
func ResponseLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wbuf := newResponseBuffer(w)

		next.ServeHTTP(wbuf, r)

		body := wbuf.Body.Bytes()

		var normBody interface{}

		if hdr, ok := w.Header()["Content-Type"]; ok && hdr[0] == "application/json" {
			if err := json.Unmarshal(body, &normBody); err != nil {
				panic(err)
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
		}

		out, err := json.Marshal(fields)
		if err != nil {
			panic(err)
		}

		log.Printf("%v", string(out))
		wbuf.flush()
	})
}
