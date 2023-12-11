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

type GaugeMetricRepository interface {
	FindByName(name string) (*GaugeMetric, bool)
	Add(metric *GaugeMetric)
}

// CounterMetric Тип counter, int64 — новое значение должно добавляться к предыдущему,
// если какое-то значение уже было известно серверу.
type CounterMetric struct {
	name  string
	value int64
}

func NewCounterMetric(name string) *CounterMetric {
	return &CounterMetric{name: name, value: 0}
}

func (obj *CounterMetric) Name() string {
	return obj.name
}

func (obj *CounterMetric) Value() int64 {
	return obj.value
}

func (obj *CounterMetric) Update(value int64) {
	obj.value = value
}

func (obj *CounterMetric) Increase() {
	obj.value++
}

type CounterMetricRepository interface {
	FindByName(name string) (*CounterMetric, bool)
	Add(metric *CounterMetric)
}
