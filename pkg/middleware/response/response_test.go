package response_test

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/snobb/susanin/pkg/middleware/response"
)

func TestResponse_Payload(t *testing.T) {
	tests := []struct {
		name       string
		payload    interface{}
		wantResult string
		wantErr    bool
	}{
		{
			name:       "should write json correctly",
			wantResult: `"test"`,
			payload:    "test",
		},
		{
			name:       "should try writing json but fail",
			payload:    func() {}, // unsupported type
			wantResult: `{"code":500,"message":"Internal Server Error","error":"json: unsupported type: func()"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			r := response.New(rec)
			if err := r.Payload(context.Background(), tt.payload); (err != nil) != tt.wantErr {
				t.Errorf("Response.Payload() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.JSONEq(t, tt.wantResult, rec.Body.String())
		})
	}
}

func TestResponse_Error(t *testing.T) {
	tests := []struct {
		name       string
		code       int
		err        error
		wantResult string
		wantErr    bool
	}{
		{
			name:       "should write json correctly",
			code:       404,
			err:        errors.New("spanner"),
			wantResult: `{"code":404, "error":"spanner","message":"Not Found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			r := response.New(rec)
			if err := r.Error(context.Background(), tt.code, tt.err); (err != nil) != tt.wantErr {
				t.Errorf("Response.Payload() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.JSONEq(t, tt.wantResult, rec.Body.String())
		})
	}
}

func TestResponse_Write(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantN   int
		wantErr bool
	}{
		{
			name:  "should Write successfully",
			data:  []byte("foobar"),
			wantN: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			r := response.New(rec)

			gotN, err := r.Write(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Response.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotN != tt.wantN {
				t.Errorf("Response.Write() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}
