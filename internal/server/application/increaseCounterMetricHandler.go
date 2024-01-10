package application

import "github.com/psfpro/metrics/internal/server/domain"

// IncreaseCounterMetricHandler Обновление counter метрики
type IncreaseCounterMetricHandler struct {
	Repository domain.CounterMetricRepository
}

func (obj IncreaseCounterMetricHandler) Handle(name string) {
	metric, exist := obj.Repository.FindByName(name)
	if exist {
		metric.Increase()
	} else {
		metric = domain.NewCounterMetric(name, 0)
		obj.Repository.Add(metric)
	}
}
