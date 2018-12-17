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

func TestLogger(t *testing.T) {
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

		g.Describe("RequestLogger middleware", func() {
			g.Before(func() {
				s.Attach(middleware.RequestLogger)
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

				Expect(len(fields)).To(Equal(6))
				Expect(fields).To(HaveKey("time"))
				Expect(fields).To(HaveKeyWithValue("type", "request"))
				Expect(fields).To(HaveKeyWithValue("method", "GET"))
				Expect(fields).To(HaveKeyWithValue("uri", "/foo/bar"))
				Expect(fields).To(HaveKeyWithValue("proto", "HTTP/1.1"))
			})

			g.It("Should log the POST request with JSON body successfully", func() {
				req, err = http.NewRequest("POST", "/foo/bar", strings.NewReader(`{"foo":"bar"}`))
				Expect(err).To(BeNil())
				req.Header.Set("content-type", "application/json")

				handler := s.Router()
				handler.ServeHTTP(rr, req)

				fields, err := helper.ParseJSONLog(&buf)
				Expect(err).To(BeNil())

				Expect(len(fields)).To(Equal(7))
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
			})

			g.It("Should log the POST request with body successfully", func() {
				req, err = http.NewRequest("POST", "/foo/bar", strings.NewReader("foo"))
				Expect(err).To(BeNil())

				handler := s.Router()
				handler.ServeHTTP(rr, req)

				fields, err := helper.ParseJSONLog(&buf)
				Expect(err).To(BeNil())

				Expect(len(fields)).To(Equal(7))
				Expect(fields).To(HaveKey("time"))
				Expect(fields).To(HaveKeyWithValue("type", "request"))
				Expect(fields).To(HaveKeyWithValue("method", "POST"))
				Expect(fields).To(HaveKeyWithValue("uri", "/foo/bar"))
				Expect(fields).To(HaveKeyWithValue("proto", "HTTP/1.1"))
				Expect(fields).To(HaveKeyWithValue("body", "foo"))
			})
		})

		g.Describe("ResponseLogger middleware", func() {
			g.Before(func() {
				s = framework.NewFramework()
				s.Attach(middleware.ResponseLogger)
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

				Expect(len(fields)).To(Equal(6))
				Expect(fields).To(HaveKey("time"))
				Expect(fields).To(HaveKeyWithValue("type", "response"))
				Expect(fields).To(HaveKeyWithValue("elapsed", BeAssignableToTypeOf("string")))
				Expect(fields).To(HaveKeyWithValue("status", float64(200)))
				Expect(fields).To(HaveKeyWithValue("body", "root"))
			})
		})
	})
}
