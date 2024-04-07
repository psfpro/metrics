package handler

import (
	"log"
	"net/http"

	"github.com/mailru/easyjson"

	"github.com/psfpro/metrics/internal/server/application"
	"github.com/psfpro/metrics/internal/server/infrastructure/api/http/model"
)

type UpdatesRequestHandler struct {
	updateGaugeMetricHandler     *application.UpdateGaugeMetricHandler
	updateCounterMetricHandler   *application.UpdateCounterMetricHandler
	increaseCounterMetricHandler *application.IncreaseCounterMetricHandler
}

func NewUpdatesRequestHandler(updateGaugeMetricHandler *application.UpdateGaugeMetricHandler, updateCounterMetricHandler *application.UpdateCounterMetricHandler, increaseCounterMetricHandler *application.IncreaseCounterMetricHandler) *UpdatesRequestHandler {
	return &UpdatesRequestHandler{updateGaugeMetricHandler: updateGaugeMetricHandler, updateCounterMetricHandler: updateCounterMetricHandler, increaseCounterMetricHandler: increaseCounterMetricHandler}
}

func (obj *UpdatesRequestHandler) HandleRequest(response http.ResponseWriter, request *http.Request) {
	log.Println("Entering handler: UpdatesRequestHandler")
	if request.Method == http.MethodPost {
		var batch model.MetricsSlice

		if err := easyjson.UnmarshalFromReader(request.Body, &batch); err != nil {
			http.Error(response, err.Error(), http.StatusBadRequest)
			return
		}
		for _, metrics := range batch {
			if metrics.MType == "gauge" {
				if metrics.ID == "" || metrics.Value == nil {
					continue
				}
				log.Printf("Update gauge %v: %v\n", metrics.ID, *metrics.Value)
				obj.updateGaugeMetricHandler.Handle(metrics.ID, *metrics.Value)
			} else if metrics.MType == "counter" {
				if metrics.ID == "" {
					continue
				}
				if metrics.Delta == nil {
					log.Printf("Increase counter %v\n", metrics.ID)
					obj.increaseCounterMetricHandler.Handle(metrics.ID)
				} else {
					log.Printf("Update counter %v: %v\n", metrics.ID, *metrics.Delta)
					obj.updateCounterMetricHandler.Handle(metrics.ID, *metrics.Delta)
				}
			}
		}
	}
	response.WriteHeader(http.StatusOK)
}
