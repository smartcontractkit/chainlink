//nolint:testifylint // inverting want and got is more succinct
package changeset_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/google/go-cmp/cmp"
	mcmsevmsdk "github.com/smartcontractkit/mcms/sdk/evm"
	mcmssolanasdk "github.com/smartcontractkit/mcms/sdk/solana"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	solanainternal "github.com/smartcontractkit/chainlink/deployment/common/changeset/internal/solana"
	mcmschangesetstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestDeployMCMSWithTimelockV2(t *testing.T) {
	// --- arrange ---
	log := logger.TestLogger(t)
	envConfig := memory.MemoryEnvironmentConfig{Chains: 2, SolChains: 1}
	env := memory.NewMemoryEnvironment(t, log, zapcore.InfoLevel, envConfig)
	evmSelectors := env.AllChainSelectors()
	solanaSelectors := env.AllChainSelectorsSolana()
	changesetConfig := map[uint64]commontypes.MCMSWithTimelockConfigV2{
		evmSelectors[0]: {
			Proposer: mcmstypes.Config{
				Quorum:  1,
				Signers: []common.Address{common.HexToAddress("0x0000000000000000000000000000000000000001")},
				GroupSigners: []mcmstypes.Config{
					{
						Quorum:       1,
						Signers:      []common.Address{common.HexToAddress("0x0000000000000000000000000000000000000002")},
						GroupSigners: []mcmstypes.Config{},
					},
				},
			},
			Canceller: mcmstypes.Config{
				Quorum:       1,
				Signers:      []common.Address{common.HexToAddress("0x0000000000000000000000000000000000000003")},
				GroupSigners: []mcmstypes.Config{},
			},
			Bypasser: mcmstypes.Config{
				Quorum:       1,
				Signers:      []common.Address{common.HexToAddress("0x0000000000000000000000000000000000000004")},
				GroupSigners: []mcmstypes.Config{},
			},
			TimelockMinDelay: big.NewInt(0),
		},
		evmSelectors[1]: {
			Proposer: mcmstypes.Config{
				Quorum:       1,
				Signers:      []common.Address{common.HexToAddress("0x0000000000000000000000000000000000000011")},
				GroupSigners: []mcmstypes.Config{},
			},
			Canceller: mcmstypes.Config{
				Quorum: 2,
				Signers: []common.Address{
					common.HexToAddress("0x0000000000000000000000000000000000000012"),
					common.HexToAddress("0x0000000000000000000000000000000000000013"),
					common.HexToAddress("0x0000000000000000000000000000000000000014"),
				},
				GroupSigners: []mcmstypes.Config{},
			},
			Bypasser: mcmstypes.Config{
				Quorum:       1,
				Signers:      []common.Address{common.HexToAddress("0x0000000000000000000000000000000000000005")},
				GroupSigners: []mcmstypes.Config{},
			},
			TimelockMinDelay: big.NewInt(1),
		},
		solanaSelectors[0]: {
			Proposer: mcmstypes.Config{
				Quorum: 1,
				Signers: []common.Address{
					common.HexToAddress("0x0000000000000000000000000000000000000021"),
					common.HexToAddress("0x0000000000000000000000000000000000000022"),
				},
				GroupSigners: []mcmstypes.Config{
					{
						Quorum: 2,
						Signers: []common.Address{
							common.HexToAddress("0x0000000000000000000000000000000000000023"),
							common.HexToAddress("0x0000000000000000000000000000000000000024"),
							common.HexToAddress("0x0000000000000000000000000000000000000025"),
						},
						GroupSigners: []mcmstypes.Config{
							{
								Quorum: 1,
								Signers: []common.Address{
									common.HexToAddress("0x0000000000000000000000000000000000000026"),
								},
								GroupSigners: []mcmstypes.Config{},
							},
						},
					},
				},
			},
			Canceller: mcmstypes.Config{
				Quorum: 1,
				Signers: []common.Address{
					common.HexToAddress("0x0000000000000000000000000000000000000027"),
				},
				GroupSigners: []mcmstypes.Config{},
			},
			Bypasser: mcmstypes.Config{
				Quorum: 1,
				Signers: []common.Address{
					common.HexToAddress("0x0000000000000000000000000000000000000028"),
					common.HexToAddress("0x0000000000000000000000000000000000000029"),
				},
				GroupSigners: []mcmstypes.Config{},
			},
			TimelockMinDelay: big.NewInt(2),
		},
	}
	changesetApplication := []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
			Config:    changesetConfig,
		},
	}
	setPreloadedSolanaAddresses(t, env, solanaSelectors[0])

	// --- act ---
	updatedEnv, err := commonchangeset.ApplyChangesets(t, env, nil, changesetApplication)
	require.NoError(t, err)

	evmState, err := mcmschangesetstate.MaybeLoadMCMSWithTimelockState(updatedEnv, evmSelectors)
	require.NoError(t, err)
	solanaState, err := mcmschangesetstate.MaybeLoadMCMSWithTimelockStateSolana(updatedEnv, solanaSelectors)
	require.NoError(t, err)

	// --- assert ---
	require.Len(t, evmState, 2)
	require.Len(t, solanaState, 1)

	// evm chain 0
	evmState0 := evmState[evmSelectors[0]]
	evmInspector := mcmsevmsdk.NewInspector(updatedEnv.Chains[evmSelectors[0]].Client)
	evmTimelockInspector := mcmsevmsdk.NewTimelockInspector(updatedEnv.Chains[evmSelectors[0]].Client)

	config, err := evmInspector.GetConfig(updatedEnv.GetContext(), evmState0.ProposerMcm.Address().Hex())
	require.NoError(t, err)
	require.Empty(t, cmp.Diff(*config, changesetConfig[evmSelectors[0]].Proposer))

	config, err = evmInspector.GetConfig(updatedEnv.GetContext(), evmState0.CancellerMcm.Address().Hex())
	require.NoError(t, err)
	require.Empty(t, cmp.Diff(*config, changesetConfig[evmSelectors[0]].Canceller))

	config, err = evmInspector.GetConfig(updatedEnv.GetContext(), evmState0.BypasserMcm.Address().Hex())
	require.NoError(t, err)
	require.Empty(t, cmp.Diff(*config, changesetConfig[evmSelectors[0]].Bypasser))

	proposers, err := evmTimelockInspector.GetProposers(updatedEnv.GetContext(), evmState0.Timelock.Address().Hex())
	require.NoError(t, err)
	require.Equal(t, proposers, []string{evmState0.ProposerMcm.Address().Hex()})

	executors, err := evmTimelockInspector.GetExecutors(updatedEnv.GetContext(), evmState0.Timelock.Address().Hex())
	require.NoError(t, err)
	require.Equal(t, executors, []string{evmState0.CallProxy.Address().Hex()})

	cancellers, err := evmTimelockInspector.GetCancellers(updatedEnv.GetContext(), evmState0.Timelock.Address().Hex())
	require.NoError(t, err)
	require.ElementsMatch(t, cancellers, []string{
		evmState0.CancellerMcm.Address().Hex(),
		evmState0.ProposerMcm.Address().Hex(),
		evmState0.BypasserMcm.Address().Hex(),
	})

	bypassers, err := evmTimelockInspector.GetBypassers(updatedEnv.GetContext(), evmState0.Timelock.Address().Hex())
	require.NoError(t, err)
	require.Equal(t, bypassers, []string{evmState0.BypasserMcm.Address().Hex()})

	// evm chain 1
	evmState1 := evmState[evmSelectors[1]]
	evmInspector = mcmsevmsdk.NewInspector(updatedEnv.Chains[evmSelectors[1]].Client)
	evmTimelockInspector = mcmsevmsdk.NewTimelockInspector(updatedEnv.Chains[evmSelectors[1]].Client)

	config, err = evmInspector.GetConfig(updatedEnv.GetContext(), evmState1.ProposerMcm.Address().Hex())
	require.NoError(t, err)
	require.Empty(t, cmp.Diff(*config, changesetConfig[evmSelectors[1]].Proposer))

	config, err = evmInspector.GetConfig(updatedEnv.GetContext(), evmState1.CancellerMcm.Address().Hex())
	require.NoError(t, err)
	require.Empty(t, cmp.Diff(*config, changesetConfig[evmSelectors[1]].Canceller))

	config, err = evmInspector.GetConfig(updatedEnv.GetContext(), evmState1.BypasserMcm.Address().Hex())
	require.NoError(t, err)
	require.Empty(t, cmp.Diff(*config, changesetConfig[evmSelectors[1]].Bypasser))

	proposers, err = evmTimelockInspector.GetProposers(updatedEnv.GetContext(), evmState1.Timelock.Address().Hex())
	require.NoError(t, err)
	require.Equal(t, proposers, []string{evmState1.ProposerMcm.Address().Hex()})

	executors, err = evmTimelockInspector.GetExecutors(updatedEnv.GetContext(), evmState1.Timelock.Address().Hex())
	require.NoError(t, err)
	require.Equal(t, executors, []string{evmState1.CallProxy.Address().Hex()})

	cancellers, err = evmTimelockInspector.GetCancellers(updatedEnv.GetContext(), evmState1.Timelock.Address().Hex())
	require.NoError(t, err)
	require.ElementsMatch(t, cancellers, []string{
		evmState1.CancellerMcm.Address().Hex(),
		evmState1.ProposerMcm.Address().Hex(),
		evmState1.BypasserMcm.Address().Hex(),
	})

	bypassers, err = evmTimelockInspector.GetBypassers(updatedEnv.GetContext(), evmState1.Timelock.Address().Hex())
	require.NoError(t, err)
	require.Equal(t, bypassers, []string{evmState1.BypasserMcm.Address().Hex()})

	// solana chain 0
	solanaState0 := solanaState[solanaSelectors[0]]
	solanaInspector := mcmssolanasdk.NewInspector(updatedEnv.SolChains[solanaSelectors[0]].Client)
	solanaTimelockInspector := mcmssolanasdk.NewTimelockInspector(updatedEnv.SolChains[solanaSelectors[0]].Client)

	addr := mcmssolanasdk.ContractAddress(solanaState0.McmProgram, mcmssolanasdk.PDASeed(solanaState0.ProposerMcmSeed))
	config, err = solanaInspector.GetConfig(updatedEnv.GetContext(), addr)
	require.NoError(t, err)
	require.Empty(t, cmp.Diff(*config, changesetConfig[solanaSelectors[0]].Proposer))

	addr = mcmssolanasdk.ContractAddress(solanaState0.McmProgram, mcmssolanasdk.PDASeed(solanaState0.CancellerMcmSeed))
	config, err = solanaInspector.GetConfig(updatedEnv.GetContext(), addr)
	require.NoError(t, err)
	require.Empty(t, cmp.Diff(*config, changesetConfig[solanaSelectors[0]].Canceller))

	addr = mcmssolanasdk.ContractAddress(solanaState0.McmProgram, mcmssolanasdk.PDASeed(solanaState0.BypasserMcmSeed))
	config, err = solanaInspector.GetConfig(updatedEnv.GetContext(), addr)
	require.NoError(t, err)
	require.Empty(t, cmp.Diff(*config, changesetConfig[solanaSelectors[0]].Bypasser))

	addr = mcmssolanasdk.ContractAddress(solanaState0.TimelockProgram, mcmssolanasdk.PDASeed(solanaState0.TimelockSeed))
	proposers, err = solanaTimelockInspector.GetProposers(updatedEnv.GetContext(), addr)
	require.NoError(t, err)
	require.Equal(t, proposers, []string{signerPDA(solanaState0.McmProgram, solanaState0.ProposerMcmSeed)})

	executors, err = solanaTimelockInspector.GetExecutors(updatedEnv.GetContext(), addr)
	require.NoError(t, err)
	require.Equal(t, executors, []string{}) // FIXME: who should be executor?

	cancellers, err = solanaTimelockInspector.GetCancellers(updatedEnv.GetContext(), addr)
	require.NoError(t, err)
	require.ElementsMatch(t, cancellers, []string{
		signerPDA(solanaState0.McmProgram, solanaState0.CancellerMcmSeed),
		signerPDA(solanaState0.McmProgram, solanaState0.ProposerMcmSeed),
		signerPDA(solanaState0.McmProgram, solanaState0.BypasserMcmSeed),
	})

	bypassers, err = solanaTimelockInspector.GetBypassers(updatedEnv.GetContext(), addr)
	require.NoError(t, err)
	require.Equal(t, bypassers, []string{signerPDA(solanaState0.McmProgram, solanaState0.BypasserMcmSeed)})
}

// ----- helpers -----

func signerPDA(programID solana.PublicKey, seed mcmschangesetstate.PDASeed) string {
	return solanainternal.GetMCMSignerPDA(programID, seed).String()
}

func setPreloadedSolanaAddresses(t *testing.T, env deployment.Environment, selector uint64) {
	typeAndVersion := deployment.NewTypeAndVersion(commontypes.ManyChainMultisigProgram, deployment.Version1_0_0)
	err := env.ExistingAddresses.Save(selector, memory.SolanaProgramIDs["mcm"], typeAndVersion)
	require.NoError(t, err)

	typeAndVersion = deployment.NewTypeAndVersion(commontypes.AccessControllerProgram, deployment.Version1_0_0)
	err = env.ExistingAddresses.Save(selector, memory.SolanaProgramIDs["access_controller"], typeAndVersion)
	require.NoError(t, err)

	typeAndVersion = deployment.NewTypeAndVersion(commontypes.RBACTimelockProgram, deployment.Version1_0_0)
	err = env.ExistingAddresses.Save(selector, memory.SolanaProgramIDs["timelock"], typeAndVersion)
	require.NoError(t, err)
}
