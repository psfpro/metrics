package main

import "github.com/psfpro/metrics/internal/infrastructure/api/http"

func main() {
	app := http.NewApp(&http.Config{
		Address: ":8080",
	})
	app.Run()
}
