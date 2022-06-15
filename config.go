package play

import (
	"go/build"
	"go/token"
	"io/fs"
)

type Config struct {
	Fset    *token.FileSet
	Fsys    fs.FS
	Context *build.Context
}

func (config *Config) HasVendor() bool {
	if config.Context == nil {
		return false
	}
	return config.Context.IsDir("vendor")
}
