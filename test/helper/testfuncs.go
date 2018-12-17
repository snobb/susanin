package helper

/**
 * @author: Alex Kozadaev
 */

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)

// HandlerFactory is a helper function useful for unit tests
func HandlerFactory(code int, message string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Write([]byte(message))
	}
}

// ParseAllJSONLog accepts buffer and finds and parses ALL json blobs.
// This is useful when mocking log.* function.
func ParseAllJSONLog(buf *bytes.Buffer) (data []map[string]interface{}, err error) {
	str := buf.String()
	line := 0

	for line < buf.Len() {
		begin := strings.IndexRune(str, '{')
		end := strings.IndexAny(str, "\n\r")

		row := make(map[string]interface{})
		err = json.Unmarshal(buf.Bytes()[line+begin:line+end], &row)
		if err != nil {
			return
		}

		data = append(data, row)

		str = str[end+1:]
		line += end + 1
	}

	return
}

// ParseJSONLog accepts buffer and finds and parses a json blob.
// This is useful when mocking log.* function.
func ParseJSONLog(buf *bytes.Buffer) (data map[string]interface{}, err error) {
	begin := strings.IndexRune(buf.String(), '{')
	err = json.Unmarshal(buf.Bytes()[begin:], &data)
	return
}
