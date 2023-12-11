package application

import (
	"fmt"
	"github.com/psfpro/metrics/internal/domain"
)

// UpdateGaugeMetricHandler Обновление gauge метрики
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

// UpdateCounterMetricHandler Обновление counter метрики
type UpdateCounterMetricHandler struct {
	Repository domain.CounterMetricRepository
}

func (obj UpdateCounterMetricHandler) Handle(name string, value int64) {
	metric, exist := obj.Repository.FindByName(name)
	if exist {
		fmt.Println(metric)
		metric.Update(value)
		fmt.Println(metric)
	} else {
		metric = domain.NewCounterMetric(name)
		obj.Repository.Add(metric)
	}
}

// IncreaseCounterMetricHandler Обновление counter метрики
type IncreaseCounterMetricHandler struct {
	Repository domain.CounterMetricRepository
}

func (obj IncreaseCounterMetricHandler) Handle(name string) {
	metric, exist := obj.Repository.FindByName(name)
	if exist {
		fmt.Println(metric)
		metric.Increase()
		fmt.Println(metric)
	} else {
		metric = domain.NewCounterMetric(name)
		obj.Repository.Add(metric)
	}
}
