//nolint:testifylint // inverting want and got is more succinct
package changeset_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/google/go-cmp/cmp"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	mcmsevmsdk "github.com/smartcontractkit/mcms/sdk/evm"
	mcmssolanasdk "github.com/smartcontractkit/mcms/sdk/solana"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/smartcontractkit/quarantine"
	"github.com/stretchr/testify/require"

	timelockBindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/timelock"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/onchain"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	mcmschangesetstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/internal/soltestutils"
)

func TestGrantRoleInTimeLock(t *testing.T) {
	selector := chain_selectors.TEST_90000001.Selector
	env, err := environment.New(t.Context(),
		environment.WithEVMSimulatedWithConfig(t, []uint64{selector}, onchain.EVMSimLoaderConfig{
			NumAdditionalAccounts: 1,
		}),
	)
	require.NoError(t, err)

	// deploy the MCMS with timelock contracts
	configuredChangeset := commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
		map[uint64]commontypes.MCMSWithTimelockConfigV2{
			selector: proposalutils.SingleGroupTimelockConfigV2(t),
		},
	)
	updatedEnv, err := commonchangeset.Apply(t, *env, configuredChangeset)
	require.NoError(t, err)
	mcmsState, err := mcmschangesetstate.MaybeLoadMCMSWithTimelockState(updatedEnv, []uint64{selector})
	require.NoError(t, err)

	// change the environment to remove proposer from the timelock, so that we can deploy new proposer
	// and then grant the role to the new proposer
	existingProposer := mcmsState[selector].ProposerMcm
	ab := cldf.NewMemoryAddressBook()
	require.NoError(t, ab.Save(selector, existingProposer.Address().String(),
		cldf.NewTypeAndVersion(commontypes.ProposerManyChainMultisig, deployment.Version1_0_0)))
	require.NoError(t, updatedEnv.ExistingAddresses.Remove(ab))

	// remove from DataStore since deployment now uses DataStore
	// Since DataStore is immutable, create a new one without the proposer
	newDataStore := datastore.NewMemoryDataStore()
	refs, err := updatedEnv.DataStore.Addresses().Fetch()
	require.NoError(t, err)

	// Copy all address refs except the proposer we want to remove
	for _, ref := range refs {
		// Skip the proposer we want to remove
		if ref.ChainSelector == selector &&
			ref.Address == existingProposer.Address().String() &&
			ref.Type == datastore.ContractType(commontypes.ProposerManyChainMultisig) {
			continue
		}
		err := newDataStore.Addresses().Add(ref)
		require.NoError(t, err)
	}

	// Replace the DataStore in the environment
	updatedEnv.DataStore = newDataStore.Seal()

	// change the deployer key, so that we can deploy proposer with a new key
	// the new deployer key will not be admin of the timelock
	// we can test granting roles through proposal
	evmChains := updatedEnv.BlockChains.EVMChains()
	chain := evmChains[selector]
	chain.DeployerKey = evmChains[selector].Users[0]

	// now deploy MCMS again so that only the proposer is new
	updatedEnv, err = commonchangeset.Apply(t, updatedEnv, configuredChangeset)
	require.NoError(t, err)
	mcmsState, err = mcmschangesetstate.MaybeLoadMCMSWithTimelockState(updatedEnv, []uint64{selector})
	require.NoError(t, err)

	require.NotEqual(t, existingProposer.Address(), mcmsState[selector].ProposerMcm.Address())
	updatedEnv, err = commonchangeset.Apply(t, updatedEnv, commonchangeset.Configure(
		commonchangeset.GrantRoleInTimeLock,
		commonchangeset.GrantRoleInput{
			ExistingProposerByChain: map[uint64]common.Address{
				selector: existingProposer.Address(),
			},
			MCMS: &proposalutils.TimelockConfig{MinDelay: 0},
		},
	))
	require.NoError(t, err)
	mcmsState, err = mcmschangesetstate.MaybeLoadMCMSWithTimelockState(updatedEnv, []uint64{selector})
	require.NoError(t, err)

	evmTimelockInspector := mcmsevmsdk.NewTimelockInspector(updatedEnv.BlockChains.EVMChains()[selector].Client)

	proposers, err := evmTimelockInspector.GetProposers(t.Context(), mcmsState[selector].Timelock.Address().Hex())
	require.NoError(t, err)
	require.Contains(t, proposers, mcmsState[selector].ProposerMcm.Address().Hex())
	require.Contains(t, proposers, existingProposer.Address().Hex())
}

