package multi

import (
	"encoding/json"
	"fmt"
	"github.com/gorigin/config"
	"github.com/gorigin/config/file"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMultiConfig(t *testing.T) {
	assert := assert.New(t)

	files := map[string]string{
		"first.ini": "[firstIni]\nvariable=334",
		"main.json": `
{
  "threshold": "!firstIni:variable!",
  "secondary": "%someValue%",
  "replaced": "foo"
}
		`,
		"app.yaml": `
repositories:
  - one
  - two
someValue: 300
tertiary: "%name%"`,
		"second.ini.asc": `
-----BEGIN PGP MESSAGE-----
Version: GnuPG v1

jA0ECQMCbx0OsFniBm5g0kYBPzfC2vJofXPD9Ogqss3Y/Di7hAKBCpjrC2A4qWRl
CRz5s258W0k3bRGgoy5b0b9l3m4n587Gf2o4IzBAqpIvnPjLzfbe
=Pb9G
-----END PGP MESSAGE-----`,
	}

	c := New(Options{
		Filenames: []string{"first.ini", "main.json", "app.yaml", "second.ini.asc"},
		Locator:   file.NoopLocator,
		Reader: func(name string) ([]byte, error) {
			t.Log("Reading", name)
			v, ok := files[name]
			if !ok {
				return nil, fmt.Errorf("File %s not found", name)
			}

			return []byte(v), nil
		},
		Prompter: func(filename string) ([]byte, error) {
			if filename == "second.ini.asc" {
				return []byte("123"), nil
			}

			return nil, fmt.Errorf("Invalid file %s", filename)
		},
		Prepend: []config.Configuration{
			config.MapConfiguration(map[string]interface{}{
				"prependedValue": 33.02,
				"someValue":      100, // Will be replaced
			}),
		},
		Append: []config.Configuration{
			config.MapConfiguration(map[string]interface{}{
				"replaced": "hello, world",
			}),
		},
	})

	assert.NotNil(c)

	// Reading
	assert.NoError(c.Test())
	q, err := c.Qualifiers()
	assert.NoError(err)

	// Asserting
	assert.Len(q, 10)
	valueEq := func(key string, value interface{}) {
		v, ok := c.Value(key)
		if ok {
			if r, k := v.(json.RawMessage); k {
				v = string(r)
			} else if b, k := v.([]byte); k {
				v = string(b)
			}
		}
		assert.True(ok)
		assert.Equal(value, v)
	}

	valueEq("someValue", "300")
	valueEq("variable", "334")
	valueEq("firstIni:variable", "334")
	valueEq("threshold", "334")
	valueEq("name", "admin")
	valueEq("tertiary", "admin")
	valueEq("replaced", "hello, world")
}
