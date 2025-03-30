package sequence

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commonstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/verification"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier_v0_5_0"
	"github.com/stretchr/testify/require"
)

func TestDeployDataStreams(t *testing.T) {
	e := testutil.NewMemoryEnv(t, false, 0)

	// Need the Link Token
	e, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			deployment.CreateLegacyChangeSet(commonChangesets.DeployLinkToken),
			[]uint64{testutil.TestChain.Selector},
		),
	)
	require.NoError(t, err)

	addresses, err := e.ExistingAddresses.AddressesForChain(testutil.TestChain.Selector)
	require.NoError(t, err)

	chain := e.Chains[testutil.TestChain.Selector]
	linkState, err := commonstate.MaybeLoadLinkTokenChainState(chain, addresses)
	require.NoError(t, err)

	// SetConfig settings
	configDigest := [32]byte{1}

	signers := []common.Address{
		common.HexToAddress("0x1111111111111111111111111111111111111111"),
		common.HexToAddress("0x2222222222222222222222222222222222222222"),
		common.HexToAddress("0x3333333333333333333333333333333333333333"),
		common.HexToAddress("0x4444444444444444444444444444444444444444"),
	}
	f := uint8(1)

	cfg := DeployDataStreamsConfig{
		ChainsToDeploy: map[uint64]DeployDataStreams{testutil.TestChain.Selector: {
			VerifierConfig: verification.SetConfig{
				ConfigDigest:               configDigest,
				Signers:                    signers,
				F:                          f,
				RecipientAddressesAndProps: []verifier_v0_5_0.CommonAddressAndWeight{},
			},
			BillingConfig: &BillingConfig{
				LinkTokenAddress:   linkState.LinkToken.Address(),
				NativeTokenAddress: common.HexToAddress("0x3e5e9111ae8eb78fe1cc3bb8915d5d461f3ef9a9"),
			},
		}}}

	resp, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(DeployDataStreamsChangeset, cfg),
	)

	require.NoError(t, err)

	var expectedContracts = []deployment.ContractType{types.VerifierProxy, types.Verifier, types.FeeManager, types.RewardManager}

	// Check the address book for all contract existence
	for _, contract := range expectedContracts {
		contractAddrHex, err := deployment.SearchAddressBook(resp.ExistingAddresses, testutil.TestChain.Selector, contract)
		require.NoError(t, err, "failed to find %s address in address book", contract)
		contractAddr := common.HexToAddress(contractAddrHex)
		require.NotEqual(t, common.Address{}, contractAddr, "%s should not be zero address", contract)
	}
}
