package play

import (
	"go/ast"
	"go/parser"
	"path"
)

func ParseDir(config *Config, dir string, mode parser.Mode) (map[string][]*ast.File, error) {
	return parseDir(config, dir, mode)
}

func parseDir(config *Config, dir string, mode parser.Mode) (_ map[string][]*ast.File, rerr error) {
	defer derr(&rerr, "ParseDir")

	fis, err := config.Context.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	pkgs := make(map[string][]*ast.File)

	for _, fi := range fis {
		if path.Ext(fi.Name()) != ".go" {
			continue
		}

		match, err := config.Context.MatchFile(dir, fi.Name())
		if err != nil {
			return nil, err
		}

		if !match {
			continue
		}

		filename := config.Context.JoinPath(dir, fi.Name())
		file, err := ParseFile(config, filename, mode)
		if err != nil {
			return nil, err
		}

		pkgs[file.Name.Name] = append(pkgs[file.Name.Name], file)
	}

	return pkgs, nil
}

func ParseFile(config *Config, filename string, mode parser.Mode) (*ast.File, error) {
	return parseFile(config, filename, mode)
}

func parseFile(config *Config, filename string, mode parser.Mode) (_ *ast.File, rerr error) {
	defer derr(&rerr, "ParseFile")

	src, err := config.Context.OpenFile(filename)
	if err != nil {
		return nil, err
	}
	defer src.Close()

	file, err := parser.ParseFile(config.Fset, filename, src, mode)
	if err != nil {
		return nil, err
	}

	return file, nil
}
