package domain

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounterMetric_Increase(t *testing.T) {
	type fields struct {
		name  string
		value int64
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{
			name: "zero",
			fields: fields{
				name:  "MetricName",
				value: 0,
			},
			want: 1,
		},
		{
			name: "negative",
			fields: fields{
				name:  "MetricName",
				value: -1,
			},
			want: 0,
		},
		{
			name: "positive",
			fields: fields{
				name:  "MetricName",
				value: 1,
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := &CounterMetric{
				name:  tt.fields.name,
				value: tt.fields.value,
			}
			obj.Increase()
			assert.Equal(t, tt.want, obj.Value())
		})
	}
}

func TestCounterMetric_Name(t *testing.T) {
	type fields struct {
		name  string
		value int64
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
			obj := &CounterMetric{
				name:  tt.fields.name,
				value: tt.fields.value,
			}
			assert.Equal(t, tt.want, obj.Name())
		})
	}
}

func TestCounterMetric_Update(t *testing.T) {
	type fields struct {
		name  string
		value int64
	}
	type args struct {
		value int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int64
	}{
		{
			name: "zero",
			fields: fields{
				name:  "MetricName",
				value: 0,
			},
			args: args{
				value: 1,
			},
			want: 1,
		},
		{
			name: "negative",
			fields: fields{
				name:  "MetricName",
				value: -1,
			},
			args: args{
				value: 2,
			},
			want: 2,
		},
		{
			name: "positive",
			fields: fields{
				name:  "MetricName",
				value: 1,
			},
			args: args{
				value: 3,
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := &CounterMetric{
				name:  tt.fields.name,
				value: tt.fields.value,
			}
			obj.Update(tt.args.value)
		})
	}
}

func TestCounterMetric_Value(t *testing.T) {
	type fields struct {
		name  string
		value int64
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{
			name: "zero",
			fields: fields{
				name:  "MetricName",
				value: 0,
			},
			want: 0,
		},
		{
			name: "negative",
			fields: fields{
				name:  "MetricName",
				value: -1,
			},
			want: -1,
		},
		{
			name: "positive",
			fields: fields{
				name:  "MetricName",
				value: 1,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := &CounterMetric{
				name:  tt.fields.name,
				value: tt.fields.value,
			}
			assert.Equal(t, tt.want, obj.Value())
		})
	}
}

func TestNewCounterMetric(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want *CounterMetric
	}{
		{
			name: "metric",
			args: args{
				name: "MetricName",
			},
			want: &CounterMetric{
				name:  "MetricName",
				value: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCounterMetric(tt.args.name, 0); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCounterMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}
