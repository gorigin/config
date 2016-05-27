package file

import (
	"os"
	"strings"
)

// exists returns true if provided file exists
func exists(name string) bool {
	s, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false
	}

	return s.Mode().IsRegular()
}

// hasPathSeparator returns true if name contains path separator
func hasPathSeparator(name string) bool {
	return strings.Index(name, "/") > -1 || strings.Index(name, "\\") > -1
}
