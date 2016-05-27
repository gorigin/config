package config

import "fmt"

// MapConfiguration is in-memory configuration for placeholders
type MapConfiguration map[string]interface{}

// Tests configuration for consistency and any errors
func (this MapConfiguration) Test() error {
	if len(this) == 0 {
		return fmt.Errorf("MapConfiguration is empty")
	}

	return nil
}

// Has returns true if configuration has provided qualifier
func (this MapConfiguration) Has(qualifier string) bool {
	if this.Test() != nil {
		return false
	}
	_, ok := this[qualifier]
	return ok
}

// Value returns configuration value and boolean flag
// which set to false when no data found
func (this MapConfiguration) Value(qualifier string) (interface{}, bool) {
	if this.Test() != nil {
		return nil, false
	}
	v, ok := this[qualifier]
	return v, ok
}

// Configure performs configuration of target using internal
// data, found by qualifier
func (this MapConfiguration) Configure(qualifier string, target interface{}) error {
	return fmt.Errorf("MapConfiguration is not supposed to be used as configurer")
}

// Qualifiers returns list of qualifiers
func (this MapConfiguration) Qualifiers() ([]string, error) {
	if err := this.Test(); err != nil {
		return []string{}, err
	}

	qa := []string{}
	for k := range this {
		qa = append(qa, k)
	}

	return qa, nil
}
