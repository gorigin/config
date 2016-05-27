package ini

import (
	"github.com/gorigin/config/file"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewIniConfigFile(t *testing.T) {
	assert := assert.New(t)

	// Building mock data
	reader := file.NewStringReader(`
# Comments are ignored
foo=bar
[numbers]
bar=10
amount=-0.22

[booleans]
baz=true
yes=yes
enabled=1
	`)

	// Building configuration reader
	ini := NewIniConfigFile("anything", file.NoopLocator, reader)

	q, err := ini.Qualifiers()

	assert.NoError(err)
	assert.Len(q, 11)
	assert.True(ini.Has("foo"))
	assert.True(ini.Has("bar"))
	assert.True(ini.Has("numbers:bar"))
	assert.True(ini.Has("amount"))
	assert.True(ini.Has("numbers:amount"))
	assert.True(ini.Has("baz"))
	assert.True(ini.Has("booleans:baz"))
	assert.True(ini.Has("yes"))
	assert.True(ini.Has("booleans:yes"))
	assert.True(ini.Has("enabled"))
	assert.True(ini.Has("booleans:enabled"))

	valueEq := func(key string, value interface{}) {
		v, ok := ini.Value(key)
		assert.True(ok)
		assert.Equal(value, v)
	}

	valueEq("foo", "bar")
	valueEq("bar", "10")
	valueEq("numbers:bar", "10")
	valueEq("amount", "-0.22")
	valueEq("numbers:amount", "-0.22")
	valueEq("baz", "true")
	valueEq("booleans:baz", "true")
	valueEq("yes", "yes")
	valueEq("booleans:yes", "yes")
	valueEq("enabled", "1")
	valueEq("booleans:enabled", "1")

	var s string
	assert.NoError(ini.Configure("foo", &s))
	assert.Equal("bar", s)

	var i int
	assert.NoError(ini.Configure("bar", &i))
	assert.Equal(10, i)
	var i64 int64
	assert.NoError(ini.Configure("bar", &i64))
	assert.Equal(int64(10), i64)

	var f float64
	assert.NoError(ini.Configure("amount", &f))
	assert.Equal(-0.22, f)

	var b bool
	assert.NoError(ini.Configure("baz", &b))
	assert.True(b)
	b = false
	assert.NoError(ini.Configure("booleans:yes", &b))
	assert.True(b)
	b = false
	assert.NoError(ini.Configure("enabled", &b))
	assert.True(b)
	b = false
}
