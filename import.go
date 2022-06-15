package play

import (
	"fmt"
	"go/types"
)

func Importer(config *Config) types.Importer {
	return &defaultImporter{
		imports: make(map[string]*types.Package),
		config:  config,
	}
}

type defaultImporter struct {
	imports map[string]*types.Package
	config  *Config
}

func (im *defaultImporter) Import(path string) (*types.Package, error) {

	if path == "unsafe" {
		return types.Unsafe, nil
	}

	pkg := im.imports[path]
	if pkg != nil {
		return pkg, nil
	}

	pkg, err := ImportFromExportData(im.config, path)
	if pkg != nil {
		im.imports[path] = pkg
		return pkg, nil
	}

	pkg, err = ImportFromVendor(im, im.config, path)
	if err != nil {
		return nil, fmt.Errorf("import %s: %w", path, err)
	}

	if pkg != nil {
		im.imports[path] = pkg
		return pkg, nil
	}

	return nil, fmt.Errorf("not found %s", path)
}
