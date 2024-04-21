package main

import (
	"github.com/go-critic/go-critic/checkers/analyzer"
	"github.com/kisielk/errcheck/errcheck"
	"github.com/psfpro/metrics/internal/staticlint/osexit"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	checks := []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		osexit.Analyzer,
		errcheck.Analyzer,
		analyzer.Analyzer,
	}

	for _, v := range staticcheck.Analyzers {
		checks = append(checks, v.Analyzer)
	}
	for _, v := range simple.Analyzers {
		checks = append(checks, v.Analyzer)
	}

	multichecker.Main(checks...)
}
