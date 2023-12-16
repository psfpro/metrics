package main

import (
	"github.com/psfpro/metrics/internal/agent"
	"time"
)

func main() {
	app := agent.NewApp(&agent.Config{
		ServerAddress:  "http://localhost:8080",
		PollInterval:   2 * time.Second,
		ReportInterval: 10 * time.Second,
	})
	app.Run()
}
