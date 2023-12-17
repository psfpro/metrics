package handler

import "net/http"

type MetricsListHandler struct {
}

func NewMetricsListHandler() *MetricsListHandler {
	return &MetricsListHandler{}
}

func (obj *MetricsListHandler) HandleRequest(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusBadRequest)
}
