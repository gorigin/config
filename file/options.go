package file

import "io/ioutil"

// Options holds minimal information, required to read a file
type Options struct {
	// Configuration file name
	Filename string

	// Locator, used to find file
	Locator FileLocator

	// Reader, used to read file contents
	Reader FileReader
}

// WithDefaults returns copy of options with default values set instead missing ones
func (o Options) WithDefaults() Options {
	fco := Options{
		Filename: o.Filename,
		Locator:  o.Locator,
		Reader:   o.Reader,
	}

	if fco.Reader == nil {
		fco.Reader = ioutil.ReadFile
	}
	if fco.Locator == nil {
		fco.Locator = LocalFolderLocator
	}

	return fco
}

// FullOptions holds full information, required to read a file
type FullOptions struct {
	Options

	// Reader, used to read byte slice (file contents) into
	// intermediate key-value map
	ByteToMapReader func([]byte) (map[string]interface{}, error)

	// Mapper used to write final configuration values
	ReflectionMapper func(source interface{}, target interface{}) error
}
