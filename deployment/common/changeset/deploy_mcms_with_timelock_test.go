//nolint:testifylint // inverting want and got is more succinct
package changeset_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	mcmsevmsdk "github.com/smartcontractkit/mcms/sdk/evm"
	mcmssolanasdk "github.com/smartcontractkit/mcms/sdk/solana"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	timelockBindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/timelock"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	mcmschangesetstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestGrantRoleInTimeLock(t *testing.T) {
	ctx := testutils.Context(t)
	env := memory.NewMemoryEnvironment(t, logger.TestLogger(t), zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains:             2,
		NumOfUsersPerChain: 2,
	})
	evmSelectors := env.AllChainSelectors()
	changesetConfig := make(map[uint64]commontypes.MCMSWithTimelockConfigV2)
	for _, chain := range evmSelectors {
		changesetConfig[chain] = proposalutils.SingleGroupTimelockConfigV2(t)
	}
	// deploy the MCMS with timelock contracts
	configuredChangeset := commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
		changesetConfig,
	)
	updatedEnv, err := commonchangeset.Apply(t, env, nil, configuredChangeset)
	require.NoError(t, err)
	mcmsState, err := mcmschangesetstate.MaybeLoadMCMSWithTimelockState(updatedEnv, evmSelectors)
	require.NoError(t, err)

	// change the environment to remove proposer from the timelock, so that we can deploy new proposer
	// and then grant the role to the new proposer
	existingProposer := mcmsState[evmSelectors[0]].ProposerMcm
	ab := deployment.NewMemoryAddressBook()
	require.NoError(t, ab.Save(evmSelectors[0], existingProposer.Address().String(),
		deployment.NewTypeAndVersion(commontypes.ProposerManyChainMultisig, deployment.Version1_0_0)))
	require.NoError(t, updatedEnv.ExistingAddresses.Remove(ab))

	// change the deployer key, so that we can deploy proposer with a new key
	// the new deployer key will not be admin of the timelock
	// we can test granting roles through proposal
	chain := updatedEnv.Chains[evmSelectors[0]]
	chain.DeployerKey = updatedEnv.Chains[evmSelectors[0]].Users[0]
	updatedEnv.Chains[evmSelectors[0]] = chain

	// now deploy MCMS again so that only the proposer is new
	updatedEnv, err = commonchangeset.Apply(t, updatedEnv, nil, configuredChangeset)
	require.NoError(t, err)
	mcmsState, err = mcmschangesetstate.MaybeLoadMCMSWithTimelockState(updatedEnv, evmSelectors)
	require.NoError(t, err)

	require.NotEqual(t, existingProposer.Address(), mcmsState[evmSelectors[0]].ProposerMcm.Address())
	updatedEnv, err = commonchangeset.Apply(t, updatedEnv, nil, commonchangeset.Configure(
		commonchangeset.GrantRoleInTimeLock,
		commonchangeset.GrantRoleInput{
			ExistingProposerByChain: map[uint64]common.Address{
				evmSelectors[0]: existingProposer.Address(),
			},
			MCMS: &proposalutils.TimelockConfig{MinDelay: 0},
		},
	))
	require.NoError(t, err)
	mcmsState, err = mcmschangesetstate.MaybeLoadMCMSWithTimelockState(updatedEnv, evmSelectors)
	require.NoError(t, err)

	evmTimelockInspector := mcmsevmsdk.NewTimelockInspector(updatedEnv.Chains[evmSelectors[0]].Client)

	proposers, err := evmTimelockInspector.GetProposers(ctx, mcmsState[evmSelectors[0]].Timelock.Address().Hex())
	require.NoError(t, err)
	require.Contains(t, proposers, mcmsState[evmSelectors[0]].ProposerMcm.Address().Hex())
	require.Contains(t, proposers, existingProposer.Address().Hex())
}

