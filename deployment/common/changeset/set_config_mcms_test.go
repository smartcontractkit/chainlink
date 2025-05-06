package changeset_test

import (
	"crypto/ecdsa"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	solanasdk "github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/config"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/mcms/sdk/evm"
	"github.com/smartcontractkit/mcms/sdk/solana"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commonchangesetsolana "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// setupSetConfigTestEnv deploys all required contracts for the setConfig MCMS contract call.
func setupSetConfigTestEnv(t *testing.T) deployment.Environment {
	lggr := logger.TestLogger(t)
	cfg := memory.MemoryEnvironmentConfig{
		Chains:    2,
		SolChains: 1,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
	chainSelector := env.AllChainSelectors()[0]
	chainSelectorSolana := env.AllChainSelectorsSolana()[0]

	commonchangeset.SetPreloadedSolanaAddresses(t, env, chainSelectorSolana)
	config := proposalutils.SingleGroupTimelockConfigV2(t)
	// Deploy MCMS and Timelock
	env, err := commonchangeset.Apply(t, env, nil,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.DeployLinkToken),
			[]uint64{chainSelector, chainSelectorSolana},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
			map[uint64]commontypes.MCMSWithTimelockConfigV2{
				chainSelector:       config,
				chainSelectorSolana: config,
			},
		),
	)
	require.NoError(t, err)
	return env
}

// TestSetConfigMCMSVariants tests the SetConfigMCMS changeset variants.
func TestSetConfigMCMSVariants(t *testing.T) {
	t.Parallel()
	// Add the timelock as a signer to check state changes
	for _, tc := range []struct {
		name       string
		changeSets func(mcmsState *commonchangeset.MCMSWithTimelockState, chainSel uint64, cfgProp, cfgCancel, cfgBypass config.Config) []commonchangeset.ConfiguredChangeSet
	}{
		{
			name: "MCMS disabled",
			changeSets: func(mcmsState *commonchangeset.MCMSWithTimelockState, chainSel uint64, cfgProp, cfgCancel, cfgBypass config.Config) []commonchangeset.ConfiguredChangeSet {
				return []commonchangeset.ConfiguredChangeSet{
					commonchangeset.Configure(
						cldf.CreateLegacyChangeSet(commonchangeset.SetConfigMCMS),
						commonchangeset.MCMSConfig{
							ConfigsPerChain: map[uint64]commonchangeset.ConfigPerRole{
								chainSel: {
									Proposer:  cfgProp,
									Canceller: cfgCancel,
									Bypasser:  cfgBypass,
								},
							},
						},
					),
				}
			},
		},
		{
			name: "MCMS enabled",
			changeSets: func(mcmsState *commonchangeset.MCMSWithTimelockState, chainSel uint64, cfgProp, cfgCancel, cfgBypass config.Config) []commonchangeset.ConfiguredChangeSet {
				return []commonchangeset.ConfiguredChangeSet{
					commonchangeset.Configure(
						cldf.CreateLegacyChangeSet(commonchangeset.TransferToMCMSWithTimelock),
						commonchangeset.TransferToMCMSWithTimelockConfig{
							ContractsByChain: map[uint64][]common.Address{
								chainSel: {mcmsState.ProposerMcm.Address(), mcmsState.BypasserMcm.Address(), mcmsState.CancellerMcm.Address()},
							},
						},
					),
					commonchangeset.Configure(
						cldf.CreateLegacyChangeSet(commonchangeset.SetConfigMCMS),
						commonchangeset.MCMSConfig{
							ProposalConfig: &proposalutils.TimelockConfig{
								MinDelay: 0,
							},
							ConfigsPerChain: map[uint64]commonchangeset.ConfigPerRole{
								chainSel: {
									Proposer:  cfgProp,
									Canceller: cfgCancel,
									Bypasser:  cfgBypass,
								},
							},
						},
					),
				}
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := tests.Context(t)

			env := setupSetConfigTestEnv(t)
			chainSelector := env.AllChainSelectors()[0]
			chain := env.Chains[chainSelector]
			addrs, err := env.ExistingAddresses.AddressesForChain(chainSelector)
			require.NoError(t, err)
			require.Len(t, addrs, 6)

			mcmsState, err := commonchangeset.MaybeLoadMCMSWithTimelockChainState(chain, addrs)
			require.NoError(t, err)
			timelockAddress := mcmsState.Timelock.Address()
			cfgProposer := proposalutils.SingleGroupMCMS(t)
			cfgProposer.Signers = append(cfgProposer.Signers, timelockAddress)
			cfgProposer.Quorum = 2 // quorum should change to 2 out of 2 signers
			timelockMap := map[uint64]*proposalutils.TimelockExecutionContracts{
				chainSelector: {
					Timelock:  mcmsState.Timelock,
					CallProxy: mcmsState.CallProxy,
				},
			}
			cfgCanceller := proposalutils.SingleGroupMCMS(t)
			cfgBypasser := proposalutils.SingleGroupMCMS(t)
			cfgBypasser.Signers = append(cfgBypasser.Signers, timelockAddress)
			cfgBypasser.Signers = append(cfgBypasser.Signers, mcmsState.ProposerMcm.Address())
			cfgBypasser.Quorum = 3 // quorum should change to 3 out of 3 signers

			// Set config on all 3 MCMS contracts
			changesetsToApply := tc.changeSets(mcmsState, chainSelector, cfgProposer, cfgCanceller, cfgBypasser)
			_, err = commonchangeset.ApplyChangesets(t, env, timelockMap, changesetsToApply)
			require.NoError(t, err)

			// Check new State
			expected := cfgProposer.ToRawConfig()
			opts := &bind.CallOpts{Context: ctx}
			newConf, err := mcmsState.ProposerMcm.GetConfig(opts)
			require.NoError(t, err)
			require.Equal(t, expected, newConf)

			expected = cfgBypasser.ToRawConfig()
			newConf, err = mcmsState.BypasserMcm.GetConfig(opts)
			require.NoError(t, err)
			require.Equal(t, expected, newConf)

			expected = cfgCanceller.ToRawConfig()
			newConf, err = mcmsState.CancellerMcm.GetConfig(opts)
			require.NoError(t, err)
			require.Equal(t, expected, newConf)
		})
	}
}

