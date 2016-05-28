package config

import (
	"os"
	"strings"
)

// ArgsConfigurationOptions is configuration for command line arguments
// configurator
type ArgsConfigurationOptions struct {
	// If provided, only options, prefixed with it will be analyzed
	OnlyWithPrefix string

	// If true, parses -v -vv -vvv -a flags and fills verbosity
	// quiet                   : -1
	// normal                  :  0
	// verbose (-v)            :  1
	// very verbose (-vv)      :  2
	// very very verbose (-vvv):  3
	VerboseAndQuiet bool

	// Prefix used for verbose and quiet flags
	VerboseAndQuietPrefix string
}

// NewCommonOsArgsConfiguration reads `os.Args` and returns Configuration
// built on it
func NewCommonOsArgsConfiguration() Configuration {
	return NewArgsConfiguration(os.Args, ArgsConfigurationOptions{
		OnlyWithPrefix:        "",
		VerboseAndQuiet:       true,
		VerboseAndQuietPrefix: "",
	})
}

// NewArgsConfiguration returns configuration, read from command line arguments
func NewArgsConfiguration(args []string, opts ArgsConfigurationOptions) Configuration {
	cnf := map[string]interface{}{}
	if opts.VerboseAndQuiet {
		cnf[opts.VerboseAndQuietPrefix+"verbosity"] = 0
	}

	for _, v := range args {
		key, value := extractKeyValue(v)

		if opts.VerboseAndQuiet {
			if key == "vvv" {
				cnf[opts.VerboseAndQuietPrefix+"verbosity"] = 3
			} else if key == "vv" {
				cnf[opts.VerboseAndQuietPrefix+"verbosity"] = 2
			} else if key == "v" || key == "verbose" {
				cnf[opts.VerboseAndQuietPrefix+"verbosity"] = 1
			} else if key == "q" || key == "quiet" {
				cnf[opts.VerboseAndQuietPrefix+"verbosity"] = -1
			}
		}

		if opts.OnlyWithPrefix != "" {
			if strings.HasPrefix(key, opts.OnlyWithPrefix) {
				cnf[key] = value
			}
		} else {
			cnf[key] = value
		}
	}

	return MapConfiguration(cnf)
}

func extractKeyValue(v string) (key string, value interface{}) {
	if len(v) > 2 && v[0:2] == "--" {
		chunks := strings.Split(v[2:], "=")
		if len(chunks) == 1 {
			key = chunks[0]
			value = true
		} else {
			key = chunks[0]
			value = chunks[1]
		}
	} else if len(v) > 1 && v[0:1] == "-" {
		key = v[1:]
		value = true
	}

	return
}
