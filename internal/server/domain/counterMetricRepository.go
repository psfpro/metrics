package domain

type CounterMetricRepository interface {
	FindAll() map[string]*CounterMetric
	FindByName(name string) (*CounterMetric, bool)
	Add(metric *CounterMetric)
}
