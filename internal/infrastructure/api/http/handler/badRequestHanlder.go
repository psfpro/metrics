package handler

import "net/http"

type BadRequestHandler struct {
}

func NewBadRequestHandler() *BadRequestHandler {
	return &BadRequestHandler{}
}

func (obj *BadRequestHandler) HandleRequest(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusBadRequest)
}