func TestDeployMCMSWithTimelockV2WithFewExistingContracts(t *testing.T) {
	ctx := testutils.Context(t)
	env := memory.NewMemoryEnvironment(t, logger.TestLogger(t), zapcore.InfoLevel, memory.MemoryEnvironmentConfig{Chains: 2})
	evmSelectors := env.AllChainSelectors()
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
	}

	// set up some dummy address in env address book for callproxy, canceller and bypasser
	// to simulate the case where they already exist
	// this is to test that the changeset will not try to deploy them again
	addrBook := deployment.NewMemoryAddressBook()
	callProxyAddress := utils.RandomAddress()
	mcmsAddress := utils.RandomAddress()
	mcmsType := deployment.NewTypeAndVersion(commontypes.ManyChainMultisig, deployment.Version1_0_0)
	// we use same address for bypasser and canceller
	mcmsType.AddLabel(commontypes.BypasserRole.String())
	mcmsType.AddLabel(commontypes.CancellerRole.String())
	require.NoError(t, addrBook.Save(evmSelectors[0], callProxyAddress.String(),
		deployment.NewTypeAndVersion(commontypes.CallProxy, deployment.Version1_0_0)))
	require.NoError(t, addrBook.Save(evmSelectors[0], mcmsAddress.String(), mcmsType))
	require.NoError(t, env.ExistingAddresses.Merge(addrBook))

	configuredChangeset := commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
		changesetConfig,
	)

	// --- act ---
	updatedEnv, err := commonchangeset.Apply(t, env, nil, configuredChangeset)
	require.NoError(t, err)

	state, err := mcmschangesetstate.MaybeLoadMCMSWithTimelockState(updatedEnv, evmSelectors)
	require.NoError(t, err)
	evmState0 := state[evmSelectors[0]]

	// --- assert ---
	require.Equal(t, callProxyAddress, evmState0.CallProxy.Address())
	require.Equal(t, mcmsAddress, evmState0.BypasserMcm.Address())
	require.Equal(t, mcmsAddress, evmState0.CancellerMcm.Address())
	// proposer should be newly deployed
	require.NotEqual(t, mcmsAddress, evmState0.ProposerMcm.Address())

	evmTimelockInspector := mcmsevmsdk.NewTimelockInspector(updatedEnv.Chains[evmSelectors[0]].Client)

	proposers, err := evmTimelockInspector.GetProposers(ctx, evmState0.Timelock.Address().Hex())
	require.NoError(t, err)
	require.Equal(t, proposers, []string{evmState0.ProposerMcm.Address().Hex()})

	executors, err := evmTimelockInspector.GetExecutors(ctx, evmState0.Timelock.Address().Hex())
	require.NoError(t, err)
	require.Equal(t, executors, []string{evmState0.CallProxy.Address().Hex()})

	cancellers, err := evmTimelockInspector.GetCancellers(ctx, evmState0.Timelock.Address().Hex())
	require.NoError(t, err)
	require.ElementsMatch(t, cancellers, []string{
		evmState0.CancellerMcm.Address().Hex(), // bypasser and canceller are same
		evmState0.ProposerMcm.Address().Hex(),
	})

	bypassers, err := evmTimelockInspector.GetBypassers(ctx, evmState0.Timelock.Address().Hex())
	require.NoError(t, err)
	require.Equal(t, bypassers, []string{evmState0.BypasserMcm.Address().Hex()})
}

