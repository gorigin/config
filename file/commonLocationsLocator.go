package file

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

var separator = string(os.PathSeparator)

// NewCommonLocationsLocator return a locator, that will search for file in current folder,
// home folder and /etc/.
// Subfolder may be configured
func NewCommonLocationsLocator(currentFolder, homeFolder, etc bool, subfolder string) Locator {
	return func(filename string) (string, error) {
		if hasPathSeparator(filename) {
			// Absolute path
			return filename, nil
		}

		// Building paths list
		paths := []string{}

		if currentFolder {
			path, err := getCwd()
			if err != nil {
				return "", err
			}

			paths = append(paths, path)
		}

		if homeFolder {
			path, err := getHomeFolder()
			if err != nil {
				return "", err
			}
			if len(path) > 0 {
				paths = append(paths, path)
			}
		}

		if etc {
			paths = append(paths, "/etc")
		}

		// Appending path separators and subfolder
		for k, v := range paths {
			if !strings.HasSuffix(v, separator) {
				paths[k] = v + separator
			}

			if len(subfolder) > 0 {
				paths[k] = paths[k] + subfolder + separator
			}
		}

		// Searching
		for _, path := range paths {
			if exists(path + filename) {
				return path + filename, nil
			}
		}

		return "", fmt.Errorf("File %s not found in %s", filename, strings.Join(paths, " "))
	}
}

func getCwd() (string, error) {
	return os.Getwd()
}

func getHomeFolder() (string, error) {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home, nil
	}
	return os.Getenv("HOME"), nil
}
