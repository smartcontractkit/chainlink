package stateview_test

import (
	"crypto/ecdsa"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/commit_store"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/ccip_home"
	onrampv16 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	cldftesthelpers "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils/testhelpers"

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
	state, err := stateview.LoadChainState(t.Context(), tenv.Env.BlockChains.EVMChains()[tenv.HomeChainSel], map[string]cldf.TypeAndVersion{
		fq1: cldf.NewTypeAndVersion(shared.FeeQuoter, deployment.Version1_0_0),
		fq2: cldf.NewTypeAndVersion(shared.FeeQuoter, deployment.Version1_2_0),
	})
	require.NoError(t, err)

	require.Equal(t, fq2, state.FeeQuoter.Address().Hex(), "expected latest fee quoter to be selected")
	require.Equal(t, deployment.Version1_2_0, *state.FeeQuoterVersion, "expected latest fee quoter version to be selected")
}

func TestLoadChainState_LegacyV15EVM2EVMDatastoreKeys(t *testing.T) {
	t.Parallel()

	srcSel := chain_selectors.TEST_90000001.Selector
	dstSel := chain_selectors.TEST_90000002.Selector

	e, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{srcSel, dstSel}),
		environment.WithLogger(logger.Test(t)),
	)
	require.NoError(t, err)

	chain := e.BlockChains.EVMChains()[srcSel]

	_, tx, onRamp, err := evm_2_evm_onramp.DeployEVM2EVMOnRamp(
		chain.DeployerKey, chain.Client,
		evm_2_evm_onramp.EVM2EVMOnRampStaticConfig{
			LinkToken:          common.HexToAddress("0x1"),
			ChainSelector:      chain.Selector,
			DestChainSelector:  dstSel,
			DefaultTxGasLimit:  10,
			MaxNopFeesJuels:    big.NewInt(10),
			PrevOnRamp:         common.Address{},
			RmnProxy:           common.HexToAddress("0x2"),
			TokenAdminRegistry: common.HexToAddress("0x3"),
		},
		evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig{
			Router:                            common.HexToAddress("0x4"),
			MaxNumberOfTokensPerMsg:           0,
			DestGasOverhead:                   0,
			DestGasPerPayloadByte:             0,
			DestDataAvailabilityOverheadGas:   0,
			DestGasPerDataAvailabilityByte:    0,
			DestDataAvailabilityMultiplierBps: 0,
			PriceRegistry:                     common.HexToAddress("0x5"),
			MaxDataBytes:                      0,
			MaxPerMsgGasLimit:                 0,
			DefaultTokenFeeUSDCents:           0,
			DefaultTokenDestGasOverhead:       0,
			EnforceOutOfOrder:                 false,
		},
		evm_2_evm_onramp.RateLimiterConfig{
			IsEnabled: true,
			Capacity:  big.NewInt(100),
			Rate:      big.NewInt(10),
		},
		[]evm_2_evm_onramp.EVM2EVMOnRampFeeTokenConfigArgs{},
		[]evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfigArgs{},
		[]evm_2_evm_onramp.EVM2EVMOnRampNopAndWeight{},
	)
	_, err = cldf.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	_, tx, cs, err := commit_store.DeployCommitStore(
		chain.DeployerKey, chain.Client, commit_store.CommitStoreStaticConfig{
			ChainSelector:       dstSel,
			SourceChainSelector: srcSel,
			OnRamp:              common.HexToAddress("0x4"),
			RmnProxy:            common.HexToAddress("0x1"),
		})
	_, err = cldf.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	offRampStatic := evm_2_evm_offramp.EVM2EVMOffRampStaticConfig{
		ChainSelector:       dstSel,
		SourceChainSelector: srcSel,
		RmnProxy:            common.HexToAddress("0x1"),
		CommitStore:         cs.Address(),
		TokenAdminRegistry:  common.HexToAddress("0x3"),
		OnRamp:              common.HexToAddress("0x4"),
	}
	rl := evm_2_evm_offramp.RateLimiterConfig{
		IsEnabled: true,
		Capacity:  big.NewInt(100),
		Rate:      big.NewInt(10),
	}
	_, tx, offRamp, err := evm_2_evm_offramp.DeployEVM2EVMOffRamp(
		chain.DeployerKey, chain.Client, offRampStatic, rl)
	_, err = cldf.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	evm2evmAddrs := map[string]cldf.TypeAndVersion{
		onRamp.Address().Hex():  cldf.NewTypeAndVersion(shared.EVM2EVMOnRamp, deployment.Version1_5_0),
		offRamp.Address().Hex(): cldf.NewTypeAndVersion(shared.EVM2EVMOffRamp, deployment.Version1_5_0),
	}

	legacyDisabled, err := stateview.LoadChainState(t.Context(), chain, evm2evmAddrs)
	require.NoError(t, err)
	require.Nil(t, legacyDisabled.EVM2EVMOnRamp)
	require.Nil(t, legacyDisabled.EVM2EVMOffRamp)

	lcs, err := stateview.LoadChainState(t.Context(), chain, evm2evmAddrs, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err)
	require.NotNil(t, lcs.EVM2EVMOnRamp[dstSel])
	require.Equal(t, onRamp.Address(), lcs.EVM2EVMOnRamp[dstSel].Address())
	require.NotNil(t, lcs.EVM2EVMOffRamp[srcSel])
	require.Equal(t, offRamp.Address(), lcs.EVM2EVMOffRamp[srcSel].Address())
	require.Equal(t, evm_2_evm_onramp.EVM2EVMOnRampABI, lcs.ABIByAddress[onRamp.Address().Hex()])
	require.Equal(t, evm_2_evm_offramp.EVM2EVMOffRampABI, lcs.ABIByAddress[offRamp.Address().Hex()])

	legacyNamesAddrs := map[string]cldf.TypeAndVersion{
		onRamp.Address().Hex():  cldf.NewTypeAndVersion(shared.OnRamp, deployment.Version1_5_0),
		offRamp.Address().Hex(): cldf.NewTypeAndVersion(shared.OffRamp, deployment.Version1_5_0),
	}
	byLegacyKeys, err := stateview.LoadChainState(t.Context(), chain, legacyNamesAddrs, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err)
	require.Equal(t, onRamp.Address(), byLegacyKeys.EVM2EVMOnRamp[dstSel].Address())
	require.Equal(t, offRamp.Address(), byLegacyKeys.EVM2EVMOffRamp[srcSel].Address())
}

