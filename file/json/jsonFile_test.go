package json

import (
	"github.com/gorigin/config/file"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewJsonConfigFile(t *testing.T) {
	assert := assert.New(t)

	// Building mock data
	reader := file.NewStringReader(`
{
  "dsn": "guest:guest@localhost",
  "credentials": [
    {"role": "admin", "enabled": true, "id": 10}
  ]
}`)
	json := NewJsonConfigFile("anything", file.NoopLocator, reader)
	q, err := json.Qualifiers()

	assert.NoError(err)
	assert.Len(q, 2)
	assert.True(json.Has("dsn"))
	assert.True(json.Has("credentials"))

	var s string
	assert.NoError(json.Configure("dsn", &s))
	assert.Equal("guest:guest@localhost", s)

	var creds []struct {
		ID       int    `json:"id"`
		Role     string `json:"role"`
		Approved bool   `json:"enabled"`
	}
	assert.NoError(json.Configure("credentials", &creds))
	assert.Len(creds, 1)
	assert.Equal(10, creds[0].ID)
	assert.True(creds[0].Approved)
	assert.Equal("admin", creds[0].Role)
}