func TestDeployMCMSWithTimelockV2WithFewExistingContracts(t *testing.T) {
	ctx := t.Context()

	selector1 := chain_selectors.TEST_90000001.Selector
	selector2 := chain_selectors.TEST_90000002.Selector
	selectors := []uint64{selector1, selector2}

	// Build a datastore with some dummy address for callproxy, canceller and bypasser
	// to simulate the case where they already exist and so the changeset will not try to deploy
	// them again
	ds := datastore.NewMemoryDataStore()

	callProxyAddress := utils.RandomAddress()
	mcmsAddress := utils.RandomAddress()
	mcmsType := cldf.NewTypeAndVersion(commontypes.ManyChainMultisig, deployment.Version1_0_0)
	// we use same address for bypasser and canceller
	mcmsType.AddLabel(commontypes.BypasserRole.String())
	mcmsType.AddLabel(commontypes.CancellerRole.String())

	// Add CallProxy for first chain only
	require.NoError(t, ds.AddressRefStore.Add(datastore.AddressRef{
		ChainSelector: selector1,
		Address:       callProxyAddress.String(),
		Type:          datastore.ContractType(commontypes.CallProxy),
		Version:       &deployment.Version1_0_0,
	}))

	// Add MCMS contract with both bypasser and canceller labels for first chain only
	require.NoError(t, ds.AddressRefStore.Add(datastore.AddressRef{
		ChainSelector: selector1,
		Address:       mcmsAddress.String(),
		Type:          datastore.ContractType(mcmsType.Type),
		Version:       &mcmsType.Version,
		Labels:        datastore.NewLabelSet(mcmsType.Labels.List()...),
	}))

	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, selectors),
		environment.WithLogger(logger.Test(t)),
		environment.WithDatastore(ds.Seal()),
	))
	require.NoError(t, err)

	chain1 := rt.Environment().BlockChains.EVMChains()[selector1]

	changesetConfig := map[uint64]commontypes.MCMSWithTimelockConfigV2{
		selector1: {
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
		selector2: {
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

	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2), changesetConfig),
	)
	require.NoError(t, err)

	state, err := mcmschangesetstate.MaybeLoadMCMSWithTimelockState(rt.Environment(), selectors)
	require.NoError(t, err)
	evmState1 := state[selector1]

	// --- assert ---
	require.Equal(t, callProxyAddress, evmState1.CallProxy.Address())
	require.Equal(t, mcmsAddress, evmState1.BypasserMcm.Address())
	require.Equal(t, mcmsAddress, evmState1.CancellerMcm.Address())
	// proposer should be newly deployed
	require.NotEqual(t, mcmsAddress, evmState1.ProposerMcm.Address())

	evmTimelockInspector := mcmsevmsdk.NewTimelockInspector(chain1.Client)

	proposers, err := evmTimelockInspector.GetProposers(ctx, evmState1.Timelock.Address().Hex())
	require.NoError(t, err)
	require.Equal(t, proposers, []string{evmState1.ProposerMcm.Address().Hex()})

	executors, err := evmTimelockInspector.GetExecutors(ctx, evmState1.Timelock.Address().Hex())
	require.NoError(t, err)
	require.Equal(t, executors, []string{evmState1.CallProxy.Address().Hex()})

	cancellers, err := evmTimelockInspector.GetCancellers(ctx, evmState1.Timelock.Address().Hex())
	require.NoError(t, err)
	require.ElementsMatch(t, cancellers, []string{
		evmState1.CancellerMcm.Address().Hex(), // bypasser and canceller are same
		evmState1.ProposerMcm.Address().Hex(),
	})

	bypassers, err := evmTimelockInspector.GetBypassers(ctx, evmState1.Timelock.Address().Hex())
	require.NoError(t, err)
	require.Equal(t, bypassers, []string{evmState1.BypasserMcm.Address().Hex()})
}

