package play

import (
	"fmt"
	"go/types"
	"path"
)

func ImportFromVendor(parent types.Importer, config *Config, p string) (*types.Package, error) {
	if !config.HasVendor() {
		return nil, nil
	}

	pkgs, err := ParseDir(config, path.Join("vendor", p), 0)
	if err != nil {
		return nil, err
	}

	for name, files := range pkgs {
		var firstHardErr error
		typesConfig := &types.Config{
			IgnoreFuncBodies: true,
			Importer:         parent,
			Error: func(err error) {
				if firstHardErr == nil && !err.(types.Error).Soft {
					firstHardErr = err
				}
			},
		}

		typesPkg, err := typesConfig.Check(name, config.Fset, files, nil)
		if err != nil {
			if firstHardErr != nil {
				typesPkg = nil
				err = firstHardErr
			}
			return typesPkg, fmt.Errorf("type-checking package %q failed: %w", name, err)
		}
		if firstHardErr != nil {
			panic("package is not safe yet no error was returned")
		}

		if typesPkg != nil {
			return typesPkg, nil
		}
	}

	return nil, nil
}
