package request_test

/**
 * @author: Alex Kozadaev
 */

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/snobb/susanin/pkg/logging"
	"github.com/snobb/susanin/pkg/middleware/request"
	"github.com/snobb/susanin/test/helper"
	"github.com/stretchr/testify/assert"
)

func loggerHandler(logger logging.Logger) http.Handler {
	mw := request.NewLogger(logger)
	return mw(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
}

func TestLogger(t *testing.T) {
	var buf bytes.Buffer

	tests := []struct {
		name        string
		args        logging.Logger
		handler     func(logging.Logger) http.Handler
		request     string
		uri         string
		body        io.Reader
		contentType string
		checker     []func(map[string]interface{})
	}{
		{
			"test of pkg/middleware/request/logger.go (GET request)",
			logging.New("logger", &buf),
			loggerHandler,
			"GET",
			"/foo/bar?filter=13",
			nil,
			"application/json",
			[]func(map[string]interface{}){
				func(v map[string]interface{}) {
					assert.Equal(t, 9, len(v))
					assert.Contains(t, v, "time")
					assert.Equal(t, "LOGGER", v["name"])
					assert.Equal(t, "trace", v["level"])
					assert.Contains(t, v, "headers")
					assert.Equal(t, "/foo/bar", v["uri"])
					assert.Equal(t, "request", v["type"])
					assert.Equal(t, "GET", v["method"])
					assert.Equal(t, "HTTP/1.1", v["proto"])
					assert.InDelta(t, os.Getpid(), v["pid"], 0)
				},
			},
		},
		{
			"test of pkg/middleware/request/logger.go (POST request)",
			logging.New("logger", &buf),
			loggerHandler,
			"POST",
			"/foo/bar",
			strings.NewReader(`{"foo":"bar"}`),
			"application/json",
			[]func(map[string]interface{}){
				func(v map[string]interface{}) {
					assert.Equal(t, 10, len(v))
					assert.Contains(t, v, "time")
					assert.Equal(t, "LOGGER", v["name"])
					assert.Equal(t, "trace", v["level"])
					assert.Contains(t, v, "headers")
					assert.Equal(t, "request", v["type"])
					assert.Equal(t, "POST", v["method"])
					assert.Equal(t, "HTTP/1.1", v["proto"])
					assert.InDelta(t, os.Getpid(), v["pid"], 0)
					assert.Equal(t, "bar", v["body"].(map[string]interface{})["foo"])

					hdrs := v["headers"].(map[string]interface{})
					assert.Contains(t, hdrs, "Content-Type")
					assert.Equal(t, "application/json", hdrs["Content-Type"].([]interface{})[0])
				},
			},
		},
		{
			"test of pkg/middleware/request/logger.go (POST request - non-json body)",
			logging.New("logger", &buf),
			loggerHandler,
			"POST",
			"/foo/bar",
			strings.NewReader("bar"),
			"text/plain",
			[]func(map[string]interface{}){
				func(v map[string]interface{}) {
					assert.Equal(t, 10, len(v))
					assert.Contains(t, v, "time")
					assert.Equal(t, "LOGGER", v["name"])
					assert.Equal(t, "trace", v["level"])
					assert.Contains(t, v, "headers")
					assert.Equal(t, "request", v["type"])
					assert.Equal(t, "POST", v["method"])
					assert.Equal(t, "HTTP/1.1", v["proto"])
					assert.InDelta(t, os.Getpid(), v["pid"], 0)
					assert.Equal(t, "bar", v["body"])

					hdrs := v["headers"].(map[string]interface{})
					assert.Contains(t, hdrs, "Content-Type")
					assert.Equal(t, "text/plain", hdrs["Content-Type"].([]interface{})[0])
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.request, tt.uri, tt.body)
			assert.Nil(t, err)

			req.Header.Set("content-type", tt.contentType)
			rr := httptest.NewRecorder()

			tt.handler(tt.args).ServeHTTP(rr, req)

			line, err := helper.ParseJSONLog(&buf)
			assert.Nil(t, err)

			tt.checker[0](line)

			buf.Reset()
		})
	}
}
