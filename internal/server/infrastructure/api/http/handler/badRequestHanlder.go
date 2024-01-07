package handler

import (
	"log"
	"net/http"
)

type BadRequestHandler struct {
}

func NewBadRequestHandler() *BadRequestHandler {
	return &BadRequestHandler{}
}

func (obj *BadRequestHandler) HandleRequest(response http.ResponseWriter, request *http.Request) {
	log.Println("Entering handler: BadRequestHandler")
	response.WriteHeader(http.StatusBadRequest)
}
