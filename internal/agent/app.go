package agent

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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
			err := obj.sendBatchMetrics(metrics, pollCount)
			if err != nil {
				log.Printf("Ошибка отправки метрик: %v", err)
			}
			lastReportTime = time.Now()
		}

		log.Printf("Ждем интервал для сбора %v\n", obj.config.PollInterval)
		time.Sleep(obj.config.PollInterval)
	}
}

func (obj *App) sendMetrics(metrics map[string]float64, pollCount int64) {
	log.Println("Отправка всех собранных метрик")
	for name, value := range metrics {
		metric := obj.gaugeMetric(name, &value)
		obj.send(metric)
	}

	log.Println("Отправка метрики PollCount")
	metric := obj.counterMetric("PollCount", &pollCount)
	obj.send(metric)
}

func (obj *App) sendBatchMetrics(metrics map[string]float64, pollCount int64) error {
	var batch []model.Metrics
	log.Println("Отправка всех собранных метрик")
	for name, value := range metrics {
		batch = append(batch, obj.gaugeMetric(name, &value))
	}
	batch = append(batch, obj.counterMetric("PollCount", &pollCount))

	var err error
	retryDelays := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}

	for _, delay := range retryDelays {
		err = obj.sendBatch(batch)
		if err == nil {
			return nil
		}

		time.Sleep(delay)
	}

	return fmt.Errorf("после нескольких попыток: %w", err)
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

func (obj *App) gaugeMetric(name string, value *float64) model.Metrics {
	return model.Metrics{ID: name, MType: "gauge", Value: value}
}

func (obj *App) counterMetric(name string, value *int64) model.Metrics {
	return model.Metrics{ID: name, MType: "counter", Delta: value}
}

func (obj *App) send(metric model.Metrics) {
	reqBytes, err := json.Marshal(metric)
	if err != nil {
		fmt.Println(err)
		return
	}
	urlString := fmt.Sprintf("%s/update", obj.config.ServerAddress)
	body, err := obj.compress(reqBytes)
	if err != nil {
		log.Printf("Error compress metric: %s\n", err)
		return
	}
	request, _ := http.NewRequest("POST", urlString, &body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Content-Encoding", "gzip")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("Error sending metric: %s\n", err)
		return
	}
	defer resp.Body.Close()
}

func (obj *App) sendBatch(metric []model.Metrics) error {
	reqBytes, err := json.Marshal(metric)
	if err != nil {
		fmt.Println(err)
		return err
	}
	urlString := fmt.Sprintf("%s/updates", obj.config.ServerAddress)
	body, err := obj.compress(reqBytes)
	if err != nil {
		return err
	}
	request, _ := http.NewRequest("POST", urlString, &body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Content-Encoding", "gzip")
	if obj.config.HashKey != "" {
		hmac := hmac.New(sha256.New, []byte(obj.config.HashKey))
		hmac.Write(reqBytes)
		signature := hex.EncodeToString(hmac.Sum(nil))
		request.Header.Set("HashSHA256", signature)
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (obj *App) compress(data []byte) (bytes.Buffer, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)

	_, err := w.Write(data)
	if err != nil {
		return b, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}
	err = w.Close()
	if err != nil {
		return b, fmt.Errorf("failed compress data: %v", err)
	}

	return b, nil
}
