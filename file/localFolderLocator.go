package file

import (
	"fmt"
	"path/filepath"
)

// LocalFolderLocator locates files only in current folder
func LocalFolderLocator(filename string) (name string, err error) {
	if !exists(filename) {
		return "", fmt.Errorf("File %s not found", filename)
	}

	return filepath.Abs(filename)
}
