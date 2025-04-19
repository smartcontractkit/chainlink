package sequence

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"

	"github.com/smartcontractkit/chainlink/deployment"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commonstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/verification"
	dsutil "github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
)

var (
	verificationCfg = verification.SetConfig{
		ConfigDigest: [32]byte{1},
		Signers: []common.Address{
			common.HexToAddress("0x1111111111111111111111111111111111111111"),
			common.HexToAddress("0x2222222222222222222222222222222222222222"),
			common.HexToAddress("0x3333333333333333333333333333333333333333"),
			common.HexToAddress("0x4444444444444444444444444444444444444444"),
		},
		F: 1,
	}
)

func TestDeployDataStreamsContracts(t *testing.T) {
	t.Parallel()
	proposalCfg := proposalutils.SingleGroupTimelockConfigV2(t)
	tests := []struct {
		name                    string
		hasExistingMcms         bool
		deployDataStreamsConfig DeployDataStreamsConfig
		expectedContracts       []deployment.ContractType
	}{
		{
			name:            "Deploy with billing and MCMS",
			hasExistingMcms: false,
			deployDataStreamsConfig: DeployDataStreamsConfig{
				ChainsToDeploy: map[uint64]DeployDataStreams{testutil.TestChain.Selector: {
					VerifierConfig: verificationCfg,
					Billing: types.BillingFeature{
						Enabled: true,
						Config: &types.BillingConfig{
							NativeTokenAddress: common.HexToAddress("0x3e5e9111ae8eb78fe1cc3bb8915d5d461f3ef9a9"),
						},
					},
					Ownership: types.OwnershipFeature{
						ShouldTransfer:     true,
						MCMSProposalConfig: &proposalutils.TimelockConfig{MinDelay: 0},
						ShouldDeployMCMS:   true,
						DeployMCMSConfig:   &proposalCfg,
					},
				}},
			},
			expectedContracts: []deployment.ContractType{types.VerifierProxy, types.Verifier, types.RewardManager, types.FeeManager,
				commontypes.ProposerManyChainMultisig, commontypes.BypasserManyChainMultisig, commontypes.CancellerManyChainMultisig},
		},
		{
			name:            "Deploy no billing and MCMS",
			hasExistingMcms: false,
			deployDataStreamsConfig: DeployDataStreamsConfig{
				ChainsToDeploy: map[uint64]DeployDataStreams{testutil.TestChain.Selector: {
					VerifierConfig: verificationCfg,
					Ownership: types.OwnershipFeature{
						ShouldTransfer:     true,
						MCMSProposalConfig: &proposalutils.TimelockConfig{MinDelay: 0},
						ShouldDeployMCMS:   true,
						DeployMCMSConfig:   &proposalCfg,
					},
				}},
			},
			expectedContracts: []deployment.ContractType{types.VerifierProxy, types.Verifier,
				commontypes.ProposerManyChainMultisig, commontypes.BypasserManyChainMultisig, commontypes.CancellerManyChainMultisig},
		},
		{
			name:            "Deploy no billing with existing MCMS",
			hasExistingMcms: true,
			deployDataStreamsConfig: DeployDataStreamsConfig{
				ChainsToDeploy: map[uint64]DeployDataStreams{testutil.TestChain.Selector: {
					VerifierConfig: verificationCfg,
					Ownership: types.OwnershipFeature{
						ShouldTransfer:     true,
						MCMSProposalConfig: &proposalutils.TimelockConfig{MinDelay: 0},
					},
				}},
			},
			expectedContracts: []deployment.ContractType{types.VerifierProxy, types.Verifier},
		},
		{
			name:            "Deploy but do not propose transfer",
			hasExistingMcms: true,
			deployDataStreamsConfig: DeployDataStreamsConfig{
				ChainsToDeploy: map[uint64]DeployDataStreams{testutil.TestChain.Selector: {
					VerifierConfig: verificationCfg,
				}},
			},
			expectedContracts: []deployment.ContractType{types.VerifierProxy, types.Verifier},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.deployDataStreamsConfig
			billingEnabled := cfg.ChainsToDeploy[testutil.TestChain.Selector].Billing.Enabled

			testEnv := testutil.NewMemoryEnvV2(t, testutil.MemoryEnvConfig{
				ShouldDeployMCMS:      tt.hasExistingMcms,
				ShouldDeployLinkToken: billingEnabled,
			})

			chain := testEnv.Environment.Chains[testutil.TestChain.Selector]

			if cfg.ChainsToDeploy[testutil.TestChain.Selector].Billing.Enabled {
				cfg.ChainsToDeploy[testutil.TestChain.Selector].Billing.Config.LinkTokenAddress = testEnv.LinkTokenState.LinkToken.Address()
			}

			resp, _, err := commonChangesets.ApplyChangesetsV2(t, testEnv.Environment, []commonChangesets.ConfiguredChangeSet{
				commonChangesets.Configure(DeployDataStreamsChainContractsChangeset, cfg),
			})
			require.NoError(t, err)

			var timelockAddr common.Address
			if tt.hasExistingMcms {
				timelockAddr = testEnv.Timelocks[testutil.TestChain.Selector].Timelock.Address()
			} else {
				addresses, err := resp.ExistingAddresses.AddressesForChain(testutil.TestChain.Selector)
				require.NoError(t, err)
				mcmsState, err := commonstate.MaybeLoadMCMSWithTimelockChainState(chain, addresses)
				require.NoError(t, err)
				timelockAddr = mcmsState.Timelock.Address()
			}

			for _, contract := range tt.expectedContracts {
				contractAddr, err := dsutil.MaybeFindEthAddress(resp.ExistingAddresses, testutil.TestChain.Selector, contract)
				require.NoError(t, err, "failed to find %s address in address book", contract)
				require.NotNil(t, "contractAddr", "address for %s was not found", contract)

				owner, _, err := commonChangesets.LoadOwnableContract(contractAddr, chain.Client)

				require.NoError(t, err)

				if cfg.ChainsToDeploy[testutil.TestChain.Selector].Ownership.ShouldTransfer {
					require.Equal(t, timelockAddr, owner, "%s contract owner should be the MCMS timelock", contract)
				} else {
					require.Equal(t, chain.DeployerKey.From, owner, "%s contract owner should be the deployer", contract)
				}
			}
		})
	}
}
