package response

/**
 * @author: Alex Kozadaev
 */

import (
	"bytes"
	"net/http"
)

// Buffer implements ResponseWriter interface and can be used to collect the response
// body in a middleware handlers.
type Buffer struct {
	Response http.ResponseWriter
	Status   int
	Body     *bytes.Buffer
}

// NewBuffer create a new response buffer
func NewBuffer(w http.ResponseWriter) *Buffer {
	return &Buffer{
		Response: w,
		Status:   200,
		Body:     new(bytes.Buffer),
	}
}

// Header returns the response header handle.
func (w *Buffer) Header() http.Header {
	return w.Response.Header() // pass the response headers
}

// Write buffers the response.
func (w *Buffer) Write(buf []byte) (int, error) {
	w.Body.Write(buf)
	return len(buf), nil
}

// WriteHeader stores the status code
func (w *Buffer) WriteHeader(status int) {
	w.Status = status
}

// Flush implements the Flusher interface
func (w *Buffer) Flush() {
	if w.Body.Len() > 0 {
		w.Response.WriteHeader(w.Status)
		if _, err := w.Response.Write(w.Body.Bytes()); err != nil {
			panic(err)
		}

		w.Body.Reset()
	}
}
