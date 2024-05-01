package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestBoolVar(t *testing.T) {
	type args struct {
		value    bool
		envKey   string
		flagName string
		usage    string
	}
	tests := []struct {
		name    string
		prepare func()
		args    args
		want    bool
	}{
		{
			name: "default",
			prepare: func() {

			},
			args: args{
				value:    true,
				envKey:   "BOOL_A",
				flagName: "bool.a",
				usage:    "a",
			},
			want: true,
		},
		{
			name: "arg parse",
			prepare: func() {
				os.Args = append(os.Args, "-bool.b=true")
			},
			args: args{
				value:    false,
				envKey:   "BOOL_B",
				flagName: "bool.b",
				usage:    "b",
			},
			want: true,
		},
		{
			name: "env parse",
			prepare: func() {
				_ = os.Setenv("BOOL_C", "true")
			},
			args: args{
				value:    false,
				envKey:   "BOOL_C",
				flagName: "bool.c",
				usage:    "c",
			},
			want: true,
		},
		{
			name: "arg and env parse",
			prepare: func() {
				os.Args = append(os.Args, "-bool.d=false")
				_ = os.Setenv("BOOL_D", "true")
			},
			args: args{
				value:    false,
				envKey:   "BOOL_D",
				flagName: "bool.d",
				usage:    "d",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare()
			BoolVar(&tt.args.value, tt.args.envKey, tt.args.flagName, tt.args.usage)
			Parse()
			assert.Equal(t, tt.want, tt.args.value)
		})
	}
}

func TestInt64Var(t *testing.T) {
	type args struct {
		value    int64
		envKey   string
		flagName string
		usage    string
	}
	tests := []struct {
		name    string
		prepare func()
		args    args
		want    int64
	}{
		{
			name: "default",
			prepare: func() {

			},
			args: args{
				value:    1,
				envKey:   "INT_A",
				flagName: "int.a",
				usage:    "a",
			},
			want: 1,
		},
		{
			name: "arg parse",
			prepare: func() {
				os.Args = append(os.Args, "-int.b=2")
			},
			args: args{
				value:    1,
				envKey:   "INT_B",
				flagName: "int.b",
				usage:    "b",
			},
			want: 2,
		},
		{
			name: "env parse",
			prepare: func() {
				_ = os.Setenv("INT_C", "3")
			},
			args: args{
				value:    1,
				envKey:   "INT_C",
				flagName: "int.c",
				usage:    "c",
			},
			want: 3,
		},
		{
			name: "arg and env parse",
			prepare: func() {
				os.Args = append(os.Args, "-int.d=2")
				_ = os.Setenv("INT_D", "3")
			},
			args: args{
				value:    1,
				envKey:   "INT_D",
				flagName: "int.d",
				usage:    "d",
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare()
			Int64Var(&tt.args.value, tt.args.envKey, tt.args.flagName, tt.args.usage)
			Parse()
			assert.Equal(t, tt.want, tt.args.value)
		})
	}
}

func TestStringVar(t *testing.T) {
	type args struct {
		value    string
		envKey   string
		flagName string
		usage    string
	}
	var tests = []struct {
		name    string
		prepare func()
		args    args
		want    string
	}{
		{
			name: "default",
			prepare: func() {

			},
			args: args{
				value:    "v",
				envKey:   "STRING_A",
				flagName: "string.a",
				usage:    "a",
			},
			want: "v",
		},
		{
			name: "arg parse",
			prepare: func() {
				os.Args = append(os.Args, "-string.b=f")
			},
			args: args{
				value:    "v",
				envKey:   "STRING_B",
				flagName: "string.b",
				usage:    "b",
			},
			want: "f",
		},
		{
			name: "env parse",
			prepare: func() {
				_ = os.Setenv("STRING_C", "e")
			},
			args: args{
				value:    "v",
				envKey:   "STRING_C",
				flagName: "string.c",
				usage:    "c",
			},
			want: "e",
		},
		{
			name: "arg and env parse",
			prepare: func() {
				os.Args = append(os.Args, "-string.d=f")
				_ = os.Setenv("STRING_D", "e")
			},
			args: args{
				value:    "v",
				envKey:   "STRING_D",
				flagName: "string.d",
				usage:    "d",
			},
			want: "e",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare()
			StringVar(&tt.args.value, tt.args.envKey, tt.args.flagName, tt.args.usage)
			Parse()
			assert.Equal(t, tt.want, tt.args.value)
		})
	}
}
