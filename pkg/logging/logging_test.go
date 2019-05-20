package logging_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/snobb/susanin/pkg/logging"
	"github.com/snobb/susanin/test/helper"
	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestRouter(t *testing.T) {
	g := goblin.Goblin(t)

	//special hook for gomega
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Default Logger", func() {
		var (
			logger logging.Logger
			buf    bytes.Buffer
		)

		g.Before(func() {
			logger = logging.New("test", &buf)
		})

		g.AfterEach(func() {
			fmt.Println(buf.String())
			buf.Reset()
		})

		g.It("Should create logger and log info message", func() {
			meta := logging.Meta{
				"msg":   "test message",
				"meta1": "value1",
				"meta2": "value2",
				"meta3": "value3",
			}

			logger.Info(meta)

			fields, err := helper.ParseJSONLog(&buf)
			Expect(err).To(BeNil())
			Expect(fields).To(HaveLen(8))
			Expect(fields).To(HaveKeyWithValue("level", "info"))
			Expect(fields).To(HaveKeyWithValue("msg", "test message"))
			Expect(fields).To(HaveKeyWithValue("meta1", "value1"))
			Expect(fields).To(HaveKeyWithValue("meta2", "value2"))
			Expect(fields).To(HaveKeyWithValue("meta3", "value3"))
			Expect(fields).To(HaveKeyWithValue("name", "TEST"))
			Expect(fields).To(HaveKeyWithValue("pid", BeEquivalentTo(os.Getpid())))
			Expect(fields).To(HaveKey("time"))
		})

		g.It("Should create logger and log error message", func() {
			meta := logging.Meta{
				"msg": "test message",
			}

			logger.Error(meta)

			fields, err := helper.ParseJSONLog(&buf)
			Expect(err).To(BeNil())
			Expect(fields).To(HaveLen(5))
			Expect(fields).To(HaveKeyWithValue("level", "error"))
			Expect(fields).To(HaveKeyWithValue("msg", "test message"))
			Expect(fields).To(HaveKeyWithValue("name", "TEST"))
			Expect(fields).To(HaveKeyWithValue("pid", BeEquivalentTo(os.Getpid())))
			Expect(fields).To(HaveKey("time"))
		})

		g.It("Should create sublogger and log info message", func() {
			l := logger.NewSubLogger("SUB")
			meta := logging.Meta{
				"msg": "test message",
				"data": map[string]string{
					"meta1": "value1",
					"meta2": "value2",
					"meta3": "value3",
				},
			}

			l.Info(meta)

			fields, err := helper.ParseJSONLog(&buf)
			Expect(err).To(BeNil())
			Expect(fields).To(HaveLen(6))
			Expect(fields).To(HaveKeyWithValue("level", "info"))
			Expect(fields).To(HaveKeyWithValue("msg", "test message"))
			Expect(fields).To(HaveKeyWithValue("name", "TEST::SUB"))
			Expect(fields).To(HaveKeyWithValue("pid", BeEquivalentTo(os.Getpid())))
			Expect(fields).To(HaveKey("time"))
			Expect(fields).To(HaveKeyWithValue("data", BeEquivalentTo(map[string]interface{}{
				"meta1": "value1",
				"meta2": "value2",
				"meta3": "value3",
			})))
		})

		g.It("Should create logger and log debug message", func() {
			meta := logging.Meta{
				"msg":   "test message",
				"debug": true,
			}

			logger.Log(logging.Debug, meta)

			fields, err := helper.ParseJSONLog(&buf)
			Expect(err).To(BeNil())
			Expect(fields).To(HaveLen(6))
			Expect(fields).To(HaveKeyWithValue("level", "debug"))
			Expect(fields).To(HaveKeyWithValue("debug", true))
			Expect(fields).To(HaveKeyWithValue("msg", "test message"))
			Expect(fields).To(HaveKeyWithValue("name", "TEST"))
			Expect(fields).To(HaveKeyWithValue("pid", BeEquivalentTo(os.Getpid())))
			Expect(fields).To(HaveKey("time"))
		})
	})
}
