package handler

import (
	"github.com/mailru/easyjson"
	"github.com/psfpro/metrics/internal/server/application"
	"github.com/psfpro/metrics/internal/server/infrastructure/api/http/model"
	"log"
	"net/http"
)

type UpdateRequestHandler struct {
	updateGaugeMetricHandler     *application.UpdateGaugeMetricHandler
	updateCounterMetricHandler   *application.UpdateCounterMetricHandler
	increaseCounterMetricHandler *application.IncreaseCounterMetricHandler
}

func NewUpdateRequestHandler(updateGaugeMetricHandler *application.UpdateGaugeMetricHandler, updateCounterMetricHandler *application.UpdateCounterMetricHandler, increaseCounterMetricHandler *application.IncreaseCounterMetricHandler) *UpdateRequestHandler {
	return &UpdateRequestHandler{updateGaugeMetricHandler: updateGaugeMetricHandler, updateCounterMetricHandler: updateCounterMetricHandler, increaseCounterMetricHandler: increaseCounterMetricHandler}
}

func (obj *UpdateRequestHandler) HandleRequest(response http.ResponseWriter, request *http.Request) {
	log.Println("Entering handler: UpdateRequestHandler")
	if request.Method == http.MethodPost {
		var metrics model.Metrics

		if err := easyjson.UnmarshalFromReader(request.Body, &metrics); err != nil {
			http.Error(response, err.Error(), http.StatusBadRequest)
			return
		}
		if metrics.MType == "gauge" {
			if metrics.ID == "" || metrics.Value == nil {
				response.WriteHeader(http.StatusNotFound)
				return
			}
			log.Printf("Update gauge %v: %v\n", metrics.ID, *metrics.Value)
			obj.updateGaugeMetricHandler.Handle(metrics.ID, *metrics.Value)
			response.WriteHeader(http.StatusOK)
			return
		} else if metrics.MType == "counter" {
			if metrics.ID == "" {
				response.WriteHeader(http.StatusNotFound)
				return
			}
			if metrics.Delta == nil {
				log.Printf("Increase counter %v\n", metrics.ID)
				obj.increaseCounterMetricHandler.Handle(metrics.ID)
			} else {
				log.Printf("Update counter %v: %v\n", metrics.ID, *metrics.Delta)
				obj.updateCounterMetricHandler.Handle(metrics.ID, *metrics.Delta)
			}
			response.WriteHeader(http.StatusOK)
			return
		}
	}
	response.WriteHeader(http.StatusNotFound)
}
