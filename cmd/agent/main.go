package main

import (
	"github.com/psfpro/metrics/internal/agent"
	"time"
)

func main() {
	parseFlags()
	app := agent.NewApp(&agent.Config{
		ServerAddress:  "http://" + flagRunAddr,
		PollInterval:   time.Duration(flagPollInterval) * time.Second,
		ReportInterval: time.Duration(flagReportInterval) * time.Second,
	})
	app.Run()
}
