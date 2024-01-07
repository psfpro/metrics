package handler

import (
	"log"
	"net/http"
)

type MetricsListHandler struct {
}

func NewMetricsListHandler() *MetricsListHandler {
	return &MetricsListHandler{}
}

func (obj *MetricsListHandler) HandleRequest(response http.ResponseWriter, request *http.Request) {
	log.Println("Entering handler: MetricsListHandler")
	response.WriteHeader(http.StatusBadRequest)
}
