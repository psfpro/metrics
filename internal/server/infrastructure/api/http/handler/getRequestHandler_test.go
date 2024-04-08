package handler

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRequestHandler_HandleRequest(t *testing.T) {
	ts := httptest.NewServer(Router())
	defer ts.Close()
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name       string
		target     string
		updateBody string
		body       string
		want       want
	}{
		{
			name:       "get gauge positive",
			target:     "/value",
			updateBody: `{"id":"1","type":"gauge","value":1.1}`,
			body:       `{"id":"1","type":"gauge"}`,
			want: want{
				code:        200,
				response:    `{"id":"1","type":"gauge","value":1.1}`,
				contentType: "application/json",
			},
		},
		{
			name:       "get counter positive",
			target:     "/value",
			updateBody: `{"id":"1","type":"counter","delta": 1}`,
			body:       `{"id":"1","type":"counter"}`,
			want: want{
				code:        200,
				response:    `{"id":"1","type":"counter","delta":1}`,
				contentType: "application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateRequest, _ := http.NewRequest(http.MethodPost, ts.URL+"/update", bytes.NewBufferString(tt.updateBody))
			resUpdate, _ := ts.Client().Do(updateRequest)
			assert.Equal(t, http.StatusOK, resUpdate.StatusCode)
			defer resUpdate.Body.Close()

			request, err := http.NewRequest(http.MethodPost, ts.URL+tt.target, bytes.NewBufferString(tt.body))
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
