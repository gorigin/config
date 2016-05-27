package config

// Configuration holds configuration information
type Configuration interface {
	// Tests configuration for consistency and any errors
	Test() error

	// Has returns true if configuration has provided qualifier
	Has(qualifier string) bool

	// Value returns configuration value and boolean flag
	// which set to false when no data found
	Value(qualifier string) (interface{}, bool)

	// Configure performs configuration of target using internal
	// data, found by qualifier
	Configure(qualifier string, target interface{}) error

	// Qualifiers returns list of qualifiers
	Qualifiers() ([]string, error)
}
