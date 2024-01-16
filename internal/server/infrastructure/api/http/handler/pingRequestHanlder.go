package handler

import (
	"database/sql"
	"log"
	"net/http"
)

type PingRequestHandler struct {
	db *sql.DB
}

func NewPingRequestHandler(db *sql.DB) *PingRequestHandler {
	return &PingRequestHandler{db: db}
}

func (obj *PingRequestHandler) HandleRequest(response http.ResponseWriter, request *http.Request) {
	log.Println("Entering handler: PingRequestHandler")

	if err := obj.db.Ping(); err != nil {
		log.Printf("db connection error: %v", err)
		response.WriteHeader(http.StatusInternalServerError)
	}

	response.WriteHeader(http.StatusOK)
}
