package aptos_test

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/aptos"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func TestDeployAptosCache(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		AptosChains: 1,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)

	chainSelector := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyAptos))[0]

	resp, err := commonChangesets.Apply(t, env, commonChangesets.Configure(
		aptos.DeployDataFeedsChangeset,
		types.DeployAptosConfig{
			ChainsToDeploy:           []uint64{chainSelector},
			OwnerAddress:             "0x0000000000000000000000000000000000000000",
			PlatformAddress:          "0x0000000000000000000000000000000000000001",
			SecondaryPlatformAddress: "0x0000000000000000000000000000000000000002",
			Qualifier:                "aptos",
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
			"aptos",
		))
	require.NoError(t, err)
	require.NotNil(t, addrs.Address)
	require.Equal(t, datastore.ContractType("DataFeedsCache"), addrs.Type)
	require.Equal(t, "aptos", addrs.Qualifier)
}
