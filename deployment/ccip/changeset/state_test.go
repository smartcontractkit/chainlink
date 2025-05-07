package changeset_test

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/ccip_home"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestSmokeState(t *testing.T) {
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(3))
	state, err := changeset.LoadOnchainState(tenv.Env)
	require.NoError(t, err)
	_, _, err = state.View(&tenv.Env, tenv.Env.AllChainSelectors())
	require.NoError(t, err)
}

func TestMCMSState(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/DX-106")
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNoJobsAndContracts())
	addressbook := deployment.NewMemoryAddressBook()
	newTv := deployment.NewTypeAndVersion(types.ManyChainMultisig, deployment.Version1_0_0)
	newTv.AddLabel(types.BypasserRole.String())
	newTv.AddLabel(types.CancellerRole.String())
	newTv.AddLabel(types.ProposerRole.String())
	addr := utils.RandomAddress()
	require.NoError(t, addressbook.Save(tenv.HomeChainSel, addr.String(), newTv))
	require.NoError(t, tenv.Env.ExistingAddresses.Merge(addressbook))
	state, err := changeset.LoadOnchainState(tenv.Env)
	require.NoError(t, err)
	require.Equal(t, addr.String(), state.Chains[tenv.HomeChainSel].BypasserMcm.Address().String())
	require.Equal(t, addr.String(), state.Chains[tenv.HomeChainSel].ProposerMcm.Address().String())
	require.Equal(t, addr.String(), state.Chains[tenv.HomeChainSel].CancellerMcm.Address().String())
}

