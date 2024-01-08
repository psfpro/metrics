package memstorage

import (
	"github.com/psfpro/metrics/internal/server/domain"
)

type GaugeMetricRepository struct {
	data map[string]*domain.GaugeMetric
}

func NewGaugeMetricRepository() *GaugeMetricRepository {
	return &GaugeMetricRepository{data: make(map[string]*domain.GaugeMetric)}
}

func (obj *GaugeMetricRepository) FindAll() map[string]*domain.GaugeMetric {
	return obj.data
}

func (obj *GaugeMetricRepository) FindByName(name string) (*domain.GaugeMetric, bool) {
	result := obj.data[name]
	if result == nil {
		return nil, false
	}

	return result, true
}
func (obj *GaugeMetricRepository) Add(metric *domain.GaugeMetric) {
	obj.data[metric.Name()] = metric
}

type CounterMetricRepository struct {
	data map[string]*domain.CounterMetric
}

func NewCounterMetricRepository() *CounterMetricRepository {
	return &CounterMetricRepository{data: make(map[string]*domain.CounterMetric)}
}

func (obj *CounterMetricRepository) FindAll() map[string]*domain.CounterMetric {
	return obj.data
}

func (obj *CounterMetricRepository) FindByName(name string) (*domain.CounterMetric, bool) {
	result := obj.data[name]
	if result == nil {
		return nil, false
	}

	return result, true
}
func (obj *CounterMetricRepository) Add(metric *domain.CounterMetric) {
	obj.data[metric.Name()] = metric
}
