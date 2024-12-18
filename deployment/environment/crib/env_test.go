package crib

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldProvideEnvironmentConfig(t *testing.T) {
	t.Parallel()
	env := NewDevspaceEnvFromStateDir("testdata/lanes-deployed-state")
	config := env.GetConfig()
	require.NotNil(t, config)
	assert.NotEmpty(t, config.NodeIDs)
	assert.NotNil(t, config.AddressBook)
	assert.NotEmpty(t, config.Chains)
}
