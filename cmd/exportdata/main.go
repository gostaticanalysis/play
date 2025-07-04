package main

import (
	"archive/zip"
	"errors"
	"flag"
	"fmt"
	"go/version"
	"io"
	"os"
	"runtime"
	"strings"

	"go.uber.org/multierr"
	"golang.org/x/tools/go/gcexportdata"
	"golang.org/x/tools/go/packages"
)

var (
	flagGOOS   string
	flagGOARCH string
	flagPkgs   string
)

func init() {
	flag.StringVar(&flagGOOS, "goos", runtime.GOOS, "GOOS")
	flag.StringVar(&flagGOARCH, "goarch", runtime.GOARCH, "GOARCH")
	flag.StringVar(&flagPkgs, "pkgs", "std", "comma separated target packages")
}

func main() {
	flag.Parse()
	if err := run(flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, "dump:", err)
		os.Exit(1)
	}
}

func run(args []string) (rerr error) {

	config := &packages.Config{
		Env:  append(os.Environ(), "GOOS="+flagGOOS, "GOARCH="+flagGOARCH),
		Mode: packages.NeedName | packages.NeedExportFile,
	}

	var patterns []string
	for _, p := range strings.Split(flagPkgs, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			patterns = append(patterns, p)
		}
	}

	if len(patterns) == 0 {
		return errors.New("pacakges must be specified")
	}

	pkgs, err := packages.Load(config, patterns...)
	if err != nil {
		return err
	}

	version := strings.ReplaceAll(version.Lang(runtime.Version()), ".", "")
	fname := strings.Join(patterns, "-") + "_" + version + "_" + flagGOOS + "_" + flagGOARCH + ".zip"
	dst, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer func() {
		rerr = multierr.Append(rerr, dst.Close())
	}()

	w := zip.NewWriter(dst)
	defer func() {
		rerr = multierr.Append(rerr, w.Close())
	}()

	for _, pkg := range pkgs {
		if err := writeExportFile(w, pkg); err != nil {
			return err
		}
	}

	return nil
}

func writeExportFile(w *zip.Writer, pkg *packages.Package) error {
	if strings.HasPrefix(pkg.PkgPath, "vendor/") ||
		strings.HasPrefix(pkg.PkgPath, "internal/") ||
		pkg.ExportFile == "" {
		return nil
	}

	src, err := os.Open(pkg.ExportFile)
	if err != nil {
		return err
	}
	defer src.Close()

	r, err := gcexportdata.NewReader(src)
	if err != nil {
		return err
	}

	dst, err := w.Create(strings.ReplaceAll(pkg.PkgPath, "/", "_"))
	if err != nil {
		return err
	}

	if _, err := io.Copy(dst, r); err != nil {
		return err
	}

	return nil
}
