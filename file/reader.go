package file

// Reader can read resolved file contents
type Reader func(string) ([]byte, error)
