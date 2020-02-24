package response_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"

	"github.com/snobb/susanin/pkg/framework"
	"github.com/snobb/susanin/pkg/logging"
	"github.com/snobb/susanin/pkg/middleware/response"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	values, ok := framework.GetValues(r)
	if !ok {
		response.WithError(r.Context(), response.Error{
			Code:    400,
			Error:   fmt.Errorf("%s", "name not found"),
			Message: http.StatusText(400),
		})
		return
	}

	response.WithPayload(r.Context(), map[string]interface{}{
		"has_name": ok,
		"name":     values["name"],
	})
}

func TestJSONEncoder(t *testing.T) {
	g := goblin.Goblin(t)

	//special hook for gomega
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Generic", func() {
		var (
			fw     *framework.Framework
			buf    bytes.Buffer
			req    *http.Request
			logger logging.Logger
			w      *httptest.ResponseRecorder
			err    error
		)

		g.Before(func() {
			logger = logging.New("logger", &buf)
			fw = framework.New()
		})

		g.JustBeforeEach(func() {
			w = httptest.NewRecorder()
		})

		g.Describe("response.JSONEncoder test", func() {
			g.Before(func() {
				fw.Attach(response.NewJSONEncoder(logger))
				fw.Get("/", Handler)
				fw.Get("/name/:name", Handler)
			})

			g.AfterEach(func() {
				buf.Reset()
			})

			g.It("Should JSON encode the response successfully", func() {
				req, err = http.NewRequest("GET", "/name/test", nil)
				Expect(err).To(BeNil())

				fw.ServeHTTP(w, req)

				body := w.Body.Bytes()
				Expect(err).To(BeNil())

				expect := `{"has_name":true,"name":"test"}`
				Expect(body).To(MatchJSON(expect))
			})

			g.It("Should return a JSON encoded error", func() {
				req, err = http.NewRequest("GET", "/", nil)
				Expect(err).To(BeNil())

				fw.ServeHTTP(w, req)

				resp := w.Result()
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).To(BeNil())

				expect := `{"code":400,"error":"name not found","message":"Bad Request"}`
				Expect(body).To(MatchJSON(expect))
			})
		})
	})
}
