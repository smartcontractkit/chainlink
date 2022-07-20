package tomlutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtils_TomlFloat32_Success_Decimal(t *testing.T) {
	t.Parallel()

	var tomlF32 Float32

	err := tomlF32.UnmarshalText([]byte("0.23"))

	assert.Nil(t, err)
	assert.Equal(t, tomlF32, Float32(0.23))
}

func TestUtils_TomlFloat32_Success_Integer(t *testing.T) {
	t.Parallel()

	var tomlF32 Float32

	err := tomlF32.UnmarshalText([]byte("13"))

	assert.Nil(t, err)
	assert.Equal(t, tomlF32, Float32(13))
}

func TestUtils_TomlFloat32_Failure(t *testing.T) {
	t.Parallel()

	var tomlF32 Float32

	err := tomlF32.UnmarshalText([]byte("1s"))

	assert.NotNil(t, err)
}

func TestUtils_TomlFloat64_Success_Decimal(t *testing.T) {
	t.Parallel()

	var tomlF64 Float64

	err := tomlF64.UnmarshalText([]byte("2.82"))

	assert.Nil(t, err)
	assert.Equal(t, tomlF64, Float64(2.82))
}

func TestUtils_TomlFloat64_Success_Integer(t *testing.T) {
	t.Parallel()

	var tomlF64 Float64

	err := tomlF64.UnmarshalText([]byte("3"))

	assert.Nil(t, err)
	assert.Equal(t, tomlF64, Float64(3))
}

func TestUtils_TomlFloat64_Failure(t *testing.T) {
	t.Parallel()

	var tomlF64 Float64

	err := tomlF64.UnmarshalText([]byte("1s"))

	assert.NotNil(t, err)
}
