package response_test

/**
 * @author: Alex Kozadaev
 */

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/snobb/susanin/pkg/logging"
	"github.com/snobb/susanin/pkg/middleware/response"
	"github.com/snobb/susanin/test/helper"
	"github.com/stretchr/testify/assert"
)

func loggerHandler(logger logging.Logger) http.Handler {
	mw := response.NewLogger(logger)
	return mw(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
}

func TestLogger(t *testing.T) {
	var buf bytes.Buffer

	tests := []struct {
		name    string
		args    logging.Logger
		handler func(logging.Logger) http.Handler
		checker []func(map[string]interface{})
	}{
		{
			"test of pkg/middleware/response/logger.go",
			logging.New("logger", &buf),
			loggerHandler,
			[]func(map[string]interface{}){
				func(v map[string]interface{}) {
					assert.Equal(t, "trace", v["level"])
					assert.Equal(t, "LOGGER", v["name"])
					assert.Equal(t, "response", v["type"])
					assert.InDelta(t, 200, v["status"], 0)
					assert.Contains(t, v, "time")
					assert.InDelta(t, os.Getpid(), v["pid"], 0)
					assert.Equal(t, "", v["body"])
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/foo/bar", nil)
			assert.Nil(t, err)
			rr := httptest.NewRecorder()

			tt.handler(tt.args).ServeHTTP(rr, req)

			lines, err := helper.ParseAllJSONLog(&buf)
			assert.Nil(t, err)

			for i, line := range lines {
				tt.checker[i](line)
			}

			buf.Reset()
		})
	}
}
