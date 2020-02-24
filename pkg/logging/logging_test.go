package logging_test

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/snobb/susanin/pkg/logging"
	"github.com/snobb/susanin/test/helper"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	levels := []string{
		"Info",
		"Warn",
		"Error",
		"Debug",
		"Trace",
	}

	for _, level := range levels {
		testRunner(t, level)
	}
}

func testRunner(t *testing.T, level string) {
	var buf bytes.Buffer

	type fields struct {
		name string
		Out  *bytes.Buffer
	}

	tests := []struct {
		name    string
		fields  fields
		args    []interface{}
		checker func(map[string]interface{})
	}{
		{
			fmt.Sprintf("create logger and log %s message", strings.ToUpper(level)),
			fields{
				level,
				&buf,
			},
			[]interface{}{
				"msg", "test message",
				"key1", "value1",
				"key2", "value2",
				"key3", "value3",
				"answer", 42,
			},
			func(v map[string]interface{}) {
				assert.Equal(t, 9, len(v))
				assert.Equal(t, strings.ToLower(level), v["level"])
				assert.Equal(t, "test message", v["msg"])
				assert.Equal(t, "value1", v["key1"])
				assert.Equal(t, "value2", v["key2"])
				assert.Equal(t, "value3", v["key3"])
				assert.InDelta(t, 42.0, v["answer"], 0)
				assert.Equal(t, strings.ToUpper(level), v["name"])
				assert.InDelta(t, os.Getpid(), v["pid"], 0)
				assert.Contains(t, v, "time")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dl := logging.New(tt.fields.name, &buf)

			method := reflect.ValueOf(dl).MethodByName(level)
			var args []reflect.Value
			for _, arg := range tt.args {
				args = append(args, reflect.ValueOf(arg))
			}
			method.Call(args)

			line, err := helper.ParseJSONLog(&buf)
			assert.Nil(t, err)

			tt.checker(line)

			buf.Reset()
		})
	}
}
