package main

import (
	"github.com/psfpro/metrics/internal/server/infrastructure/api/http"
)

func main() {
	app := http.NewApp(&http.Config{
		Address: ":8080",
	})
	app.Run()
}
