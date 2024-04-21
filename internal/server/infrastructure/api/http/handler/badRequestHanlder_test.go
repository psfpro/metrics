package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBadRequestHandler_HandleRequest(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name   string
		target string
		want   want
	}{
		{
			name:   "positive test",
			target: "/request",
			want: want{
				code:        400,
				response:    "",
				contentType: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.target, nil)
			w := httptest.NewRecorder()
			NewBadRequestHandler().HandleRequest(w, request)

			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			defer func() {
				err := res.Body.Close()
				assert.NoError(t, err)
			}()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, tt.want.response, string(resBody))
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
