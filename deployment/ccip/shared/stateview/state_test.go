package stateview_test

import (
	"context"
	"crypto/ecdsa"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/ccip_home"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

func TestLoadChainState_MultipleFeeQuoters(t *testing.T) {
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(3))
	fq1 := utils.RandomAddress().Hex()
	fq2 := utils.RandomAddress().Hex()
	state, err := stateview.LoadChainState(context.Background(), tenv.Env.BlockChains.EVMChains()[tenv.HomeChainSel], map[string]cldf.TypeAndVersion{
		fq1: cldf.NewTypeAndVersion(shared.FeeQuoter, deployment.Version1_0_0),
		fq2: cldf.NewTypeAndVersion(shared.FeeQuoter, deployment.Version1_2_0),
	})
	require.NoError(t, err)

	require.Equal(t, fq2, state.FeeQuoter.Address().Hex(), "expected latest fee quoter to be selected")
	require.Equal(t, deployment.Version1_2_0, *state.FeeQuoterVersion, "expected latest fee quoter version to be selected")
}

func TestSmokeState(t *testing.T) {
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(3))
	state, err := stateview.LoadOnchainState(tenv.Env)
	require.NoError(t, err)
	_, err = state.View(&tenv.Env, tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM)))
	require.NoError(t, err)
}

func TestMCMSState(t *testing.T) {
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNoJobsAndContracts())
	addressbook := cldf.NewMemoryAddressBook()
	newTv := cldf.NewTypeAndVersion(types.ManyChainMultisig, deployment.Version1_0_0)
	newTv.AddLabel(types.BypasserRole.String())
	newTv.AddLabel(types.CancellerRole.String())
	newTv.AddLabel(types.ProposerRole.String())
	addr := utils.RandomAddress()
	require.NoError(t, addressbook.Save(tenv.HomeChainSel, addr.String(), newTv))
	require.NoError(t, tenv.Env.ExistingAddresses.Merge(addressbook))
	state, err := stateview.LoadOnchainState(tenv.Env)
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

			homeChainSelector := chain_selectors.TEST_90000001.Selector
			lggr := logger.Test(t)
			rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
				environment.WithEVMSimulated(t, []uint64{homeChainSelector}),
				environment.WithLogger(lggr),
			))
			require.NoError(t, err)

			evmChains := rt.Environment().BlockChains.EVMChains()

			if test.DeployCCIPHome {
				_, err = cldf.DeployContract(lggr, evmChains[homeChainSelector], rt.State().AddressBook,
					func(chain cldf_evm.Chain) cldf.ContractDeploy[*ccip_home.CCIPHome] {
						address, tx2, contract, err2 := ccip_home.DeployCCIPHome(
							chain.DeployerKey,
							chain.Client,
							utils.RandomAddress(), // We don't need a real contract address here, just a random one to satisfy the constructor.
						)
						return cldf.ContractDeploy[*ccip_home.CCIPHome]{
							Address: address, Contract: contract, Tx: tx2, Tv: cldf.NewTypeAndVersion(shared.CCIPHome, deployment.Version1_6_0), Err: err2,
						}
					})
				require.NoError(t, err, "failed to deploy CCIP home")
			}

			if test.DeployCapReg {
				_, err = cldf.DeployContract(lggr, evmChains[homeChainSelector], rt.State().AddressBook,
					func(chain cldf_evm.Chain) cldf.ContractDeploy[*capabilities_registry.CapabilitiesRegistry] {
						address, tx2, contract, err2 := capabilities_registry.DeployCapabilitiesRegistry(
							chain.DeployerKey,
							chain.Client,
						)
						return cldf.ContractDeploy[*capabilities_registry.CapabilitiesRegistry]{
							Address: address, Contract: contract, Tx: tx2, Tv: cldf.NewTypeAndVersion(shared.CapabilitiesRegistry, deployment.Version1_0_0), Err: err2,
						}
					})
				require.NoError(t, err, "failed to deploy capability registry")
			}

			if test.DeployMCMS {
				err = rt.Exec(runtime.ChangesetTask(cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2), map[uint64]types.MCMSWithTimelockConfigV2{
					homeChainSelector: proposalutils.SingleGroupTimelockConfigV2(t),
				}))
				require.NoError(t, err, "failed to deploy MCMS")

				state, err := stateview.LoadOnchainState(rt.Environment())
				require.NoError(t, err, "failed to load onchain state")

				addrs := make([]common.Address, 0, 2)
				if test.TransferCCIPHomeToMCMS {
					addrs = append(addrs, state.Chains[homeChainSelector].CCIPHome.Address())
				}
				if test.TransferCapRegToMCMS {
					addrs = append(addrs, state.Chains[homeChainSelector].CapabilityRegistry.Address())
				}
				if len(addrs) > 0 {
					err = rt.Exec(
						runtime.ChangesetTask(cldf.CreateLegacyChangeSet(commonchangeset.TransferToMCMSWithTimelockV2), commonchangeset.TransferToMCMSWithTimelockConfig{
							ContractsByChain: map[uint64][]common.Address{
								homeChainSelector: addrs,
							},
							MCMSConfig: proposalutils.TimelockConfig{
								MinDelay: 0 * time.Second,
							},
						}),
						runtime.SignAndExecuteProposalsTask([]*ecdsa.PrivateKey{proposalutils.TestXXXMCMSSigner}),
					)

					require.NoError(t, err, "failed to transfer contracts to MCMS")
				}
			}

			state, err := stateview.LoadOnchainState(rt.Environment())
			require.NoError(t, err, "failed to load onchain state")

			err = state.EnforceMCMSUsageIfProd(t.Context(), test.MCMSConfig)
			if test.ExpectedErr != "" {
				require.Error(t, err, "expected error but got nil")
				require.ErrorContains(t, err, test.ExpectedErr, "error message mismatch")
				return
			}
			require.NoError(t, err, "failed to validate MCMS config")
		})
	}
}

// TODO: add solana state test
