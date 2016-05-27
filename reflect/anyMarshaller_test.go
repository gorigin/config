package reflect

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAnyMarshaller(t *testing.T) {
	assert := assert.New(t)

	si32 := int32(5000)
	var i32 int32
	assert.NoError(AnyMarshaller(si32, &i32))
	assert.Equal(si32, i32)

	var i int
	assert.NoError(AnyMarshaller("-18", &i))
	assert.Equal(-18, i)
}