func TestSetConfigMCMSV2EVM(t *testing.T) {
	t.Parallel()
	// Add the timelock as a signer to check state changes
	for _, tc := range []struct {
		name       string
		changeSets func(mcmsState *commonchangeset.MCMSWithTimelockState, chainSel uint64, cfgProp, cfgCancel, cfgBypass mcmstypes.Config) []commonchangeset.ConfiguredChangeSet
	}{
		{
			name: "MCMS disabled",
			changeSets: func(mcmsState *commonchangeset.MCMSWithTimelockState, chainSel uint64, cfgProp, cfgCancel, cfgBypass mcmstypes.Config) []commonchangeset.ConfiguredChangeSet {
				return []commonchangeset.ConfiguredChangeSet{
					commonchangeset.Configure(
						cldf.CreateLegacyChangeSet(commonchangeset.SetConfigMCMSV2),
						commonchangeset.MCMSConfigV2{
							ConfigsPerChain: map[uint64]commonchangeset.ConfigPerRoleV2{
								chainSel: {
									Proposer:  cfgProp,
									Canceller: cfgCancel,
									Bypasser:  cfgBypass,
								},
							},
						},
					),
				}
			},
		},
		{
			name: "MCMS enabled",
			changeSets: func(mcmsState *commonchangeset.MCMSWithTimelockState, chainSel uint64, cfgProp, cfgCancel, cfgBypass mcmstypes.Config) []commonchangeset.ConfiguredChangeSet {
				return []commonchangeset.ConfiguredChangeSet{
					commonchangeset.Configure(
						cldf.CreateLegacyChangeSet(commonchangeset.TransferToMCMSWithTimelockV2),
						commonchangeset.TransferToMCMSWithTimelockConfig{
							ContractsByChain: map[uint64][]common.Address{
								chainSel: {mcmsState.ProposerMcm.Address(), mcmsState.BypasserMcm.Address(), mcmsState.CancellerMcm.Address()},
							},
						},
					),
					commonchangeset.Configure(
						cldf.CreateLegacyChangeSet(commonchangeset.SetConfigMCMSV2),
						commonchangeset.MCMSConfigV2{
							ProposalConfig: &proposalutils.TimelockConfig{
								MinDelay: 0,
							},
							ConfigsPerChain: map[uint64]commonchangeset.ConfigPerRoleV2{
								chainSel: {
									Proposer:  cfgProp,
									Canceller: cfgCancel,
									Bypasser:  cfgBypass,
								},
							},
						},
					),
				}
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := tests.Context(t)

			env := setupSetConfigTestEnv(t)
			chainSelector := env.AllChainSelectors()[0]
			chain := env.Chains[chainSelector]
			addrs, err := env.ExistingAddresses.AddressesForChain(chainSelector)
			require.NoError(t, err)
			require.Len(t, addrs, 6)

			mcmsState, err := commonchangeset.MaybeLoadMCMSWithTimelockChainState(chain, addrs)
			require.NoError(t, err)
			timelockAddress := mcmsState.Timelock.Address()
			cfgProposer := proposalutils.SingleGroupMCMSV2(t)
			cfgProposer.Signers = append(cfgProposer.Signers, timelockAddress)
			cfgProposer.Quorum = 2 // quorum should change to 2 out of 2 signers
			timelockMap := map[uint64]*proposalutils.TimelockExecutionContracts{
				chainSelector: {
					Timelock:  mcmsState.Timelock,
					CallProxy: mcmsState.CallProxy,
				},
			}
			cfgCanceller := proposalutils.SingleGroupMCMSV2(t)
			cfgBypasser := proposalutils.SingleGroupMCMSV2(t)
			cfgBypasser.Signers = append(cfgBypasser.Signers, timelockAddress)
			cfgBypasser.Signers = append(cfgBypasser.Signers, mcmsState.ProposerMcm.Address())
			cfgBypasser.Quorum = 3 // quorum should change to 3 out of 3 signers

			// Set config on all 3 MCMS contracts
			changesetsToApply := tc.changeSets(mcmsState, chainSelector, cfgProposer, cfgCanceller, cfgBypasser)
			_, err = commonchangeset.ApplyChangesets(t, env, timelockMap, changesetsToApply)
			require.NoError(t, err)

			// Check new State
			inspector := evm.NewInspector(chain.Client)
			newConf, err := inspector.GetConfig(ctx, mcmsState.ProposerMcm.Address().Hex())
			require.NoError(t, err)
			require.ElementsMatch(t, cfgProposer.Signers, newConf.Signers)
			require.Equal(t, cfgProposer.Quorum, newConf.Quorum)

			newConf, err = inspector.GetConfig(ctx, mcmsState.BypasserMcm.Address().Hex())
			require.NoError(t, err)
			require.ElementsMatch(t, cfgBypasser.Signers, newConf.Signers)
			require.Equal(t, cfgBypasser.Quorum, newConf.Quorum)

			newConf, err = inspector.GetConfig(ctx, mcmsState.CancellerMcm.Address().Hex())
			require.NoError(t, err)
			require.ElementsMatch(t, cfgCanceller.Signers, newConf.Signers)
			require.Equal(t, cfgCanceller.Quorum, newConf.Quorum)
		})
	}
}

