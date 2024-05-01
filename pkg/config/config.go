package config

import (
	"flag"
	"os"
	"strconv"
)

var set []param

func StringVar(value *string, envKey string, flagName string, usage string) {
	flagValue := flag.String(flagName, "", usage)
	set = append(set, &stringParam{value: value, flagValue: flagValue, envKey: envKey})
}

func Int64Var(value *int64, envKey string, flagName string, usage string) {
	flagValue := flag.String(flagName, "", usage)
	set = append(set, &int64Param{value: value, flagValue: flagValue, envKey: envKey})
}

func BoolVar(value *bool, envKey string, flagName string, usage string) {
	flagValue := flag.String(flagName, "", usage)
	set = append(set, &boolParam{value: value, flagValue: flagValue, envKey: envKey})
}

func Parse() {
	flag.Parse()
	for _, v := range set {
		v.GetEnv()
	}
}

type param interface {
	GetEnv()
}

type stringParam struct {
	value     *string
	flagValue *string
	envKey    string
}

func (p *stringParam) GetEnv() {
	if envValue := os.Getenv(p.envKey); envValue != "" {
		*p.flagValue = envValue
	}
	if *p.flagValue != "" {
		*p.value = *p.flagValue
	}
}

type int64Param struct {
	value     *int64
	flagValue *string
	envKey    string
}

func (p *int64Param) GetEnv() {
	if envValue := os.Getenv(p.envKey); envValue != "" {
		*p.flagValue = envValue
	}
	if *p.flagValue != "" {
		v, _ := strconv.ParseInt(*p.flagValue, 10, 64)
		*p.value = v
	}
}

type boolParam struct {
	value     *bool
	flagValue *string
	envKey    string
}

func (p *boolParam) GetEnv() {
	if envValue := os.Getenv(p.envKey); envValue != "" {
		*p.flagValue = envValue
	}
	if *p.flagValue != "" {
		v, _ := strconv.ParseBool(*p.flagValue)
		*p.value = v
	}
}