func TestEnforceMCMSUsageIfProd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Msg                    string
		DeployCCIPHome         bool
		DeployCapReg           bool
		DeployMCMS             bool
		TransferCCIPHomeToMCMS bool
		TransferCapRegToMCMS   bool
		ExpectedErr            string
		MCMSConfig             *proposalutils.TimelockConfig
	}{
		{
			Msg:                    "CCIPHome & CapReg ownership mismatch",
			DeployCCIPHome:         true,
			DeployCapReg:           true,
			DeployMCMS:             true,
			TransferCCIPHomeToMCMS: true,
			TransferCapRegToMCMS:   false,
			MCMSConfig:             &proposalutils.TimelockConfig{},
			ExpectedErr:            "CCIPHome and CapabilitiesRegistry owners do not match",
		},
		{
			Msg:                    "CCIPHome MCMS owned & MCMS config provided",
			DeployCCIPHome:         true,
			DeployCapReg:           true,
			DeployMCMS:             true,
			TransferCCIPHomeToMCMS: true,
			TransferCapRegToMCMS:   true,
			MCMSConfig:             &proposalutils.TimelockConfig{},
			ExpectedErr:            "",
		},
		{
			Msg:                    "CCIPHome MCMS owned & MCMS config not provided",
			DeployCCIPHome:         true,
			DeployCapReg:           true,
			DeployMCMS:             true,
			TransferCCIPHomeToMCMS: true,
			TransferCapRegToMCMS:   true,
			MCMSConfig:             nil,
			ExpectedErr:            "MCMS is enforced for environment",
		},
		{
			Msg:                    "CCIPHome not MCMS owned & MCMS config provided",
			DeployCCIPHome:         true,
			DeployCapReg:           true,
			DeployMCMS:             true,
			TransferCCIPHomeToMCMS: false,
			TransferCapRegToMCMS:   false,
			MCMSConfig:             &proposalutils.TimelockConfig{},
			ExpectedErr:            "",
		},
		{
			Msg:                    "CCIPHome not MCMS owned & MCMS config not provided",
			DeployCCIPHome:         true,
			DeployCapReg:           true,
			DeployMCMS:             true,
			TransferCCIPHomeToMCMS: false,
			TransferCapRegToMCMS:   false,
			MCMSConfig:             nil,
			ExpectedErr:            "",
		},
		{
			Msg:                    "CCIPHome not deployed & MCMS config provided",
			DeployCCIPHome:         false,
			DeployCapReg:           true,
			DeployMCMS:             false,
			TransferCCIPHomeToMCMS: false,
			TransferCapRegToMCMS:   false,
			MCMSConfig:             &proposalutils.TimelockConfig{},
			ExpectedErr:            "",
		},
		{
			Msg:                    "CCIPHome not deployed & MCMS config not provided",
			DeployCCIPHome:         false,
			DeployCapReg:           true,
			DeployMCMS:             false,
			TransferCCIPHomeToMCMS: false,
			TransferCapRegToMCMS:   false,
			MCMSConfig:             nil,
			ExpectedErr:            "",
		},
		{
			Msg:                    "MCMS not deployed & MCMS config provided",
			DeployCCIPHome:         true,
			DeployCapReg:           true,
			DeployMCMS:             false,
			TransferCCIPHomeToMCMS: false,
			TransferCapRegToMCMS:   false,
			MCMSConfig:             &proposalutils.TimelockConfig{},
			ExpectedErr:            "",
		},
		{
			Msg:                    "MCMS not deployed & MCMS config not provided",
			DeployCCIPHome:         true,
			DeployCapReg:           true,
			DeployMCMS:             false,
			TransferCCIPHomeToMCMS: false,
			TransferCapRegToMCMS:   false,
			MCMSConfig:             nil,
			ExpectedErr:            "",
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			var err error
			lggr := logger.TestLogger(t)
			e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
				Chains: 1,
			})
			homeChainSelector := e.AllChainSelectors()[0]

			if test.DeployCCIPHome {
				_, err = cldf.DeployContract(e.Logger, e.Chains[homeChainSelector], e.ExistingAddresses,
					func(chain deployment.Chain) cldf.ContractDeploy[*ccip_home.CCIPHome] {
						address, tx2, contract, err2 := ccip_home.DeployCCIPHome(
							chain.DeployerKey,
							chain.Client,
							utils.RandomAddress(), // We don't need a real contract address here, just a random one to satisfy the constructor.
						)
						return cldf.ContractDeploy[*ccip_home.CCIPHome]{
							Address: address, Contract: contract, Tx: tx2, Tv: deployment.NewTypeAndVersion(changeset.CCIPHome, deployment.Version1_6_0), Err: err2,
						}
					})
				require.NoError(t, err, "failed to deploy CCIP home")
			}

			if test.DeployCapReg {
				_, err = cldf.DeployContract(e.Logger, e.Chains[homeChainSelector], e.ExistingAddresses,
					func(chain deployment.Chain) cldf.ContractDeploy[*capabilities_registry.CapabilitiesRegistry] {
						address, tx2, contract, err2 := capabilities_registry.DeployCapabilitiesRegistry(
							chain.DeployerKey,
							chain.Client,
						)
						return cldf.ContractDeploy[*capabilities_registry.CapabilitiesRegistry]{
							Address: address, Contract: contract, Tx: tx2, Tv: deployment.NewTypeAndVersion(changeset.CapabilitiesRegistry, deployment.Version1_0_0), Err: err2,
						}
					})
				require.NoError(t, err, "failed to deploy capability registry")
			}

			if test.DeployMCMS {
				e, err = commonchangeset.Apply(t, e, nil,
					commonchangeset.Configure(cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2), map[uint64]types.MCMSWithTimelockConfigV2{
						homeChainSelector: proposalutils.SingleGroupTimelockConfigV2(t),
					}),
				)
				require.NoError(t, err, "failed to deploy MCMS")
				state, err := changeset.LoadOnchainState(e)
				require.NoError(t, err, "failed to load onchain state")

				addrs := make([]common.Address, 0, 2)
				if test.TransferCCIPHomeToMCMS {
					addrs = append(addrs, state.Chains[homeChainSelector].CCIPHome.Address())
				}
				if test.TransferCapRegToMCMS {
					addrs = append(addrs, state.Chains[homeChainSelector].CapabilityRegistry.Address())
				}
				if len(addrs) > 0 {
					e, err = commonchangeset.Apply(t, e,
						map[uint64]*proposalutils.TimelockExecutionContracts{
							homeChainSelector: &proposalutils.TimelockExecutionContracts{
								Timelock:  state.Chains[homeChainSelector].Timelock,
								CallProxy: state.Chains[homeChainSelector].CallProxy,
							},
						},
						commonchangeset.Configure(
							cldf.CreateLegacyChangeSet(commonchangeset.TransferToMCMSWithTimelockV2),
							commonchangeset.TransferToMCMSWithTimelockConfig{
								ContractsByChain: map[uint64][]common.Address{
									homeChainSelector: addrs,
								},
								MCMSConfig: proposalutils.TimelockConfig{
									MinDelay: 0 * time.Second,
								},
							},
						),
					)
					require.NoError(t, err, "failed to transfer contracts to MCMS")
				}
			}

			state, err := changeset.LoadOnchainState(e)
			require.NoError(t, err, "failed to load onchain state")

			err = state.EnforceMCMSUsageIfProd(e.GetContext(), test.MCMSConfig)
			if test.ExpectedErr != "" {
				require.Error(t, err, "expected error but got nil")
				require.Contains(t, err.Error(), test.ExpectedErr, "error message mismatch")
				return
			}
			require.NoError(t, err, "failed to validate MCMS config")
		})
	}
}

// TODO: add solana state test