func TestSetConfigMCMSV2Solana(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name       string
		changeSets func(chainSel uint64, cfgs map[uint64]commonchangeset.ConfigPerRoleV2) []commonchangeset.ConfiguredChangeSet
	}{
		{
			name: "MCMS disabled",
			changeSets: func(chainSel uint64, cfgs map[uint64]commonchangeset.ConfigPerRoleV2) []commonchangeset.ConfiguredChangeSet {
				return []commonchangeset.ConfiguredChangeSet{
					commonchangeset.Configure(
						cldf.CreateLegacyChangeSet(commonchangeset.SetConfigMCMSV2),
						commonchangeset.MCMSConfigV2{
							ConfigsPerChain: cfgs,
						},
					),
				}
			},
		},
		{
			name: "MCMS enabled",
			changeSets: func(chainSel uint64, cfgs map[uint64]commonchangeset.ConfigPerRoleV2) []commonchangeset.ConfiguredChangeSet {
				return []commonchangeset.ConfiguredChangeSet{
					commonchangeset.Configure(&commonchangesetsolana.TransferMCMSToTimelockSolana{}, commonchangesetsolana.TransferMCMSToTimelockSolanaConfig{
						Chains:  []uint64{chainSel},
						MCMSCfg: proposalutils.TimelockConfig{MinDelay: time.Second * 1},
					}),
					commonchangeset.Configure(
						cldf.CreateLegacyChangeSet(commonchangeset.SetConfigMCMSV2),
						commonchangeset.MCMSConfigV2{
							ProposalConfig: &proposalutils.TimelockConfig{
								MinDelay: time.Second * 1,
							},
							ConfigsPerChain: cfgs,
						},
					),
				}
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := tests.Context(t)

			env := setupSetConfigTestEnv(t)
			chainSelectorSolana := env.AllChainSelectorsSolana()[0]
			solChain := env.SolChains[chainSelectorSolana]

			addrs, err := env.ExistingAddresses.AddressesForChain(chainSelectorSolana)
			require.NoError(t, err)
			solMcmsState, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(solChain, addrs)
			require.NoError(t, err)

			fundSignerPDAs(t, env, chainSelectorSolana, solMcmsState)

			// new signer
			key, err := crypto.GenerateKey()
			require.NoError(t, err)
			publicKey := key.Public().(*ecdsa.PublicKey)
			newSigner := crypto.PubkeyToAddress(*publicKey)

			newCfgProposer := proposalutils.SingleGroupMCMSV2(t)
			newCfgProposer.Signers = append(newCfgProposer.Signers, newSigner)
			newCfgProposer.Quorum = 2
			newCfgCanceller := proposalutils.SingleGroupMCMSV2(t)
			newCfgBypasser := proposalutils.SingleGroupMCMSV2(t)
			newCfgBypasser.Signers = append(newCfgBypasser.Signers, newSigner)
			newCfgBypasser.Quorum = 2

			// act
			changesetsToApply := tc.changeSets(
				chainSelectorSolana,
				map[uint64]commonchangeset.ConfigPerRoleV2{
					chainSelectorSolana: {
						Proposer:  newCfgProposer,
						Canceller: newCfgCanceller,
						Bypasser:  newCfgBypasser,
					},
				})
			_, _, err = commonchangeset.ApplyChangesetsV2(t, env, changesetsToApply)
			require.NoError(t, err)

			// assert
			inspectorSolana := solana.NewInspector(solChain.Client)
			confs, err := inspectorSolana.GetConfig(ctx, solana.ContractAddress(solMcmsState.McmProgram, solana.PDASeed(solMcmsState.ProposerMcmSeed)))
			require.NoError(t, err)
			require.ElementsMatch(t, newCfgProposer.Signers, confs.Signers)
			require.Equal(t, newCfgProposer.Quorum, confs.Quorum)

			confs, err = inspectorSolana.GetConfig(ctx, solana.ContractAddress(solMcmsState.McmProgram, solana.PDASeed(solMcmsState.BypasserMcmSeed)))
			require.NoError(t, err)
			require.ElementsMatch(t, newCfgBypasser.Signers, confs.Signers)
			require.Equal(t, newCfgBypasser.Quorum, confs.Quorum)

			confs, err = inspectorSolana.GetConfig(ctx, solana.ContractAddress(solMcmsState.McmProgram, solana.PDASeed(solMcmsState.CancellerMcmSeed)))
			require.NoError(t, err)
			require.ElementsMatch(t, newCfgCanceller.Signers, confs.Signers)
			require.Equal(t, newCfgCanceller.Quorum, confs.Quorum)
		})
	}
}

