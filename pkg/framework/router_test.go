package framework_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"runtime"
	"testing"

	"github.com/snobb/susanin/pkg/framework"
	"github.com/snobb/susanin/test/helper"
	"github.com/stretchr/testify/assert"
)

var dummy http.HandlerFunc = helper.HandlerFactory(200, "dummy")
var dynamic http.HandlerFunc = helper.HandlerFactory(200, "dynamic")
var dynamic1 http.HandlerFunc = helper.HandlerFactory(200, "dynamic1")
var static http.HandlerFunc = helper.HandlerFactory(200, "static")
var static1 http.HandlerFunc = helper.HandlerFactory(200, "static1")
var static2 http.HandlerFunc = helper.HandlerFactory(200, "static2")
var splat http.HandlerFunc = helper.HandlerFactory(200, "splat")
var fallback http.HandlerFunc = helper.HandlerFactory(200, "fallback")

type response struct {
	code int
	msg  string
}

func TestRouter_NewRouter(t *testing.T) {
	tests := map[string]struct {
		what *framework.Router
	}{
		"should instantiate correct router": {
			what: framework.NewRouter(nil),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.NotNil(t, tt.what)
		})
	}
}

func TestRouter_Handle(t *testing.T) {
	tests := map[string]struct {
		path    string
		handler http.HandlerFunc
		wantErr bool
	}{
		"should return an error if splat is in the middle of the path": {
			path:    "/test/*/hello",
			handler: dummy,
			wantErr: true,
		},
		"should return an error if splat is there is more than one splats": {
			path:    "/test/hello/*/*",
			handler: dummy,
			wantErr: true,
		},
		"should successfully add a correct pass handler with a variable": {
			path:    "/test/:param1/hello",
			handler: dummy,
		},
		"should return an error path with different variable already exists": {
			path:    "/test/:param2/hello",
			handler: dummy,
			wantErr: true,
		},
	}

	r := framework.NewRouter(nil)
	assert.NotNil(t, r)

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := r.Handle(tt.path, tt.handler)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRouter_Lookup(t *testing.T) {
	tests := map[string]struct {
		path         string
		wantHandler  http.HandlerFunc
		wantValues   map[string]string
		wantResponse response
		wantErr      bool
	}{
		"should find a static handler": {
			path:         "/hello/test",
			wantHandler:  static,
			wantResponse: response{code: 200, msg: "static"},
		},
		"should find a static1 handler": {
			path:         "/hello/test/all",
			wantHandler:  static1,
			wantResponse: response{code: 200, msg: "static1"},
		},
		"should find a static2 handler": {
			path:         "/test/all",
			wantHandler:  static2,
			wantResponse: response{code: 200, msg: "static2"},
		},
		"should find a dynamic handler": {
			path:         "/hello/alex",
			wantHandler:  dynamic,
			wantValues:   map[string]string{"name": "alex"},
			wantResponse: response{code: 200, msg: "dynamic"},
		},
		"should find a dynamic handler with suffix": {
			path:         "/hello/alex/by-name",
			wantHandler:  dynamic1,
			wantValues:   map[string]string{"name": "alex"},
			wantResponse: response{code: 200, msg: "dynamic1"},
		},
		"should find a splat handler": {
			path:         "/hello/alex/nonexistant",
			wantHandler:  splat,
			wantValues:   map[string]string{"name": "alex"},
			wantResponse: response{code: 200, msg: "splat"},
		},
		"should fallback to generic splat on no match": {
			path:         "/foobar",
			wantHandler:  fallback,
			wantValues:   nil,
			wantResponse: response{code: 200, msg: "fallback"},
		},
		"should match static and return 2 variables": {
			path:         "/by-name/john/doe",
			wantHandler:  dynamic,
			wantValues:   map[string]string{"fname": "john", "lname": "doe"},
			wantResponse: response{code: 200, msg: "dynamic"},
		},
		"should fallback to splat if no longer matching the line but fill the values": {
			path:        "/by-name/john",
			wantHandler: fallback,
			// currently it still fills the values during the longest match search
			wantValues:   map[string]string{"fname": "john"},
			wantResponse: response{code: 200, msg: "fallback"},
		},
	}

	r := framework.NewRouter(nil)
	assert.NoError(t, r.Handle("/hello/:name", dynamic))
	assert.NoError(t, r.Handle("/hello/:name/by-name", dynamic1))
	assert.NoError(t, r.Handle("/hello/*", splat))
	assert.NoError(t, r.Handle("/hello/test", static))
	assert.NoError(t, r.Handle("/hello/test/all", static1))
	assert.NoError(t, r.Handle("/test/all", static2))
	assert.NoError(t, r.Handle("/by-name/:fname/:lname", dynamic))
	assert.NoError(t, r.Handle("/*", fallback))

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			handler, values := r.Lookup(tt.path)
			if tt.wantErr {
				assert.Nil(t, handler)
				return
			}

			assert.NotNil(t, handler)
			assert.Equal(t, tt.wantValues, values)

			wantHandler := runtime.FuncForPC(reflect.ValueOf(tt.wantHandler).Pointer()).Name()
			gotHandler := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
			assert.Equal(t, wantHandler, gotHandler)

			rec := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, tt.path, nil)
			assert.NoError(t, err)
			handler.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponse.code, rec.Code)
			assert.Equal(t, tt.wantResponse.msg, rec.Body.String())
		})
	}
}
