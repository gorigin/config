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

func (this *fileConfiguration) Test() error {
	if this.err != nil {
		return this.err
	}
	if this.values == nil {
		err := this.Reload()
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *fileConfiguration) Reload() error {
	// Locating file
	filename, err := this.Locator(this.Filename)
	if err != nil {
		this.err = err
		return err
	}

	// Reading file
	bts, err := this.Reader(filename)
	if err != nil {
		this.err = err
		return err
	}

	// Filling values
	this.values, err = this.ByteToMapReader(bts)
	if err != nil {
		this.err = err
		return err
	}

	return nil
}

func (this *fileConfiguration) Has(qualifier string) bool {
	if this.Test() != nil {
		return false
	}

	_, ok := this.values[qualifier]
	return ok
}

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
