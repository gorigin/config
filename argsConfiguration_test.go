package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewArgsConfiguration(t *testing.T) {
	assert := assert.New(t)

	// Default
	c := NewArgsConfiguration([]string{"-foo", "--bar=baz"}, ArgsConfigurationOptions{})
	q, err := c.Qualifiers()
	assert.NoError(err)
	assert.Len(q, 2)

	// Custom prefix options
	c = NewArgsConfiguration([]string{"-configure.foo", "--configure.bar=baz", "--baz=something"}, ArgsConfigurationOptions{Prefix: "configure."})
	q, err = c.Qualifiers()
	assert.NoError(err)
	assert.Len(q, 2)
	assert.True(c.Has("foo"))
	assert.True(c.Has("bar"))
	assert.False(c.Has("bazm"))

	valueEq := func(key string, value interface{}) {
		v, ok := c.Value(key)
		assert.True(ok)
		assert.Equal(value, v)
	}

	valueEq("foo", true)
	valueEq("bar", "baz")

	// Verbose&quiet options
	c = NewArgsConfiguration([]string{"-foo", "--bar=baz", "-q"}, ArgsConfigurationOptions{
		VerboseAndQuiet:       true,
		VerboseAndQuietPrefix: "default.",
	})

	q, err = c.Qualifiers()
	assert.NoError(err)
	assert.Len(q, 4)

	valueEq("foo", true)
	valueEq("bar", "baz")
	valueEq("q", true)
	valueEq("default.verbosity", -1)

	// VerboseAndQuietPrefix
	c = NewArgsConfiguration([]string{"-foo", "--somebar=baz", "-vv"}, ArgsConfigurationOptions{
		VerboseAndQuiet:       true,
		VerboseAndQuietPrefix: "default.",
		Prefix:                "notexist",
	})

	q, err = c.Qualifiers()
	assert.NoError(err)
	assert.Len(q, 1)

	assert.False(c.Has("foo"))
	assert.False(c.Has("vv"))
	valueEq("default.verbosity", 2)
}
