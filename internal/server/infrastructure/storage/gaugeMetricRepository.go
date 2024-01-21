package storage

import "github.com/psfpro/metrics/internal/server/domain"

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
