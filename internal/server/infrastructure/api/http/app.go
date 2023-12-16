package http

import (
	"github.com/psfpro/metrics/internal/server/application"
	handler2 "github.com/psfpro/metrics/internal/server/infrastructure/api/http/handler"
	"github.com/psfpro/metrics/internal/server/infrastructure/storage/memstorage"
	"net/http"
)

type App struct {
	config *Config
	mux    *http.ServeMux
}

func NewApp(config *Config) *App {
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

	badRequestHandler := handler2.NewBadRequestHandler()
	metricsRequestHandler := handler2.NewMetricsRequestHandler(gaugeMetricRepository, counterMetricRepository)
	updateGaugeRequestHandler := handler2.NewUpdateGaugeRequestHandler(updateGaugeMetricHandler)
	updateCounterRequestHandler := handler2.NewUpdateCounterRequestHandler(updateCounterMetricHandler, increaseCounterMetricHandler)

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, badRequestHandler.HandleRequest)
	mux.HandleFunc(`/metrics`, metricsRequestHandler.HandleRequest)
	mux.HandleFunc(`/update/gauge/`, updateGaugeRequestHandler.HandleRequest)
	mux.HandleFunc(`/update/counter/`, updateCounterRequestHandler.HandleRequest)

	return &App{
		config: config,
		mux:    mux,
	}
}

func (obj *App) Run() {
	err := http.ListenAndServe(obj.config.Address, obj.mux)
	if err != nil {
		panic(err)
	}
}
