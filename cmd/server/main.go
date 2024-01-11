package main

import (
	"github.com/psfpro/metrics/internal/server"
)

func main() {
	container := server.NewContainer()
	app := container.App()
	app.Run()
}
