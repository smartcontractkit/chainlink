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

	/////////////////////////////////////
	// START Setup initial chains
	/////////////////////////////////////
	e = setupChain(
		t,
		e,
		tEnv,
		toDeploy,
		true,  // deployJobs
		true,  // deployHomeChain
		false, // mcmsEnabled
	)

	state, err := LoadOnchainState(e.Env)
	require.NoError(t, err)

	// check RMNRemote is up and RMNProxy is correctly wired.
	assertRMNRemoteAndProxyState(t, toDeploy, state)

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
	transferToMCMSAndRenounceTimelockDeployer(
		t,
		e,
		toDeploy,
		state,
		false, // onlyChainContracts
	)

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

	/////////////////////////////////////
	// END Setup initial chains
	/////////////////////////////////////

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

	/////////////////////////////////////
	// START Deploy to the remaining chain.
	/////////////////////////////////////

	// MCMS needs to be enabled because the home chain contracts have been
	// transferred to MCMS.
	e = setupChain(
		t,
		e,
		tEnv,
		[]uint64{remainingChain},
		false, // deployJobs
		false, // deployHomeChain
		true,  // mcmsEnabled
	)

	state, err = LoadOnchainState(e.Env)
	require.NoError(t, err)

	assertRMNRemoteAndProxyState(t, []uint64{remainingChain}, state)

	// TODO: wait for gas price of new chain to be updated on all other chains.

	/////////////////////////////////////
	// END Deploy to the remaining chain.
	/////////////////////////////////////

	/////////////////////////////////////
	// START Wire up toDeploy -> remainingChain outbound.
	/////////////////////////////////////
	e = setupOutboundWiring(
		t,
		e,
		toDeploy,
		[]uint64{remainingChain},
		true, // testRouterEnabled
		true, // mcmsEnabled
	)

	state, err = LoadOnchainState(e.Env)
	require.NoError(t, err)

	assertChainWiringOutbound(
		t,
		state,
		remainingChain,
		toDeploy,
		true, // testRouterEnabled
	)

	// At this point we can send messages from the test router to the new chain.
	// These won't be processed yet because the offRamp is not aware of these new sources.

	/////////////////////////////////////
	// END Wire up toDeploy -> remainingChain outbound.
	/////////////////////////////////////

	/////////////////////////////////////
	// START Wire up toDeploy -> remainingChain inbound on remainingChain.
	/////////////////////////////////////

	// UpdateOffRampSourcesConfig called on the new chain to enable existing sources. Also with the test router.
	// UpdateRouterRampsConfig to enable the existing sources on the new test router.
	// This means we can send messages from toDeploy to remainingChain.
	// NOTE: not using MCMS since haven't transferred to timelock yet.
	e = setupInboundWiring(t, e, toDeploy, []uint64{remainingChain}, true, false)

	assertChainWiringInbound(
		t,
		state,
		remainingChain,
		toDeploy,
		true, // testRouterEnabled
	)

	// Send messages from toDeploy to the newly added chain thru the test router.
	sendMsgs(toDeploy, []uint64{remainingChain}, true)

	/////////////////////////////////////
	// END Wire up toDeploy -> remainingChain inbound on remainingChain.
	/////////////////////////////////////

	/////////////////////////////////////
	// START Wire up remainingChain -> toDeploy outbound.
	/////////////////////////////////////

	// Now we switch to testing outbound from the new chain.
	// This amounts to enabling a new lane on the existing OCR instances
	// (assuming by default we want to enable all chains).
	e = setupOutboundWiring(
		t,
		e,
		[]uint64{remainingChain},
		toDeploy,
		true,  // testRouterEnabled
		false, // mcmsEnabled
	)

	// sanity check that everything is correctly set up.
	for _, chain := range toDeploy {
		assertChainWiringOutbound(t, state, chain, []uint64{remainingChain}, true)
	}

	/////////////////////////////////////
	// END Wire up remainingChain -> toDeploy outbound.
	/////////////////////////////////////

	/////////////////////////////////////
	// START Wire up remainingChain -> toDeploy inbound on toDeploy.
	/////////////////////////////////////

	e = setupInboundWiring(t, e, []uint64{remainingChain}, toDeploy, true, true)

	for _, chain := range toDeploy {
		assertChainWiringInbound(t, state, chain, []uint64{remainingChain}, true)
	}

	/////////////////////////////////////
	// END Wire up remainingChain -> toDeploy inbound on toDeploy.
	/////////////////////////////////////

	/////////////////////////////////////
	// START send messages from remainingChain to toDeploy.
	/////////////////////////////////////

	// Send messages from remainingChain to toDeploy thru the test router.
	sendMsgs([]uint64{remainingChain}, toDeploy, true)

	/////////////////////////////////////
	// END send messages from remainingChain to toDeploy.
	/////////////////////////////////////

	// Transfer the new chain contracts to the timelock ownership.
	transferToMCMSAndRenounceTimelockDeployer(
		t,
		e,
		[]uint64{remainingChain},
		state,
		true, // onlyChainContracts
	)

	// Once verified with the test routers, the last step is to whitelist the new chain on the real routers everywhere with
	// UpdateRouterRamps, UpdateOffRampSources and UpdateOnRampDests with TestRouter=False.
	mcmsConfig := &MCMSConfig{
		MinDelay: 0,
	}
	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, e.TimelockContracts(t), []commonchangeset.ChangesetApplication{
		// OffRamp updates
		// TODO: combine into one changeset call w/ the appropriate config.
		{
			// This updates the offRamps in toDeploy to enable the new chain as a source through the real router.
			Changeset: commonchangeset.WrapChangeSet(UpdateOffRampSources),
			Config: UpdateOffRampSourcesConfig{
				UpdatesByChain: offRampSourceUpdates(t, toDeploy, []uint64{remainingChain}, false),
				MCMS:           mcmsConfig,
			},
		},
		{
			// This updates the offRamp in remainingChain to enable the toDeploy chains as sources through the real router.
			Changeset: commonchangeset.WrapChangeSet(UpdateOffRampSources),
			Config: UpdateOffRampSourcesConfig{
				UpdatesByChain: offRampSourceUpdates(t, []uint64{remainingChain}, toDeploy, false),
				MCMS:           mcmsConfig,
			},
		},

		// OnRamp updates
		// TODO: combine into one changeset call w/ the appropriate config.
		{
			// This updates the onRamps in toDeploy to enable the new chain as a dest through the real router.
			Changeset: commonchangeset.WrapChangeSet(UpdateOnRampsDests),
			Config: UpdateOnRampDestsConfig{
				UpdatesByChain: onRampDestUpdates(t, []uint64{remainingChain}, toDeploy, false),
				MCMS:           mcmsConfig,
			},
		},
		{
			// This updates the onRamp in remainingChain to enable the toDeploy chains as dests through the real router.
			Changeset: commonchangeset.WrapChangeSet(UpdateOnRampsDests),
			Config: UpdateOnRampDestsConfig{
				UpdatesByChain: onRampDestUpdates(t, toDeploy, []uint64{remainingChain}, false),
				MCMS:           mcmsConfig,
			},
		},

		// Router updates
		// TODO: combine into one changeset call w/ the appropriate config.
		{
			// This updates the real router on remainingChain to enable the toDeploy chains as destinations.
			Changeset: commonchangeset.WrapChangeSet(UpdateRouterRamps),
			Config: UpdateRouterRampsConfig{
				TestRouter:     false,
				UpdatesByChain: routerOnRampUpdates(t, toDeploy, []uint64{remainingChain}),
				MCMS:           mcmsConfig,
			},
		},
		{
			// This updates the real router on toDeploy chains to enable the new chain as a destination.
			Changeset: commonchangeset.WrapChangeSet(UpdateRouterRamps),
			Config: UpdateRouterRampsConfig{
				TestRouter:     false,
				UpdatesByChain: routerOnRampUpdates(t, []uint64{remainingChain}, toDeploy),
				MCMS:           mcmsConfig,
			},
		},
		{
			// This updates the offRamps in toDeploy to enable the new chain as a source through the real router.
			Changeset: commonchangeset.WrapChangeSet(UpdateRouterRamps),
			Config: UpdateRouterRampsConfig{
				TestRouter:     false,
				UpdatesByChain: routerOffRampUpdates(t, toDeploy, []uint64{remainingChain}),
				MCMS:           mcmsConfig,
			},
		},
		{
			// This updates the offRamps in remaining chain to enable the toDeploy chains sources through the real router.
			Changeset: commonchangeset.WrapChangeSet(UpdateRouterRamps),
			Config: UpdateRouterRampsConfig{
				TestRouter:     false,
				UpdatesByChain: routerOffRampUpdates(t, []uint64{remainingChain}, toDeploy),
				MCMS:           mcmsConfig,
			},
		},
	})
	require.NoError(t, err)

	// assert wiring is correct but with real router now.
	// this checks that the toDeploy chains can be requested from remainingChain.
	for _, chain := range toDeploy {
		assertChainWiringOutbound(t, state, chain, []uint64{remainingChain}, false)
	}
	// this checks that the remainingChain can be requested from the toDeploy chains.
	assertChainWiringOutbound(t, state, remainingChain, toDeploy, false)

	// this checks that the toDeploy chains have remainingChain enabled as as a source.
	for _, chain := range toDeploy {
		assertChainWiringInbound(t, state, chain, []uint64{remainingChain}, false)
	}
	// this checks that the remainingChain has the toDeploy chains enabled as sources.
	assertChainWiringInbound(t, state, remainingChain, toDeploy, false)

	// Finally, at this point, we can request messages from the real routers.
	// sendMsgs(toDeploy, []uint64{remainingChain}, false)
	// sendMsgs([]uint64{remainingChain}, toDeploy, false)
}

