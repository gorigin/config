package file

// Locator is function, able to found file in different locations
type Locator func(string) (string, error)
