package collector

import (
	"github.com/psfpro/metrics/internal/agent/model"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"log"
	mathRand "math/rand"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type Worker struct {
	pollCount    int64
	collectGroup sync.WaitGroup
}

func NewWorker() *Worker {
	return &Worker{}
}

func (w *Worker) Run(collectJobs chan int, collectResults chan []model.Metrics) {
	for i := 1; i <= 3; i++ {
		w.collectGroup.Add(1)
		go w.collect(i, collectJobs, collectResults)
	}
	w.collectGroup.Wait()
	close(collectResults)
}

func (w *Worker) collectMetrics() map[string]float64 {
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

func (w *Worker) gaugeMetric(name string, value *float64) model.Metrics {
	return model.Metrics{ID: name, MType: "gauge", Value: value}
}

func (w *Worker) counterMetric(name string, value *int64) model.Metrics {
	return model.Metrics{ID: name, MType: "counter", Delta: value}
}

func (w *Worker) collect(id int, jobs <-chan int, results chan<- []model.Metrics) {
	for range jobs {
		log.Println("collect", id)
		metrics := w.collectMetrics()
		metrics["RandomValue"] = mathRand.Float64()
		w.pollCount++
		var batch []model.Metrics
		for name, value := range metrics {
			batch = append(batch, w.gaugeMetric(name, &value))
		}
		batch = append(batch, w.counterMetric("PollCount", &w.pollCount))
		results <- batch
	}
	w.collectGroup.Done()
	log.Printf("collect %d stopping\n", id)
}

func (w *Worker) psutil(id int, jobs <-chan int, results chan<- []model.Metrics) {
	for j := range jobs {
		log.Println("psutil", id, "start", j)
		v, _ := mem.VirtualMemory()
		totalMemory := float64(v.Total)
		freeMemory := float64(v.Free)
		cpuPercentages, _ := cpu.Percent(time.Second, true)
		var batch []model.Metrics
		batch = append(batch, w.gaugeMetric("TotalMemory", &totalMemory))
		batch = append(batch, w.gaugeMetric("FreeMemory", &freeMemory))
		for i, percent := range cpuPercentages {
			batch = append(batch, w.gaugeMetric("CPUtilization"+strconv.Itoa(i), &percent))
		}
		log.Println("psutil", id, "end", j)
		results <- batch
	}
}
