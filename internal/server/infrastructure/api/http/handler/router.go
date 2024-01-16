package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/psfpro/metrics/internal/server/application"
	"github.com/psfpro/metrics/internal/server/infrastructure/storage/memstorage"
)

func Router() *chi.Mux {
	gaugeMetricRepository := memstorage.NewGaugeMetricRepository()
	counterMetricRepository := memstorage.NewCounterMetricRepository()
	updateGaugeMetricHandler := &application.UpdateGaugeMetricHandler{
		Repository: gaugeMetricRepository,
	}
	updateCounterMetricHandler := &application.UpdateCounterMetricHandler{
		Repository: counterMetricRepository,
	}
	increaseCounterMetricHandler := &application.IncreaseCounterMetricHandler{
		Repository: counterMetricRepository,
	}

	badRequestHandler := NewBadRequestHandler()
	notFoundHandler := NewNotFoundRequestHandler()
	metricsRequestHandler := NewMetricsRequestHandler(gaugeMetricRepository, counterMetricRepository)
	metricsListRequestHandler := NewMetricsListRequestHandler(gaugeMetricRepository, counterMetricRepository)
	updateGaugeRequestHandler := NewUpdateGaugeRequestHandler(updateGaugeMetricHandler)
	updateCounterRequestHandler := NewUpdateCounterRequestHandler(updateCounterMetricHandler, increaseCounterMetricHandler)
	updateRequestHandler := NewUpdateRequestHandler(updateGaugeMetricHandler, updateCounterMetricHandler, increaseCounterMetricHandler)
	getMetricValueRequestHandler := NewGetMetricValueRequestHandler(gaugeMetricRepository, counterMetricRepository)
	getRequestHandler := NewGetRequestHandler(gaugeMetricRepository, counterMetricRepository)

	router := chi.NewRouter()
	router.Use(middleware.RealIP, Compressor, Logger, middleware.Logger, middleware.Recoverer)
	router.Get(`/`, metricsListRequestHandler.HandleRequest)
	router.Get(`/metrics`, metricsRequestHandler.HandleRequest)
	router.Route(`/update`, func(router chi.Router) {
		router.Post(`/`, updateRequestHandler.HandleRequest)
		router.Post(`/gauge/{name}/{value}`, updateGaugeRequestHandler.HandleRequest)
		router.Post(`/counter/{name}`, updateCounterRequestHandler.HandleRequest)
		router.Post(`/counter/{name}/{value}`, updateCounterRequestHandler.HandleRequest)
		router.Post(`/{type}/{name}/{value}`, badRequestHandler.HandleRequest)
	})
	router.Post(`/value`, getRequestHandler.HandleRequest)
	router.Post(`/value/`, getRequestHandler.HandleRequest)
	router.Get(`/value/{type}/{name}`, getMetricValueRequestHandler.HandleRequest)
	router.NotFound(notFoundHandler.HandleRequest)

	return router
}
