package tron_test

import (
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cldf_tron "github.com/smartcontractkit/chainlink-deployments-framework/chain/tron"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
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

		deployChangeset := commonchangeset.Configure(tron.DeployForwarder{},
			&tron.DeployForwarderRequest{
				ChainSelectors: []uint64{registrySel},
				Qualifier:      "my-test-forwarder",
				DeployOptions:  deployOptions,
			},
		)

		// deploy
		var err error
		_, resp, err := commonchangeset.ApplyChangesets(t, env, []commonchangeset.ConfiguredChangeSet{deployChangeset})
		require.NoError(t, err)
		require.NotNil(t, resp)
	})
}
