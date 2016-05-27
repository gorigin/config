package file

// FileReader can read resolved file contents
type FileReader func(string) ([]byte, error)
