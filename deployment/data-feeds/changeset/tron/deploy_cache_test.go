package tron_test

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/tron"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func TestDeployCache(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		TronChains: 1,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)

	chainSelector := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyTron))[0]

	resp, err := commonChangesets.Apply(t, env,
		commonChangesets.Configure(
			tron.DeployCacheChangeset,
			types.DeployTronConfig{
				ChainsToDeploy: []uint64{chainSelector},
				Qualifier:      "tron",
			},
		),
	)
	require.NoError(t, err)
	require.NotNil(t, resp)

	addrs, err := resp.DataStore.Addresses().Get(
		datastore.NewAddressRefKey(
			chainSelector,
			"DataFeedsCache",
			semver.MustParse("1.0.0"),
			"tron",
		))
	require.NoError(t, err)
	require.NotNil(t, addrs.Address)
	require.Equal(t, datastore.ContractType("DataFeedsCache"), addrs.Type)
	require.Equal(t, "tron", addrs.Qualifier)
}
