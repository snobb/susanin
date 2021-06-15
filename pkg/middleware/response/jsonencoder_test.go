package response_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/snobb/susanin/pkg/middleware/response"
	"github.com/stretchr/testify/assert"
)

func getPayloadHandler(args interface{}) http.Handler {
	return response.JSONEncoder(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		response.WithPayload(r.Context(), args)
	}))
}

func getErrorHandler(args interface{}) http.Handler {
	return response.JSONEncoder(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		response.WithError(r.Context(), args.(response.Error))
	}))
}

func TestJSONEncoder(t *testing.T) {
	type args interface{}

	tests := []struct {
		name    string
		args    args
		handler func(interface{}) http.Handler
		want    string
	}{
		{
			"test if the response body gets encoded",
			map[string]interface{}{"result": 0},
			getPayloadHandler,
			`{"result":0}`,
		},
		{
			"test if the nested response body gets encoded",
			map[string]interface{}{
				"result": 0,
				"test": map[string]interface{}{
					"foo": 1,
					"bar": "string",
				},
			},
			getPayloadHandler,
			`{"result":0,"test":{"foo":1,"bar":"string"}}`,
		},
		{
			"test if the error is reported",
			response.Error{
				Code:    500,
				Error:   fmt.Errorf("spanner"),
				Message: "spanner thrown",
			},
			getErrorHandler,
			`{"code":500,"error":"spanner","message":"spanner thrown"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/", nil)
			assert.Nil(t, err)
			rec := httptest.NewRecorder()

			tt.handler(tt.args).ServeHTTP(rec, req)

			bytes, err := ioutil.ReadAll(rec.Result().Body)
			assert.Nil(t, err)

			body := strings.Trim(string(bytes), "\n")

			assert.JSONEq(t, tt.want, body)
		})
	}
}
