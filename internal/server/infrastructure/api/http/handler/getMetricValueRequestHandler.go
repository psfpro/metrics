package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/psfpro/metrics/internal/server/domain"
	"log"
	"net/http"
	"strconv"
)

type GetMetricValueRequestHandler struct {
	gaugeMetricRepository   domain.GaugeMetricRepository
	counterMetricRepository domain.CounterMetricRepository
}

func NewGetMetricValueRequestHandler(gaugeMetricRepository domain.GaugeMetricRepository, counterMetricRepository domain.CounterMetricRepository) *GetMetricValueRequestHandler {
	return &GetMetricValueRequestHandler{gaugeMetricRepository: gaugeMetricRepository, counterMetricRepository: counterMetricRepository}
}

func (obj *GetMetricValueRequestHandler) HandleRequest(response http.ResponseWriter, request *http.Request) {
	log.Println("Entering handler: GetMetricValueRequestHandler")
	metricType := chi.URLParam(request, "type")
	metricName := chi.URLParam(request, "name")

	if metricType == "counter" && metricName != "" {
		metric, exist := obj.counterMetricRepository.FindByName(metricName)
		if exist {
			body := strconv.FormatInt(metric.Value(), 10)
			log.Printf("Get %v %v: %v\n", metricType, metricName, body)
			response.Write([]byte(body))
			return
		}
	} else if metricType == "gauge" && metricName != "" {
		metric, exist := obj.gaugeMetricRepository.FindByName(metricName)
		if exist {
			body := strconv.FormatFloat(metric.Value(), 'f', -1, 64)
			log.Printf("Get %v %v: %v\n", metricType, metricName, body)
			response.Write([]byte(body))
			return
		}
	}

	response.WriteHeader(http.StatusNotFound)
}
