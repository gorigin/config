package yaml

import (
	"fmt"
	"github.com/gorigin/config"
	"github.com/gorigin/config/file"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
)

func YamlFileValuesFiller(bts []byte) (map[string]interface{}, error) {
	var target yaml.MapSlice
	err := yaml.Unmarshal(bts, &target)
	if err != nil {
		return nil, err
	}

	// Converting to map[string]interface{}
	response := map[string]interface{}{}
	for _, item := range target {
		out, err := yaml.Marshal(item.Value)
		if err != nil {
			return nil, err
		}
		if lo := len(out); lo > 0 && out[lo-1] == 10 {
			out = out[0 : lo-1]
		}
		response[fmt.Sprintf("%v", item.Key)] = out
	}

	return response, nil
}

func YamlFileValuesMapper(source interface{}, target interface{}) error {
	raw, ok := source.([]byte)
	if !ok {
		return fmt.Errorf("Value must be byte slice")
	}

	return yaml.Unmarshal(raw, target)
}

func NewYamlConfigFile(filename string, l file.FileLocator, r file.FileReader) config.Configuration {
	return file.NewFileConfiguration(filename, l, r, YamlFileValuesFiller, YamlFileValuesMapper)
}

func NewLocalCommonYamlConfigFile(filename string) config.Configuration {
	return NewYamlConfigFile(filename, file.LocalFolderLocator, ioutil.ReadFile)
}

func NewLocalYamlConfigWithPlaceholders(filename string, props config.Configuration) config.Configuration {
	return NewYamlConfigFile(filename, file.LocalFolderLocator, file.NewPlaceholdersReplacerReader(ioutil.ReadFile, props))
}

func NewCommonYamlConfigFile(filename string, subfolder string) config.Configuration {
	return NewYamlConfigFile(filename, file.GetCommonLocationsLocator(true, true, true, subfolder), ioutil.ReadFile)
}

func NewYamlConfigFileWithPlaceholders(filename string, subfolder string, props config.Configuration) config.Configuration {
	return NewYamlConfigFile(filename, file.GetCommonLocationsLocator(true, true, true, subfolder), file.NewPlaceholdersReplacerReader(ioutil.ReadFile, props))
}
