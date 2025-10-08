package aptos_test

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/aptos-labs/aptos-go-sdk"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	aptosCS "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/aptos"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func TestDeployAptosCache(t *testing.T) {
	t.Parallel()

	selector := chain_selectors.APTOS_LOCALNET.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithAptosContainer(t, []uint64{selector}),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	chain := rt.Environment().BlockChains.AptosChains()[selector]

	// deploy platform
	platform1, err := aptosCS.DeployPlatform(chain, aptos.AccountAddress{}, []string{})
	require.NoError(t, err)
	platform2, err := aptosCS.DeployPlatformSecondary(chain, aptos.AccountAddress{}, []string{})
	require.NoError(t, err)

	// deploy cache
	err = rt.Exec(
		runtime.ChangesetTask(aptosCS.DeployDataFeedsChangeset, types.DeployAptosConfig{
			ChainsToDeploy:           []uint64{selector},
			PlatformAddress:          platform1.Address.String(),
			SecondaryPlatformAddress: platform2.Address.String(),
			Qualifier:                "aptos",
		}),
	)
	require.NoError(t, err)

	addrs, err := rt.State().DataStore.Addresses().Get(
		datastore.NewAddressRefKey(
			selector,
			"DataFeedsCache",
			semver.MustParse("1.0.0"),
			"aptos",
		))
	require.NoError(t, err)
	require.NotNil(t, addrs.Address)
	require.Equal(t, datastore.ContractType("DataFeedsCache"), addrs.Type)
	require.Equal(t, "aptos", addrs.Qualifier)
}
