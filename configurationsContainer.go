package config

import "fmt"

func NewConfigurationsContainer(configs ...Configuration) Configuration {
	return ConfigurationsContainer(configs)
}

type ConfigurationsContainer []Configuration

func (this ConfigurationsContainer) Test() error {
	for _, conf := range this {
		if err := conf.Test(); err != nil {
			return err
		}
	}

	return nil
}

func (this ConfigurationsContainer) Has(qualifier string) bool {
	for _, c := range this {
		if c.Has(qualifier) {
			return true
		}
	}

	return false
}

func (this ConfigurationsContainer) Value(qualifier string) (interface{}, bool) {
	for i := len(this) - 1; i >= 0; i-- {
		c := this[i]
		if c.Has(qualifier) {
			return c.Value(qualifier)
		}
	}

	return nil, false
}

func (this ConfigurationsContainer) Configure(qualifier string, target interface{}) error {
	for i := len(this) - 1; i >= 0; i-- {
		c := this[i]
		if c.Has(qualifier) {
			return c.Configure(qualifier, target)
		}
	}

	return fmt.Errorf("No config with qualifier %s found", qualifier)
}

func (this ConfigurationsContainer) Qualifiers() ([]string, error) {
	unique := map[string]bool{}
	for _, c := range this {
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