func TestDeployMCMSWithTimelockV2(t *testing.T) {
	quarantine.Flaky(t, "DX-1719")
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

	evmChain := rt.Environment().BlockChains.EVMChains()[evmSelector]
	solChain := rt.Environment().BlockChains.SolanaChains()[solSelector]

	changesetConfig := map[uint64]commontypes.MCMSWithTimelockConfigV2{
		evmSelector: {
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
		solSelector: {
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

	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2), changesetConfig),
	)
	require.NoError(t, err)

	// Load the MCMS contractsstate
	evmMCMSState, err := mcmschangesetstate.MaybeLoadMCMSWithTimelockState(rt.Environment(), []uint64{evmSelector})
	require.NoError(t, err)
	require.Len(t, evmMCMSState, 1)

	solMCMSState := soltestutils.GetMCMSStateFromAddressBook(t, rt.State().AddressBook, solChain)

	// --- assert ---

	ctx := t.Context()

	// evm chain
	evmState := evmMCMSState[evmSelector]
	evmInspector := mcmsevmsdk.NewInspector(evmChain.Client)
	evmTimelockInspector := mcmsevmsdk.NewTimelockInspector(evmChain.Client)

	config, err := evmInspector.GetConfig(ctx, evmState.ProposerMcm.Address().Hex())
	require.NoError(t, err)
	require.Empty(t, cmp.Diff(*config, changesetConfig[evmSelector].Proposer))

	config, err = evmInspector.GetConfig(ctx, evmState.CancellerMcm.Address().Hex())
	require.NoError(t, err)
	require.Empty(t, cmp.Diff(*config, changesetConfig[evmSelector].Canceller))

	config, err = evmInspector.GetConfig(ctx, evmState.BypasserMcm.Address().Hex())
	require.NoError(t, err)
	require.Empty(t, cmp.Diff(*config, changesetConfig[evmSelector].Bypasser))

	proposers, err := evmTimelockInspector.GetProposers(ctx, evmState.Timelock.Address().Hex())
	require.NoError(t, err)
	require.Equal(t, proposers, []string{evmState.ProposerMcm.Address().Hex()})

	executors, err := evmTimelockInspector.GetExecutors(ctx, evmState.Timelock.Address().Hex())
	require.NoError(t, err)
	require.Equal(t, executors, []string{evmState.CallProxy.Address().Hex()})

	cancellers, err := evmTimelockInspector.GetCancellers(ctx, evmState.Timelock.Address().Hex())
	require.NoError(t, err)
	require.ElementsMatch(t, cancellers, []string{
		evmState.CancellerMcm.Address().Hex(),
		evmState.ProposerMcm.Address().Hex(),
		evmState.BypasserMcm.Address().Hex(),
	})

	bypassers, err := evmTimelockInspector.GetBypassers(ctx, evmState.Timelock.Address().Hex())
	require.NoError(t, err)
	require.Equal(t, bypassers, []string{evmState.BypasserMcm.Address().Hex()})

	// solana chain
	solanaInspector := mcmssolanasdk.NewInspector(solChain.Client)
	solanaTimelockInspector := mcmssolanasdk.NewTimelockInspector(solChain.Client)

	addr := mcmssolanasdk.ContractAddress(solMCMSState.McmProgram, mcmssolanasdk.PDASeed(solMCMSState.ProposerMcmSeed))
	config, err = solanaInspector.GetConfig(ctx, addr)
	require.NoError(t, err)
	require.Empty(t, cmp.Diff(*config, changesetConfig[solSelector].Proposer))

	addr = mcmssolanasdk.ContractAddress(solMCMSState.McmProgram, mcmssolanasdk.PDASeed(solMCMSState.CancellerMcmSeed))
	config, err = solanaInspector.GetConfig(ctx, addr)
	require.NoError(t, err)
	require.Empty(t, cmp.Diff(*config, changesetConfig[solSelector].Canceller))

	addr = mcmssolanasdk.ContractAddress(solMCMSState.McmProgram, mcmssolanasdk.PDASeed(solMCMSState.BypasserMcmSeed))
	config, err = solanaInspector.GetConfig(ctx, addr)
	require.NoError(t, err)
	require.Empty(t, cmp.Diff(*config, changesetConfig[solSelector].Bypasser))

	addr = mcmssolanasdk.ContractAddress(solMCMSState.TimelockProgram, mcmssolanasdk.PDASeed(solMCMSState.TimelockSeed))
	proposers, err = solanaTimelockInspector.GetProposers(ctx, addr)
	require.NoError(t, err)
	require.Equal(t, proposers, []string{mcmSignerPDA(solMCMSState.McmProgram, solMCMSState.ProposerMcmSeed)})

	executors, err = solanaTimelockInspector.GetExecutors(ctx, addr)
	require.NoError(t, err)
	require.Equal(t, executors, []string{solChain.DeployerKey.PublicKey().String()})

	cancellers, err = solanaTimelockInspector.GetCancellers(ctx, addr)
	require.NoError(t, err)
	require.ElementsMatch(t, cancellers, []string{
		mcmSignerPDA(solMCMSState.McmProgram, solMCMSState.CancellerMcmSeed),
		mcmSignerPDA(solMCMSState.McmProgram, solMCMSState.ProposerMcmSeed),
		mcmSignerPDA(solMCMSState.McmProgram, solMCMSState.BypasserMcmSeed),
	})

	bypassers, err = solanaTimelockInspector.GetBypassers(ctx, addr)
	require.NoError(t, err)
	require.Equal(t, bypassers, []string{mcmSignerPDA(solMCMSState.McmProgram, solMCMSState.BypasserMcmSeed)})

	timelockConfig := solanaTimelockConfig(ctx, t, solChain, solMCMSState.TimelockProgram, solMCMSState.TimelockSeed)
	require.NoError(t, err)
	require.Equal(t, timelockConfig.ProposedOwner.String(), "11111111111111111111111111111111")
}

