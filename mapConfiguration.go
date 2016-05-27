package config

import (
	"fmt"
	"github.com/gorigin/config/reflect"
)

// MapConfiguration is in-memory configuration for placeholders
type MapConfiguration map[string]interface{}

// Test method tests configuration for consistency and any errors
func (mc MapConfiguration) Test() error {
	if len(mc) == 0 {
		return fmt.Errorf("MapConfiguration is empty")
	}

	return nil
}

// Has returns true if configuration has provided qualifier
func (mc MapConfiguration) Has(qualifier string) bool {
	if mc.Test() != nil {
		return false
	}
	_, ok := mc[qualifier]
	return ok
}

// Value returns configuration value and boolean flag
// which set to false when no data found
func (mc MapConfiguration) Value(qualifier string) (interface{}, bool) {
	if mc.Test() != nil {
		return nil, false
	}
	v, ok := mc[qualifier]
	return v, ok
}

// Configure performs configuration of target using internal
// data, found by qualifier
func (mc MapConfiguration) Configure(qualifier string, target interface{}) error {
	v, ok := mc.Value(qualifier)
	if !ok {
		return fmt.Errorf("Qualifier %s not found", qualifier)
	}

	return reflect.AnyMarshaller(v, target)
}

// Qualifiers returns list of qualifiers
func (mc MapConfiguration) Qualifiers() ([]string, error) {
	if err := mc.Test(); err != nil {
		return []string{}, err
	}

	qa := []string{}
	for k := range mc {
		qa = append(qa, k)
	}

	return qa, nil
}
