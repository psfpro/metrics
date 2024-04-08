package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mailru/easyjson"

	"github.com/psfpro/metrics/internal/server/domain"
	"github.com/psfpro/metrics/internal/server/infrastructure/api/http/model"
)

type GetRequestHandler struct {
	gaugeMetricRepository   domain.GaugeMetricRepository
	counterMetricRepository domain.CounterMetricRepository
}

func NewGetRequestHandler(gaugeMetricRepository domain.GaugeMetricRepository, counterMetricRepository domain.CounterMetricRepository) *GetRequestHandler {
	return &GetRequestHandler{gaugeMetricRepository: gaugeMetricRepository, counterMetricRepository: counterMetricRepository}
}

func (obj *GetRequestHandler) HandleRequest(response http.ResponseWriter, request *http.Request) {
	log.Println("Entering handler: GetRequestHandler")
	if request.Method == http.MethodPost {
		var metrics model.Metrics

		if err := easyjson.UnmarshalFromReader(request.Body, &metrics); err != nil {
			http.Error(response, err.Error(), http.StatusBadRequest)
			return
		}
		if metrics.ID == "" {
			response.WriteHeader(http.StatusNotFound)
			return
		}
		if metrics.MType == "gauge" {
			metric, exist := obj.gaugeMetricRepository.FindByName(metrics.ID)
			if exist {
				val := metric.Value()
				metrics.Value = &val
				jsonData, err := json.Marshal(metrics)
				log.Printf("Get %v %v: %v\n", "gauge", metrics.ID, jsonData)
				if err != nil {
					log.Fatalf("Error occurred during marshaling. Error: %s", err.Error())
				}

				response.Header().Set("Content-Type", "application/json")
				response.WriteHeader(http.StatusOK)
				response.Write(jsonData)
				return
			}
		} else if metrics.MType == "counter" {
			metric, exist := obj.counterMetricRepository.FindByName(metrics.ID)
			if exist {
				val := metric.Value()
				metrics.Delta = &val
				jsonData, err := json.Marshal(metrics)
				log.Printf("Get %v %v: %v\n", "counter", metrics.ID, jsonData)
				if err != nil {
					log.Fatalf("Error occurred during marshaling. Error: %s", err.Error())
				}

				response.Header().Set("Content-Type", "application/json")
				response.WriteHeader(http.StatusOK)
				response.Write(jsonData)
				return
			}
		}
		log.Printf("Get %v %v: %v\n", metrics.MType, metrics.ID, "NOT FOUND")
	}
	response.WriteHeader(http.StatusNotFound)
}
