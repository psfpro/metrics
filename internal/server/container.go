package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/psfpro/metrics/internal/server/application"
	"github.com/psfpro/metrics/internal/server/infrastructure/api/http"
	"github.com/psfpro/metrics/internal/server/infrastructure/api/http/handler"
	"github.com/psfpro/metrics/internal/server/infrastructure/storage/filestorage"
)

type Container struct {
	app *http.App
}

func (c Container) App() *http.App {
	return c.app
}

func NewContainer() *Container {
	config := NewConfig()

	entityManager := filestorage.NewEntityManager(config.fileStoragePath)
	storageMiddleware := filestorage.NewMiddleware(entityManager)
	gaugeMetricRepository := filestorage.NewGaugeMetricRepository(entityManager)
	counterMetricRepository := filestorage.NewCounterMetricRepository(entityManager)
	updateGaugeMetricHandler := &application.UpdateGaugeMetricHandler{
		Repository: gaugeMetricRepository,
	}
	updateCounterMetricHandler := &application.UpdateCounterMetricHandler{
		Repository: counterMetricRepository,
	}
	increaseCounterMetricHandler := &application.IncreaseCounterMetricHandler{
		Repository: counterMetricRepository,
	}

	badRequestHandler := handler.NewBadRequestHandler()
	notFoundHandler := handler.NewNotFoundRequestHandler()
	metricsRequestHandler := handler.NewMetricsRequestHandler(gaugeMetricRepository, counterMetricRepository)
	metricsListRequestHandler := handler.NewMetricsListRequestHandler(gaugeMetricRepository, counterMetricRepository)
	updateGaugeRequestHandler := handler.NewUpdateGaugeRequestHandler(updateGaugeMetricHandler)
	updateCounterRequestHandler := handler.NewUpdateCounterRequestHandler(updateCounterMetricHandler, increaseCounterMetricHandler)
	updateRequestHandler := handler.NewUpdateRequestHandler(updateGaugeMetricHandler, updateCounterMetricHandler, increaseCounterMetricHandler)
	getMetricValueRequestHandler := handler.NewGetMetricValueRequestHandler(gaugeMetricRepository, counterMetricRepository)
	getRequestHandler := handler.NewGetRequestHandler(gaugeMetricRepository, counterMetricRepository)

	router := chi.NewRouter()
	router.Use(middleware.RealIP, handler.Compressor, handler.Logger, middleware.Logger, storageMiddleware.Handle, middleware.Recoverer)
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

	app := http.NewApp(config.serverAddress.String(), router)

	return &Container{
		app: app,
	}
}
