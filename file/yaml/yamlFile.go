package yaml

import (
	"fmt"
	"github.com/gorigin/config"
	"github.com/gorigin/config/file"
	"gopkg.in/yaml.v2"
)

func yamlFileValuesFiller(bts []byte) (map[string]interface{}, error) {
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

func yamlFileValuesMapper(source interface{}, target interface{}) error {
	raw, ok := source.([]byte)
	if !ok {
		return fmt.Errorf("Value must be byte slice")
	}

	return yaml.Unmarshal(raw, target)
}

// New returns new .yaml (and .yml) file configuration source
func New(options file.Options) config.Configuration {
	return file.NewFileConfiguration(
		file.FullOptions{
			Options:          options.WithDefaults(),
			ByteToMapReader:  yamlFileValuesFiller,
			ReflectionMapper: yamlFileValuesMapper,
		},
	)
}
