package domain

// CounterMetric Тип counter, int64 — новое значение должно добавляться к предыдущему,
// если какое-то значение уже было известно серверу.
type CounterMetric struct {
	name  string
	value int64
}

func NewCounterMetric(name string, value int64) *CounterMetric {
	return &CounterMetric{name: name, value: value}
}

func (obj *CounterMetric) Name() string {
	return obj.name
}

func (obj *CounterMetric) Value() int64 {
	return obj.value
}

func (obj *CounterMetric) Update(value int64) {
	obj.value += value
}

func (obj *CounterMetric) Increase() {
	obj.value++
}