func TestDeployMCMSWithTimelockV2(t *testing.T) {
	t.Parallel()
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
	configuredChangeset := commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
		changesetConfig,
	)
	commonchangeset.SetPreloadedSolanaAddresses(t, env, solanaSelectors[0])

	// --- act ---
	updatedEnv, err := commonchangeset.Apply(t, env, nil, configuredChangeset)
	require.NoError(t, err)

	evmState, err := mcmschangesetstate.MaybeLoadMCMSWithTimelockState(updatedEnv, evmSelectors)
	require.NoError(t, err)
	solanaState, err := mcmschangesetstate.MaybeLoadMCMSWithTimelockStateSolana(updatedEnv, solanaSelectors)
	require.NoError(t, err)

	// --- assert ---
	require.Len(t, evmState, 2)
	require.Len(t, solanaState, 1)
	ctx := updatedEnv.GetContext()

	// evm chain 0
	evmState0 := evmState[evmSelectors[0]]
	evmInspector := mcmsevmsdk.NewInspector(updatedEnv.Chains[evmSelectors[0]].Client)
	evmTimelockInspector := mcmsevmsdk.NewTimelockInspector(updatedEnv.Chains[evmSelectors[0]].Client)

	config, err := evmInspector.GetConfig(ctx, evmState0.ProposerMcm.Address().Hex())
	require.NoError(t, err)
	require.Empty(t, cmp.Diff(*config, changesetConfig[evmSelectors[0]].Proposer))

	config, err = evmInspector.GetConfig(ctx, evmState0.CancellerMcm.Address().Hex())
	require.NoError(t, err)
	require.Empty(t, cmp.Diff(*config, changesetConfig[evmSelectors[0]].Canceller))

	config, err = evmInspector.GetConfig(ctx, evmState0.BypasserMcm.Address().Hex())
	require.NoError(t, err)
	require.Empty(t, cmp.Diff(*config, changesetConfig[evmSelectors[0]].Bypasser))

	proposers, err := evmTimelockInspector.GetProposers(ctx, evmState0.Timelock.Address().Hex())
	require.NoError(t, err)
	require.Equal(t, proposers, []string{evmState0.ProposerMcm.Address().Hex()})

	executors, err := evmTimelockInspector.GetExecutors(ctx, evmState0.Timelock.Address().Hex())
	require.NoError(t, err)
	require.Equal(t, executors, []string{evmState0.CallProxy.Address().Hex()})

	cancellers, err := evmTimelockInspector.GetCancellers(ctx, evmState0.Timelock.Address().Hex())
	require.NoError(t, err)
	require.ElementsMatch(t, cancellers, []string{
		evmState0.CancellerMcm.Address().Hex(),
		evmState0.ProposerMcm.Address().Hex(),
		evmState0.BypasserMcm.Address().Hex(),
	})

	bypassers, err := evmTimelockInspector.GetBypassers(ctx, evmState0.Timelock.Address().Hex())
	require.NoError(t, err)
	require.Equal(t, bypassers, []string{evmState0.BypasserMcm.Address().Hex()})

	// evm chain 1
	evmState1 := evmState[evmSelectors[1]]
	evmInspector = mcmsevmsdk.NewInspector(updatedEnv.Chains[evmSelectors[1]].Client)
	evmTimelockInspector = mcmsevmsdk.NewTimelockInspector(updatedEnv.Chains[evmSelectors[1]].Client)

	config, err = evmInspector.GetConfig(ctx, evmState1.ProposerMcm.Address().Hex())
	require.NoError(t, err)
	require.Empty(t, cmp.Diff(*config, changesetConfig[evmSelectors[1]].Proposer))

	config, err = evmInspector.GetConfig(ctx, evmState1.CancellerMcm.Address().Hex())
	require.NoError(t, err)
	require.Empty(t, cmp.Diff(*config, changesetConfig[evmSelectors[1]].Canceller))

	config, err = evmInspector.GetConfig(ctx, evmState1.BypasserMcm.Address().Hex())
	require.NoError(t, err)
	require.Empty(t, cmp.Diff(*config, changesetConfig[evmSelectors[1]].Bypasser))

	proposers, err = evmTimelockInspector.GetProposers(ctx, evmState1.Timelock.Address().Hex())
	require.NoError(t, err)
	require.Equal(t, proposers, []string{evmState1.ProposerMcm.Address().Hex()})

	executors, err = evmTimelockInspector.GetExecutors(ctx, evmState1.Timelock.Address().Hex())
	require.NoError(t, err)
	require.Equal(t, executors, []string{evmState1.CallProxy.Address().Hex()})

	cancellers, err = evmTimelockInspector.GetCancellers(ctx, evmState1.Timelock.Address().Hex())
	require.NoError(t, err)
	require.ElementsMatch(t, cancellers, []string{
		evmState1.CancellerMcm.Address().Hex(),
		evmState1.ProposerMcm.Address().Hex(),
		evmState1.BypasserMcm.Address().Hex(),
	})

	bypassers, err = evmTimelockInspector.GetBypassers(ctx, evmState1.Timelock.Address().Hex())
	require.NoError(t, err)
	require.Equal(t, bypassers, []string{evmState1.BypasserMcm.Address().Hex()})

	// solana chain 0
	solanaState0 := solanaState[solanaSelectors[0]]
	solanaChain0 := updatedEnv.SolChains[solanaSelectors[0]]
	solanaInspector := mcmssolanasdk.NewInspector(solanaChain0.Client)
	solanaTimelockInspector := mcmssolanasdk.NewTimelockInspector(solanaChain0.Client)

	addr := mcmssolanasdk.ContractAddress(solanaState0.McmProgram, mcmssolanasdk.PDASeed(solanaState0.ProposerMcmSeed))
	config, err = solanaInspector.GetConfig(ctx, addr)
	require.NoError(t, err)
	require.Empty(t, cmp.Diff(*config, changesetConfig[solanaSelectors[0]].Proposer))

	addr = mcmssolanasdk.ContractAddress(solanaState0.McmProgram, mcmssolanasdk.PDASeed(solanaState0.CancellerMcmSeed))
	config, err = solanaInspector.GetConfig(ctx, addr)
	require.NoError(t, err)
	require.Empty(t, cmp.Diff(*config, changesetConfig[solanaSelectors[0]].Canceller))

	addr = mcmssolanasdk.ContractAddress(solanaState0.McmProgram, mcmssolanasdk.PDASeed(solanaState0.BypasserMcmSeed))
	config, err = solanaInspector.GetConfig(ctx, addr)
	require.NoError(t, err)
	require.Empty(t, cmp.Diff(*config, changesetConfig[solanaSelectors[0]].Bypasser))

	addr = mcmssolanasdk.ContractAddress(solanaState0.TimelockProgram, mcmssolanasdk.PDASeed(solanaState0.TimelockSeed))
	proposers, err = solanaTimelockInspector.GetProposers(ctx, addr)
	require.NoError(t, err)
	require.Equal(t, proposers, []string{mcmSignerPDA(solanaState0.McmProgram, solanaState0.ProposerMcmSeed)})

	executors, err = solanaTimelockInspector.GetExecutors(ctx, addr)
	require.NoError(t, err)
	require.Equal(t, executors, []string{solanaChain0.DeployerKey.PublicKey().String()})

	cancellers, err = solanaTimelockInspector.GetCancellers(ctx, addr)
	require.NoError(t, err)
	require.ElementsMatch(t, cancellers, []string{
		mcmSignerPDA(solanaState0.McmProgram, solanaState0.CancellerMcmSeed),
		mcmSignerPDA(solanaState0.McmProgram, solanaState0.ProposerMcmSeed),
		mcmSignerPDA(solanaState0.McmProgram, solanaState0.BypasserMcmSeed),
	})

	bypassers, err = solanaTimelockInspector.GetBypassers(ctx, addr)
	require.NoError(t, err)
	require.Equal(t, bypassers, []string{mcmSignerPDA(solanaState0.McmProgram, solanaState0.BypasserMcmSeed)})

	timelockConfig := solanaTimelockConfig(ctx, t, solanaChain0, solanaState0.TimelockProgram, solanaState0.TimelockSeed)
	require.NoError(t, err)
	require.Equal(t, timelockConfig.ProposedOwner.String(), "11111111111111111111111111111111")
}

