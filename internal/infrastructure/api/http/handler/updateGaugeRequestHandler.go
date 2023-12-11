package handler

import (
	"github.com/psfpro/metrics/internal/application"
	"net/http"
	"strconv"
	"strings"
)

type UpdateGaugeRequestHandler struct {
	updateGaugeMetricHandler *application.UpdateGaugeMetricHandler
}

func NewUpdateGaugeRequestHandler(updateGaugeMetricHandler *application.UpdateGaugeMetricHandler) *UpdateGaugeRequestHandler {
	return &UpdateGaugeRequestHandler{updateGaugeMetricHandler: updateGaugeMetricHandler}
}

func (obj *UpdateGaugeRequestHandler) HandleRequest(response http.ResponseWriter, request *http.Request) {
	parts := strings.Split(request.RequestURI, "/")
	if len(parts) == 5 && parts[3] != "" && parts[4] != "" {
		name := parts[3]
		value, err := strconv.ParseFloat(parts[4], 64)
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			return
		}
		obj.updateGaugeMetricHandler.Handle(name, value)
	} else {
		response.WriteHeader(http.StatusNotFound)
	}
}
