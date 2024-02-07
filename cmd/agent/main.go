package main

import (
	"flag"
	"github.com/psfpro/metrics/internal/agent"
	"os"
	"strconv"
	"time"
)

func main() {
	parseFlags()
	app := agent.NewApp(&agent.Config{
		HashKey:        hashKey,
		ServerAddress:  "http://" + flagRunAddr,
		PollInterval:   time.Duration(flagPollInterval) * time.Second,
		ReportInterval: time.Duration(flagReportInterval) * time.Second,
	})
	app.Run()
}

// Не экспортированная переменная flagRunAddr содержит адрес и порт для запуска сервера
var hashKey string
var flagRunAddr string
var flagReportInterval int
var flagPollInterval int

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func parseFlags() {
	flag.StringVar(&hashKey, "k", "", "hash key")
	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
	flag.IntVar(&flagReportInterval, "r", 10, "frequency of sending metrics to the server")
	flag.IntVar(&flagPollInterval, "p", 2, "metrics polling rate")
	flag.Parse()
	if envHashKey := os.Getenv("KEY"); envHashKey != "" {
		hashKey = envHashKey
	}
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		flagReportInterval, _ = strconv.Atoi(envReportInterval)
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		flagPollInterval, _ = strconv.Atoi(envPollInterval)
	}
}
