package memstorage

import (
	"github.com/psfpro/metrics/internal/server/domain"
)

type GaugeMetricRepository struct {
	Data map[string]*domain.GaugeMetric
}

func NewGaugeMetricRepository() *GaugeMetricRepository {
	return &GaugeMetricRepository{Data: make(map[string]*domain.GaugeMetric)}
}

func (obj GaugeMetricRepository) FindByName(name string) (*domain.GaugeMetric, bool) {
	result := obj.Data[name]
	if result == nil {
		return nil, false
	}

	return result, true
}
func (obj GaugeMetricRepository) Add(metric *domain.GaugeMetric) {
	obj.Data[metric.Name()] = metric
}

type CounterMetricRepository struct {
	Data map[string]*domain.CounterMetric
}

func NewCounterMetricRepository() *CounterMetricRepository {
	return &CounterMetricRepository{Data: make(map[string]*domain.CounterMetric)}
}

func (obj CounterMetricRepository) FindByName(name string) (*domain.CounterMetric, bool) {
	result := obj.Data[name]
	if result == nil {
		return nil, false
	}

	return result, true
}
func (obj CounterMetricRepository) Add(metric *domain.CounterMetric) {
	obj.Data[metric.Name()] = metric
}
