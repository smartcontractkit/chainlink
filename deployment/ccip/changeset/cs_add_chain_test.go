package changeset

import (
	"math/big"
	"slices"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	changesetcommon "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"
)

func Test_AddChain(t *testing.T) {
	t.Parallel()

	const (
		numChains     = 4
		usersPerChain = 2
	)

	// Set up an env with 4 chains but initially
	// only deploy and configure 3 of them.
	e, tEnv := NewMemoryEnvironment(
		t,
		WithChains(numChains),
		WithNodes(4),
		WithPrerequisiteDeployment(),
		WithUsersPerChain(usersPerChain),
		WithNoJobsAndContracts(),
	)

	allChains := maps.Keys(e.Env.Chains)
	slices.Sort(allChains)
	toDeploy := e.Env.AllChainSelectorsExcluding([]uint64{allChains[0]})
	require.Len(t, toDeploy, numChains-1)
	remainingChain := allChains[0]
	t.Log("initially deploying chains:", toDeploy, "and afterwards adding chain", remainingChain)

	e = AddCCIPContractsToEnvironment(
		t,
		toDeploy,
		tEnv,
		true,  // deployJobs
		true,  // deployHomeChain
		false, // mcmsEnabled
	)

	// Need to update what the RMNProxy is pointing to, otherwise plugin will not work.
	var err error
	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, e.TimelockContracts(t), []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(SetRMNRemoteOnRMNProxy),
			Config: SetRMNRemoteOnRMNProxyConfig{
				ChainSelectors: toDeploy,
			},
		},
	})
	require.NoError(t, err)

	state, err := LoadOnchainState(e.Env)
	require.NoError(t, err)

	// Setup densely connected lanes between all chains.
	for _, source := range toDeploy {
		for _, dest := range toDeploy {
			if source == dest {
				continue
			}
			AddLaneWithDefaultPricesAndFeeQuoterConfig(
				t,
				&e,
				state,
				source,
				dest,
				false, // isTestRouter
			)
		}
	}

	// Transfer ownership of all contracts to the MCMS and renounce the timelock deployer.
	transferToMCMSAndRenounceTimelockDeployer(t, e, toDeploy, state)

	// check RMNRemote is up and RMNProxy is correctly wired.
	assertRMNRemoteAndProxyState(t, toDeploy, state)

	// At this stage we can send some requests and confirm the setup is working.
	sendMsgs := func(
		sources []uint64,
		dests []uint64,
		testRouter bool,
	) (gasPricePreUpdate map[SourceDestPair]*big.Int, startBlocks map[uint64]*uint64) {
		startBlocks = make(map[uint64]*uint64)
		gasPricePreUpdate = make(map[SourceDestPair]*big.Int)
		var (
			expectedSeqNum     = make(map[SourceDestPair]uint64)
			expectedSeqNumExec = make(map[SourceDestPair][]uint64)
		)
		for _, source := range sources {
			for _, dest := range dests {
				if source == dest {
					continue
				}

				gp, err := state.Chains[source].FeeQuoter.GetDestinationChainGasPrice(&bind.CallOpts{
					Context: tests.Context(t),
				}, dest)
				require.NoError(t, err)
				gasPricePreUpdate[SourceDestPair{
					SourceChainSelector: source,
					DestChainSelector:   dest,
				}] = gp.Value

				latesthdr, err := e.Env.Chains[dest].Client.HeaderByNumber(testcontext.Get(t), nil)
				require.NoError(t, err)
				block := latesthdr.Number.Uint64()
				msgSentEvent := TestSendRequest(t, e.Env, state, source, dest, testRouter, router.ClientEVM2AnyMessage{
					Receiver:     common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32),
					Data:         []byte("hello world"),
					TokenAmounts: nil,
					FeeToken:     common.HexToAddress("0x0"),
					ExtraArgs:    nil,
				})

				startBlocks[dest] = &block
				expectedSeqNum[SourceDestPair{
					SourceChainSelector: source,
					DestChainSelector:   dest,
				}] = msgSentEvent.SequenceNumber
				expectedSeqNumExec[SourceDestPair{
					SourceChainSelector: source,
					DestChainSelector:   dest,
				}] = append(expectedSeqNumExec[SourceDestPair{
					SourceChainSelector: source,
					DestChainSelector:   dest,
				}], msgSentEvent.SequenceNumber)
			}
		}

		// Confirm execution of the message
		ConfirmCommitForAllWithExpectedSeqNums(t, e.Env, state, expectedSeqNum, startBlocks)
		ConfirmExecWithSeqNrsForAll(t, e.Env, state, expectedSeqNumExec, startBlocks)
		return gasPricePreUpdate, startBlocks
	}

	// wait for plugins to come up.
	time.Sleep(30 * time.Second)
	sendMsgs(toDeploy, toDeploy, false)

	// TODO: Not working. Need to fix/figure out why.
	// gasPricePreUpdate, startBlocks := sendMsgs(toDeploy)
	// for sourceDestPair, preUpdateGp := range gasPricePreUpdate {
	// 	// check that each chain's fee quoter has updated its gas price
	// 	// for all dests.
	// 	err := ConfirmGasPriceUpdated(
	// 		t,
	// 		e.Env.Chains[sourceDestPair.DestChainSelector],
	// 		state.Chains[sourceDestPair.SourceChainSelector].FeeQuoter,
	// 		*startBlocks[sourceDestPair.DestChainSelector],
	// 		preUpdateGp,
	// 	)
	// 	require.NoError(t, err)
	// }

	// Deploy to the remaining chain.
	// MCMS needs to be enabled because the home chain contracts have been
	// transferred to MCMS.
	e = AddCCIPContractsToEnvironment(
		t,
		[]uint64{remainingChain},
		tEnv,
		false, // deployJobs
		false, // deployHomeChain
		true,  // mcmsEnabled
	)

	// Need to update what the RMNProxy is pointing to, otherwise plugin will not work.
	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, e.TimelockContracts(t), []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(SetRMNRemoteOnRMNProxy),
			Config: SetRMNRemoteOnRMNProxyConfig{
				ChainSelectors: []uint64{remainingChain},
			},
		},
	})
	require.NoError(t, err)

	state, err = LoadOnchainState(e.Env)
	require.NoError(t, err)

	assertRMNRemoteAndProxyState(t, []uint64{remainingChain}, state)

	// TODO: wait for gas price of new chain to be updated on all other chains.

	// UpdateOnRampDestsConfig on the existing chains with TestRouter=true and MCMS config enabled.
	// UpdateFeeQuoterPrices (new CS) to set initial FQ prices to the new destination.
	// UpdateFeeQuoterDestsConfig to enable quoting for the new destination
	// UpdateRouterRampsConfig on the test router to enable the new destination
	mcmsConfig := &MCMSConfig{
		MinDelay: 0,
	}
	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, e.TimelockContracts(t), []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(UpdateOnRampsDests),
			Config: UpdateOnRampDestsConfig{
				UpdatesByChain: onRampDestUpdates(t, []uint64{remainingChain}, toDeploy, true),
				MCMS:           mcmsConfig,
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(UpdateFeeQuoterPricesCS),
			Config: UpdateFeeQuoterPricesConfig{
				PricesByChain: feeQuoterPricesByChain(t, []uint64{remainingChain}, toDeploy),
				MCMS:          mcmsConfig,
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(UpdateFeeQuoterDests),
			Config: UpdateFeeQuoterDestsConfig{
				UpdatesByChain: feeQuoterDestUpdates(t, []uint64{remainingChain}, toDeploy),
				MCMS:           mcmsConfig,
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(UpdateRouterRamps),
			Config: UpdateRouterRampsConfig{
				TestRouter:     true,
				UpdatesByChain: routerOnRampUpdates(t, []uint64{remainingChain}, toDeploy),
				MCMS:           mcmsConfig,
			},
		},
	})
	require.NoError(t, err)

	state, err = LoadOnchainState(e.Env)
	require.NoError(t, err)

	assertExistingChainsWiring(
		t,
		state,
		remainingChain,
		toDeploy,
		true, // testRouterEnabled
	)

	// At this point we can send messages from the test router to the new chain.
	// These won't be processed yet because the offRamp is not aware of these new sources.
	// TODO: do we want to send a request before enabling on the offRamp? Don't think this would
	// be a practical use case.

	// UpdateOffRampSourcesConfig called on the new chain to enable existing sources. Also with the test router.
	// UpdateRouterRampsConfig to enable the existing sources on the new test router.
	// This means we can send messages from toDeploy to remainingChain.
	// NOTE: not using MCMS since haven't transferred to timelock yet.
	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, e.TimelockContracts(t), []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(UpdateOffRampSources),
			Config: UpdateOffRampSourcesConfig{
				UpdatesByChain: offRampSourceUpdates(t, remainingChain, toDeploy, true),
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(UpdateRouterRamps),
			Config: UpdateRouterRampsConfig{
				TestRouter:     true,
				UpdatesByChain: routerOffRampUpdates(t, remainingChain, toDeploy),
			},
		},
	})
	require.NoError(t, err)

	assertNewChainWiring(
		t,
		state,
		remainingChain,
		toDeploy,
		true, // testRouterEnabled
	)

	// Send messages from toDeploy to the newly added chain thru the test router.
	sendMsgs(toDeploy, []uint64{remainingChain}, true)

	// Now we switch to testing outbound from the new chain.
	// This amounts to enabling a new lane on the existing OCR instances
	// (assuming by default we want to enable all chains).

	// UpdateOnRampDestsConfig on the new chain with TestRouter=true and MCMS config enabled.
	// UpdateFeeQuoterPrices on the new chain to enable quoting for all existing chains
	// UpdateFeeQuoterDestsConfig on the new chain to set quoting params for the existing chains.
	// UpdateRouterRampsConfig with TestRouter=true to set the new chains onramp for all existing chain destinations
	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, e.TimelockContracts(t), []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(UpdateOnRampsDests),
			Config: UpdateOnRampDestsConfig{
				UpdatesByChain: onRampDestUpdates(t, toDeploy, []uint64{remainingChain}, true),
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(UpdateFeeQuoterPricesCS),
			Config: UpdateFeeQuoterPricesConfig{
				PricesByChain: feeQuoterPricesByChain(t, toDeploy, []uint64{remainingChain}),
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(UpdateFeeQuoterDests),
			Config: UpdateFeeQuoterDestsConfig{
				UpdatesByChain: feeQuoterDestUpdates(t, toDeploy, []uint64{remainingChain}),
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(UpdateRouterRamps),
			Config: UpdateRouterRampsConfig{
				TestRouter:     true,
				UpdatesByChain: routerOnRampUpdates(t, toDeploy, []uint64{remainingChain}),
			},
		},
	})
	require.NoError(t, err)

	// Send messages from remainingChain to toDeploy thru the test router.
	// TODO: not working due to discovery bug.
	// sendMsgs([]uint64{remainingChain}, toDeploy, true)
}

