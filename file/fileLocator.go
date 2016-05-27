package file

// FileLocator is function, able to found file in different locations
type FileLocator func(string) (string, error)
