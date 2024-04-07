package domain

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGaugeMetric_Name(t *testing.T) {
	type fields struct {
		name  string
		value float64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "zero",
			fields: fields{
				name:  "MetricName",
				value: 0,
			},
			want: "MetricName",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := &GaugeMetric{
				name:  tt.fields.name,
				value: tt.fields.value,
			}
			assert.Equal(t, tt.want, obj.Name())
		})
	}
}

func TestGaugeMetric_Update(t *testing.T) {
	type fields struct {
		name  string
		value float64
	}
	type args struct {
		value float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			name: "zero",
			fields: fields{
				name:  "MetricName",
				value: 0.0,
			},
			args: args{
				value: 1.1,
			},
			want: 1.1,
		},
		{
			name: "negative",
			fields: fields{
				name:  "MetricName",
				value: -1.1,
			},
			args: args{
				value: 2.2,
			},
			want: 2.2,
		},
		{
			name: "positive",
			fields: fields{
				name:  "MetricName",
				value: 1.1,
			},
			args: args{
				value: 3.3,
			},
			want: 3.3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := &GaugeMetric{
				name:  tt.fields.name,
				value: tt.fields.value,
			}
			obj.Update(tt.args.value)
			assert.Equal(t, tt.want, obj.Value())
		})
	}
}

func TestGaugeMetric_Value(t *testing.T) {
	type fields struct {
		name  string
		value float64
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "zero",
			fields: fields{
				name:  "MetricName",
				value: 0.0,
			},
			want: 0.0,
		},
		{
			name: "negative",
			fields: fields{
				name:  "MetricName",
				value: -1.1,
			},
			want: -1.1,
		},
		{
			name: "positive",
			fields: fields{
				name:  "MetricName",
				value: 1.1,
			},
			want: 1.1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := &GaugeMetric{
				name:  tt.fields.name,
				value: tt.fields.value,
			}
			assert.Equal(t, tt.want, obj.Value())
		})
	}
}

func TestNewGaugeMetric(t *testing.T) {
	type args struct {
		name  string
		value float64
	}
	tests := []struct {
		name string
		args args
		want *GaugeMetric
	}{
		{
			name: "metric",
			args: args{
				name:  "MetricName",
				value: 1.1,
			},
			want: &GaugeMetric{
				name:  "MetricName",
				value: 1.1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGaugeMetric(tt.args.name, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGaugeMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}
