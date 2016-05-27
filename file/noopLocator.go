package file

func NoopLocator(name string) (string, error) {
	return name, nil
}