func assertNewChainWiring(
	t *testing.T,
	state CCIPOnChainState,
	newChain uint64,
	existingChains []uint64,
	testRouterEnabled bool,
) {
	for _, existingChain := range existingChains {
		var rtr *router.Router
		if testRouterEnabled {
			rtr = state.Chains[newChain].TestRouter
		} else {
			rtr = state.Chains[newChain].Router
		}

		// check that the onRamp has the new chain enabled as a dest.
		dcc, err := state.Chains[newChain].OffRamp.GetSourceChainConfig(&bind.CallOpts{
			Context: tests.Context(t),
		}, existingChain)
		require.NoError(t, err)
		require.Equal(t, rtr.Address(), dcc.Router)

		// check that the router has the new chain enabled as a dest.
		routerOffRamps, err := rtr.GetOffRamps(&bind.CallOpts{
			Context: tests.Context(t),
		})
		require.NoError(t, err)

		var found bool
		for _, offRamp := range routerOffRamps {
			if offRamp.SourceChainSelector == existingChain {
				require.Equal(t, state.Chains[newChain].OffRamp.Address(), offRamp.OffRamp)
				found = true
				break
			}
		}
		require.True(t, found)
	}
}

// assertExistingChainsWiring asserts that the following changes are applied correctly on all existing chains:
// UpdateOnRampDestsConfig on the existing chains with TestRouter=true and MCMS config enabled.
// UpdateFeeQuoterPrices to set initial FQ prices to the new destination.
// UpdateFeeQuoterDestsConfig to enable quoting for the new destination
// UpdateRouterRampsConfig on the test router to enable the new destination
func assertExistingChainsWiring(
	t *testing.T,
	state CCIPOnChainState,
	newChain uint64,
	existingChains []uint64,
	testRouterEnabled bool,
) {
	for _, existingChain := range existingChains {
		var rtr *router.Router
		if testRouterEnabled {
			rtr = state.Chains[existingChain].TestRouter
		} else {
			rtr = state.Chains[existingChain].Router
		}

		// check that the onRamp has the new chain enabled as a dest.
		dcc, err := state.Chains[existingChain].OnRamp.GetDestChainConfig(&bind.CallOpts{
			Context: tests.Context(t),
		}, newChain)
		require.NoError(t, err)
		require.Equal(t, rtr.Address(), dcc.Router)

		// check that the feeQuoter has the new chain enabled as a dest.
		fqdcc, err := state.Chains[existingChain].FeeQuoter.GetDestChainConfig(&bind.CallOpts{
			Context: tests.Context(t),
		}, newChain)
		require.NoError(t, err)
		require.Equal(t, true, fqdcc.IsEnabled)

		// check that the router has the new chain enabled as a dest.
		routerOnRamp, err := rtr.GetOnRamp(&bind.CallOpts{
			Context: tests.Context(t),
		}, newChain)
		require.NoError(t, err)
		require.Equal(t, state.Chains[existingChain].OnRamp.Address(), routerOnRamp)
	}
}

