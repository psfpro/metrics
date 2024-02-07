package handler

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateGaugeRequestHandler_HandleRequest(t *testing.T) {
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
			name:   "not found",
			target: "/request",
			want: want{
				code:        404,
				response:    "",
				contentType: "",
			},
		},
		{
			name:   "metric update",
			target: "/update/gauge/Metric/1",
			want: want{
				code:        200,
				response:    "",
				contentType: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(Router())
			defer ts.Close()
			request, err := http.NewRequest(http.MethodPost, ts.URL+tt.target, nil)
			require.NoError(t, err)

			res, _ := ts.Client().Do(request)
			assert.Equal(t, tt.want.code, res.StatusCode)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, tt.want.response, string(resBody))
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
