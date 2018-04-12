package file

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// ErrGoPath ...
var ErrGoPath = errors.New("GOPATH isn't set")

// ErrPath ...
var ErrPath = errors.New("path doesn't have GOPATH as sub-path")

// IsGo ...
func IsGo(p string) bool {
	return filepath.Ext(p) == ".go" && !strings.HasSuffix(p, "_test.go")
}

// CopyPackage ...
func CopyPackage(pkg, to string) error {
	from, err := packageDir(pkg)
	if err != nil {
		return err
	}
	to = filepath.Join(to, filepath.Base(from))
	if err := os.MkdirAll(to, os.ModePerm); err != nil {
		return err
	}
	fs, err := PackageFiles(pkg)
	if err != nil {
		return err
	}
	for old, bs := range fs {
		new := filepath.Join(to, filepath.Base(old))
		if err := Write(new, bs); err != nil {
			return err
		}
	}
	return nil
}

// Package ...
func Package(path string) (string, error) {
	base, err := goPath()
	if err != nil {
		return "", err
	}
	split := strings.Split(path, filepath.Join(base, "src")+"/")
	if len(split) != 2 {
		return "", ErrPath
	}
	return split[1], nil
}

// PackageFiles ...
func PackageFiles(pkg string) (map[string][]byte, error) {
	dir, err := packageDir(pkg)
	if err != nil {
		return nil, err
	}
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	files := make(map[string][]byte)
	for _, fi := range fis {
		if !IsGo(fi.Name()) {
			continue
		}
		path := filepath.Join(dir, fi.Name())
		bs, err := Read(path)
		if err != nil {
			return nil, err
		}
		files[path] = bs
	}
	return files, nil
}

func packageDir(pkg string) (string, error) {
	path, err := goPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(path, "src", pkg), nil
}

func goPath() (string, error) {
	path := os.ExpandEnv("$GOPATH")
	if path == "" {
		return "", ErrGoPath
	}
	return path, nil
}
