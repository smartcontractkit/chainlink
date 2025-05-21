package sequence

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"

	"github.com/smartcontractkit/chainlink/deployment"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commonstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/verification"
	dsutil "github.com/smartcontractkit/chainlink/deployment/data-streams/utils"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
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
		expectedContracts       []cldf.TypeAndVersion
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
			expectedContracts: []cldf.TypeAndVersion{
				cldf.NewTypeAndVersion(types.VerifierProxy, deployment.Version0_5_0),
				cldf.NewTypeAndVersion(types.Verifier, deployment.Version0_5_0),
				cldf.NewTypeAndVersion(types.RewardManager, deployment.Version0_5_0),
				cldf.NewTypeAndVersion(types.FeeManager, deployment.Version0_5_0),
				cldf.NewTypeAndVersion(commontypes.ProposerManyChainMultisig, deployment.Version1_0_0),
				cldf.NewTypeAndVersion(commontypes.BypasserManyChainMultisig, deployment.Version1_0_0),
				cldf.NewTypeAndVersion(commontypes.CancellerManyChainMultisig, deployment.Version1_0_0),
			},
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
			expectedContracts: []cldf.TypeAndVersion{
				cldf.NewTypeAndVersion(types.VerifierProxy, deployment.Version0_5_0),
				cldf.NewTypeAndVersion(types.Verifier, deployment.Version0_5_0),
				cldf.NewTypeAndVersion(commontypes.ProposerManyChainMultisig, deployment.Version1_0_0),
				cldf.NewTypeAndVersion(commontypes.BypasserManyChainMultisig, deployment.Version1_0_0),
				cldf.NewTypeAndVersion(commontypes.CancellerManyChainMultisig, deployment.Version1_0_0),
			},
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
			expectedContracts: []cldf.TypeAndVersion{
				cldf.NewTypeAndVersion(types.VerifierProxy, deployment.Version0_5_0),
				cldf.NewTypeAndVersion(types.Verifier, deployment.Version0_5_0),
			},
		},
		{
			name:            "Deploy but do not propose transfer",
			hasExistingMcms: true,
			deployDataStreamsConfig: DeployDataStreamsConfig{
				ChainsToDeploy: map[uint64]DeployDataStreams{testutil.TestChain.Selector: {
					VerifierConfig: verificationCfg,
				}},
			},
			expectedContracts: []cldf.TypeAndVersion{
				cldf.NewTypeAndVersion(types.VerifierProxy, deployment.Version0_5_0),
				cldf.NewTypeAndVersion(types.Verifier, deployment.Version0_5_0),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.deployDataStreamsConfig
			chainSel := testutil.TestChain.Selector
			billingEnabled := cfg.ChainsToDeploy[chainSel].Billing.Enabled

			testEnv := testutil.NewMemoryEnvV2(t, testutil.MemoryEnvConfig{
				ShouldDeployMCMS:      tt.hasExistingMcms,
				ShouldDeployLinkToken: billingEnabled,
			})

			chain := testEnv.Environment.Chains[chainSel]

			if cfg.ChainsToDeploy[chainSel].Billing.Enabled {
				cfg.ChainsToDeploy[chainSel].Billing.Config.LinkTokenAddress = testEnv.LinkTokenState.LinkToken.Address()
			}

			env, _, err := commonChangesets.ApplyChangesetsV2(t, testEnv.Environment, []commonChangesets.ConfiguredChangeSet{
				commonChangesets.Configure(DeployDataStreamsChainContractsChangeset, cfg),
			})
			require.NoError(t, err)

			var timelockAddr common.Address
			if tt.hasExistingMcms {
				// Deployed by test env
				timelockAddr = testEnv.Timelocks[chainSel].Timelock.Address()
			} else {
				// Deployed by changeset
				addresses, err := dsutil.EnvironmentAddresses(env)
				require.NoError(t, err)
				mcmsState, err := commonstate.MaybeLoadMCMSWithTimelockChainState(chain, addresses)
				require.NoError(t, err)
				timelockAddr = mcmsState.Timelock.Address()
			}

			for _, contract := range tt.expectedContracts {
				record, err := env.DataStore.Addresses().Get(
					datastore.NewAddressRefKey(testutil.TestChain.Selector, datastore.ContractType(contract.Type), &contract.Version, ""),
				)
				require.NoError(t, err)
				contractAddress := common.HexToAddress(record.Address)
				owner, _, err := commonChangesets.LoadOwnableContract(contractAddress, chain.Client)

				require.NoError(t, err)

				if cfg.ChainsToDeploy[chainSel].Ownership.ShouldTransfer {
					require.Equal(t, timelockAddr, owner, "%s contract owner should be the MCMS timelock", contract)
				} else {
					require.Equal(t, chain.DeployerKey.From, owner, "%s contract owner should be the deployer", contract)
				}
			}
		})
	}
}
