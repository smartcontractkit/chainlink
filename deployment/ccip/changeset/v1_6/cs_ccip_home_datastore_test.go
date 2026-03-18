package v1_6_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
)

func TestSetCandidateAcceptsUSDCTokenPoolProxyFromDataStore(t *testing.T) {
	t.Parallel()

	tenv, _ := testhelpers.NewMemoryEnvironment(t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithNumOfNodes(4))
	dest := remoteChainSelector(
		tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM)),
		tenv.HomeChainSel,
	)
	proxyAddress := common.HexToAddress("0x1000000000000000000000000000000000000001")

	ds := datastore.NewMemoryDataStore()
	require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: dest,
		Address:       proxyAddress.Hex(),
		Type:          datastore.ContractType(shared.USDCTokenPoolProxy),
		Version:       &deployment.Version1_7_0,
		Qualifier:     "proxy-only-in-datastore",
	}))
	tenv.Env.DataStore = ds.Seal()

	_, err := commonchangeset.Apply(t, tenv.Env,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.SetCandidateChangeset),
			setCandidateExecConfig(tenv.HomeChainSel, tenv.FeedChainSel, dest, proxyAddress.Hex()),
		),
	)
	require.NoError(t, err)
}

func TestSetCandidateErrorsOnDuplicateUSDCTokenPoolProxyInDataStore(t *testing.T) {
	t.Parallel()

	tenv, _ := testhelpers.NewMemoryEnvironment(t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithNumOfNodes(4))
	dest := remoteChainSelector(
		tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM)),
		tenv.HomeChainSel,
	)

	ds := datastore.NewMemoryDataStore()
	for i, ref := range []struct {
		address   string
		qualifier string
	}{
		{address: "0x1000000000000000000000000000000000000001", qualifier: "proxy-a"},
		{address: "0x2000000000000000000000000000000000000002", qualifier: "proxy-b"},
	} {
		require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
			ChainSelector: dest,
			Address:       common.HexToAddress(ref.address).Hex(),
			Type:          datastore.ContractType(shared.USDCTokenPoolProxy),
			Version:       &deployment.Version1_7_0,
			Qualifier:     ref.qualifier,
		}), "add datastore ref %d", i)
	}
	tenv.Env.DataStore = ds.Seal()

	_, err := commonchangeset.Apply(t, tenv.Env,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.SetCandidateChangeset),
			setCandidateExecConfig(tenv.HomeChainSel, tenv.FeedChainSel, dest, "0x1000000000000000000000000000000000000001"),
		),
	)
	require.Error(t, err)
	require.ErrorContains(t, err, "multiple datastore entries found for USDCTokenPoolProxy 1.7.0")
}

func setCandidateExecConfig(homeChainSel, feedChainSel, dest uint64, sourcePoolAddress string) v1_6.SetCandidateChangesetConfig {
	return v1_6.SetCandidateChangesetConfig{
		SetCandidateConfigBase: v1_6.SetCandidateConfigBase{
			HomeChainSelector: homeChainSel,
			FeedChainSelector: feedChainSel,
		},
		PluginInfo: []v1_6.SetCandidatePluginInfo{
			{
				PluginType: types.PluginTypeCCIPExec,
				OCRConfigPerRemoteChainSelector: map[uint64]v1_6.CCIPOCRParams{
					dest: v1_6.DeriveOCRParamsForExec(v1_6.SimulationTest, []pluginconfig.TokenDataObserverConfig{
						{
							Type:    pluginconfig.USDCCCTPHandlerType,
							Version: "1.0",
							USDCCCTPObserverConfig: &pluginconfig.USDCCCTPObserverConfig{
								AttestationConfig: pluginconfig.AttestationConfig{
									AttestationAPI: "http://example.com",
								},
								Tokens: map[ccipocr3.ChainSelector]pluginconfig.USDCCCTPTokenConfig{
									ccipocr3.ChainSelector(dest): {
										SourcePoolAddress:            sourcePoolAddress,
										SourceMessageTransmitterAddr: common.HexToAddress("0x3000000000000000000000000000000000000003").Hex(),
									},
								},
							},
						},
					}, nil),
				},
			},
		},
	}
}

func remoteChainSelector(selectors []uint64, exclude uint64) uint64 {
	for _, selector := range selectors {
		if selector != exclude {
			return selector
		}
	}
	return 0
}
