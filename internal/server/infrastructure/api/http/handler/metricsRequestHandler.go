package handler

import (
	"encoding/json"
	"github.com/psfpro/metrics/internal/server/domain"
	"github.com/psfpro/metrics/internal/server/infrastructure/storage/memstorage"
	"log"
	"net/http"
)

type MetricsRequestHandler struct {
	gaugeMetricRepository   *memstorage.GaugeMetricRepository
	counterMetricRepository *memstorage.CounterMetricRepository
}

func NewMetricsRequestHandler(gaugeMetricRepository *memstorage.GaugeMetricRepository, counterMetricRepository *memstorage.CounterMetricRepository) *MetricsRequestHandler {
	return &MetricsRequestHandler{gaugeMetricRepository: gaugeMetricRepository, counterMetricRepository: counterMetricRepository}
}

type GaugeMetric struct {
	Name  string
	Value float64
}

func GaugeMetricFromModel(metric *domain.GaugeMetric) *GaugeMetric {
	return &GaugeMetric{
		Name:  metric.Name(),
		Value: metric.Value(),
	}
}

type CounterMetric struct {
	Name  string
	Value int64
}

func CounterMetricFromModel(metric *domain.CounterMetric) *CounterMetric {
	return &CounterMetric{
		Name:  metric.Name(),
		Value: metric.Value(),
	}
}

type CombinedMetrics struct {
	GaugeMetrics   map[string]*GaugeMetric
	CounterMetrics map[string]*CounterMetric
}

func (obj *MetricsRequestHandler) HandleRequest(response http.ResponseWriter, request *http.Request) {
	log.Println("Entering handler: MetricsRequestHandler")
	gaugeData := make(map[string]*GaugeMetric)
	counterData := make(map[string]*CounterMetric)

	for k, v := range obj.gaugeMetricRepository.FindAll() {
		gaugeData[k] = GaugeMetricFromModel(v)
	}
	for k, v := range obj.counterMetricRepository.FindAll() {
		counterData[k] = CounterMetricFromModel(v)
	}

	combinedMetrics := CombinedMetrics{
		GaugeMetrics:   gaugeData,
		CounterMetrics: counterData,
	}
	jsonData, err := json.Marshal(combinedMetrics)
	if err != nil {
		log.Fatalf("Error occurred during marshaling. Error: %s", err.Error())
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(jsonData)
}
