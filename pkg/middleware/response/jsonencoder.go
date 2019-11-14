package response

/*
 * this middleware was inspired by work of Rafał Lorenz
 * https://rafallorenz.com/go/go-middleware-parsing-http-response-as-json/
 */

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/snobb/susanin/pkg/logging"
	"github.com/snobb/susanin/pkg/middleware"
)

type payloadKey struct{}

type encode struct {
	payload interface{}
}

// Error is a struct to store a JSON error representation
type Error struct {
	Code    int    `json:"code"`
	Error   error  `json:"error"`
	Message string `json:"message"`
}

// MarshalJSON is the implementation of Marshaler interface for Error
func (e Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"code":    e.Code,
		"error":   e.Error.Error(),
		"message": e.Message,
	})
}

func (enc *encode) write(payload interface{}) {
	enc.payload = payload
}

func fromContext(ctx context.Context) (*encode, bool) {
	enc, ok := ctx.Value(payloadKey{}).(*encode)
	return enc, ok
}

func contextWithPayload(ctx context.Context) context.Context {
	return context.WithValue(ctx, payloadKey{}, &encode{})
}

// WithPayload adds payload to the given context
func WithPayload(ctx context.Context, payload interface{}) {
	enc, ok := fromContext(ctx)
	if !ok {
		panic("no response in the context - please use response middleware")
	}

	enc.write(payload)
}

// WithError add an Error struct to the given context
func WithError(ctx context.Context, err Error) {
	WithPayload(ctx, err)
}

// NewJSONEncoder creates a new response middleware
func NewJSONEncoder(logger logging.Logger) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("content-type", "application/json")
			ctx := contextWithPayload(r.Context())

			next.ServeHTTP(w, r.WithContext(ctx))

			enc, ok := fromContext(ctx)
			if !ok {
				// no response - nothing to do
				return
			}

			// set the error code in case of an error
			switch t := enc.payload.(type) {
			case Error:
				w.WriteHeader(t.Code)
			case *Error:
				w.WriteHeader(t.Code)
			}

			encoder := json.NewEncoder(w)

			if enc.payload != nil {
				if err := encoder.Encode(enc.payload); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					encoder.Encode(map[string]interface{}{
						"code":    http.StatusInternalServerError,
						"error":   err.Error(),
						"message": http.StatusText(http.StatusInternalServerError),
					})
				}
			}
		})
	}
}