// setupInboundWiring sets up the newChain to be able to receive ccip messages from the provided sources.
// This only touches the newChain and does not touch the sources.
func setupInboundWiring(
	t *testing.T,
	e DeployedEnv,
	sources []uint64,
	newChains []uint64,
	testRouterEnabled,
	mcmsEnabled bool,
) DeployedEnv {
	var mcmsConfig *MCMSConfig
	if mcmsEnabled {
		mcmsConfig = &MCMSConfig{
			MinDelay: 0,
		}
	}

	var err error
	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, e.TimelockContracts(t), []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(UpdateOffRampSources),
			Config: UpdateOffRampSourcesConfig{
				UpdatesByChain: offRampSourceUpdates(t, newChains, sources, testRouterEnabled),
				MCMS:           mcmsConfig,
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(UpdateRouterRamps),
			Config: UpdateRouterRampsConfig{
				TestRouter:     testRouterEnabled,
				UpdatesByChain: routerOffRampUpdates(t, newChains, sources),
				MCMS:           mcmsConfig,
			},
		},
	})
	require.NoError(t, err)

	return e
}

// setupOutboundWiring sets up the given sources to be able to request ccip messages to the newChain.
// This only touches the sources and does not touch newChain.
// Therefore any requests to newChain will not be processed until the inbound wiring is set up.
func setupOutboundWiring(
	t *testing.T,
	e DeployedEnv,
	sources []uint64,
	newChains []uint64,
	testRouterEnabled,
	mcmsEnabled bool,
) DeployedEnv {
	var mcmsConfig *MCMSConfig
	if mcmsEnabled {
		mcmsConfig = &MCMSConfig{
			MinDelay: 0,
		}
	}

	var err error
	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, e.TimelockContracts(t), []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(UpdateOnRampsDests),
			Config: UpdateOnRampDestsConfig{
				UpdatesByChain: onRampDestUpdates(t, newChains, sources, true),
				MCMS:           mcmsConfig,
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(UpdateFeeQuoterPricesCS),
			Config: UpdateFeeQuoterPricesConfig{
				PricesByChain: feeQuoterPricesByChain(t, newChains, sources),
				MCMS:          mcmsConfig,
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(UpdateFeeQuoterDests),
			Config: UpdateFeeQuoterDestsConfig{
				UpdatesByChain: feeQuoterDestUpdates(t, newChains, sources),
				MCMS:           mcmsConfig,
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(UpdateRouterRamps),
			Config: UpdateRouterRampsConfig{
				TestRouter:     testRouterEnabled,
				UpdatesByChain: routerOnRampUpdates(t, newChains, sources),
				MCMS:           mcmsConfig,
			},
		},
	})
	require.NoError(t, err)

	return e
}

