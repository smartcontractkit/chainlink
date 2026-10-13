package tomlutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUtils_TomlFloat32_Success_Decimal(t *testing.T) {
	t.Parallel()

	var tomlF32 Float32

	err := tomlF32.UnmarshalText([]byte("0.23"))

	require.NoError(t, err)
	assert.InEpsilon(t, float32(0.23), float32(tomlF32), 1e-9)
}

func TestUtils_TomlFloat32_Success_Integer(t *testing.T) {
	t.Parallel()

	var tomlF32 Float32

	err := tomlF32.UnmarshalText([]byte("13"))

	require.NoError(t, err)
	assert.InEpsilon(t, float32(13), float32(tomlF32), 1e-9)
}

func TestUtils_TomlFloat32_Failure(t *testing.T) {
	t.Parallel()

	var tomlF32 Float32

	err := tomlF32.UnmarshalText([]byte("1s"))

	assert.Error(t, err)
}

func TestUtils_TomlFloat64_Success_Decimal(t *testing.T) {
	t.Parallel()

	var tomlF64 Float64

	err := tomlF64.UnmarshalText([]byte("2.82"))

	require.NoError(t, err)
	assert.InEpsilon(t, float64(2.82), float64(tomlF64), 1e-9)
}

func TestUtils_TomlFloat64_Success_Integer(t *testing.T) {
	t.Parallel()

	var tomlF64 Float64

	err := tomlF64.UnmarshalText([]byte("3"))

	require.NoError(t, err)
	assert.InEpsilon(t, float64(3), float64(tomlF64), 1e-9)
}

func TestUtils_TomlFloat64_Failure(t *testing.T) {
	t.Parallel()

	var tomlF64 Float64

	err := tomlF64.UnmarshalText([]byte("1s"))

	assert.Error(t, err)
}
