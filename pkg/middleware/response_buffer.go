package middleware

/**
 * @author: Alex Kozadaev
 */

import (
	"bytes"
	"net/http"
)

// responseBuffer implements ResponseWriter interface and can be used to collect the response
// body in a middleware handlers.
type responseBuffer struct {
	Response http.ResponseWriter
	Status   int
	Body     *bytes.Buffer
}

func newResponseBuffer(w http.ResponseWriter) *responseBuffer {
	return &responseBuffer{
		Response: w,
		Status:   200,
		Body:     &bytes.Buffer{},
	}
}

func (w *responseBuffer) Header() http.Header {
	return w.Response.Header() // pass the response headers
}

func (w *responseBuffer) Write(buf []byte) (int, error) {
	w.Body.Write(buf)
	return len(buf), nil
}

func (w *responseBuffer) WriteHeader(status int) {
	w.Status = status
}

// flush needs to be called in order to deliver the intercepted response
func (w *responseBuffer) flush() {
	if w.Body.Len() > 0 {
		w.Response.WriteHeader(w.Status)
		if _, err := w.Response.Write(w.Body.Bytes()); err != nil {
			panic(err)
		}

		w.Body.Reset()
	}
}
