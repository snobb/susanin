package response_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"

	"github.com/snobb/susanin/pkg/framework"
	"github.com/snobb/susanin/pkg/logging"
	"github.com/snobb/susanin/pkg/middleware/response"
	"github.com/snobb/susanin/test/helper"
)

func TestLogger(t *testing.T) {
	g := goblin.Goblin(t)

	//special hook for gomega
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Generic", func() {
		var (
			fw     *framework.Framework
			buf    bytes.Buffer
			req    *http.Request
			logger logging.Logger
			rr     *httptest.ResponseRecorder
			err    error
		)

		g.Before(func() {
			logger = logging.New("logger", &buf)
			fw = framework.New()
		})

		g.JustBeforeEach(func() {
			rr = httptest.NewRecorder()
		})

		g.Describe("response.Logger middleware", func() {
			g.Before(func() {
				fw.Attach(response.NewLogger(logger))
				fw.Get("/*", helper.HandlerFactory(200, "root"))
				fw.Post("/*", helper.HandlerFactory(200, "root"))
			})

			g.AfterEach(func() {
				buf.Reset()
			})

			g.It("Should log the response successfully", func() {
				req, err = http.NewRequest("GET", "/foo/bar?filter=13", nil)
				Expect(err).To(BeNil())
				req.Header.Set("content-type", "application/json")

				fw.ServeHTTP(rr, req)

				fields, err := helper.ParseJSONLog(&buf)
				Expect(err).To(BeNil())

				Expect(len(fields)).To(Equal(9))
				Expect(fields).To(HaveKey("time"))
				Expect(fields).To(HaveKeyWithValue("type", "response"))
				Expect(fields).To(HaveKeyWithValue("elapsed", BeNumerically(">", 100)))
				Expect(fields).To(HaveKeyWithValue("status", float64(200)))
				Expect(fields).To(HaveKeyWithValue("body", "root"))

				Expect(fields).To(HaveKeyWithValue("level", "trace"))
				Expect(fields).To(HaveKeyWithValue("name", "LOGGER"))
				Expect(fields).To(HaveKeyWithValue("pid", BeEquivalentTo(os.Getpid())))
			})
		})
	})
}
