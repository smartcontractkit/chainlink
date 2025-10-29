package changeset_test

import (
	"crypto/ecdsa"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	solanago "github.com/gagliardetto/solana-go"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/mcms/sdk/evm"
	"github.com/smartcontractkit/mcms/sdk/solana"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commonchangesetsolana "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/internal/soltestutils"
)

func TestSetConfigMCMSV2EVM(t *testing.T) {
	t.Parallel()

	selector1 := chain_selectors.TEST_90000001.Selector
	selector2 := chain_selectors.TEST_90000002.Selector

	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector1, selector2}),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	chain1 := rt.Environment().BlockChains.EVMChains()[selector1]
	chain2 := rt.Environment().BlockChains.EVMChains()[selector2]

	config := proposalutils.SingleGroupTimelockConfigV2(t)

	// Deploy MCMS and Timelock for selector1 & selector2
	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2), map[uint64]commontypes.MCMSWithTimelockConfigV2{
			selector1: config,
			selector2: config,
		}),
	)
	require.NoError(t, err)

	// Transfer MCMS contracts to timelock for selector2 for testing setConfig on MCMS owned contracts
	chain2Addrs, err := rt.State().AddressBook.AddressesForChain(selector2)
	require.NoError(t, err)
	require.Len(t, chain2Addrs, 5)
	chain2MCMSState, err := commonchangeset.MaybeLoadMCMSWithTimelockChainState(chain2, chain2Addrs)
	require.NoError(t, err)

	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(commonchangeset.TransferToMCMSWithTimelockV2), commonchangeset.TransferToMCMSWithTimelockConfig{
			ContractsByChain: map[uint64][]common.Address{
				selector2: {
					chain2MCMSState.ProposerMcm.Address(),
					chain2MCMSState.BypasserMcm.Address(),
					chain2MCMSState.CancellerMcm.Address(),
				},
			},
		}),
		runtime.SignAndExecuteProposalsTask([]*ecdsa.PrivateKey{proposalutils.TestXXXMCMSSigner}),
	)
	require.NoError(t, err)

	// Add the timelock as a signer to check state changes
	for _, tt := range []struct {
		name       string
		chain      cldf_evm.Chain
		changeSets func(selector uint64, cfgProp, cfgCancel, cfgBypass mcmstypes.Config) []runtime.Executable
	}{
		{
			name:  "MCMS disabled",
			chain: chain1,
			changeSets: func(selector uint64, cfgProp, cfgCancel, cfgBypass mcmstypes.Config) []runtime.Executable {
				return []runtime.Executable{
					runtime.ChangesetTask(cldf.CreateLegacyChangeSet(commonchangeset.SetConfigMCMSV2), commonchangeset.MCMSConfigV2{
						ConfigsPerChain: map[uint64]commonchangeset.ConfigPerRoleV2{
							selector: {
								Proposer:  cfgProp,
								Canceller: cfgCancel,
								Bypasser:  cfgBypass,
							},
						},
					}),
				}
			},
		},
		{
			name:  "MCMS enabled",
			chain: chain2,
			changeSets: func(selector uint64, cfgProp, cfgCancel, cfgBypass mcmstypes.Config) []runtime.Executable {
				return []runtime.Executable{
					runtime.ChangesetTask(cldf.CreateLegacyChangeSet(commonchangeset.SetConfigMCMSV2), commonchangeset.MCMSConfigV2{
						ProposalConfig: &proposalutils.TimelockConfig{
							MinDelay: 0,
						},
						ConfigsPerChain: map[uint64]commonchangeset.ConfigPerRoleV2{
							selector: {
								Proposer:  cfgProp,
								Canceller: cfgCancel,
								Bypasser:  cfgBypass,
							},
						},
					}),
					runtime.SignAndExecuteProposalsTask([]*ecdsa.PrivateKey{proposalutils.TestXXXMCMSSigner}),
				}
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// Get the mcms addresses for the chain
			addrs, err := rt.State().AddressBook.AddressesForChain(tt.chain.Selector)
			require.NoError(t, err)

			// Check new State
			mcmsState, err := commonchangeset.MaybeLoadMCMSWithTimelockChainState(tt.chain, addrs)
			require.NoError(t, err)
			timelockAddress := mcmsState.Timelock.Address()

			// Create new configs for the MCMS contracts
			cfgProposer := proposalutils.SingleGroupMCMSV2(t)
			cfgProposer.Signers = append(cfgProposer.Signers, timelockAddress)
			cfgProposer.Quorum = 2                             // quorum should change to 2 out of 2 signers
			cfgCanceller := proposalutils.SingleGroupMCMSV2(t) // quorum should not change
			cfgBypasser := proposalutils.SingleGroupMCMSV2(t)
			cfgBypasser.Signers = append(cfgBypasser.Signers, timelockAddress)
			cfgBypasser.Signers = append(cfgBypasser.Signers, mcmsState.ProposerMcm.Address())
			cfgBypasser.Quorum = 3 // quorum should change to 3 out of 3 signers

			// Set config on all 3 MCMS contracts
			err = rt.Exec(tt.changeSets(tt.chain.Selector, cfgProposer, cfgCanceller, cfgBypasser)...)
			require.NoError(t, err)

			inspector := evm.NewInspector(tt.chain.Client)
			newConf, err := inspector.GetConfig(t.Context(), mcmsState.ProposerMcm.Address().Hex())
			require.NoError(t, err)
			require.ElementsMatch(t, cfgProposer.Signers, newConf.Signers)
			require.Equal(t, cfgProposer.Quorum, newConf.Quorum)

			newConf, err = inspector.GetConfig(t.Context(), mcmsState.BypasserMcm.Address().Hex())
			require.NoError(t, err)
			require.ElementsMatch(t, cfgBypasser.Signers, newConf.Signers)
			require.Equal(t, cfgBypasser.Quorum, newConf.Quorum)

			newConf, err = inspector.GetConfig(t.Context(), mcmsState.CancellerMcm.Address().Hex())
			require.NoError(t, err)
			require.ElementsMatch(t, cfgCanceller.Signers, newConf.Signers)
			require.Equal(t, cfgCanceller.Quorum, newConf.Quorum)
		})
	}
}

