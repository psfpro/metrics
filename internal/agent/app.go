package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/psfpro/metrics/internal/agent/model"
	"log"
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
			log.Println("Отправка всех собранных метрик")
			for name, value := range metrics {
				obj.sendGaugeMetric(name, &value)
			}

			log.Println("Отправка метрики PollCount")
			obj.sendCounterMetric("PollCount", &pollCount)
			lastReportTime = time.Now()
		}

		log.Printf("Ждем интервал для сбора %v\n", obj.config.PollInterval)
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
		log.Printf("Error sending metric: %s\n", err)
		return
	}
	defer resp.Body.Close()
}

func (obj *App) sendGaugeMetric(name string, value *float64) {
	urlString := fmt.Sprintf("%s/update", obj.config.ServerAddress)
	metric := model.Metrics{ID: name, MType: "gauge", Value: value}
	reqBytes, err := json.Marshal(metric)
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := http.Post(urlString, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		log.Printf("Error sending metric: %s\n", err)
		return
	}
	defer resp.Body.Close()
}

func (obj *App) sendCounterMetric(name string, value *int64) {
	urlString := fmt.Sprintf("%s/update", obj.config.ServerAddress)
	metric := model.Metrics{ID: name, MType: "counter", Delta: value}
	reqBytes, err := json.Marshal(metric)
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := http.Post(urlString, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		log.Printf("Error sending metric: %s\n", err)
		return
	}
	defer resp.Body.Close()
}
