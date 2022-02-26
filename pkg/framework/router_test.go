package framework_test

import (
	"net/http"
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

func TestRouter_NewRouter(t *testing.T) {
	tests := []struct {
		name string
		what *framework.Router
	}{
		{
			"should instantiate correct router",
			framework.NewRouter(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.what)
		})
	}
}

func TestRouter_Handle(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			"should return an error if splat is in the middle of the path",
			"/test/*/hello",
			dummy,
			true,
		},
		{
			"should return an error if splat is in the middle of the path",
			"/test/hello/*/*",
			dummy,
			true,
		},
		{
			"should successfully add a correct pass handler with a variable",
			"/test/:param1/hello",
			dummy,
			false,
		},
		{
			"should return an error path with different variable already exists",
			"/test/:param2/hello",
			dummy,
			true,
		},
	}

	r := framework.NewRouter(nil)
	assert.NotNil(t, r)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
	tests := []struct {
		name        string
		path        string
		wantHandler http.HandlerFunc
		wantValues  map[string]string
		wantErr     bool
	}{
		{
			"should find a static handler",
			"/hello/test",
			static,
			nil,
			false,
		},
		{
			"should find a static1 handler",
			"/hello/test/all",
			static1,
			nil,
			false,
		},
		{
			"should find a static2 handler",
			"/test/all",
			static2,
			nil,
			false,
		},
		{
			"should find a dynamic handler",
			"/hello/alex",
			dynamic,
			map[string]string{"name": "alex"},
			false,
		},
		{
			"should find a dynamic handler",
			"/hello/alex/by-name",
			dynamic,
			map[string]string{"name": "alex"},
			false,
		},
		{
			"should find a splat handler",
			"/hello/alex/nonexistant",
			splat,
			map[string]string{"name": "alex"},
			false,
		},
		{
			"should fallback to generic splat on no match",
			"/foobar",
			fallback,
			nil,
			false,
		},
		{
			"should match static and return 2 variables",
			"/by-name/john/doe",
			dynamic,
			map[string]string{"fname": "john", "lname": "doe"},
			false,
		},
		{
			"should fallback to splat if no longer matching the line but fill the values",
			"/by-name/john",
			fallback,
			// currently it still fills the values during the longest match search
			map[string]string{"fname": "john"},
			false,
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
		})
	}
}
