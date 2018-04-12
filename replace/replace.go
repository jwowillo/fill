package replace

import (
	"errors"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jwowillo/fill/file"

	"golang.org/x/tools/imports"
)

// ErrFileType ...
var ErrFileType = errors.New("file-type is bad")

// ErrType ...
var ErrType = errors.New("type doesn't exist")

// ErrTest ...
var ErrTest = errors.New("go test failed")

// PackageTypes ...
func PackageTypes(pkg string, types map[string]string) error {
	paths, err := file.PackageFiles(pkg)
	if err != nil {
		return err
	}
	var found bool
	for path := range paths {
		if err := fileTypes(path, types); err != nil {
			if err != ErrType {
				return err
			}
		} else {
			found = true
		}
	}
	if !found {
		return ErrType
	}
	return goTest(pkg)
}

// PackageType ...
func PackageType(pkg, old, new string) error {
	return PackageTypes(pkg, map[string]string{old: new})
}

// FileTypes ...
func FileTypes(path string, types map[string]string) error {
	if err := fileTypes(path, types); err != nil {
		return err
	}
	return goTest(filepath.Dir(path))
}

// FileType ...
func FileType(path, old, new string) error {
	if err := fileType(path, old, new); err != nil {
		return err
	}
	return goTest(filepath.Dir(path))
}

func fileTypes(path string, types map[string]string) error {
	for old, new := range types {
		if err := fileType(path, old, new); err != nil {
			return err
		}
	}
	return nil
}

func fileType(path, old, new string) error {
	if !file.IsGo(path) {
		return ErrFileType
	}
	fset := token.NewFileSet()
	pf, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	terr := removeType(pf, old)
	replaceTypeInNode(pf, old, new)
	bs, err := fileToBytes(fset, pf)
	if err != nil {
		return err
	}
	bs, err = imports.Process(path, bs, nil)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(path, bs, os.ModePerm); err != nil {
		return err
	}
	return terr
}

func goTest(pkg string) error {
	_, err := exec.Command("go", "test", "-c", pkg).Output()
	if err != nil {
		return ErrTest
	}
	return nil
}
