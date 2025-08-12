package tron_test

import (
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	cldf_tron "github.com/smartcontractkit/chainlink-deployments-framework/chain/tron"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/tron"
)

func TestDeployForwarder(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		TronChains: 1,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)

	registrySel := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyTron))[0]

	t.Run("should deploy forwarder", func(t *testing.T) {
		ab := cldf.NewMemoryAddressBook()

		// deploy forwarder
		env.ExistingAddresses = ab

		deployOptions := cldf_tron.DefaultDeployOptions()
		deployOptions.FeeLimit = 1_000_000_000

		configuredChangeset := commonchangeset.Configure(tron.DeployForwarder{},
			&tron.DeployForwarderRequest{
				ChainSelectors: []uint64{registrySel},
				Qualifier:      "my-test-forwarder",
				DeployOptions:  deployOptions,
			},
		)

		// deploy
		var err error
		_, resp, err := commonchangeset.ApplyChangesets(t, env, []commonchangeset.ConfiguredChangeSet{configuredChangeset})
		require.NoError(t, err)
		require.NotNil(t, resp)

		// registry, ocr3, forwarder should be deployed on registry chain
		addrs, err := resp[0].AddressBook.AddressesForChain(registrySel)
		require.NoError(t, err)
		require.Len(t, addrs, 1)
		fa := resp[0].DataStore.Addresses().Filter(datastore.AddressRefByQualifier("my-test-forwarder"))
		require.Len(t, fa, 1, "expected to find 'my-test-forwarder' qualifier")
		l := fa[0].Labels.List()
		require.Len(t, l, 2, "expected exactly 2 labels")
		require.Contains(t, l[0], internal.DeploymentBlockLabel)
		require.Contains(t, l[1], internal.DeploymentHashLabel)
	})
}
