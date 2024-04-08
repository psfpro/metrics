package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/psfpro/metrics/internal/server/application"
)

type UpdateGaugeRequestHandler struct {
	updateGaugeMetricHandler *application.UpdateGaugeMetricHandler
}

func NewUpdateGaugeRequestHandler(updateGaugeMetricHandler *application.UpdateGaugeMetricHandler) *UpdateGaugeRequestHandler {
	return &UpdateGaugeRequestHandler{updateGaugeMetricHandler: updateGaugeMetricHandler}
}

func (obj *UpdateGaugeRequestHandler) HandleRequest(response http.ResponseWriter, request *http.Request) {
	log.Println("Entering handler: UpdateGaugeRequestHandler")
	name := chi.URLParam(request, "name")
	value := chi.URLParam(request, "value")
	if name == "" || value == "" {
		response.WriteHeader(http.StatusNotFound)
		return
	}
	valueFloat, err := strconv.ParseFloat(value, 64)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("Update gauge %v: %v\n", name, valueFloat)
	obj.updateGaugeMetricHandler.Handle(name, valueFloat)
}