func TestSetConfigMCMSV2Solana(t *testing.T) {
	t.Parallel()

	selector := chain_selectors.TEST_22222222222222222222222222222222222222222222.Selector

	programsPath, programIDs, ab := soltestutils.PreloadMCMS(t, selector)

	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithSolanaContainer(t, []uint64{selector}, programsPath, programIDs),
		environment.WithAddressBook(ab),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	chain := rt.Environment().BlockChains.SolanaChains()[selector]

	// Deploy MCMS and Timelock
	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2), map[uint64]commontypes.MCMSWithTimelockConfigV2{
			selector: proposalutils.SingleGroupTimelockConfigV2(t),
		}),
	)
	require.NoError(t, err)

	// Load the MCMS state
	addrs, err := rt.State().AddressBook.AddressesForChain(selector)
	require.NoError(t, err)
	mcmsState, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addrs)
	require.NoError(t, err)

	// Fund the signer PDAs for the MCMS contracts
	soltestutils.FundSignerPDAs(t, chain, mcmsState)

	// Used to check the chain state after the changeset is applied
	inspector := solana.NewInspector(chain.Client)

	// Create some signers to set into the config
	signer1Key, signer1Addr := createSolSigner(t)
	_, signer2Addr := createSolSigner(t)

	newCfgProposer := proposalutils.SingleGroupMCMSV2(t)
	newCfgProposer.Signers = append(newCfgProposer.Signers, signer1Addr)
	newCfgProposer.Quorum = 2
	newCfgCanceller := proposalutils.SingleGroupMCMSV2(t)
	newCfgBypasser := proposalutils.SingleGroupMCMSV2(t)
	newCfgBypasser.Signers = append(newCfgBypasser.Signers, signer1Addr)
	newCfgBypasser.Quorum = 2

	t.Run("MCMS disabled", func(t *testing.T) {
		err = rt.Exec(
			runtime.ChangesetTask(cldf.CreateLegacyChangeSet(commonchangeset.SetConfigMCMSV2), commonchangeset.MCMSConfigV2{
				ConfigsPerChain: map[uint64]commonchangeset.ConfigPerRoleV2{
					selector: {
						Proposer:  newCfgProposer,
						Canceller: newCfgCanceller,
						Bypasser:  newCfgBypasser,
					},
				},
			}),
		)
		require.NoError(t, err)

		assertSolConfigEquals(t, inspector, mcmsState.McmProgram, mcmsState.ProposerMcmSeed, newCfgProposer)
		assertSolConfigEquals(t, inspector, mcmsState.McmProgram, mcmsState.BypasserMcmSeed, newCfgBypasser)
		assertSolConfigEquals(t, inspector, mcmsState.McmProgram, mcmsState.CancellerMcmSeed, newCfgCanceller)
	})

	t.Run("MCMS enabled", func(t *testing.T) {
		// Now we transfer the MCMS contracts to the timelock for testing setConfig on MCMS owned contracts
		err = rt.Exec(
			runtime.ChangesetTask(commonchangesetsolana.TransferMCMSToTimelockSolana{}, commonchangesetsolana.TransferMCMSToTimelockSolanaConfig{
				Chains:  []uint64{selector},
				MCMSCfg: proposalutils.TimelockConfig{MinDelay: time.Second * 1},
			}),
			// We must sign with an additional signer since we changed the config quorum previously.
			runtime.SignAndExecuteProposalsTask([]*ecdsa.PrivateKey{proposalutils.TestXXXMCMSSigner, signer1Key}),
		)
		require.NoError(t, err)

		// Update the configs with yet another additional signer
		newCfgProposer.Signers = append(newCfgProposer.Signers, signer2Addr)
		newCfgProposer.Quorum = 3
		newCfgBypasser.Signers = append(newCfgBypasser.Signers, signer2Addr)
		newCfgBypasser.Quorum = 3

		err = rt.Exec(
			runtime.ChangesetTask(cldf.CreateLegacyChangeSet(commonchangeset.SetConfigMCMSV2), commonchangeset.MCMSConfigV2{
				ProposalConfig: &proposalutils.TimelockConfig{
					MinDelay: time.Second * 1,
				},
				ConfigsPerChain: map[uint64]commonchangeset.ConfigPerRoleV2{
					selector: {
						Proposer:  newCfgProposer,
						Canceller: newCfgCanceller,
						Bypasser:  newCfgBypasser,
					},
				},
			}),
			runtime.SignAndExecuteProposalsTask([]*ecdsa.PrivateKey{proposalutils.TestXXXMCMSSigner, signer1Key}),
		)
		require.NoError(t, err)

		assertSolConfigEquals(t, inspector, mcmsState.McmProgram, mcmsState.ProposerMcmSeed, newCfgProposer)
		assertSolConfigEquals(t, inspector, mcmsState.McmProgram, mcmsState.BypasserMcmSeed, newCfgBypasser)
		assertSolConfigEquals(t, inspector, mcmsState.McmProgram, mcmsState.CancellerMcmSeed, newCfgCanceller)
	})
}

