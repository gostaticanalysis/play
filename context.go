package play

import (
	"go/build"
	"io"
	"io/fs"
	"path"
	"path/filepath"
	"strings"
)

func Context(fsys fs.FS) *build.Context {
	ctx := build.Default

	ctx.OpenFile = func(p string) (io.ReadCloser, error) {
		p = filepath.FromSlash(p)
		if strings.HasPrefix(p, "/") {
			p = p[1:]
		}
		return fsys.Open(p)
	}

	ctx.JoinPath = func(elem ...string) string {
		return path.Join(elem...)
	}

	ctx.SplitPathList = func(list string) []string {
		return strings.Split(list, "/")
	}

	ctx.IsAbsPath = func(p string) bool {
		return path.IsAbs(p)
	}

	ctx.IsDir = func(p string) bool {
		fi, err := fs.Stat(fsys, p)
		if err != nil {
			return false
		}
		return fi.IsDir()
	}

	ctx.HasSubdir = func(root, dir string) (rel string, ok bool) {
		root = path.Clean(root)
		dir = path.Clean(dir)
		if strings.HasPrefix(dir, root) {
			return dir[len(root):], true
		}
		return "", false
	}

	ctx.ReadDir = func(dir string) ([]fs.FileInfo, error) {
		des, err := fs.ReadDir(fsys, dir)
		if err != nil {
			return nil, err
		}

		fis := make([]fs.FileInfo, len(des))
		for i := range des {
			fi, err := des[i].Info()
			if err != nil {
				return nil, err
			}
			fis[i] = fi
		}

		return fis, nil
	}

	return &ctx
}
