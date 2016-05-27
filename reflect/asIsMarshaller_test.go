package reflect

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAsIsMarshallerPtr(t *testing.T) {
	assert := assert.New(t)

	v := "hello"

	var s string
	assert.Error(AsIsMarshaller(v, s))
	assert.Error(AsIsMarshaller(&v, &s))
	assert.NoError(AsIsMarshaller(v, &s))
}

func TestAsIsMarshaller(t *testing.T) {
	assert := assert.New(t)

	var s string
	assert.NoError(AsIsMarshaller("hello", &s))
	assert.Equal("hello", s)

	var i int
	assert.NoError(AsIsMarshaller(8012, &i))
	assert.Equal(8012, i)
}
