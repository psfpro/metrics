package application

import "github.com/psfpro/metrics/internal/server/domain"

// UpdateGaugeMetricHandler implements gauge metric update.
type UpdateGaugeMetricHandler struct {
	Repository domain.GaugeMetricRepository
}

func (obj UpdateGaugeMetricHandler) Handle(name string, value float64) {
	metric, exist := obj.Repository.FindByName(name)
	if exist {
		metric.Update(value)
	} else {
		metric = domain.NewGaugeMetric(name, value)
		obj.Repository.Add(metric)
	}
}