func TestValidate(t *testing.T) {
	t.Parallel()
	env := setupSetConfigTestEnv(t)

	chainSelector := env.AllChainSelectors()[0]
	chain := env.Chains[chainSelector]
	addrs, err := env.ExistingAddresses.AddressesForChain(chainSelector)
	require.NoError(t, err)
	require.Len(t, addrs, 6)
	mcmsState, err := commonchangeset.MaybeLoadMCMSWithTimelockChainState(chain, addrs)
	require.NoError(t, err)
	cfg := proposalutils.SingleGroupMCMS(t)
	timelockAddress := mcmsState.Timelock.Address()
	// Add the timelock as a signer to check state changes
	cfg.Signers = append(cfg.Signers, timelockAddress)
	cfg.Quorum = 2 // quorum

	cfgInvalid := proposalutils.SingleGroupMCMS(t)
	cfgInvalid.Quorum = 0
	require.NoError(t, err)
	tests := []struct {
		name     string
		cfg      commonchangeset.MCMSConfig
		errorMsg string
	}{
		{
			name: "valid config",
			cfg: commonchangeset.MCMSConfig{
				ProposalConfig: &proposalutils.TimelockConfig{
					MinDelay: 0,
				},
				ConfigsPerChain: map[uint64]commonchangeset.ConfigPerRole{
					chainSelector: {
						Proposer:  cfg,
						Canceller: cfg,
						Bypasser:  cfg,
					},
				},
			},
		},
		{
			name: "valid non mcms config",
			cfg: commonchangeset.MCMSConfig{
				ConfigsPerChain: map[uint64]commonchangeset.ConfigPerRole{
					chainSelector: {
						Proposer:  cfg,
						Canceller: cfg,
						Bypasser:  cfg,
					},
				},
			},
		},
		{
			name: "no chain configurations",
			cfg: commonchangeset.MCMSConfig{
				ConfigsPerChain: map[uint64]commonchangeset.ConfigPerRole{},
			},
			errorMsg: "no chain configs provided",
		},
		{
			name: "non evm chain",
			cfg: commonchangeset.MCMSConfig{
				ConfigsPerChain: map[uint64]commonchangeset.ConfigPerRole{
					chain_selectors.APTOS_MAINNET.Selector: {
						Proposer:  cfg,
						Canceller: cfg,
						Bypasser:  cfg,
					},
				},
			},
			errorMsg: "chain selector: 4741433654826277614 is not an ethereum chain",
		},
		{
			name: "chain selector not found in environment",
			cfg: commonchangeset.MCMSConfig{
				ConfigsPerChain: map[uint64]commonchangeset.ConfigPerRole{
					123: {
						Proposer:  cfg,
						Canceller: cfg,
						Bypasser:  cfg,
					},
				},
			},
			errorMsg: "unknown chain selector 123",
		},
		{
			name: "invalid proposer config",
			cfg: commonchangeset.MCMSConfig{
				ProposalConfig: &proposalutils.TimelockConfig{
					MinDelay: 0,
				},
				ConfigsPerChain: map[uint64]commonchangeset.ConfigPerRole{
					chainSelector: {
						Proposer:  cfgInvalid,
						Canceller: cfg,
						Bypasser:  cfg,
					},
				},
			},
			errorMsg: "invalid MCMS config: Quorum must be greater than 0",
		},
		{
			name: "invalid canceller config",
			cfg: commonchangeset.MCMSConfig{
				ProposalConfig: &proposalutils.TimelockConfig{
					MinDelay: 0,
				},
				ConfigsPerChain: map[uint64]commonchangeset.ConfigPerRole{
					chainSelector: {
						Proposer:  cfg,
						Canceller: cfgInvalid,
						Bypasser:  cfg,
					},
				},
			},
			errorMsg: "invalid MCMS config: Quorum must be greater than 0",
		},
		{
			name: "invalid bypasser config",
			cfg: commonchangeset.MCMSConfig{
				ProposalConfig: &proposalutils.TimelockConfig{
					MinDelay: 0,
				},
				ConfigsPerChain: map[uint64]commonchangeset.ConfigPerRole{
					chainSelector: {
						Proposer:  cfg,
						Canceller: cfg,
						Bypasser:  cfgInvalid,
					},
				},
			},
			errorMsg: "invalid MCMS config: Quorum must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selectors := []uint64{chainSelector}

			err := tt.cfg.Validate(env, selectors)
			if tt.errorMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateV2(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/DX-439")

	t.Parallel()
	env := setupSetConfigTestEnv(t)

	chainSelector := env.AllChainSelectors()[0]
	chainSelectorSolana := env.AllChainSelectorsSolana()[0]

	cfg := proposalutils.SingleGroupMCMSV2(t)
	cfgInvalid := proposalutils.SingleGroupMCMSV2(t)
	cfgInvalid.Quorum = 0

	tests := []struct {
		name     string
		cfg      commonchangeset.MCMSConfigV2
		errorMsg string
	}{
		{
			name: "valid config",
			cfg: commonchangeset.MCMSConfigV2{
				ProposalConfig: &proposalutils.TimelockConfig{
					MinDelay: 0,
				},
				ConfigsPerChain: map[uint64]commonchangeset.ConfigPerRoleV2{
					chainSelector: {
						Proposer:  cfg,
						Canceller: cfg,
						Bypasser:  cfg,
					},
					chainSelectorSolana: {
						Proposer:  cfg,
						Canceller: cfg,
						Bypasser:  cfg,
					},
				},
			},
		},
		{
			name: "valid non mcms config",
			cfg: commonchangeset.MCMSConfigV2{
				ConfigsPerChain: map[uint64]commonchangeset.ConfigPerRoleV2{
					chainSelector: {
						Proposer:  cfg,
						Canceller: cfg,
						Bypasser:  cfg,
					},
					chainSelectorSolana: {
						Proposer:  cfg,
						Canceller: cfg,
						Bypasser:  cfg,
					},
				},
			},
		},
		{
			name: "no chain configurations",
			cfg: commonchangeset.MCMSConfigV2{
				ConfigsPerChain: map[uint64]commonchangeset.ConfigPerRoleV2{},
			},
			errorMsg: "no chain configs provided",
		},
		{
			name: "chain selector not found in environment",
			cfg: commonchangeset.MCMSConfigV2{
				ConfigsPerChain: map[uint64]commonchangeset.ConfigPerRoleV2{
					123: {
						Proposer:  cfg,
						Canceller: cfg,
						Bypasser:  cfg,
					},
				},
			},
			errorMsg: "unknown chain selector 123",
		},
		{
			name: "invalid proposer config",
			cfg: commonchangeset.MCMSConfigV2{
				ProposalConfig: &proposalutils.TimelockConfig{
					MinDelay: 0,
				},
				ConfigsPerChain: map[uint64]commonchangeset.ConfigPerRoleV2{
					chainSelector: {
						Proposer:  cfgInvalid,
						Canceller: cfg,
						Bypasser:  cfg,
					},
					chainSelectorSolana: {
						Proposer:  cfg,
						Canceller: cfg,
						Bypasser:  cfg,
					},
				},
			},
			errorMsg: "invalid MCMS config: Quorum must be greater than 0",
		},
		{
			name: "invalid canceller config",
			cfg: commonchangeset.MCMSConfigV2{
				ProposalConfig: &proposalutils.TimelockConfig{
					MinDelay: 0,
				},
				ConfigsPerChain: map[uint64]commonchangeset.ConfigPerRoleV2{
					chainSelector: {
						Proposer:  cfg,
						Canceller: cfgInvalid,
						Bypasser:  cfg,
					},
					chainSelectorSolana: {
						Proposer:  cfg,
						Canceller: cfg,
						Bypasser:  cfg,
					},
				},
			},
			errorMsg: "invalid MCMS config: Quorum must be greater than 0",
		},
		{
			name: "invalid bypasser config",
			cfg: commonchangeset.MCMSConfigV2{
				ProposalConfig: &proposalutils.TimelockConfig{
					MinDelay: 0,
				},
				ConfigsPerChain: map[uint64]commonchangeset.ConfigPerRoleV2{
					chainSelector: {
						Proposer:  cfg,
						Canceller: cfg,
						Bypasser:  cfgInvalid,
					},
					chainSelectorSolana: {
						Proposer:  cfg,
						Canceller: cfg,
						Bypasser:  cfg,
					},
				},
			},
			errorMsg: "invalid MCMS config: Quorum must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selectors := []uint64{chainSelector, chainSelectorSolana}

			err := tt.cfg.Validate(env, selectors)
			if tt.errorMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func fundSignerPDAs(
	t *testing.T, env deployment.Environment, chainSelector uint64, chainState *state.MCMSWithTimelockStateSolana,
) {
	t.Helper()
	solChain := env.SolChains[chainSelector]
	timelockSignerPDA := state.GetTimelockSignerPDA(chainState.TimelockProgram, chainState.TimelockSeed)
	mcmSignerPDA := state.GetMCMSignerPDA(chainState.McmProgram, chainState.ProposerMcmSeed)
	signerPDAs := []solanasdk.PublicKey{timelockSignerPDA, mcmSignerPDA}
	memory.FundSolanaAccounts(env.GetContext(), t, signerPDAs, 1, solChain.Client)
}
