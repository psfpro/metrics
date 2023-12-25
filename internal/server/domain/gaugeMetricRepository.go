package domain

type GaugeMetricRepository interface {
	FindAll() map[string]*GaugeMetric
	FindByName(name string) (*GaugeMetric, bool)
	Add(metric *GaugeMetric)
}
