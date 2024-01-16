package filestorage

import (
	"github.com/psfpro/metrics/internal/server/domain"
)

type GaugeMetricRepository struct {
	entityManager *EntityManager
}

func NewGaugeMetricRepository(entityManager *EntityManager) *GaugeMetricRepository {
	return &GaugeMetricRepository{entityManager: entityManager}
}

func (obj *GaugeMetricRepository) FindAll() map[string]*domain.GaugeMetric {
	return obj.entityManager.findAllGaugeMetrics()
}

func (obj *GaugeMetricRepository) FindByName(name string) (*domain.GaugeMetric, bool) {
	return obj.entityManager.findGaugeMetric(name)
}
func (obj *GaugeMetricRepository) Add(metric *domain.GaugeMetric) {
	obj.entityManager.persistGaugeMetric(metric)
}

type CounterMetricRepository struct {
	entityManager *EntityManager
}

func NewCounterMetricRepository(entityManager *EntityManager) *CounterMetricRepository {
	return &CounterMetricRepository{entityManager: entityManager}
}

func (obj *CounterMetricRepository) FindAll() map[string]*domain.CounterMetric {
	return obj.entityManager.findAllCounterMetrics()
}

func (obj *CounterMetricRepository) FindByName(name string) (*domain.CounterMetric, bool) {
	return obj.entityManager.findCounterMetric(name)
}
func (obj *CounterMetricRepository) Add(metric *domain.CounterMetric) {
	obj.entityManager.persistCounterMetric(metric)
}
