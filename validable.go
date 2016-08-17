package config

// Validable is interface for structs, able to validate its contents
type Validable interface {
	// Validate validates its contents
	Validate() error
}
