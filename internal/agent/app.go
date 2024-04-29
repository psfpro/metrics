package agent

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	mathRand "math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/psfpro/metrics/internal/agent/model"
)

type App struct {
	config    *Config
	pollCount int64
}

func NewApp(config *Config) *App {
	return &App{
		config: config,
	}
}

func (obj *App) Run() {
	collectJobs := make(chan int, obj.config.RateLimit)
	collectResults := make(chan []model.Metrics, obj.config.RateLimit)
	sendResults := make(chan error, obj.config.RateLimit)

	defer close(collectJobs)

	// создаем и запускаем 3 воркера, это и есть пул,
	// передаем id, это для наглядности, канал задач и канал результатов
	for w := 1; w <= 3; w++ {
		go obj.collect(w, collectJobs, collectResults)
	}
	for w := 1; w <= 3; w++ {
		go obj.send(w, collectResults, sendResults)
	}

	ticker := time.NewTicker(obj.config.PollInterval)
	for {
		select {
		case <-ticker.C:
			collectJobs <- 1
		case err := <-sendResults:
			if err != nil {
				log.Printf("Ошибка отправки метрик: %v", err)
			}
		}
	}
}

func (obj *App) sendBatchMetrics(batch []model.Metrics) error {
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

func (obj *App) gaugeMetric(name string, value *float64) model.Metrics {
	return model.Metrics{ID: name, MType: "gauge", Value: value}
}

func (obj *App) counterMetric(name string, value *int64) model.Metrics {
	return model.Metrics{ID: name, MType: "counter", Delta: value}
}

func (obj *App) sendBatch(metric []model.Metrics) (err error) {
	reqBytes, err := json.Marshal(metric)
	if err != nil {
		fmt.Println(err)
		return err
	}
	urlString := fmt.Sprintf("%s/updates", obj.config.ServerAddress)
	log.Println("crypto key: ", obj.config.CryptoKey)
	if obj.config.CryptoKey != "" {
		rsaPublicKey, err := loadRSAPublicKey(obj.config.CryptoKey)
		if err != nil {
			return err
		}
		// Generate a new AES key
		aesKey := make([]byte, aes.BlockSize)
		if _, err := rand.Read(aesKey); err != nil {
			return err
		}
		// Encrypt the data with AES
		blockCipher, err := aes.NewCipher(aesKey)
		if err != nil {
			return err
		}
		gcm, err := cipher.NewGCM(blockCipher)
		if err != nil {
			return err
		}
		nonce := make([]byte, gcm.NonceSize())
		if _, err := rand.Read(nonce); err != nil {
			return err
		}
		encryptedData := gcm.Seal(nil, nonce, reqBytes, nil)
		// Encrypt the AES key with RSA
		encryptedKey, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, aesKey)
		if err != nil {
			return err
		}
		reqBytes = append(encryptedData, encryptedKey...)
		reqBytes = append(reqBytes, nonce...)
	}
	body, err := obj.compress(reqBytes)
	if err != nil {
		return err
	}
	request, _ := http.NewRequest("POST", urlString, &body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Content-Encoding", "gzip")
	if obj.config.HashKey != "" {
		hash := hmac.New(sha256.New, []byte(obj.config.HashKey))
		hash.Write(reqBytes)
		signature := hex.EncodeToString(hash.Sum(nil))
		request.Header.Set("HashSHA256", signature)
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer func() {
		err = resp.Body.Close()
	}()
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

func (obj *App) collect(id int, jobs <-chan int, results chan<- []model.Metrics) {
	for j := range jobs {
		log.Println("collect", id, "start", j)
		metrics := obj.collectMetrics()
		metrics["RandomValue"] = mathRand.Float64()
		obj.pollCount++
		var batch []model.Metrics
		for name, value := range metrics {
			batch = append(batch, obj.gaugeMetric(name, &value))
		}
		batch = append(batch, obj.counterMetric("PollCount", &obj.pollCount))
		log.Println("collect", id, "end", j)
		results <- batch
	}
}
func (obj *App) psutil(id int, jobs <-chan int, results chan<- []model.Metrics) {
	for j := range jobs {
		log.Println("psutil", id, "start", j)
		v, _ := mem.VirtualMemory()
		totalMemory := float64(v.Total)
		freeMemory := float64(v.Free)
		cpuPercentages, _ := cpu.Percent(time.Second, true)
		var batch []model.Metrics
		batch = append(batch, obj.gaugeMetric("TotalMemory", &totalMemory))
		batch = append(batch, obj.gaugeMetric("FreeMemory", &freeMemory))
		for i, percent := range cpuPercentages {
			batch = append(batch, obj.gaugeMetric("CPUtilization"+strconv.Itoa(i), &percent))
		}
		log.Println("psutil", id, "end", j)
		results <- batch
	}
}

func (obj *App) send(id int, jobs <-chan []model.Metrics, results chan<- error) {
	dataForSend := make(map[string]model.Metrics)
	ticker := time.NewTicker(obj.config.ReportInterval)
	defer ticker.Stop()

	for {
		select {
		case j, ok := <-jobs:
			if !ok {
				log.Printf("send %d stopping\n", id)
				return
			}
			log.Printf("send %d starting task\n", id)
			for _, value := range j {
				dataForSend[value.ID] = value
			}
		case <-ticker.C:
			if len(dataForSend) == 0 {
				continue
			}
			log.Printf("send %d performing action\n", id)
			var values []model.Metrics
			for _, value := range dataForSend {
				values = append(values, value)
			}
			err := obj.sendBatchMetrics(values)
			if err == nil {
				dataForSend = make(map[string]model.Metrics)
			}

			results <- err
		}
	}
}

func loadRSAPublicKey(filename string) (*rsa.PublicKey, error) {
	// Read the file containing the public key
	keyData, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to read public key file: %w", err)
	}

	// Decode the PEM data
	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	// Parse the public key
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse public key: %w", err)
	}

	// Assert the type to *rsa.PublicKey
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}

	return rsaPub, nil
}
