package main

import (
	"flag"
)

// Не экспортированная переменная flagRunAddr содержит адрес и порт для запуска сервера
var flagRunAddr string
var flagReportInterval int
var flagPollInterval int

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
	flag.IntVar(&flagReportInterval, "r", 10, "frequency of sending metrics to the server")
	flag.IntVar(&flagPollInterval, "p", 2, "metrics polling rate")
	flag.Parse()
}
