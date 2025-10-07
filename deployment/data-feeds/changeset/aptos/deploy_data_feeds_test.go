package aptos_test

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/aptos-labs/aptos-go-sdk"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	aptosCS "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/aptos"
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
	chain := env.BlockChains.AptosChains()[chainSelector]
	platform1, err := aptosCS.DeployPlatform(chain, aptos.AccountAddress{}, []string{})
	require.NoError(t, err)
	platform2, err := aptosCS.DeployPlatformSecondary(chain, aptos.AccountAddress{}, []string{})
	require.NoError(t, err)

	resp, err := commonChangesets.Apply(t, env, commonChangesets.Configure(
		aptosCS.DeployDataFeedsChangeset,
		types.DeployAptosConfig{
			ChainsToDeploy:           []uint64{chainSelector},
			PlatformAddress:          platform1.Address.String(),
			SecondaryPlatformAddress: platform2.Address.String(),
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