// TestDeployMCMSWithTimelockV2SkipInit tests calling the deploy changeset when accounts have already been initialized
func TestDeployMCMSWithTimelockV2SkipInitSolana(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/DX-438")

	t.Parallel()
	// --- arrange ---
	log := logger.TestLogger(t)
	envConfig := memory.MemoryEnvironmentConfig{Chains: 0, SolChains: 1}
	env := memory.NewMemoryEnvironment(t, log, zapcore.InfoLevel, envConfig)
	solanaSelectors := env.AllChainSelectorsSolana()
	changesetConfig := map[uint64]commontypes.MCMSWithTimelockConfigV2{
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
	configuredChangeset := commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
		changesetConfig,
	)
	commonchangeset.SetPreloadedSolanaAddresses(t, env, solanaSelectors[0])
	// --- act ---
	updatedEnv, err := commonchangeset.Apply(t, env, nil, configuredChangeset)
	require.NoError(t, err)

	solanaState, err := mcmschangesetstate.MaybeLoadMCMSWithTimelockStateSolana(updatedEnv, solanaSelectors)
	require.NoError(t, err)

	// Call deploy again, seeds and addresses from original state should not change
	updatedEnvReTriggered, err := commonchangeset.Apply(t, updatedEnv, nil, configuredChangeset)
	require.NoError(t, err)
	solanaStateNew, err := mcmschangesetstate.MaybeLoadMCMSWithTimelockStateSolana(updatedEnvReTriggered, solanaSelectors)
	require.NoError(t, err)

	// --- assert ---
	require.Len(t, solanaState, 1)
	stateOld := solanaState[solanaSelectors[0]]
	stateNew := solanaStateNew[solanaSelectors[0]]
	require.Equal(t, stateOld.TimelockSeed, stateNew.TimelockSeed)
	require.Equal(t, stateOld.TimelockProgram, stateNew.TimelockProgram)
	require.Equal(t, stateOld.BypasserAccessControllerAccount, stateNew.BypasserAccessControllerAccount)
	require.Equal(t, stateOld.CancellerAccessControllerAccount, stateNew.CancellerAccessControllerAccount)
	require.Equal(t, stateOld.ExecutorAccessControllerAccount, stateNew.ExecutorAccessControllerAccount)
	require.Equal(t, stateOld.ProposerAccessControllerAccount, stateNew.ProposerAccessControllerAccount)
	require.Equal(t, stateOld.McmProgram, stateNew.McmProgram)
	require.Equal(t, stateOld.BypasserMcmSeed, stateNew.BypasserMcmSeed)
	require.Equal(t, stateOld.CancellerMcmSeed, stateNew.CancellerMcmSeed)
	require.Equal(t, stateOld.ProposerMcmSeed, stateNew.ProposerMcmSeed)
	require.Equal(t, stateOld.AccessControllerProgram, stateNew.AccessControllerProgram)
}

// ----- helpers -----

func mcmSignerPDA(programID solana.PublicKey, seed mcmschangesetstate.PDASeed) string {
	return mcmschangesetstate.GetMCMSignerPDA(programID, seed).String()
}

func timelockSignerPDA(programID solana.PublicKey, seed mcmschangesetstate.PDASeed) string {
	return mcmschangesetstate.GetTimelockSignerPDA(programID, seed).String()
}

func solanaTimelockConfig(
	ctx context.Context, t *testing.T, chain deployment.SolChain, programID solana.PublicKey, seed mcmschangesetstate.PDASeed,
) timelockBindings.Config {
	t.Helper()

	var data timelockBindings.Config
	err := chain.GetAccountDataBorshInto(ctx, mcmschangesetstate.GetTimelockConfigPDA(programID, seed), &data)
	require.NoError(t, err)

	return data
}
