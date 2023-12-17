package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/psfpro/metrics/internal/server/application"
	"log"
	"net/http"
	"strconv"
)

type UpdateCounterRequestHandler struct {
	updateCounterMetricHandler   *application.UpdateCounterMetricHandler
	increaseCounterMetricHandler *application.IncreaseCounterMetricHandler
}

func NewUpdateCounterRequestHandler(updateCounterMetricHandler *application.UpdateCounterMetricHandler, increaseCounterMetricHandler *application.IncreaseCounterMetricHandler) *UpdateCounterRequestHandler {
	return &UpdateCounterRequestHandler{updateCounterMetricHandler: updateCounterMetricHandler, increaseCounterMetricHandler: increaseCounterMetricHandler}
}

func (obj *UpdateCounterRequestHandler) HandleRequest(response http.ResponseWriter, request *http.Request) {
	name := chi.URLParam(request, "name")
	value := chi.URLParam(request, "value")
	if name == "" {
		response.WriteHeader(http.StatusNotFound)
		return
	}
	if value == "" {
		log.Printf("Increase counter %v\n", name)
		obj.increaseCounterMetricHandler.Handle(name)
	} else {
		valueInt, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Update counter %v: %v\n", name, valueInt)
		obj.updateCounterMetricHandler.Handle(name, valueInt)
	}
}