// routerOffRampUpdates adds the provided sources to the router of the provided dest chain.
func routerOffRampUpdates(t *testing.T, dest uint64, sources []uint64) (updates map[uint64]RouterUpdates) {
	updates = make(map[uint64]RouterUpdates)
	for _, source := range sources {
		require.NotEqual(t, source, dest)
		if _, ok := updates[dest]; ok {
			updates[dest].OffRampUpdates[source] = true
		} else {
			updates[dest] = RouterUpdates{
				OffRampUpdates: map[uint64]bool{
					source: true,
				},
			}
		}
	}
	return
}

// routerOnRampUpdates adds the sets each dest selector in the given dest chains slice on the router
// to point to the local onramp on each source chain.
func routerOnRampUpdates(t *testing.T, dests []uint64, sources []uint64) (updates map[uint64]RouterUpdates) {
	updates = make(map[uint64]RouterUpdates)
	for _, existing := range sources {
		for _, newChain := range dests {
			require.NotEqual(t, existing, newChain)
			if _, ok := updates[existing]; !ok {
				updates[existing] = RouterUpdates{
					OnRampUpdates: map[uint64]bool{
						newChain: true,
					},
				}
			} else {
				updates[existing].OnRampUpdates[newChain] = true
			}
		}
	}
	return
}

