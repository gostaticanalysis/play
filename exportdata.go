package play

import (
	"archive/zip"
	"bytes"
	"fmt"
	"go/types"
	"strings"

	"golang.org/x/tools/go/gcexportdata"
)

//go:generate go run github.com/gostaticanalysis/play/cmd/exportdata

var exportFiles *zip.Reader

func init() {
	var err error
	exportFiles, err = zip.NewReader(bytes.NewReader(exportFilesZIP), int64(len(exportFilesZIP)))
	if err != nil {
		panic("play: " + err.Error())
	}
}

func ImportFromExportData(config *Config, path string) (*types.Package, error) {
	filename := strings.ReplaceAll(path, "/", "_")
	f, err := exportFiles.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("read export data of %q: %w", path, err)
	}
	defer f.Close()

	pkg, err := gcexportdata.Read(f, config.Fset, make(map[string]*types.Package), path)
	if err != nil {
		return nil, fmt.Errorf("read export data of %q: %w", path, err)
	}

	return pkg, nil
}
