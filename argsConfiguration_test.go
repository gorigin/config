package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewArgsConfiguration(t *testing.T) {
	assert := assert.New(t)

	// Default options
	c := NewArgsConfiguration([]string{"-foo", "--bar=baz", "-q"}, ArgsConfigurationOptions{})
	q, err := c.Qualifiers()
	assert.NoError(err)
	assert.Len(q, 3)
	assert.True(c.Has("foo"))
	assert.True(c.Has("bar"))
	assert.True(c.Has("q"))

	valueEq := func(key string, value interface{}) {
		v, ok := c.Value(key)
		assert.True(ok)
		assert.Equal(value, v)
	}

	valueEq("foo", true)
	valueEq("bar", "baz")
	valueEq("q", true)

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

	// Only with prefix + defaults
	c = NewArgsConfiguration([]string{"-foo", "--somebar=baz", "-vv"}, ArgsConfigurationOptions{
		VerboseAndQuiet:       true,
		VerboseAndQuietPrefix: "default.",
		OnlyWithPrefix:        "some",
	})

	q, err = c.Qualifiers()
	assert.NoError(err)
	assert.Len(q, 2)

	assert.False(c.Has("foo"))
	assert.False(c.Has("vv"))
	valueEq("somebar", "baz")
	valueEq("default.verbosity", 2)
}
