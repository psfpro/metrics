package domain

// GaugeMetric Тип gauge, float64 — новое значение должно замещать предыдущее.
type GaugeMetric struct {
	name  string
	value float64
}

func NewGaugeMetric(name string, value float64) *GaugeMetric {
	return &GaugeMetric{name: name, value: value}
}

func (obj *GaugeMetric) Name() string {
	return obj.name
}

func (obj *GaugeMetric) Value() float64 {
	return obj.value
}

func (obj *GaugeMetric) Update(value float64) {
	obj.value = value
}
