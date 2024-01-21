package server

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/psfpro/metrics/internal/server/infrastructure/storage"
	"log"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/psfpro/metrics/internal/server/application"
	"github.com/psfpro/metrics/internal/server/infrastructure/api/http"
	"github.com/psfpro/metrics/internal/server/infrastructure/api/http/handler"
)

type Container struct {
	app *http.App
}

func (c Container) App() *http.App {
	return c.app
}

func NewContainer() *Container {
	config := NewConfig()
	ctx := context.Background()

	gaugeMetricRepository := storage.NewGaugeMetricRepository()
	counterMetricRepository := storage.NewCounterMetricRepository()

	var storageAdapter storage.Adapter
	storageAdapter = storage.NewDumbAdapter()

	file, err := os.OpenFile(config.fileStoragePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("Storage adapter File error: %v", err)
	} else {
		log.Printf("Storage adapter File")
		storageAdapter = storage.NewFileAdapter(file, counterMetricRepository, gaugeMetricRepository)
	}

	db, err := sql.Open("pgx", config.databaseDsn.String())
	pingErr := db.Ping()
	if err != nil {
		log.Printf("Storage adapter DB error: %v", err)
	}
	if pingErr != nil {
		log.Printf("Storage adapter DB error: %v", pingErr)
	} else {
		log.Printf("Storage adapter DB")
		storageAdapter = storage.NewDbAdapter(db, counterMetricRepository, gaugeMetricRepository)
	}

	err = storageAdapter.Restore(ctx)
	if err != nil {
		log.Printf("Storage restore error: %v", err)
	}
	storageMiddleware := storage.NewMiddleware(storageAdapter)

	updateGaugeMetricHandler := &application.UpdateGaugeMetricHandler{
		Repository: gaugeMetricRepository,
	}
	updateCounterMetricHandler := &application.UpdateCounterMetricHandler{
		Repository: counterMetricRepository,
	}
	increaseCounterMetricHandler := &application.IncreaseCounterMetricHandler{
		Repository: counterMetricRepository,
	}

	pingRequestHandler := handler.NewPingRequestHandler(db)
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
	router.Get(`/ping`, pingRequestHandler.HandleRequest)
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
