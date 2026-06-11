package modulemap

import (
	"sort"
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

func TestModulesSorted(t *testing.T) {
	t.Parallel()
	mods := Modules()
	assert.True(t, sort.StringsAreSorted(mods), "Modules() output must be sorted; got %v", mods)
	// calling twice must return same order
	assert.Equal(t, mods, Modules())
}
