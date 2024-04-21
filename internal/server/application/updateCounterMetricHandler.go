package application

import "github.com/psfpro/metrics/internal/server/domain"

// UpdateCounterMetricHandler implements counter metric update.
type UpdateCounterMetricHandler struct {
	Repository domain.CounterMetricRepository
}

func (obj UpdateCounterMetricHandler) Handle(name string, value int64) {
	metric, exist := obj.Repository.FindByName(name)
	if exist {
		metric.Update(value)
	} else {
		metric = domain.NewCounterMetric(name, value)
		obj.Repository.Add(metric)
	}
}
