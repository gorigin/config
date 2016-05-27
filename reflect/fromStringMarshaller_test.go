package reflect

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFromStringMarshallerNotPointer(t *testing.T) {
	assert := assert.New(t)

	var to string
	assert.Error(FromStringMarshaller("foo", to))
	assert.NoError(FromStringMarshaller("foo", &to))
}

func TestFromStringMarshallerSuccess(t *testing.T) {
	assert := assert.New(t)

	var s string
	assert.NoError(FromStringMarshaller("foo", &s))
	assert.Equal("foo", s)

	for _, v := range []string{"true", "True", "tRue", "yes", "on", "1"} {
		b := false
		assert.NoError(FromStringMarshaller(v, &b))
		assert.True(b, "Must be casted to TRUE: "+v)
	}
	for _, v := range []string{"false", "", "0", "no", "off"} {
		b := true
		assert.NoError(FromStringMarshaller(v, &b))
		assert.False(b, "Must be casted to FALSE: "+v)
	}

	var i int
	assert.NoError(FromStringMarshaller("12", &i))
	assert.Equal(12, i)
	assert.NoError(FromStringMarshaller("0", &i))
	assert.Equal(0, i)
	assert.NoError(FromStringMarshaller("-8", &i))
	assert.Equal(-8, i)
	assert.Error(FromStringMarshaller("", &i))

	var i64 int64
	assert.NoError(FromStringMarshaller("12", &i64))
	assert.Equal(int64(12), i64)
	assert.NoError(FromStringMarshaller("0", &i64))
	assert.Equal(int64(0), i64)
	assert.NoError(FromStringMarshaller("-8", &i64))
	assert.Equal(int64(-8), i64)
	assert.Error(FromStringMarshaller("", &i64))

	var f64 float64
	assert.NoError(FromStringMarshaller(".21", &f64))
	assert.Equal(0.21, f64)
	assert.NoError(FromStringMarshaller("21.00006", &f64))
	assert.Equal(21.00006, f64)
	assert.NoError(FromStringMarshaller("-0.0000000000001", &f64))
	assert.Equal(-0.0000000000001, f64)
}
