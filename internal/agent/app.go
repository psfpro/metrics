package agent

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

type App struct {
	config *Config
}

func NewApp(config *Config) *App {
	return &App{
		config: config,
	}
}

func (obj *App) Run() {
	pollCount := int64(0)
	lastReportTime := time.Now()

	for {
		metrics := obj.collectMetrics()
		metrics["RandomValue"] = rand.Float64()
		pollCount++
		if time.Since(lastReportTime) >= obj.config.ReportInterval {
			fmt.Println("Отправка всех собранных метрик")
			for name, value := range metrics {
				obj.sendMetric("gauge", name, value)
			}

			fmt.Println("Отправка метрики PollCount")
			obj.sendMetric("counter", "PollCount", pollCount)
			lastReportTime = time.Now()
		}

		fmt.Printf("Ждем интервал для сбора %v\n", obj.config.PollInterval)
		time.Sleep(obj.config.PollInterval)
	}
}

func (obj *App) collectMetrics() map[string]float64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return map[string]float64{
		"Alloc":         float64(m.Alloc),
		"BuckHashSys":   float64(m.BuckHashSys),
		"Frees":         float64(m.Frees),
		"GCCPUFraction": m.GCCPUFraction,
		"GCSys":         float64(m.GCSys),
		"HeapAlloc":     float64(m.HeapAlloc),
		"HeapIdle":      float64(m.HeapIdle),
		"HeapInuse":     float64(m.HeapInuse),
		"HeapObjects":   float64(m.HeapObjects),
		"HeapReleased":  float64(m.HeapReleased),
		"HeapSys":       float64(m.HeapSys),
		"LastGC":        float64(m.LastGC),
		"Lookups":       float64(m.Lookups),
		"MCacheInuse":   float64(m.MCacheInuse),
		"MCacheSys":     float64(m.MCacheSys),
		"MSpanInuse":    float64(m.MSpanInuse),
		"MSpanSys":      float64(m.MSpanSys),
		"Mallocs":       float64(m.Mallocs),
		"NextGC":        float64(m.NextGC),
		"NumForcedGC":   float64(m.NumForcedGC),
		"NumGC":         float64(m.NumGC),
		"OtherSys":      float64(m.OtherSys),
		"PauseTotalNs":  float64(m.PauseTotalNs),
		"StackInuse":    float64(m.StackInuse),
		"StackSys":      float64(m.StackSys),
		"Sys":           float64(m.Sys),
		"TotalAlloc":    float64(m.TotalAlloc),
	}
}

func (obj *App) sendMetric(metricType, name string, value interface{}) {
	urlString := fmt.Sprintf("%s/update/%s/%s/%v", obj.config.ServerAddress, metricType, name, value)
	resp, err := http.Post(urlString, "text/plain", nil)
	if err != nil {
		fmt.Printf("Error sending metric: %s\n", err)
		return
	}
	defer resp.Body.Close()
}
