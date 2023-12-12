package http

import (
	"github.com/psfpro/metrics/internal/application"
	"github.com/psfpro/metrics/internal/domain"
	"github.com/psfpro/metrics/internal/infrastructure/api/http/handler"
	"github.com/psfpro/metrics/internal/infrastructure/storage/memstorage"
	"net/http"
)

type App struct {
	config *Config
	mux    *http.ServeMux
}

func NewApp(config *Config) *App {
	gaugeMetricRepository := &memstorage.GaugeMetricRepository{Data: make(map[string]*domain.GaugeMetric)}
	counterMetricRepository := &memstorage.CounterMetricRepository{Data: make(map[string]*domain.CounterMetric)}
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
	metricsRequestHandler := handler.NewMetricsRequestHandler(gaugeMetricRepository, counterMetricRepository)
	updateGaugeRequestHandler := handler.NewUpdateGaugeRequestHandler(updateGaugeMetricHandler)
	updateCounterRequestHandler := handler.NewUpdateCounterRequestHandler(updateCounterMetricHandler, increaseCounterMetricHandler)

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

type Config struct {
	Address string
}