// TestDeployMCMSWithTimelockV2SkipInit tests calling the deploy changeset when accounts have already been initialized
func TestDeployMCMSWithTimelockV2SkipInitSolana(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/DX-438")

	t.Parallel()

	selector := chain_selectors.TEST_22222222222222222222222222222222222222222222.Selector
	programsPath, programIDs, ab := soltestutils.PreloadMCMS(t, selector)

	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithSolanaContainer(t, []uint64{selector}, programsPath, programIDs),
		environment.WithAddressBook(ab),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	changesetConfig := map[uint64]commontypes.MCMSWithTimelockConfigV2{
		selector: {
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

	// --- act ---
	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2), changesetConfig),
	)
	require.NoError(t, err)

	solanaState, err := mcmschangesetstate.MaybeLoadMCMSWithTimelockStateSolana(rt.Environment(), []uint64{selector})
	require.NoError(t, err)

	// Call deploy again, seeds and addresses from original state should not change
	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2), changesetConfig),
	)
	require.NoError(t, err)

	solanaStateNew, err := mcmschangesetstate.MaybeLoadMCMSWithTimelockStateSolana(rt.Environment(), []uint64{selector})
	require.NoError(t, err)

	// --- assert ---
	require.Len(t, solanaState, 1)
	stateOld := solanaState[selector]
	stateNew := solanaStateNew[selector]
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

func solanaTimelockConfig(
	ctx context.Context, t *testing.T, chain cldf_solana.Chain, programID solana.PublicKey, seed mcmschangesetstate.PDASeed,
) timelockBindings.Config {
	t.Helper()

	var data timelockBindings.Config
	err := chain.GetAccountDataBorshInto(ctx, mcmschangesetstate.GetTimelockConfigPDA(programID, seed), &data)
	require.NoError(t, err)

	return data
}
