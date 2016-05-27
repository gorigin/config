package file

import (
	"github.com/gorigin/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPlaceholdersReplacerReader(t *testing.T) {
	assert := assert.New(t)

	original := `foo="%bar%"
bar="!baz!"
baz=none`

	m := map[string]interface{}{
		"bar": "string",
		"baz": 100500,
	}

	r := NewPlaceholdersReplacerReader(NewStringReader(original), config.MapConfiguration(m))
	bts, err := r("anythinf")
	assert.NoError(err)
	assert.Equal(`foo="string"
bar=100500
baz=none`, string(bts))
}
