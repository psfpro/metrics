package domain

// GaugeMetricRepository defines the interface for a repository that manages gauge metrics.
type GaugeMetricRepository interface {
	FindAll() map[string]*GaugeMetric
	FindByName(name string) (*GaugeMetric, bool)
	Add(metric *GaugeMetric)
}
