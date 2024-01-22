package storage

import (
	"bytes"
	"context"
	"os"

	"encoding/json"

	"github.com/psfpro/metrics/internal/server/domain"
)

type FileAdapter struct {
	file                    *os.File
	counterMetricRepository *CounterMetricRepository
	gaugeMetricRepository   *GaugeMetricRepository
}

func NewFileAdapter(file *os.File, counterMetricRepository *CounterMetricRepository, gaugeMetricRepository *GaugeMetricRepository) *FileAdapter {
	return &FileAdapter{file: file, counterMetricRepository: counterMetricRepository, gaugeMetricRepository: gaugeMetricRepository}
}

func (obj *FileAdapter) Flush(ctx context.Context) error {
	gaugeData := make(map[string]*GaugeMetric)
	counterData := make(map[string]*CounterMetric)

	for k, v := range obj.gaugeMetricRepository.data {
		gaugeData[k] = GaugeMetricFromModel(v)
	}
	for k, v := range obj.counterMetricRepository.data {
		counterData[k] = CounterMetricFromModel(v)
	}

	combinedMetrics := CombinedMetrics{
		GaugeMetrics:   gaugeData,
		CounterMetrics: counterData,
	}
	jsonData, err := json.Marshal(combinedMetrics)
	if err != nil {
		return err
	}

	err = obj.file.Truncate(0)
	if err != nil {
		return err
	}
	_, err = obj.file.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = obj.file.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

func (obj *FileAdapter) Restore(ctx context.Context) error {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(obj.file)
	if err != nil {
		return err
	}
	var combinedMetrics CombinedMetrics

	if err = json.Unmarshal(buf.Bytes(), &combinedMetrics); err != nil {
		return err
	}

	obj.gaugeMetricRepository.data = make(map[string]*domain.GaugeMetric)
	obj.counterMetricRepository.data = make(map[string]*domain.CounterMetric)
	for k, v := range combinedMetrics.GaugeMetrics {
		obj.gaugeMetricRepository.data[k] = domain.NewGaugeMetric(v.Name, v.Value)
	}
	for k, v := range combinedMetrics.CounterMetrics {
		obj.counterMetricRepository.data[k] = domain.NewCounterMetric(v.Name, v.Value)
	}

	return nil
}

// Mapping structures

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