// setupChain will deploy the ccip chain contracts to the provided chains.
// Based on the flags provided, it will also deploy the jobs and home chain contracts.
// mcmsEnabled should be set to true if the home chain contracts have been transferred to MCMS.
func setupChain(t *testing.T, e DeployedEnv, tEnv TestEnvironment, chains []uint64, deployJobs, deployHomeChain, mcmsEnabled bool) DeployedEnv {
	e = AddCCIPContractsToEnvironment(
		t,
		chains,
		tEnv,
		deployJobs,
		deployHomeChain,
		mcmsEnabled,
	)

	// Need to update what the RMNProxy is pointing to, otherwise plugin will not work.
	var err error
	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, e.TimelockContracts(t), []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(SetRMNRemoteOnRMNProxy),
			Config: SetRMNRemoteOnRMNProxyConfig{
				ChainSelectors: chains,
			},
		},
	})
	require.NoError(t, err)

	return e
}

// assertChainWiringInbound checks that the new chain has the existing chains enabled as sources.
// It only checks the inbound wiring on the new chain.
// It doesn't check that the existing chains have the new chain enabled as a dest.
func assertChainWiringInbound(
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

		// check that the offRamp has the existing chain enabled as a source.
		dcc, err := state.Chains[newChain].OffRamp.GetSourceChainConfig(&bind.CallOpts{
			Context: tests.Context(t),
		}, existingChain)
		require.NoError(t, err)
		require.Equal(t, rtr.Address(), dcc.Router)

		// check that the router has the existing chain enabled as a source.
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

// assertChainWiringOutbound checks that the new chain can be requested from the existing chains.
// This only checks the outbound wiring on the existing chains.
// It doesn't check that the new chain can process the requests.
func assertChainWiringOutbound(
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
		require.True(t, fqdcc.IsEnabled)

		// check that the router has the new chain enabled as a dest.
		routerOnRamp, err := rtr.GetOnRamp(&bind.CallOpts{
			Context: tests.Context(t),
		}, newChain)
		require.NoError(t, err)
		require.Equal(t, state.Chains[existingChain].OnRamp.Address(), routerOnRamp)
	}
}

