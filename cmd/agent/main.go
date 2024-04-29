package main

import (
	"fmt"
	"github.com/psfpro/metrics/internal/agent"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	fmt.Printf(
		"Build version: %s\nBuild date: %s\nBuild commit: %s\n\n",
		getString(buildVersion), getString(buildDate), getString(buildCommit),
	)
	app := agent.NewApp(agent.NewConfig())
	app.Run()
}

func getString(s string) string {
	if s == "" {
		return "N/A"
	}
	return s
}
