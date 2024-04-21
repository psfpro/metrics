package handler

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/psfpro/metrics/internal/server/application"
	"github.com/psfpro/metrics/internal/server/infrastructure/storage"
)

func TestUpdateCounterRequestHandler_HandleRequest(t *testing.T) {
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
			request, err := http.NewRequest(http.MethodPost, ts.URL+tt.target, nil)
			require.NoError(t, err)

			res, _ := ts.Client().Do(request)
			assert.Equal(t, tt.want.code, res.StatusCode)
			defer func() {
				err = res.Body.Close()
				assert.NoError(t, err)
			}()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, tt.want.response, string(resBody))
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func BenchmarkUpdateCounterRequestHandler_HandleRequest(b *testing.B) {
	const triesN = 1000
	counterMetricRepository := storage.NewCounterMetricRepository()
	updateCounterMetricHandler := &application.UpdateCounterMetricHandler{
		Repository: counterMetricRepository,
	}
	increaseCounterMetricHandler := &application.IncreaseCounterMetricHandler{
		Repository: counterMetricRepository,
	}
	updateRequestHandler := NewUpdateCounterRequestHandler(updateCounterMetricHandler, increaseCounterMetricHandler)
	slice := make([]*http.Request, triesN)
	for i := 0; i < triesN; i++ {
		request, _ := http.NewRequest(http.MethodPost, "/update/counter/Metric/1", bytes.NewBufferString(""))
		rctx := chi.NewRouteContext()
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))
		rctx.URLParams.Add("name", "Metric")
		rctx.URLParams.Add("value", "1")
		slice[i] = request
	}
	w := httptest.NewRecorder()
	b.ResetTimer()

	for i := 0; i < triesN; i++ {
		updateRequestHandler.HandleRequest(w, slice[i])
	}
}
