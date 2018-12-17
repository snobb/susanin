package middleware_test

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/snobb/susanin/pkg/framework"
	"github.com/snobb/susanin/pkg/middleware"
	"github.com/snobb/susanin/test/helper"
	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestTimer(t *testing.T) {
	g := goblin.Goblin(t)

	//special hook for gomega
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Generic", func() {
		var (
			s   *framework.Framework
			buf bytes.Buffer
			req *http.Request
			rr  *httptest.ResponseRecorder
			err error
		)

		g.Before(func() {
			log.SetOutput(&buf)
			s = framework.NewFramework()
		})

		g.JustBeforeEach(func() {
			rr = httptest.NewRecorder()
		})

		g.Describe("Timer middleware", func() {
			g.Before(func() {
				s.Attach(middleware.Timer)
				s.Get("/*", helper.HandlerFactory(200, "root"))
			})

			g.AfterEach(func() {
				fmt.Println(buf.String())
				buf.Reset()
			})

			g.It("Should log the timing info successfully", func() {
				req, err = http.NewRequest("GET", "/foo/bar?filter=13", nil)
				Expect(err).To(BeNil())
				req.Header.Set("content-type", "application/json")

				handler := s.Router()
				handler.ServeHTTP(rr, req)

				bufString := buf.String()

				idx := strings.Index(bufString, "accepting a request")
				Expect(idx).To(Equal(20))

				idx = strings.Index(bufString, "/foo/bar")
				Expect(idx).To(Equal(46))

				idx = strings.Index(bufString, "elapsed time")
				Expect(idx).To(Equal(76))
			})
		})
	})
}
