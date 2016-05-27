package yaml

import (
	"github.com/gorigin/config/file"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewYamlConfigFile(t *testing.T) {
	assert := assert.New(t)

	// Building mock data
	reader := file.NewStringReader(`
threshold: -8.21
debug: true
dsn: guest:guest@localhost
credentials:
  -
    role: admin
    # Comments are ok
    id: 10 # More comments
    enabled: true
`)
	yaml := New(file.Options{Reader: reader, Locator: file.NoopLocator})
	q, err := yaml.Qualifiers()

	assert.NoError(err)
	assert.Len(q, 4)
	assert.True(yaml.Has("dsn"))
	assert.True(yaml.Has("credentials"))
	assert.True(yaml.Has("threshold"))
	assert.True(yaml.Has("debug"))

	var s string
	assert.NoError(yaml.Configure("dsn", &s))
	assert.Equal("guest:guest@localhost", s)

	var creds []struct {
		ID       int    `yaml:"id"`
		Role     string `yaml:"role"`
		Approved bool   `yaml:"enabled"`
	}
	assert.NoError(yaml.Configure("credentials", &creds))
	assert.Len(creds, 1)
	assert.Equal(10, creds[0].ID)
	assert.True(creds[0].Approved)
	assert.Equal("admin", creds[0].Role)

	var f float64
	assert.NoError(yaml.Configure("threshold", &f))
	assert.Equal(-8.21, f)

	var b bool
	assert.NoError(yaml.Configure("debug", &b))
	assert.True(b)
}
