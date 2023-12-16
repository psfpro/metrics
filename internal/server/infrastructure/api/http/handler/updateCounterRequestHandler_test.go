package handler

import (
	"github.com/psfpro/metrics/internal/server/application"
	"github.com/psfpro/metrics/internal/server/infrastructure/storage/memstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateCounterRequestHandler_HandleRequest(t *testing.T) {
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
			target: "/update/counter/Metric/1",
			want: want{
				code:        200,
				response:    "",
				contentType: "",
			},
		},
		{
			name:   "metric increase",
			target: "/update/counter/Metric",
			want: want{
				code:        200,
				response:    "",
				contentType: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.target, nil)
			w := httptest.NewRecorder()
			repository := memstorage.NewCounterMetricRepository()
			NewUpdateCounterRequestHandler(
				&application.UpdateCounterMetricHandler{Repository: repository},
				&application.IncreaseCounterMetricHandler{Repository: repository},
			).HandleRequest(w, request)

			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, tt.want.response, string(resBody))
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
