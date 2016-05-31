package file

// NewStringReader returns FileReader, that reads data from provided
// string. This is useful in unit tests
func NewStringReader(data string) Reader {
	return func(string) ([]byte, error) {
		return []byte(data), nil
	}
}
