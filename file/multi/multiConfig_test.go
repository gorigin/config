package multi

import (
	"encoding/json"
	"fmt"
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
  "secondary": "%someValue%"
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

	reader := func(name string) ([]byte, error) {
		t.Log("Reading", name)
		v, ok := files[name]
		if !ok {
			return nil, fmt.Errorf("File %s not found", name)
		}

		return []byte(v), nil
	}

	c, err := NewMultifileConfig([]string{"first.ini", "main.json", "app.yaml", "second.ini.asc"}, file.NoopLocator)
	assert.NoError(err)
	assert.NotNil(c)

	// Replacing reader
	rc, ok := c.(*mfc)
	assert.True(ok)
	rc.reader = reader
	rc.passwordPrompter = func(filename string) ([]byte, error) {
		if filename == "second.ini.asc" {
			return []byte("123"), nil
		}

		return nil, fmt.Errorf("Invalid file %s", filename)
	}

	// Reading
	assert.NoError(rc.Reload())
	q, err := rc.Qualifiers()

	// Asserting
	assert.Len(q, 8)
	valueEq := func(key string, value interface{}) {
		v, ok := rc.Value(key)
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

	valueEq("variable", "334")
	valueEq("firstIni:variable", "334")
	valueEq("threshold", "334")
	valueEq("name", "admin")
	valueEq("tertiary", "admin")
}
