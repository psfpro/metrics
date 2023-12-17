package domain

type CounterMetricRepository interface {
	FindByName(name string) (*CounterMetric, bool)
	Add(metric *CounterMetric)
}
