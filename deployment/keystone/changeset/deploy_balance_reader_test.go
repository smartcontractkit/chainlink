package changeset_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

func TestDeployBalanceReader(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Chains: 2,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)

	registrySel := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]

	t.Run("should deploy balancereader", func(t *testing.T) {
		ab := cldf.NewMemoryAddressBook()
		qualifier := "my-balance-reader-qualifier"

		// deploy balancereader
		env.ExistingAddresses = ab
		resp, err := changeset.DeployBalanceReader(env, changeset.DeployBalanceReaderRequest{
			Qualifier: qualifier,
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		// registry, ocr3, balancereader should be deployed on registry chain
		addrs, err := resp.AddressBook.AddressesForChain(registrySel)
		require.NoError(t, err)
		assert.Len(t, addrs, 1)

		dsAddrs, err := resp.DataStore.Addresses().Fetch()
		require.NoError(t, err)
		assert.Len(t, dsAddrs, 2) // 2 balance readers, one per chain
		assert.Equal(t, qualifier, dsAddrs[0].Qualifier)
		assert.Equal(t, qualifier, dsAddrs[1].Qualifier)

		// only balancereader on chain 1
		require.NotEqual(t, registrySel, env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[1])
		oaddrs, err := resp.AddressBook.AddressesForChain(env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[1])
		require.NoError(t, err)
		assert.Len(t, oaddrs, 1)
	})
}
