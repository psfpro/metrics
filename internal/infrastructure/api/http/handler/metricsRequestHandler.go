package handler

import (
	"fmt"
	"github.com/psfpro/metrics/internal/infrastructure/storage/memstorage"
	"net/http"
)

type MetricsRequestHandler struct {
	gaugeMetricRepository   *memstorage.GaugeMetricRepository
	counterMetricRepository *memstorage.CounterMetricRepository
}

func NewMetricsRequestHandler(gaugeMetricRepository *memstorage.GaugeMetricRepository, counterMetricRepository *memstorage.CounterMetricRepository) *MetricsRequestHandler {
	return &MetricsRequestHandler{gaugeMetricRepository: gaugeMetricRepository, counterMetricRepository: counterMetricRepository}
}

func (obj *MetricsRequestHandler) HandleRequest(response http.ResponseWriter, request *http.Request) {
	body := fmt.Sprintf("Method: %s\r\n", request.Method)
	body += "Gauge metrics =================\r\n"
	for k, v := range obj.gaugeMetricRepository.Data {
		body += fmt.Sprintf("%s: %v\r\n", k, v)
	}
	body += "Counter metrics ===============\r\n"
	for k, v := range obj.counterMetricRepository.Data {
		body += fmt.Sprintf("%s: %v\r\n", k, v)
	}
	response.Write([]byte(body))
}
