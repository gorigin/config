package file

func NewStringReader(data string) FileReader {
	return func(string) ([]byte, error) {
		return []byte(data), nil
	}
}
