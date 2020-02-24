package framework_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/snobb/susanin/pkg/framework"
	"github.com/snobb/susanin/test/helper"
)

func TestFramework_New(t *testing.T) {
	tests := []struct {
		name    string
		what    *framework.Framework
		wantNil bool
	}{
		{
			"should instantiate a new framework",
			framework.New(),
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.what)
		})
	}
}

func TestFramework_Get_Head_Delete_Options(t *testing.T) {
	methods := []string{"GET", "HEAD", "DELETE", "OPTIONS"}
	varHandler := func(w http.ResponseWriter, r *http.Request) {
		values, ok := framework.GetValues(r.Context())
		assert.True(t, ok)
		body, err := json.Marshal(values)
		assert.NoError(t, err)

		w.WriteHeader(200)
		w.Write(body)
	}

	for _, m := range methods {
		tests := []struct {
			name       string
			path       string
			wantValues *map[string]string
			wantBody   string
			wantCode   int
		}{
			{
				"should match /short endpoint and return HTTP 200 with 'short' body",
				"/short",
				nil,
				"short",
				200,
			},
			{
				"should match / endpoint and return HTTP 200 with 'root' body",
				"/",
				nil,
				"root",
				200,
			},
			{
				"should match /home endpoint and return HTTP 200 with 'home' body",
				"/home/foobar",
				nil,
				"home",
				200,
			},
			{
				"should match /hello/<vars> endpoint and return HTTP 200 with json body",
				"/hello/john/doe",
				nil,
				`{"fname":"john","lname":"doe"}`,
				200,
			},
			{
				"should not match any endpoind and return HTTP 400",
				"/foobar",
				nil,
				"{\"code\":404,\"msg\":\"Endpoint is not found\"}\n",
				404,
			},
		}

		fw := framework.New()

		routes := []struct {
			route   string
			handler http.HandlerFunc
		}{
			{"/", helper.HandlerFactory(200, "root")},
			{"/short", helper.HandlerFactory(200, "short")},
			{"/home/*", helper.HandlerFactory(200, "home")},
			{"/hello/:fname/:lname/", varHandler},
		}

		for _, r := range routes {
			switch m {
			case "GET":
				fw.Get(r.route, r.handler)
			case "HEAD":
				fw.Head(r.route, r.handler)
			case "DELETE":
				fw.Delete(r.route, r.handler)
			case "OPTIONS":
				fw.Options(r.route, r.handler)
			}
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				rr := httptest.NewRecorder()

				req, err := http.NewRequest(m, tt.path, nil)
				assert.NoError(t, err)

				fw.ServeHTTP(rr, req)
				assert.Equal(t, tt.wantBody, rr.Body.String())
				assert.Equal(t, tt.wantCode, rr.Code)
			})
		}

		fw.Clear()
		rr := httptest.NewRecorder()
		req, err := http.NewRequest(m, "/short", nil)
		assert.NoError(t, err)
		fw.ServeHTTP(rr, req)
		assert.Equal(t, 404, rr.Code)
	}
}

func TestFramework_Post_Put_Patch(t *testing.T) {
	methods := []string{"POST", "PUT", "PATCH"}

	for _, m := range methods {
		handler := func(code int, name string) func(http.ResponseWriter, *http.Request) {
			return func(w http.ResponseWriter, r *http.Request) {
				body, err := ioutil.ReadAll(r.Body)
				assert.NoError(t, err)

				assert.Equal(t, name, string(body))
				w.WriteHeader(code)
				w.Write([]byte(name))
			}
		}

		tests := []struct {
			name     string
			path     string
			wantBody string
			wantCode int
		}{
			{
				"should route to POST / handler and have a body",
				"/",
				"root",
				200,
			},
			{
				"should route to POST /short handler and have a body",
				"/short",
				"short",
				200,
			},
			{
				"should route to POST /home splat handler and have a body",
				"/home/test",
				"home",
				200,
			},
		}

		fw := framework.New()

		routes := []struct {
			route   string
			handler http.HandlerFunc
		}{
			{"/", handler(200, "root")},
			{"/short", handler(200, "short")},
			{"/home/*", handler(200, "home")},
		}

		for _, r := range routes {
			switch m {
			case "POST":
				fw.Post(r.route, r.handler)
			case "PUT":
				fw.Put(r.route, r.handler)
			case "PATCH":
				fw.Patch(r.route, r.handler)
			}
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				rr := httptest.NewRecorder()

				req, err := http.NewRequest(m, tt.path, strings.NewReader(tt.wantBody))
				assert.NoError(t, err)

				fw.ServeHTTP(rr, req)
				assert.Equal(t, tt.wantBody, rr.Body.String())
				assert.Equal(t, tt.wantCode, rr.Code)
			})
		}

		fw.Clear()
		rr := httptest.NewRecorder()
		req, err := http.NewRequest(m, "/short", nil)
		assert.NoError(t, err)
		fw.ServeHTTP(rr, req)
		assert.Equal(t, 404, rr.Code)
	}
}