// feeQuoterDestUpdates adds a fee quoter configuration for the provided dest chains on the fee quoters on the provided sources.
func feeQuoterDestUpdates(t *testing.T, dests []uint64, sources []uint64) (updates map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig) {
	updates = make(map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig)
	for _, source := range sources {
		for _, dest := range dests {
			require.NotEqual(t, source, dest)
			if _, ok := updates[source]; !ok {
				updates[source] = make(map[uint64]fee_quoter.FeeQuoterDestChainConfig)
			}
			updates[source][dest] = DefaultFeeQuoterDestChainConfig()
		}
	}
	return
}

// feeQuoterPricesByChain sets the gas price for the provided dests on the fee quoters in the provided sources.
func feeQuoterPricesByChain(t *testing.T, dests []uint64, sources []uint64) (prices map[uint64]FeeQuoterPriceUpdatePerSource) {
	prices = make(map[uint64]FeeQuoterPriceUpdatePerSource)
	for _, source := range sources {
		prices[source] = FeeQuoterPriceUpdatePerSource{
			GasPrices: make(map[uint64]*big.Int),
		}
		for _, dest := range dests {
			require.NotEqual(t, source, dest)
			prices[source].GasPrices[dest] = DefaultGasPrice
		}
	}
	return
}

// onRampDestUpdates adds the provided dests as destination chains to the onRamps on the provided sources.
func onRampDestUpdates(t *testing.T, dests []uint64, sources []uint64, testRouterEnabled bool) (updates map[uint64]map[uint64]OnRampDestinationUpdate) {
	updates = make(map[uint64]map[uint64]OnRampDestinationUpdate)
	for _, source := range sources {
		for _, dest := range dests {
			require.NotEqual(t, source, dest)
			if _, ok := updates[source]; !ok {
				updates[source] = map[uint64]OnRampDestinationUpdate{
					dest: {
						IsEnabled:  true,
						TestRouter: testRouterEnabled,
					},
				}
			} else {
				updates[source][dest] = OnRampDestinationUpdate{
					IsEnabled:  true,
					TestRouter: testRouterEnabled,
				}
			}
		}
	}
	return
}

