package http

import (
	"github.com/psfpro/metrics/internal/server/infrastructure/api/http/handler"
	"net/http"
)

type App struct {
	config *Config
	router http.Handler
}

func NewApp(config *Config) *App {
	return &App{
		config: config,
		router: handler.Router(),
	}
}

func (obj *App) Run() {
	err := http.ListenAndServe(obj.config.Address, obj.router)
	if err != nil {
		panic(err)
	}
}
