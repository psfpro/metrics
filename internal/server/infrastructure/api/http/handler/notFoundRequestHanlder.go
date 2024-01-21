package handler

import (
	"log"
	"net/http"
)

type NotFoundRequestHandler struct {
}

func NewNotFoundRequestHandler() *NotFoundRequestHandler {
	return &NotFoundRequestHandler{}
}

func (obj *NotFoundRequestHandler) HandleRequest(response http.ResponseWriter, request *http.Request) {
	log.Println("Entering handler: MetricsListHandler")
	response.WriteHeader(http.StatusNotFound)
}
