package ini

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gorigin/config"
	"github.com/gorigin/config/file"
	"github.com/gorigin/config/reflect"
	"io/ioutil"
	"strings"
)

// IniFileValuesReader reads byte contents into values map
func IniFileValuesFiller(bts []byte) (map[string]interface{}, error) {
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

func NewIniConfigFile(filename string, l file.FileLocator, r file.FileReader) config.Configuration {
	return file.NewFileConfiguration(filename, l, r, IniFileValuesFiller, reflect.AnyMarshaller)
}

func NewLocalCommonIniConfigFile(filename string) config.Configuration {
	return NewIniConfigFile(filename, file.LocalFolderLocator, ioutil.ReadFile)
}

func NewCommonIniConfigFile(filename string, subfolder string) config.Configuration {
	return NewIniConfigFile(filename, file.NewCommonLocationsLocator(true, true, true, subfolder), ioutil.ReadFile)
}