// offRampSourceUpdates adds the provided sources to the offRamp on the provided dest chain.
func offRampSourceUpdates(t *testing.T, dest uint64, sources []uint64, testRouterEnabled bool) (updates map[uint64]map[uint64]OffRampSourceUpdate) {
	updates = make(map[uint64]map[uint64]OffRampSourceUpdate)
	for _, source := range sources {
		require.NotEqual(t, source, dest)
		if _, ok := updates[dest]; !ok {
			updates[dest] = make(map[uint64]OffRampSourceUpdate)
		}
		updates[dest][source] = OffRampSourceUpdate{
			IsEnabled:  true,
			TestRouter: testRouterEnabled,
		}
	}
	return
}

func assertRMNRemoteAndProxyState(t *testing.T, chains []uint64, state CCIPOnChainState) {
	for _, chain := range chains {
		require.NotEqual(t, common.Address{}, state.Chains[chain].RMNRemote.Address())
		_, err := state.Chains[chain].RMNRemote.GetCursedSubjects(&bind.CallOpts{
			Context: tests.Context(t),
		})
		require.NoError(t, err)

		// check which address RMNProxy is pointing to
		rmnAddress, err := state.Chains[chain].RMNProxy.GetARM(&bind.CallOpts{
			Context: tests.Context(t),
		})
		require.NoError(t, err)
		require.Equal(t, state.Chains[chain].RMNRemote.Address(), rmnAddress)

		t.Log("RMNRemote address for chain", chain, "is:", state.Chains[chain].RMNRemote.Address().Hex())
		t.Log("RMNProxy address for chain", chain, "is:", state.Chains[chain].RMNProxy.Address().Hex())
	}
}

func transferToMCMSAndRenounceTimelockDeployer(
	t *testing.T,
	e DeployedEnv,
	chains []uint64,
	state CCIPOnChainState,
) {
	var apps []changesetcommon.ChangesetApplication
	apps = append(apps, changesetcommon.ChangesetApplication{
		Changeset: changesetcommon.WrapChangeSet(changesetcommon.TransferToMCMSWithTimelock),
		Config:    genTestTransferOwnershipConfig(e, chains, state),
	})
	for _, chain := range chains {
		apps = append(apps, changesetcommon.ChangesetApplication{
			Changeset: changesetcommon.WrapChangeSet(changesetcommon.RenounceTimelockDeployer),
			Config: changesetcommon.RenounceTimelockDeployerConfig{
				ChainSel: chain,
			},
		})
	}
	var err error
	e.Env, err = changesetcommon.ApplyChangesets(t, e.Env, e.TimelockContracts(t), apps)
	require.NoError(t, err)
}
