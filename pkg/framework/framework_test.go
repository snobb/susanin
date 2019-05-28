package framework_test

/**
 * @author: Alex Kozadaev
 */

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/snobb/susanin/pkg/framework"
	"github.com/snobb/susanin/test/helper"

	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestFrameWork(t *testing.T) {
	g := Goblin(t)

	//special hook for gomega
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("The Router function", func() {
		var s *framework.Framework
		var rr *httptest.ResponseRecorder

		g.Before(func() {
			s = framework.NewFramework()
			s.Get("/", helper.HandlerFactory(200, "root"))
			s.Get("/short", helper.HandlerFactory(200, "short"))
			s.Get("/home/*", helper.HandlerFactory(200, "home"))

			s.Get("/hello/:fname/:lname/", func(w http.ResponseWriter, r *http.Request) {
				values, ok := framework.GetValues(r)
				Expect(ok).To(BeTrue())
				message := fmt.Sprintf("%s %s", values["fname"], values["lname"])

				w.WriteHeader(200)
				w.Write([]byte(message))
			})

			s.Post("/post/*", func(w http.ResponseWriter, r *http.Request) {
				body, err := ioutil.ReadAll(r.Body)
				Expect(err).To(BeNil())
				w.Write(body)
			})
		})

		g.BeforeEach(func() {
			rr = httptest.NewRecorder()
		})

		g.It("should route to root handler", func() {
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			handler := s.Router()
			handler.ServeHTTP(rr, req)
			Expect(rr.Body.String()).To(Equal("root"))
		})

		g.It("should route to short handler", func() {
			req, err := http.NewRequest("GET", "/short/", nil)
			if err != nil {
				t.Fatal(err)
			}

			handler := s.Router()
			handler.ServeHTTP(rr, req)
			Expect(rr.Body.String()).To(Equal("short"))
		})

		g.It("should route to home handler", func() {
			req, err := http.NewRequest("GET", "/home/test", nil)
			if err != nil {
				t.Fatal(err)
			}

			handler := s.Router()
			handler.ServeHTTP(rr, req)
			Expect(rr.Body.String()).To(Equal("home"))
		})

		g.It("should route to hello handler", func() {
			req, err := http.NewRequest("GET", "/hello/john/doe/", nil)
			if err != nil {
				t.Fatal(err)
			}

			handler := s.Router()
			handler.ServeHTTP(rr, req)
			Expect(rr.Body.String()).To(Equal("john doe"))
		})

		g.It("should route to hello handler (no trailing backslash)", func() {
			req, err := http.NewRequest("GET", "/hello/john/doe", nil)
			if err != nil {
				t.Fatal(err)
			}

			handler := s.Router()
			handler.ServeHTTP(rr, req)
			Expect(rr.Body.String()).To(Equal("john doe"))
		})

		g.It("should route to POST post handler", func() {
			req, err := http.NewRequest("POST", "/post/handler",
				strings.NewReader("HELLO WORLD"))
			if err != nil {
				t.Fatal(err)
			}

			handler := s.Router()
			handler.ServeHTTP(rr, req)
			Expect(rr.Body.String()).To(Equal("HELLO WORLD"))
		})

		g.It("should fail to route to non-existing handler (router)", func() {
			req, err := http.NewRequest("GET", "/does/not/exist", nil)
			if err != nil {
				t.Fatal(err)
			}

			handler := s.Router()
			handler.ServeHTTP(rr, req)

			var body map[string]interface{}
			err = json.Unmarshal(rr.Body.Bytes(), &body)
			Expect(err).To(BeNil())
			Expect(body).To(HaveKeyWithValue("msg", "Endpoint is not found"))
			Expect(body).To(HaveKeyWithValue("code", float64(404)))
		})

		g.It("should fail to route to non-existing handler (framework)", func() {
			req, err := http.NewRequest("PATCH", "/does/not/exist", nil)
			if err != nil {
				t.Fatal(err)
			}

			handler := s.Router()
			handler.ServeHTTP(rr, req)

			var body map[string]interface{}
			err = json.Unmarshal(rr.Body.Bytes(), &body)
			Expect(err).To(BeNil())
			Expect(body).To(HaveKeyWithValue("msg", "Method is not found"))
			Expect(body).To(HaveKeyWithValue("code", float64(404)))
		})

		g.It("should route to the fallback handler", func() {
			s.Get("/*", helper.HandlerFactory(200, "fallback"))
			req, err := http.NewRequest("GET", "/does/not/exist", nil)
			if err != nil {
				t.Fatal(err)
			}

			handler := s.Router()
			handler.ServeHTTP(rr, req)
			Expect(rr.Body.String()).To(Equal("fallback"))
		})
	})
}
