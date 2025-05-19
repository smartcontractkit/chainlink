package verification

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"

	ds "github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/verifier_proxy_v0_5_0"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

func TestInitializeVerifier(t *testing.T) {
	t.Parallel()
	e := testutil.NewMemoryEnv(t, true, 0)

	chainSelector := e.AllChainSelectors()[0]
	e, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			DeployVerifierProxyChangeset,
			DeployVerifierProxyConfig{
				ChainsToDeploy: map[uint64]DeployVerifierProxy{
					chainSelector: {AccessControllerAddress: common.Address{}},
				},
				Version: deployment.Version0_5_0,
			},
		),
	)
	require.NoError(t, err)

	// Ensure the VerifierProxy was deployed
	record, err := e.DataStore.Addresses().Get(ds.NewAddressRefKey(chainSelector, ds.ContractType(types.VerifierProxy), &deployment.Version0_5_0, ""))
	require.NoError(t, err)
	verifierProxyAddr := common.HexToAddress(record.Address)

	// Deploy Verifier
	e, err = commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			DeployVerifierChangeset,
			DeployVerifierConfig{
				ChainsToDeploy: map[uint64]DeployVerifier{
					chainSelector: {VerifierProxyAddress: verifierProxyAddr},
				},
			},
		),
	)
	require.NoError(t, err)

	// Ensure the Verifier was deployed
	record, err = e.DataStore.Addresses().Get(ds.NewAddressRefKey(chainSelector, ds.ContractType(types.Verifier), &deployment.Version0_5_0, ""))
	require.NoError(t, err)
	verifierAddr := common.HexToAddress(record.Address)

	e, err = commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			InitializeVerifierChangeset,
			VerifierProxyInitializeVerifierConfig{
				ConfigPerChain: map[uint64][]InitializeVerifierConfig{
					chainSelector: {
						{VerifierAddress: verifierAddr, VerifierProxyAddress: verifierProxyAddr},
					},
				},
			},
		),
	)
	require.NoError(t, err)

	chain := e.Chains[chainSelector]

	vp, err := verifier_proxy_v0_5_0.NewVerifierProxy(verifierProxyAddr, chain.Client)
	require.NoError(t, err)
	logIterator, err := vp.FilterVerifierInitialized(nil)
	require.NoError(t, err)

	foundExpected := false
	for logIterator.Next() {
		if verifierAddr == logIterator.Event.VerifierAddress {
			foundExpected = true
			break
		}
	}
	require.True(t, foundExpected)
	err = logIterator.Close()
	if err != nil {
		t.Errorf("Error closing log iterator: %v", err)
	}
}