// routerOffRampUpdates adds the provided sources to the router of the provided dest chain.
func routerOffRampUpdates(t *testing.T, dests []uint64, sources []uint64) (updates map[uint64]RouterUpdates) {
	updates = make(map[uint64]RouterUpdates)
	for _, source := range sources {
		for _, dest := range dests {
			require.NotEqual(t, source, dest)
			if _, ok := updates[dest]; !ok {
				updates[dest] = RouterUpdates{
					OffRampUpdates: map[uint64]bool{
						source: true,
					},
				}
			} else {
				updates[dest].OffRampUpdates[source] = true
			}
		}
	}
	return
}

// routerOnRampUpdates sets each dest selector in the given dest chains slice on the router
// to point to the local onramp on each source chain.
func routerOnRampUpdates(t *testing.T, dests []uint64, sources []uint64) (updates map[uint64]RouterUpdates) {
	updates = make(map[uint64]RouterUpdates)
	for _, source := range sources {
		for _, dest := range dests {
			require.NotEqual(t, source, dest)
			if _, ok := updates[source]; !ok {
				updates[source] = RouterUpdates{
					OnRampUpdates: map[uint64]bool{
						dest: true,
					},
				}
			} else {
				updates[source].OnRampUpdates[dest] = true
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

// offRampSourceUpdates adds the provided sources to the offRamp on the provided dest chains.
func offRampSourceUpdates(t *testing.T, dests []uint64, sources []uint64, testRouterEnabled bool) (updates map[uint64]map[uint64]OffRampSourceUpdate) {
	updates = make(map[uint64]map[uint64]OffRampSourceUpdate)
	for _, source := range sources {
		for _, dest := range dests {
			require.NotEqual(t, source, dest)
			if _, ok := updates[dest]; !ok {
				updates[dest] = make(map[uint64]OffRampSourceUpdate)
			}
			updates[dest][source] = OffRampSourceUpdate{
				IsEnabled:  true,
				TestRouter: testRouterEnabled,
			}
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
	onlyChainContracts bool,
) {
	apps := make([]changesetcommon.ChangesetApplication, 0, len(chains)+1)
	cfg := genTestTransferOwnershipConfig(e, chains, state)
	if onlyChainContracts {
		// filter out the home chain contracts from e.HomeChainSel
		var homeChainContracts = map[common.Address]struct{}{
			state.Chains[e.HomeChainSel].CapabilityRegistry.Address(): {},
			state.Chains[e.HomeChainSel].CCIPHome.Address():           {},
			state.Chains[e.HomeChainSel].RMNHome.Address():            {},
		}
		var chainContracts []common.Address
		for _, contract := range cfg.ContractsByChain[e.HomeChainSel] {
			if _, ok := homeChainContracts[contract]; !ok {
				chainContracts = append(chainContracts, contract)
			}
		}
		cfg.ContractsByChain[e.HomeChainSel] = chainContracts
	}
	apps = append(apps, changesetcommon.ChangesetApplication{
		Changeset: changesetcommon.WrapChangeSet(changesetcommon.TransferToMCMSWithTimelock),
		Config:    cfg,
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
