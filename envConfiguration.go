package config

import (
	"os"
	"strings"
)

// NewEnvConfiguration creates configuration from ENVIRONMENT variables
// If prefix provided, only keys with it will be added
func NewEnvConfiguration(Prefix string) Configuration {
	values := map[string]interface{}{}
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		if len(pair) == 2 {
			if Prefix != "" && strings.HasPrefix(pair[0], Prefix) {
				values[pair[0][len(Prefix):]] = pair[1]
			} else {
				values[pair[0]] = pair[1]
			}
		}
	}

	return MapConfiguration(values)
}
