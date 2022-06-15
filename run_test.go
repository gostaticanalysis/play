package play_test

import (
	"embed"
	"fmt"
	"testing"

	"github.com/gostaticanalysis/play"
	"golang.org/x/tools/go/analysis"
)

//go:embed testdata/*
var testdata embed.FS

var Analyzer = &analysis.Analyzer{
	Name: "test",
	Doc:  "test",
	Run: func(pass *analysis.Pass) (any, error) {
		fmt.Println(pass.Pkg)
		return nil, nil
	},
}

func Test(t *testing.T) {
	rs, err := play.Run(testdata, Analyzer, "a")
	fmt.Println(rs, err)
}
