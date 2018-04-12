package get

import (
	"os/exec"

	"github.com/jwowillo/fill/file"
)

// To ...
func To(pkg, path string) error {
	if err := goGet(pkg); err != nil {
		return err
	}
	return file.CopyPackage(pkg, path)
}

func goGet(pkg string) error {
	_, err := exec.Command("go", "get", "-u", pkg).Output()
	return err
}
