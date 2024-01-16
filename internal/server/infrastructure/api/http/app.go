package http

import (
	"log"
	"net/http"
)

type App struct {
	addr   string
	router http.Handler
}

func NewApp(addr string, router http.Handler) *App {
	return &App{
		addr:   addr,
		router: router,
	}
}

func (obj *App) Run() {
	log.Printf("Start server addr: %v", obj.addr)
	err := http.ListenAndServe(obj.addr, obj.router)
	if err != nil {
		panic(err)
	}
}
