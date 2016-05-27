package json

import (
	"encoding/json"
	"fmt"
	"github.com/gorigin/config"
	"github.com/gorigin/config/file"
	"io/ioutil"
)

func JsonFileValuesFiller(bts []byte) (map[string]interface{}, error) {
	var target map[string]json.RawMessage
	err := json.Unmarshal(bts, &target)

	if err != nil {
		return nil, err
	}

	// Converting to map[string]interface{}
	response := map[string]interface{}{}
	for k, v := range target {
		response[k] = v
	}

	return response, nil
}

func JsonFileValuesMapper(source interface{}, target interface{}) error {
	raw, ok := source.(json.RawMessage)
	if !ok {
		return fmt.Errorf("Value must be RawMessage")
	}

	return json.Unmarshal(raw, target)
}

func NewJsonConfigFile(filename string, l file.FileLocator, r file.FileReader) config.Configuration {
	return file.NewFileConfiguration(filename, l, r, JsonFileValuesFiller, JsonFileValuesMapper)
}

func NewLocalCommonJsonConfigFile(filename string) config.Configuration {
	return NewJsonConfigFile(filename, file.LocalFolderLocator, ioutil.ReadFile)
}

func NewLocalJsonConfigWithPlaceholders(filename string, props config.Configuration) config.Configuration {
	return NewJsonConfigFile(filename, file.LocalFolderLocator, file.NewPlaceholdersReplacerReader(ioutil.ReadFile, props))
}

func NewCommonJsonConfigFile(filename string, subfolder string) config.Configuration {
	return NewJsonConfigFile(filename, file.GetCommonLocationsLocator(true, true, true, subfolder), ioutil.ReadFile)
}

func NewJsonConfigFileWithPlaceholders(filename string, subfolder string, props config.Configuration) config.Configuration {
	return NewJsonConfigFile(filename, file.GetCommonLocationsLocator(true, true, true, subfolder), file.NewPlaceholdersReplacerReader(ioutil.ReadFile, props))
}
