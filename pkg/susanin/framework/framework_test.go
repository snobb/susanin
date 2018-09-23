package framework

/**
 * @author: Alex Kozadaev
 */

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"susanin/pkg/susanin/helpers"
	"susanin/pkg/susanin/router"

	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestHelloWorld(t *testing.T) {
	g := Goblin(t)

	//special hook for gomega
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Generic framework", func() {
		g.It("should test route binding", func() {
			f := NewFramework()
			dummy := helpers.HandlerFactory(200, "dummy")

			Expect(f.methods[mGet]).To(BeNil())
			Expect(f.methods[mPut]).To(BeNil())
			Expect(f.methods[mPost]).To(BeNil())
			Expect(f.methods[mDelete]).To(BeNil())
			Expect(f.methods[mPatch]).To(BeNil())

			f.Get("/", dummy)
			Expect(f.methods[mGet]).NotTo(BeNil())

			f.Put("/", dummy)
			Expect(f.methods[mPut]).NotTo(BeNil())

			f.Post("/", dummy)
			Expect(f.methods[mPost]).NotTo(BeNil())

			f.Delete("/", dummy)
			Expect(f.methods[mDelete]).NotTo(BeNil())

			f.Patch("/", dummy)
			Expect(f.methods[mPatch]).NotTo(BeNil())
		})

	})

	g.Describe("The Router function", func() {
		var s *Framework
		var rr *httptest.ResponseRecorder

		g.Before(func() {
			s = NewFramework()
			s.Get("/", helpers.HandlerFactory(200, "root"))
			s.Get("/short", helpers.HandlerFactory(200, "short"))
			s.Get("/home/*", helpers.HandlerFactory(200, "home"))

			s.Get("/hello/:fname/:lname/", func(w http.ResponseWriter, r *http.Request) {
				values, ok := router.GetValues(r)
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

			handler := http.HandlerFunc(s.Router())
			handler(rr, req)
			Expect(rr.Body.String()).To(Equal("root"))
		})

		g.It("should route to short handler", func() {
			req, err := http.NewRequest("GET", "/short/", nil)
			if err != nil {
				t.Fatal(err)
			}

			handler := http.HandlerFunc(s.Router())
			handler(rr, req)
			Expect(rr.Body.String()).To(Equal("short"))
		})

		g.It("should route to home handler", func() {
			req, err := http.NewRequest("GET", "/home/test", nil)
			if err != nil {
				t.Fatal(err)
			}

			handler := http.HandlerFunc(s.Router())
			handler(rr, req)
			Expect(rr.Body.String()).To(Equal("home"))
		})

		g.It("should route to hello handler", func() {
			req, err := http.NewRequest("GET", "/hello/john/doe/", nil)
			if err != nil {
				t.Fatal(err)
			}

			handler := http.HandlerFunc(s.Router())
			handler(rr, req)
			Expect(rr.Body.String()).To(Equal("john doe"))
		})

		g.It("should route to hello handler (no trailing backslash)", func() {
			req, err := http.NewRequest("GET", "/hello/john/doe", nil)
			if err != nil {
				t.Fatal(err)
			}

			handler := http.HandlerFunc(s.Router())
			handler(rr, req)
			Expect(rr.Body.String()).To(Equal("john doe"))
		})

		g.It("should route to POST post handler", func() {
			req, err := http.NewRequest("POST", "/post/handler",
				strings.NewReader("HELLO WORLD"))
			if err != nil {
				t.Fatal(err)
			}

			handler := http.HandlerFunc(s.Router())
			handler(rr, req)
			Expect(rr.Body.String()).To(Equal("HELLO WORLD"))
		})

		g.It("should fail to route to non-existing handler", func() {
			req, err := http.NewRequest("GET", "/does/not/exist", nil)
			if err != nil {
				t.Fatal(err)
			}

			handler := http.HandlerFunc(s.Router())
			handler(rr, req)
			Expect(rr.Body.String()).To(Equal("not found\n"))
		})

		g.It("should route to the fallback handler", func() {
			s.Get("/*", helpers.HandlerFactory(200, "fallback"))
			req, err := http.NewRequest("GET", "/does/not/exist", nil)
			if err != nil {
				t.Fatal(err)
			}

			handler := http.HandlerFunc(s.Router())
			handler(rr, req)
			Expect(rr.Body.String()).To(Equal("fallback"))
		})
	})
}
