package domain

// CounterMetricRepository defines the interface for a repository that manages counter metrics.
type CounterMetricRepository interface {
	FindAll() map[string]*CounterMetric
	FindByName(name string) (*CounterMetric, bool)
	Add(metric *CounterMetric)
}
