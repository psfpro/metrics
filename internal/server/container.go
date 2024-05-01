package server

import (
	"context"
	"database/sql"
	"log"
	http2 "net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/psfpro/metrics/internal/server/application"
	"github.com/psfpro/metrics/internal/server/infrastructure/api/http"
	"github.com/psfpro/metrics/internal/server/infrastructure/api/http/handler"
	"github.com/psfpro/metrics/internal/server/infrastructure/storage"
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

	file, err := os.OpenFile(config.StoragePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("Storage adapter File error: %v", err)
	} else {
		log.Printf("Storage adapter File")
		storageAdapter = storage.NewFileAdapter(file, counterMetricRepository, gaugeMetricRepository)
	}

	db, err := sql.Open("pgx", config.DatabaseDsn)
	if err != nil {
		log.Printf("Storage adapter DB error: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Printf("Storage adapter DB error: %v", err)
	} else {
		log.Printf("Storage adapter DB")
		storageAdapter = storage.NewDBAdapter(db, counterMetricRepository, gaugeMetricRepository)
	}

	if config.Restore {
		err = storageAdapter.Restore(ctx)
		if err != nil {
			log.Printf("Storage restore error: %v", err)
		} else {
			log.Printf("Storage restore")
		}
	}
	storageMiddleware := storage.NewMiddleware(storageAdapter)
	hashCheckerMiddleware := handler.NewHashChecker(config.HashKey)
	cryptoDecoderMiddleware := handler.NewCryptoDecoder(config.CryptoKey)

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
	updatesRequestHandler := handler.NewUpdatesRequestHandler(updateGaugeMetricHandler, updateCounterMetricHandler, increaseCounterMetricHandler)
	getMetricValueRequestHandler := handler.NewGetMetricValueRequestHandler(gaugeMetricRepository, counterMetricRepository)
	getRequestHandler := handler.NewGetRequestHandler(gaugeMetricRepository, counterMetricRepository)

	router := chi.NewRouter()
	router.Use(
		middleware.RealIP, handler.Compressor, handler.Logger, middleware.Logger, cryptoDecoderMiddleware.Decode,
		hashCheckerMiddleware.Check, storageMiddleware.Handle, middleware.Recoverer,
	)
	router.Mount("/debug", middleware.Profiler())
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
	router.Post(`/updates`, updatesRequestHandler.HandleRequest)
	router.Post(`/updates/`, updatesRequestHandler.HandleRequest)
	router.Post(`/value`, getRequestHandler.HandleRequest)
	router.Post(`/value/`, getRequestHandler.HandleRequest)
	router.Get(`/value/{type}/{name}`, getMetricValueRequestHandler.HandleRequest)
	router.NotFound(notFoundHandler.HandleRequest)
	srv := &http2.Server{Addr: config.Address, Handler: router}

	app := http.NewApp(srv)

	return &Container{
		app: app,
	}
}
