package filestorage

import (
	"bytes"
	"encoding/json"
	"github.com/psfpro/metrics/internal/server/domain"
	"log"
	"os"
)

type EntityManager struct {
	file              *os.File
	dataGaugeMetric   map[string]*domain.GaugeMetric
	dataCounterMetric map[string]*domain.CounterMetric
}

func NewEntityManager(filePath string) *EntityManager {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	em := &EntityManager{
		file:              file,
		dataGaugeMetric:   make(map[string]*domain.GaugeMetric),
		dataCounterMetric: make(map[string]*domain.CounterMetric),
	}
	em.Restore()

	return em
}

func (obj *EntityManager) Flush() {
	gaugeData := make(map[string]*GaugeMetric)
	counterData := make(map[string]*CounterMetric)

	for k, v := range obj.dataGaugeMetric {
		gaugeData[k] = GaugeMetricFromModel(v)
	}
	for k, v := range obj.dataCounterMetric {
		counterData[k] = CounterMetricFromModel(v)
	}

	combinedMetrics := CombinedMetrics{
		GaugeMetrics:   gaugeData,
		CounterMetrics: counterData,
	}
	jsonData, err := json.Marshal(combinedMetrics)
	if err != nil {
		log.Fatalf("Error occurred during marshaling. Error: %s", err.Error())
	}

	obj.file.Truncate(0)
	obj.file.Seek(0, 0)
	obj.file.Write(jsonData)
}

func (obj *EntityManager) Restore() {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(obj.file)
	if err != nil {
		log.Printf("Error occurred during read file. Error: %s", err.Error())
		return
	}
	var combinedMetrics CombinedMetrics

	if err = json.Unmarshal(buf.Bytes(), &combinedMetrics); err != nil {
		log.Printf("Error occurred during unmarshaling. Error: %s", err.Error())
		return
	}

	obj.dataGaugeMetric = make(map[string]*domain.GaugeMetric)
	obj.dataCounterMetric = make(map[string]*domain.CounterMetric)
	for k, v := range combinedMetrics.GaugeMetrics {
		obj.dataGaugeMetric[k] = domain.NewGaugeMetric(v.Name, v.Value)
	}
	for k, v := range combinedMetrics.CounterMetrics {
		obj.dataCounterMetric[k] = domain.NewCounterMetric(v.Name, v.Value)
	}
}

func (obj *EntityManager) findAllGaugeMetrics() map[string]*domain.GaugeMetric {
	return obj.dataGaugeMetric
}

func (obj *EntityManager) findGaugeMetric(name string) (*domain.GaugeMetric, bool) {
	result := obj.dataGaugeMetric[name]
	if result == nil {
		return nil, false
	}

	return result, true
}

func (obj *EntityManager) persistGaugeMetric(metric *domain.GaugeMetric) {
	obj.dataGaugeMetric[metric.Name()] = metric
}

func (obj *EntityManager) findAllCounterMetrics() map[string]*domain.CounterMetric {
	return obj.dataCounterMetric
}

func (obj *EntityManager) findCounterMetric(name string) (*domain.CounterMetric, bool) {
	result := obj.dataCounterMetric[name]
	if result == nil {
		return nil, false
	}

	return result, true
}

func (obj *EntityManager) persistCounterMetric(metric *domain.CounterMetric) {
	obj.dataCounterMetric[metric.Name()] = metric
}

// Mapping classes

type GaugeMetric struct {
	Name  string
	Value float64
}

func GaugeMetricFromModel(metric *domain.GaugeMetric) *GaugeMetric {
	return &GaugeMetric{
		Name:  metric.Name(),
		Value: metric.Value(),
	}
}

type CounterMetric struct {
	Name  string
	Value int64
}

func CounterMetricFromModel(metric *domain.CounterMetric) *CounterMetric {
	return &CounterMetric{
		Name:  metric.Name(),
		Value: metric.Value(),
	}
}

type CombinedMetrics struct {
	GaugeMetrics   map[string]*GaugeMetric
	CounterMetrics map[string]*CounterMetric
}
