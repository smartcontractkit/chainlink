package modulemap

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModulePath(t *testing.T) {
	t.Parallel()
	mod, err := ModulePath("mockery")
	require.NoError(t, err)
	assert.Equal(t, "github.com/vektra/mockery/v2", mod)

	_, err = ModulePath("protoc")
	require.Error(t, err)
}
