package helper

/**
 * @author: Alex Kozadaev
 */

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	begin := 0

	for begin < len(str) {
		for begin < len(str) && str[begin] != '{' {
			begin++
		}

		end := begin + 1
		for end < len(str) && str[end] != '\n' && str[end] != '\r' {
			end++
		}

		row := make(map[string]interface{})
		err = json.Unmarshal(buf.Bytes()[begin:end], &row)
		if err != nil {
			return
		}
		data = append(data, row)

		begin = end + 1
	}

	return
}

// ParseJSONLog accepts buffer and finds and parses a json blob.
// This is useful when mocking log.* function.
func ParseJSONLog(buf *bytes.Buffer) (data map[string]interface{}, err error) {
	begin := strings.IndexRune(buf.String(), '{')
	if begin == -1 {
		return nil, fmt.Errorf("invalid json")
	}
	err = json.Unmarshal(buf.Bytes()[begin:], &data)
	return
}
