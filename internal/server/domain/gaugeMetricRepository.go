package domain

type GaugeMetricRepository interface {
	FindByName(name string) (*GaugeMetric, bool)
	Add(metric *GaugeMetric)
}
