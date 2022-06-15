package play

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"io/fs"
	"path"

	"golang.org/x/tools/go/analysis"
)

type Result struct {
	Pass        *analysis.Pass
	Diagnostics []analysis.Diagnostic
	Facts       map[types.Object][]analysis.Fact
	Result      any
	Err         error
}

func Run(testdata fs.FS, a *analysis.Analyzer, pkgs ...string) ([]*Result, error) {
	if len(pkgs) == 0 {
		return nil, nil
	}

	var results []*Result
	for _, pkg := range pkgs {
		rs, err := run(testdata, a, pkg)
		if err != nil {
			return nil, err
		}
		results = append(results, rs...)
	}

	return results, nil
}

func run(testdata fs.FS, a *analysis.Analyzer, pkg string) (results []*Result, rerr error) {
	defer derr(&rerr, "Run")

	fsys, err := fs.Sub(testdata, path.Join("testdata", "src", pkg))
	if err != nil {
		return nil, err
	}

	config := &Config{
		Fset:    token.NewFileSet(),
		Context: Context(fsys),
		Fsys:    fsys,
	}

	pkgs, err := ParseDir(config, ".", parser.ParseComments)
	if err != nil {
		return nil, err
	}

	for name, files := range pkgs {
		typesConfig := &types.Config{
			FakeImportC: true,
			Importer:    Importer(config),
		}

		info := &types.Info{
			Types:      make(map[ast.Expr]types.TypeAndValue),
			Instances:  make(map[*ast.Ident]types.Instance),
			Defs:       make(map[*ast.Ident]types.Object),
			Uses:       make(map[*ast.Ident]types.Object),
			Implicits:  make(map[ast.Node]types.Object),
			Selections: make(map[*ast.SelectorExpr]*types.Selection),
			Scopes:     make(map[ast.Node]*types.Scope),
		}

		typesPkg, err := typesConfig.Check(name, config.Fset, files, info)
		if err != nil {
			return nil, err
		}

		var result Result

		result.Pass = &analysis.Pass{
			Analyzer:   a,
			Fset:       config.Fset,
			Files:      files,
			Pkg:        typesPkg,
			TypesInfo:  info,
			TypesSizes: types.SizesFor(config.Context.Compiler, config.Context.GOARCH),
			Report: func(d analysis.Diagnostic) {
				result.Diagnostics = append(result.Diagnostics, d)
			},
			// FIXME: dependency
			ResultOf: make(map[*analysis.Analyzer]any),
			// FIXME: for fact
		}

		result.Result, result.Err = a.Run(result.Pass)

		results = append(results, &result)
	}

	return results, nil
}
