package storage

import "github.com/psfpro/metrics/internal/server/domain"

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
