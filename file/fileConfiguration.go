package file

import (
	"fmt"
	"github.com/gorigin/config"
)

type fileConfiguration struct {
	filename string

	locate FileLocator
	read   FileReader
	fill   func([]byte) (map[string]interface{}, error)
	mapper func(source interface{}, target interface{}) error

	values map[string]interface{}
	err    error
}

func NewFileConfiguration(filename string, l FileLocator, r FileReader, f func([]byte) (map[string]interface{}, error), m func(interface{}, interface{}) error) config.Configuration {
	return &fileConfiguration{
		filename: filename,
		locate:   l,
		read:     r,
		fill:     f,
		mapper:   m,
	}
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
	filename, err := this.locate(this.filename)
	if err != nil {
		this.err = err
		return err
	}

	// Reading file
	bts, err := this.read(filename)
	if err != nil {
		this.err = err
		return err
	}

	// Filling values
	this.values, err = this.fill(bts)
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
	return this.mapper(val, target)
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
