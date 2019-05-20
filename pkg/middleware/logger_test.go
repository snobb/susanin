package middleware_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/snobb/susanin/pkg/framework"
	"github.com/snobb/susanin/pkg/logging"
	"github.com/snobb/susanin/pkg/middleware"
	"github.com/snobb/susanin/test/helper"
	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestLogger(t *testing.T) {
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
			logger = logging.New("logger", &buf)
			s = framework.NewFramework()
		})

		g.JustBeforeEach(func() {
			rr = httptest.NewRecorder()
		})

		g.Describe("RequestLogger middleware", func() {
			g.Before(func() {
				s.Attach(middleware.RequestLogger(logger))
				s.Get("/*", helper.HandlerFactory(200, "root"))
				s.Post("/*", helper.HandlerFactory(200, "root"))
			})

			g.AfterEach(func() {
				fmt.Println(buf.String())
				buf.Reset()
			})

			g.It("Should log the request successfully", func() {
				req, err = http.NewRequest("GET", "/foo/bar?filter=13", nil)
				Expect(err).To(BeNil())
				req.Header.Set("content-type", "application/json")

				handler := s.Router()
				handler.ServeHTTP(rr, req)

				fields, err := helper.ParseJSONLog(&buf)
				Expect(err).To(BeNil())

				Expect(len(fields)).To(Equal(9))
				Expect(fields).To(HaveKey("time"))
				Expect(fields).To(HaveKey("headers"))
				Expect(fields).To(HaveKeyWithValue("type", "request"))
				Expect(fields).To(HaveKeyWithValue("method", "GET"))
				Expect(fields).To(HaveKeyWithValue("uri", "/foo/bar"))
				Expect(fields).To(HaveKeyWithValue("proto", "HTTP/1.1"))

				Expect(fields).To(HaveKeyWithValue("level", "info"))
				Expect(fields).To(HaveKeyWithValue("name", "LOGGER"))
				Expect(fields).To(HaveKeyWithValue("pid", BeEquivalentTo(os.Getpid())))
			})

			g.It("Should log the POST request with JSON body successfully", func() {
				req, err = http.NewRequest("POST", "/foo/bar", strings.NewReader(`{"foo":"bar"}`))
				Expect(err).To(BeNil())
				req.Header.Set("content-type", "application/json")

				handler := s.Router()
				handler.ServeHTTP(rr, req)

				fields, err := helper.ParseJSONLog(&buf)
				Expect(err).To(BeNil())

				Expect(len(fields)).To(Equal(10))
				Expect(fields).To(HaveKey("time"))
				Expect(fields).To(HaveKeyWithValue("type", "request"))
				Expect(fields).To(HaveKeyWithValue("method", "POST"))
				Expect(fields).To(HaveKeyWithValue("uri", "/foo/bar"))
				Expect(fields).To(HaveKeyWithValue("proto", "HTTP/1.1"))
				Expect(fields).To(HaveKey("headers"))
				Expect(fields).To(HaveKeyWithValue("body",
					BeEquivalentTo(map[string]interface{}{"foo": "bar"})))
				Expect(fields).To(HaveKey("headers"))
				hdrs := fields["headers"].(map[string]interface{})
				Expect(hdrs).To(HaveKey("Content-Type"))
				Expect(hdrs["Content-Type"].([]interface{})[0]).To(Equal("application/json"))

				Expect(fields).To(HaveKeyWithValue("level", "info"))
				Expect(fields).To(HaveKeyWithValue("name", "LOGGER"))
				Expect(fields).To(HaveKeyWithValue("pid", BeEquivalentTo(os.Getpid())))
			})

			g.It("Should log the POST request with body successfully", func() {
				req, err = http.NewRequest("POST", "/foo/bar", strings.NewReader("foo"))
				Expect(err).To(BeNil())

				handler := s.Router()
				handler.ServeHTTP(rr, req)

				fields, err := helper.ParseJSONLog(&buf)
				Expect(err).To(BeNil())

				Expect(len(fields)).To(Equal(10))
				Expect(fields).To(HaveKey("time"))
				Expect(fields).To(HaveKeyWithValue("type", "request"))
				Expect(fields).To(HaveKeyWithValue("method", "POST"))
				Expect(fields).To(HaveKeyWithValue("uri", "/foo/bar"))
				Expect(fields).To(HaveKeyWithValue("proto", "HTTP/1.1"))
				Expect(fields).To(HaveKeyWithValue("body", "foo"))

				Expect(fields).To(HaveKeyWithValue("level", "info"))
				Expect(fields).To(HaveKeyWithValue("name", "LOGGER"))
				Expect(fields).To(HaveKeyWithValue("pid", BeEquivalentTo(os.Getpid())))
			})
		})

		g.Describe("ResponseLogger middleware", func() {
			g.Before(func() {
				s = framework.NewFramework()
				s.Attach(middleware.ResponseLogger(logger))
				s.Get("/*", helper.HandlerFactory(200, "root"))
				s.Post("/*", helper.HandlerFactory(200, "root"))
			})

			g.AfterEach(func() {
				fmt.Println(buf.String())
				buf.Reset()
			})

			g.It("Should log the response successfully", func() {
				req, err = http.NewRequest("GET", "/foo/bar?filter=13", nil)
				Expect(err).To(BeNil())
				req.Header.Set("content-type", "application/json")

				handler := s.Router()
				handler.ServeHTTP(rr, req)

				fields, err := helper.ParseJSONLog(&buf)
				Expect(err).To(BeNil())

				Expect(len(fields)).To(Equal(9))
				Expect(fields).To(HaveKey("time"))
				Expect(fields).To(HaveKeyWithValue("type", "response"))
				Expect(fields).To(HaveKeyWithValue("elapsed", BeNumerically(">", 100)))
				Expect(fields).To(HaveKeyWithValue("status", float64(200)))
				Expect(fields).To(HaveKeyWithValue("body", "root"))

				Expect(fields).To(HaveKeyWithValue("level", "info"))
				Expect(fields).To(HaveKeyWithValue("name", "LOGGER"))
				Expect(fields).To(HaveKeyWithValue("pid", BeEquivalentTo(os.Getpid())))
			})
		})
	})
}
