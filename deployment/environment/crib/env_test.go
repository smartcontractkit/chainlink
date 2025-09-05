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
	config, err := env.GetConfig(DeployerKeys{
		EVMKey:   "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		SolKey:   "57qbvFjTChfNwQxqkFZwjHp7xYoPZa7f9ow6GA59msfCH1g6onSjKUTrrLp4w1nAwbwQuit8YgJJ2AwT9BSwownC",
		AptosKey: "0x906b8a983b434318ca67b7eff7300f91b02744c84f87d243d2fbc3e528414366",
	})
	require.NoError(t, err)
	require.NotNil(t, config)
	assert.NotEmpty(t, config.NodeIDs)
	assert.NotNil(t, config.AddressBook)
	assert.NotEmpty(t, config.Chains)
}
