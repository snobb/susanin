package middleware_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/snobb/susanin/pkg/framework"
	"github.com/snobb/susanin/pkg/logging"
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
			s      *framework.Framework
			buf    bytes.Buffer
			req    *http.Request
			logger logging.Logger
			rr     *httptest.ResponseRecorder
			err    error
		)

		g.Before(func() {
			logger = logging.New("timer", &buf)
			s = framework.NewFramework()
		})

		g.JustBeforeEach(func() {
			rr = httptest.NewRecorder()
		})

		g.Describe("Timer middleware", func() {
			g.Before(func() {
				s.Attach(middleware.Timer(logger))
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

				lines, err := helper.ParseAllJSONLog(&buf)
				Expect(err).To(BeNil())

				Expect(lines[0]).To(HaveKeyWithValue("level", "trace"))
				Expect(lines[0]).To(HaveKeyWithValue("name", "TIMER"))
				Expect(lines[0]).To(HaveKey("time"))
				Expect(lines[0]).To(HaveKeyWithValue("pid", BeEquivalentTo(os.Getpid())))
				Expect(lines[0]).To(HaveKeyWithValue("uri", "/foo/bar"))
				Expect(lines[0]).To(HaveKeyWithValue("msg", "accepted connection"))

				Expect(lines[1]).To(HaveKeyWithValue("level", "trace"))
				Expect(lines[1]).To(HaveKeyWithValue("name", "TIMER"))
				Expect(lines[1]).To(HaveKey("time"))
				Expect(lines[1]).To(HaveKeyWithValue("pid", BeEquivalentTo(os.Getpid())))
				Expect(lines[1]).To(HaveKeyWithValue("elapsed", BeNumerically(">", 100)))
				Expect(lines[1]).To(HaveKey("elapsed_str"))
			})
		})
	})
}
