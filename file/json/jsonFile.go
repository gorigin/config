package json

import (
	"encoding/json"
	"fmt"
	"github.com/gorigin/config"
	"github.com/gorigin/config/file"
)

func jsonFileValuesFiller(bts []byte) (map[string]interface{}, error) {
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

func jsonFileValuesMapper(source interface{}, target interface{}) error {
	raw, ok := source.(json.RawMessage)
	if !ok {
		return fmt.Errorf("Value must be RawMessage")
	}

	return json.Unmarshal(raw, target)
}

// New returns new .json file configuration source
func New(options file.Options) config.Configuration {
	return file.NewFileConfiguration(
		file.FullOptions{
			Options:          options.WithDefaults(),
			ByteToMapReader:  jsonFileValuesFiller,
			ReflectionMapper: jsonFileValuesMapper,
		},
	)
}