func TestValidateV2(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/DX-439")

	t.Parallel()

	evmSelector := chain_selectors.TEST_90000001.Selector
	solSelector := chain_selectors.TEST_22222222222222222222222222222222222222222222.Selector

	programsPath, programIDs, ab := soltestutils.PreloadMCMS(t, solSelector)

	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{evmSelector}),
		environment.WithSolanaContainer(t, []uint64{solSelector}, programsPath, programIDs),
		environment.WithAddressBook(ab),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	config := proposalutils.SingleGroupTimelockConfigV2(t)

	// Deploy MCMS and Timelock
	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(commonchangeset.DeployLinkToken), []uint64{evmSelector}),
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2), map[uint64]commontypes.MCMSWithTimelockConfigV2{
			evmSelector: config,
			solSelector: config,
		}),
	)
	require.NoError(t, err)

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
					evmSelector: {
						Proposer:  cfg,
						Canceller: cfg,
						Bypasser:  cfg,
					},
					solSelector: {
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
					evmSelector: {
						Proposer:  cfg,
						Canceller: cfg,
						Bypasser:  cfg,
					},
					solSelector: {
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
					evmSelector: {
						Proposer:  cfgInvalid,
						Canceller: cfg,
						Bypasser:  cfg,
					},
					solSelector: {
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
					evmSelector: {
						Proposer:  cfg,
						Canceller: cfgInvalid,
						Bypasser:  cfg,
					},
					solSelector: {
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
					evmSelector: {
						Proposer:  cfg,
						Canceller: cfg,
						Bypasser:  cfgInvalid,
					},
					solSelector: {
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
			selectors := []uint64{evmSelector, solSelector}

			err := tt.cfg.Validate(rt.Environment(), selectors)
			if tt.errorMsg != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func assertSolConfigEquals(
	t *testing.T, inspector *solana.Inspector, programID solanago.PublicKey, seed state.PDASeed, want mcmstypes.Config,
) {
	t.Helper()

	cfg, err := inspector.GetConfig(t.Context(), solana.ContractAddress(programID, solana.PDASeed(seed)))
	require.NoError(t, err)
	require.ElementsMatch(t, want.Signers, cfg.Signers)
	require.Equal(t, want.Quorum, cfg.Quorum)
}

// createSolSigner creates a new Solana signer and returns the private key and address
func createSolSigner(t *testing.T) (*ecdsa.PrivateKey, common.Address) {
	t.Helper()

	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	publicKey := key.Public().(*ecdsa.PublicKey)

	return key, crypto.PubkeyToAddress(*publicKey)
}
