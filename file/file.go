package file

import (
	"io/ioutil"
	"os"
)

// Read ...
func Read(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

// Write ...
func Write(path string, bs []byte) error {
	return ioutil.WriteFile(path, bs, os.ModePerm)
}
