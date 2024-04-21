package osexit

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "osexit",
	Doc:  "os exit in main function",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			fn, ok := node.(*ast.FuncDecl)
			if !ok {
				return true
			}
			// Check if the function is named "main" and it's in the "main" package
			if fn.Name.Name == "main" && pass.Pkg.Name() == "main" {
				ast.Inspect(fn.Body, func(node ast.Node) bool {
					callExpr, ok := node.(*ast.CallExpr)
					if !ok {
						return true
					}
					if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
						// Check if the call is to os.Exit
						if pkgIdent, ok := selExpr.X.(*ast.Ident); ok && pkgIdent.Name == "os" && selExpr.Sel.Name == "Exit" {
							pass.Reportf(callExpr.Pos(), "call to os.Exit found in main function")
						}
					}
					return true
				})
			}
			return true
		})
	}

	return nil, nil
}
