package response

import (
	"context"
	"encoding/json"
	"net/http"
)

// Response is a generic way of responding with JSON from the HTTP endpoints
type Response struct {
	Writer http.ResponseWriter
}

// New returns a pointer to new Response writer
func New(w http.ResponseWriter) *Response {
	return &Response{Writer: w}
}

// Payload responds with a JSON pyaload
func (r *Response) Payload(ctx context.Context, payload interface{}) error {
	r.Writer.Header().Add("content-type", "application/json")

	encoder := json.NewEncoder(r.Writer)

	if err := encoder.Encode(payload); err != nil {
		r.Writer.WriteHeader(http.StatusInternalServerError)

		_ = encoder.Encode(map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"error":   err.Error(),
			"message": http.StatusText(http.StatusInternalServerError),
		})

		return err
	}

	return nil
}

// Write writes raw data, implementing io.Writer interface.
func (r *Response) Write(data []byte) (n int, err error) {
	return r.Writer.Write(data)
}

// Error responds with an Error
func (r *Response) Error(ctx context.Context, code int, err error) error {
	r.Writer.WriteHeader(code)

	errPayload := map[string]interface{}{
		"code":    code,
		"error":   err.Error(),
		"message": http.StatusText(code),
	}

	return r.Payload(ctx, errPayload)
}
