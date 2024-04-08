package application

import (
	"testing"

	"github.com/psfpro/metrics/internal/server/domain"
	"github.com/psfpro/metrics/internal/server/infrastructure/storage"
)

func TestIncreaseCounterMetricHandler_Handle(t *testing.T) {
	type fields struct {
		Repository domain.CounterMetricRepository
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *domain.CounterMetric
	}{
		{
			name: "metric",
			fields: fields{
				Repository: storage.NewCounterMetricRepository(),
			},
			args: args{
				name: "MetricName",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := IncreaseCounterMetricHandler{
				Repository: tt.fields.Repository,
			}
			obj.Handle(tt.args.name)
		})
	}
}

func TestUpdateCounterMetricHandler_Handle(t *testing.T) {
	type fields struct {
		Repository domain.CounterMetricRepository
	}
	type args struct {
		name  string
		value int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "metric",
			fields: fields{
				Repository: storage.NewCounterMetricRepository(),
			},
			args: args{
				name:  "MetricName",
				value: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := UpdateCounterMetricHandler{
				Repository: tt.fields.Repository,
			}
			obj.Handle(tt.args.name, tt.args.value)
		})
	}
}

func TestUpdateGaugeMetricHandler_Handle(t *testing.T) {
	type fields struct {
		Repository domain.GaugeMetricRepository
	}
	type args struct {
		name  string
		value float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "metric",
			fields: fields{
				Repository: storage.NewGaugeMetricRepository(),
			},
			args: args{
				name:  "MetricName",
				value: 1.1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := UpdateGaugeMetricHandler{
				Repository: tt.fields.Repository,
			}
			obj.Handle(tt.args.name, tt.args.value)
		})
	}
}
