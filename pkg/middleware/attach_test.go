package middleware_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/snobb/susanin/pkg/framework"
	"github.com/snobb/susanin/pkg/logging"
	"github.com/snobb/susanin/pkg/middleware"
	"github.com/snobb/susanin/test/helper"
	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestAttach(t *testing.T) {
	g := goblin.Goblin(t)

	//special hook for gomega
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Generic", func() {
		var (
			s      *framework.Framework
			buf    bytes.Buffer
			logger logging.Logger
			req    *http.Request
			rr     *httptest.ResponseRecorder
			err    error
		)

		g.Before(func() {
			logger = logging.New("attach", &buf)
			s = framework.NewFramework()
		})

		g.JustBeforeEach(func() {
			rr = httptest.NewRecorder()
		})

		g.Describe("Attach middleware function", func() {
			g.Before(func() {
				handler := http.Handler(helper.HandlerFactory(200, "root"))
				mwHandler := middleware.Attach(handler,
					middleware.RequestLogger(logger),
					middleware.ResponseLogger(logger))
				s.Get("/*", http.HandlerFunc(mwHandler))
			})

			g.AfterEach(func() {
				fmt.Println(buf.String())
				buf.Reset()
			})

			g.It("Should attach request and response middleware successfully", func() {
				req, err = http.NewRequest("GET", "/foo/bar?filter=13", nil)
				Expect(err).To(BeNil())
				req.Header.Set("content-type", "application/json")

				handler := s.Router()
				handler.ServeHTTP(rr, req)

				lines, err := helper.ParseAllJSONLog(&buf)
				Expect(err).To(BeNil())

				Expect(len(lines)).To(Equal(2))
				Expect(lines[0]).To(HaveKeyWithValue("type", "request"))
				Expect(lines[1]).To(HaveKeyWithValue("type", "response"))
			})
		})
	})
}
