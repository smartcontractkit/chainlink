package crib

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func TestShouldProvideEnvironmentConfig(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/DX-291")

	t.Parallel()
	env := NewDevspaceEnvFromStateDir(nil, "testdata/lanes-deployed-state")
	config, err := env.GetConfig("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	require.NoError(t, err)
	require.NotNil(t, config)
	assert.NotEmpty(t, config.NodeIDs)
	assert.NotNil(t, config.AddressBook)
	assert.NotEmpty(t, config.Chains)
}
