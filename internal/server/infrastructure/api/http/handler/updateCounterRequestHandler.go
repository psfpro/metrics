package handler

import (
	"fmt"
	"github.com/psfpro/metrics/internal/server/application"
	"net/http"
	"strconv"
	"strings"
)

type UpdateCounterRequestHandler struct {
	updateCounterMetricHandler   *application.UpdateCounterMetricHandler
	increaseCounterMetricHandler *application.IncreaseCounterMetricHandler
}

func NewUpdateCounterRequestHandler(updateCounterMetricHandler *application.UpdateCounterMetricHandler, increaseCounterMetricHandler *application.IncreaseCounterMetricHandler) *UpdateCounterRequestHandler {
	return &UpdateCounterRequestHandler{updateCounterMetricHandler: updateCounterMetricHandler, increaseCounterMetricHandler: increaseCounterMetricHandler}
}

func (obj *UpdateCounterRequestHandler) HandleRequest(response http.ResponseWriter, request *http.Request) {
	parts := strings.Split(request.RequestURI, "/")
	if len(parts) == 4 && parts[3] != "" {
		name := parts[3]
		fmt.Printf("Increase counter %v\n", name)
		obj.increaseCounterMetricHandler.Handle(name)
	} else if len(parts) == 5 && parts[3] != "" && parts[4] != "" {
		name := parts[3]
		value, err := strconv.ParseInt(parts[4], 10, 64)
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Printf("Update counter %v: %v\n", name, value)
		obj.updateCounterMetricHandler.Handle(name, value)
	} else {
		response.WriteHeader(http.StatusNotFound)
	}
}
