package sequence

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commonstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/verification"
	dsutil "github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier_v0_5_0"
	"github.com/stretchr/testify/require"
)

func TestDeployDataStreams(t *testing.T) {
	testEnv := testutil.NewMemoryEnvV2(t, testutil.MemoryEnvConfig{DeployMCMS: true})

	// Need the Link Token
	e, err := commonChangesets.Apply(t, testEnv.Environment, nil,
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
			Billing: BillingFeature{
				Enabled: true,
				Config: &BillingConfig{
					LinkTokenAddress:   linkState.LinkToken.Address(),
					NativeTokenAddress: common.HexToAddress("0x3e5e9111ae8eb78fe1cc3bb8915d5d461f3ef9a9"),
				},
			},
			Ownership: types.OwnershipFeature{
				Transfer: true,
				MCMSProposalConfig: &proposalutils.TimelockConfig{
					MinDelay: 0,
				},
				DeployMCMS: false, // assumes existing MCMS deployment
			},
		}}}

	resp, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(DeployDataStreamsChangeset, cfg),
	)

	require.NoError(t, err)

	var expectedContracts = []deployment.ContractType{types.VerifierProxy, types.Verifier, types.RewardManager, types.FeeManager}

	// Check the address book for all contract existence + ownership
	for _, contract := range expectedContracts {
		contractAddr, err := dsutil.MaybeFindEthAddress(resp.ExistingAddresses, testutil.TestChain.Selector, contract)
		require.NoError(t, err, "failed to find %s address in address book", contract)

		owner, _, err := commonChangesets.LoadOwnableContract(contractAddr, chain.Client)

		require.NoError(t, err)
		require.Equal(t, testEnv.Timelocks[testutil.TestChain.Selector].Timelock.Address(), owner, "%s contract owner should be the MCMS timelock", contract)
	}
}

func TestDeployDataStreamsV2(t *testing.T) {
	testEnv := testutil.NewMemoryEnvV2(t, testutil.MemoryEnvConfig{DeployMCMS: false})

	// Need the Link Token
	e, err := commonChangesets.Apply(t, testEnv.Environment, nil,
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
			Billing: BillingFeature{
				Enabled: true,
				Config: &BillingConfig{
					LinkTokenAddress:   linkState.LinkToken.Address(),
					NativeTokenAddress: common.HexToAddress("0x3e5e9111ae8eb78fe1cc3bb8915d5d461f3ef9a9"),
				},
			},
			Ownership: types.OwnershipFeature{
				Transfer: true,
				MCMSProposalConfig: &proposalutils.TimelockConfig{
					MinDelay: 0,
				},
				DeployMCMS:       true,
				DeployMCMSConfig: proposalutils.SingleGroupTimelockConfigV2(t),
			},
		}}}

	result, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(DeployDataStreamsChangeset, cfg),
	)

	require.NoError(t, err)

	addresses, err = result.ExistingAddresses.AddressesForChain(testutil.TestChain.Selector)
	require.NoError(t, err)
	mcmsState, err := commonstate.MaybeLoadMCMSWithTimelockChainState(chain, addresses)
	require.NoError(t, err)
	require.NotNil(t, mcmsState.Timelock, "MCMS Timelock should be deployed")
	require.NotNil(t, mcmsState.CallProxy, "MCMS CallProxy should be deployed")
	require.NotNil(t, mcmsState.CancellerMcm, "MCMS Canceller should be deployed")
	require.NotNil(t, mcmsState.BypasserMcm, "MCMS Bypasser should be deployed")
	require.NotNil(t, mcmsState.ProposerMcm, "MCMS Proposer should be deployed")

	var expectedContracts = []deployment.ContractType{types.VerifierProxy, types.Verifier, types.RewardManager, types.FeeManager}
	// Check the address book for all contract existence + ownership
	for _, contract := range expectedContracts {
		contractAddr, err := dsutil.MaybeFindEthAddress(result.ExistingAddresses, testutil.TestChain.Selector, contract)
		require.NoError(t, err, "failed to find %s address in address book", contract)

		owner, _, err := commonChangesets.LoadOwnableContract(contractAddr, chain.Client)

		require.NoError(t, err)
		require.Equal(t, mcmsState.Timelock.Address(), owner, "%s contract owner should be the MCMS timelock", contract)
	}
}
