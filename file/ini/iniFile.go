package ini

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gorigin/config"
	"github.com/gorigin/config/file"
	"io/ioutil"
	"reflect"
	"strconv"
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

func IniFileValuesMapper(source interface{}, target interface{}) error {
	sv, ok := source.(string)
	if !ok {
		return fmt.Errorf("Ini value must be string")
	}

	kind := reflect.TypeOf(target).Elem().Kind()
	val := reflect.ValueOf(target).Elem()
	switch kind {
	case reflect.String:
		val.SetString(sv)
	case reflect.Int, reflect.Int64:
		if iv, err := strconv.Atoi(sv); err != nil {
			return err
		} else {
			val.SetInt(int64(iv))
		}
	case reflect.Float64:
		if fv, err := strconv.ParseFloat(sv, 64); err != nil {
			return err
		} else {
			val.SetFloat(fv)
		}
	case reflect.Bool:
		sv = strings.ToUpper(sv)
		val.SetBool(sv == "TRUE" || sv == "T" || sv == "1" || sv == "YES" || sv == "Y" || sv == "ON")
	default:
		return fmt.Errorf("Unsupported kind %s", kind)
	}

	return nil
}

func NewIniConfigFile(filename string, l file.FileLocator, r file.FileReader) config.Configuration {
	return file.NewFileConfiguration(filename, l, r, IniFileValuesFiller, IniFileValuesMapper)
}

func NewLocalCommonIniConfigFile(filename string) config.Configuration {
	return NewIniConfigFile(filename, file.LocalFolderLocator, ioutil.ReadFile)
}

func NewCommonIniConfigFile(filename string, subfolder string) config.Configuration {
	return NewIniConfigFile(filename, file.GetCommonLocationsLocator(true, true, true, subfolder), ioutil.ReadFile)
}
