package handler

import (
	"bytes"
	"github.com/psfpro/metrics/internal/server/application"
	"github.com/psfpro/metrics/internal/server/infrastructure/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateRequestHandler_HandleRequest(t *testing.T) {
	ts := httptest.NewServer(Router())
	defer ts.Close()
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name   string
		target string
		body   string
		want   want
	}{
		{
			name:   "not found",
			target: "/request",
			body:   "",
			want: want{
				code:        404,
				response:    "",
				contentType: "",
			},
		},
		{
			name:   "gauge metric update",
			target: "/update",
			body:   `{"id": "1", "type": "gauge", "value": 1.1}`,
			want: want{
				code:        200,
				response:    "",
				contentType: "",
			},
		},
		{
			name:   "counter metric update",
			target: "/update",
			body:   `{"id": "1", "type": "counter", "delta": 1}`,
			want: want{
				code:        200,
				response:    "",
				contentType: "",
			},
		},
		{
			name:   "counter metric increase",
			target: "/update",
			body:   `{"id": "1", "type": "counter"}`,
			want: want{
				code:        200,
				response:    "",
				contentType: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func BenchmarkUpdateRequestHandler_HandleRequest(b *testing.B) {
	const triesN = 1000
	gaugeMetricRepository := storage.NewGaugeMetricRepository()
	counterMetricRepository := storage.NewCounterMetricRepository()
	updateGaugeMetricHandler := &application.UpdateGaugeMetricHandler{
		Repository: gaugeMetricRepository,
	}
	updateCounterMetricHandler := &application.UpdateCounterMetricHandler{
		Repository: counterMetricRepository,
	}
	increaseCounterMetricHandler := &application.IncreaseCounterMetricHandler{
		Repository: counterMetricRepository,
	}
	updateRequestHandler := NewUpdateRequestHandler(updateGaugeMetricHandler, updateCounterMetricHandler, increaseCounterMetricHandler)
	slice := make([]*http.Request, triesN)
	for i := 0; i < triesN; i++ {
		request, _ := http.NewRequest(http.MethodPost, "/update", bytes.NewBufferString(`{"id": "1", "type": "gauge", "value": 1.1}`))
		slice[i] = request
	}
	w := httptest.NewRecorder()
	b.ResetTimer()

	for i := 0; i < triesN; i++ {
		updateRequestHandler.HandleRequest(w, slice[i])
	}
}
