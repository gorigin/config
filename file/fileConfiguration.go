package file

import (
	"fmt"
	"github.com/gorigin/config"
)

type fileConfiguration struct {
	FullOptions

	values map[string]interface{}
	err    error
}

// NewFileConfiguration creates and returns new file-based configuration
func NewFileConfiguration(opts FullOptions) config.Configuration {
	return &fileConfiguration{FullOptions: opts}
}

// Test method tests configuration for consistency and any errors
func (fc *fileConfiguration) Test() error {
	if fc.err != nil {
		return fc.err
	}
	if fc.values == nil {
		err := fc.Reload()
		if err != nil {
			return err
		}
	}

	return nil
}

// Reload reloads file data from disk
func (fc *fileConfiguration) Reload() error {
	// Locating file
	filename, err := fc.Locator(fc.Filename)
	if err != nil {
		fc.err = err
		return err
	}

	// Reading file
	bts, err := fc.Reader(filename)
	if err != nil {
		fc.err = err
		return err
	}

	// Filling values
	fc.values, err = fc.ByteToMapReader(bts)
	if err != nil {
		fc.err = err
		return err
	}

	return nil
}

// Has returns true if configuration has provided qualifier
func (this *fileConfiguration) Has(qualifier string) bool {
	if this.Test() != nil {
		return false
	}

	_, ok := this.values[qualifier]
	return ok
}

// Value returns configuration value and boolean flag
// which set to false when no data found
func (this *fileConfiguration) Value(qualifier string) (interface{}, bool) {
	if this.Test() != nil {
		return nil, false
	}
	if this.values == nil {
		err := this.Reload()
		if err != nil {
			return nil, false
		}
	}

	v, ok := this.values[qualifier]
	return v, ok
}

// Configure performs configuration of target using internal
// data, found by qualifier
func (this *fileConfiguration) Configure(qualifier string, target interface{}) error {
	if err := this.Test(); err != nil {
		return err
	}
	if this.values == nil {
		err := this.Reload()
		if err != nil {
			return err
		}
	}

	val, ok := this.values[qualifier]
	if !ok {
		return fmt.Errorf("Qualifier %s not found in configuration", qualifier)
	}
	return this.ReflectionMapper(val, target)
}

// Qualifiers returns list of qualifiers
func (this *fileConfiguration) Qualifiers() ([]string, error) {
	if err := this.Test(); err != nil {
		return nil, err
	}
	if this.values == nil {
		err := this.Reload()
		if err != nil {
			return nil, err
		}
	}

	names := []string{}
	for n := range this.values {
		names = append(names, n)
	}

	return names, nil
}
