package ini

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gorigin/config"
	"github.com/gorigin/config/file"
	"github.com/gorigin/config/reflect"
	"strings"
)

func iniFileValuesFiller(bts []byte) (map[string]interface{}, error) {
	props := map[string]interface{}{}
	lineNumber := 0
	scanner := bufio.NewScanner(bytes.NewReader(bts))
	section := ""
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || line[0] == '#' || line[0] == '/' {
			continue
		}
		if line[0] == '[' {
			section = line[1:strings.Index(line, "]")]
			continue
		}
		chunks := strings.SplitN(line, "=", 2)
		if len(chunks) != 2 {
			return nil, fmt.Errorf("Error parsing line %d %s", lineNumber, line)
		}
		props[strings.TrimSpace(chunks[0])] = strings.TrimSpace(chunks[1])
		if section != "" {
			props[section+":"+strings.TrimSpace(chunks[0])] = strings.TrimSpace(chunks[1])
		}
	}

	return props, nil
}

// New returns new .ini (or .cnf) configuration source
func New(options file.Options) config.Configuration {
	return file.NewFileConfiguration(
		file.FullOptions{
			Options:          options.WithDefaults(),
			ByteToMapReader:  iniFileValuesFiller,
			ReflectionMapper: reflect.AnyMarshaller,
		},
	)
}
