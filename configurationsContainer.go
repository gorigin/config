package config

import "fmt"

// NewConfigurationsContainer returns new configurations container, build
// on provided configurations
func NewConfigurationsContainer(configs ...Configuration) Configuration {
	return ConfigurationsContainer(configs)
}

// ConfigurationsContainer represents slice of configurations and implements
// Configuration interface
type ConfigurationsContainer []Configuration

// Test tests configuration for consistency and any errors
func (cc ConfigurationsContainer) Test() error {
	for _, conf := range cc {
		if err := conf.Test(); err != nil {
			return err
		}
	}

	return nil
}

// Has returns true if configuration has provided qualifier
func (cc ConfigurationsContainer) Has(qualifier string) bool {
	for _, c := range cc {
		if c.Has(qualifier) {
			return true
		}
	}

	return false
}

// Value returns configuration value and boolean flag
// which set to false when no data found
func (cc ConfigurationsContainer) Value(qualifier string) (interface{}, bool) {
	for i := len(cc) - 1; i >= 0; i-- {
		c := cc[i]
		if c.Has(qualifier) {
			return c.Value(qualifier)
		}
	}

	return nil, false
}

// Configure performs configuration of target using internal
// data, found by qualifier
func (cc ConfigurationsContainer) Configure(qualifier string, target interface{}) error {
	for i := len(cc) - 1; i >= 0; i-- {
		c := cc[i]
		if c.Has(qualifier) {
			return c.Configure(qualifier, target)
		}
	}

	return fmt.Errorf("No config with qualifier %s found", qualifier)
}

// ConfigureValidate performs configuration of target and validated
// configuration if target implements validable interface
func (cc ConfigurationsContainer) ConfigureValidate(qualifier string, target interface{}) error {
	err := cc.Configure(qualifier, target)
	if err == nil {
		if vt, ok := target.(Validable); ok {
			err = vt.Validate()
		}
	}

	return err
}

// Qualifiers returns list of qualifiers
func (cc ConfigurationsContainer) Qualifiers() ([]string, error) {
	unique := map[string]bool{}
	for _, c := range cc {
		qf, err := c.Qualifiers()
		if err != nil {
			return nil, err
		}
		for _, q := range qf {
			unique[q] = true
		}
	}

	names := []string{}
	for n := range unique {
		names = append(names, n)
	}

	return names, nil
}
