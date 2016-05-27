package file

// NoopLocator is locator that does nothing, just returns
// provided string as response with nil error
func NoopLocator(name string) (string, error) {
	return name, nil
}
