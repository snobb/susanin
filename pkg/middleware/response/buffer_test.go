package response_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/snobb/susanin/pkg/middleware/response"
	"github.com/stretchr/testify/assert"
)

func TestBuffer_Header(t *testing.T) {
	type fields struct {
		Response http.ResponseWriter
		Status   int
		Body     *bytes.Buffer
	}
	tests := []struct {
		name   string
		fields fields
		want   http.Header
	}{
		{
			"check that Header is wrapped",
			fields{recorder(), 0, new(bytes.Buffer)},
			http.Header{"Content-Type": []string{"application/json"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &response.Buffer{
				Response: tt.fields.Response,
				Status:   tt.fields.Status,
				Body:     tt.fields.Body,
			}
			if got := buf.Header(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Buffer.Header() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuffer_Write(t *testing.T) {
	type fields struct {
		Response http.ResponseWriter
		Status   int
		Body     *bytes.Buffer
	}
	type args struct {
		body []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			"check that Write buffers the body",
			fields{recorder(), 0, &bytes.Buffer{}},
			args{[]byte(`{"ima": "pc"}`)},
			13,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &response.Buffer{
				Response: tt.fields.Response,
				Status:   tt.fields.Status,
				Body:     tt.fields.Body,
			}
			got, err := buf.Write(tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Buffer.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Buffer.Write() = %v, want %v", got, tt.want)
			}

			assert.JSONEq(t, string(tt.args.body), string(buf.Body.Bytes()))
		})
	}
}

func TestBuffer_WriteHeader(t *testing.T) {
	type fields struct {
		Response http.ResponseWriter
		Status   int
		Body     *bytes.Buffer
	}
	type args struct {
		status int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"check that WriteHeader stores the status",
			fields{recorder(), 0, &bytes.Buffer{}},
			args{200},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &response.Buffer{
				Response: tt.fields.Response,
				Status:   tt.fields.Status,
				Body:     tt.fields.Body,
			}
			buf.WriteHeader(tt.args.status)

			assert.Equal(t, tt.args.status, buf.Status)
		})
	}
}

func TestBuffer_Flush(t *testing.T) {
	type fields struct {
		Response http.ResponseWriter
		Status   int
		Body     *bytes.Buffer
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			"check that Flush writes header/body and resets the buffer",
			fields{recorder(), 200, bytes.NewBufferString(`{"ima": "pc"}`)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &response.Buffer{
				Response: tt.fields.Response,
				Status:   tt.fields.Status,
				Body:     tt.fields.Body,
			}
			buf.Flush()

			recorder, _ := tt.fields.Response.(*httptest.ResponseRecorder)
			assert.Equal(t, 200, recorder.Code)
			assert.JSONEq(t, `{"ima": "pc"}`, string(recorder.Body.Bytes()))
			assert.Equal(t, 0, buf.Body.Len())
		})
	}
}

func recorder() *httptest.ResponseRecorder {
	return &httptest.ResponseRecorder{
		HeaderMap: http.Header{"Content-Type": []string{"application/json"}},
		Body:      &bytes.Buffer{},
	}
}