// TestLoadChainState_LegacyEVM2EVMAndV16OnRampCoexist verifies that a chain can list both v1.5 EVM2EVM ramps
// (loaded only with WithLoadLegacyContracts) and a v1.6 OnRamp in the address map without either clobbering the other.
func TestLoadChainState_LegacyEVM2EVMAndV16OnRampCoexist(t *testing.T) {
	t.Parallel()

	srcSel := chain_selectors.TEST_90000001.Selector
	dstSel := chain_selectors.TEST_90000002.Selector

	e, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{srcSel, dstSel}),
		environment.WithLogger(logger.Test(t)),
	)
	require.NoError(t, err)

	chain := e.BlockChains.EVMChains()[srcSel]

	_, tx, evmOnRamp, err := evm_2_evm_onramp.DeployEVM2EVMOnRamp(
		chain.DeployerKey, chain.Client,
		evm_2_evm_onramp.EVM2EVMOnRampStaticConfig{
			LinkToken:          common.HexToAddress("0x1"),
			ChainSelector:      chain.Selector,
			DestChainSelector:  dstSel,
			DefaultTxGasLimit:  10,
			MaxNopFeesJuels:    big.NewInt(10),
			PrevOnRamp:         common.Address{},
			RmnProxy:           common.HexToAddress("0x2"),
			TokenAdminRegistry: common.HexToAddress("0x3"),
		},
		evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig{
			Router:                            common.HexToAddress("0x4"),
			MaxNumberOfTokensPerMsg:           0,
			DestGasOverhead:                   0,
			DestGasPerPayloadByte:             0,
			DestDataAvailabilityOverheadGas:   0,
			DestGasPerDataAvailabilityByte:    0,
			DestDataAvailabilityMultiplierBps: 0,
			PriceRegistry:                     common.HexToAddress("0x5"),
			MaxDataBytes:                      0,
			MaxPerMsgGasLimit:                 0,
			DefaultTokenFeeUSDCents:           0,
			DefaultTokenDestGasOverhead:       0,
			EnforceOutOfOrder:                 false,
		},
		evm_2_evm_onramp.RateLimiterConfig{
			IsEnabled: true,
			Capacity:  big.NewInt(100),
			Rate:      big.NewInt(10),
		},
		[]evm_2_evm_onramp.EVM2EVMOnRampFeeTokenConfigArgs{},
		[]evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfigArgs{},
		[]evm_2_evm_onramp.EVM2EVMOnRampNopAndWeight{},
	)
	_, err = cldf.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	_, tx, cs, err := commit_store.DeployCommitStore(
		chain.DeployerKey, chain.Client, commit_store.CommitStoreStaticConfig{
			ChainSelector:       dstSel,
			SourceChainSelector: srcSel,
			OnRamp:              common.HexToAddress("0x4"),
			RmnProxy:            common.HexToAddress("0x1"),
		})
	_, err = cldf.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	offRampStatic := evm_2_evm_offramp.EVM2EVMOffRampStaticConfig{
		ChainSelector:       dstSel,
		SourceChainSelector: srcSel,
		RmnProxy:            common.HexToAddress("0x1"),
		CommitStore:         cs.Address(),
		TokenAdminRegistry:  common.HexToAddress("0x3"),
		OnRamp:              common.HexToAddress("0x4"),
	}
	rl := evm_2_evm_offramp.RateLimiterConfig{
		IsEnabled: true,
		Capacity:  big.NewInt(100),
		Rate:      big.NewInt(10),
	}
	_, tx, evmOffRamp, err := evm_2_evm_offramp.DeployEVM2EVMOffRamp(
		chain.DeployerKey, chain.Client, offRampStatic, rl)
	_, err = cldf.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	v16Static := onrampv16.OnRampStaticConfig{
		ChainSelector:      chain.Selector,
		RmnRemote:          common.HexToAddress("0xb0"),
		NonceManager:       common.HexToAddress("0xb1"),
		TokenAdminRegistry: common.HexToAddress("0xb2"),
	}
	v16Dynamic := onrampv16.OnRampDynamicConfig{
		FeeQuoter:              common.HexToAddress("0xb3"),
		ReentrancyGuardEntered: false,
		MessageInterceptor:     common.Address{},
		FeeAggregator:          common.HexToAddress("0xb4"),
		AllowlistAdmin:         common.HexToAddress("0xb5"),
	}
	v16DestArgs := []onrampv16.OnRampDestChainConfigArgs{
		{
			DestChainSelector: dstSel,
			Router:            common.HexToAddress("0xb6"),
			AllowlistEnabled:  false,
		},
	}
	v16Addr, tx, _, err := onrampv16.DeployOnRamp(chain.DeployerKey, chain.Client, v16Static, v16Dynamic, v16DestArgs)
	_, err = cldf.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	combined := map[string]cldf.TypeAndVersion{
		evmOnRamp.Address().Hex():  cldf.NewTypeAndVersion(shared.EVM2EVMOnRamp, deployment.Version1_5_0),
		evmOffRamp.Address().Hex(): cldf.NewTypeAndVersion(shared.EVM2EVMOffRamp, deployment.Version1_5_0),
		v16Addr.Hex():              cldf.NewTypeAndVersion(shared.OnRamp, deployment.Version1_6_0),
	}

	withoutLegacy, err := stateview.LoadChainState(t.Context(), chain, combined)
	require.NoError(t, err)
	require.Nil(t, withoutLegacy.EVM2EVMOnRamp)
	require.Nil(t, withoutLegacy.EVM2EVMOffRamp)
	require.NotNil(t, withoutLegacy.OnRamp)
	require.Equal(t, v16Addr, withoutLegacy.OnRamp.Address())

	withLegacy, err := stateview.LoadChainState(t.Context(), chain, combined, stateview.WithLoadLegacyContracts(true))
	require.NoError(t, err)
	require.NotNil(t, withLegacy.OnRamp)
	require.Equal(t, v16Addr, withLegacy.OnRamp.Address())
	require.NotNil(t, withLegacy.EVM2EVMOnRamp[dstSel])
	require.Equal(t, evmOnRamp.Address(), withLegacy.EVM2EVMOnRamp[dstSel].Address())
	require.NotNil(t, withLegacy.EVM2EVMOffRamp[srcSel])
	require.Equal(t, evmOffRamp.Address(), withLegacy.EVM2EVMOffRamp[srcSel].Address())
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
				err = rt.Exec(runtime.ChangesetTask(cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2), map[uint64]cldfproposalutils.MCMSWithTimelockConfig{
					homeChainSelector: cldftesthelpers.SingleGroupTimelockConfig(t),
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
						runtime.SignAndExecuteProposalsTask([]*ecdsa.PrivateKey{cldftesthelpers.TestXXXMCMSSigner}),
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
