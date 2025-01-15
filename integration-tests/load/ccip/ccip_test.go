package ccip

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	ccipchangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	crib "github.com/smartcontractkit/chainlink/deployment/environment/crib"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/stretchr/testify/require"
	"math/big"
	"sync"
	"testing"
	"time"
)

var (
	CommonTestLabels = map[string]string{
		"branch": "ccip_load_crib",
		"commit": "ccip_load_crib",
	}
	wg                 sync.WaitGroup
	SIM_CHAIN_TEST_KEY = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
)

const CRIB_DIRECTORY = "/Users/austin.wang/ccip-core/repos/crib/deployments/ccip-v2/.tmp"

// step 1: setup
// Parse the test config, initialize CRIB with configurations defined
// step 2: load
// Use wasp to initiate load
// step 3: parse logs
// Parse all events from the simulated chains, send to Loki
// step 4: teardown
// Stop the chains, cleanup the environment
func TestCCIPLoad_RPS(t *testing.T) {
	ctx := tests.Context(t)
	lggr := logger.Test(t)

	config, err := tc.GetConfig([]string{"Load"}, tc.CCIP)
	require.NoError(t, err)
	lggr.Infof("loaded ccip test config: %+v", config.CCIP.Load)
	userOverrides := config.CCIP.Load

	timeout, err := time.ParseDuration(*userOverrides.TestTimeout)
	if err != nil {
		require.NoError(t, err)
	}
	t.Cleanup(func() {
		time.AfterFunc(timeout, func() {
			t.Fatalf("Test passed timeout after %v", timeout)
		})
	})

	cribEnv := crib.NewDevspaceEnvFromStateDir(CRIB_DIRECTORY)

	cribDeployOutput, err := cribEnv.GetConfig(SIM_CHAIN_TEST_KEY)
	require.NoError(t, err)
	env, err := crib.NewDeployEnvironmentFromCribOutput(lggr, cribDeployOutput)
	require.NoError(t, err)
	require.NotNil(t, env)

	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	state, err := ccipchangeset.LoadOnchainState(*env)
	require.NoError(t, err)

	// Parse all events from the simulated chains, send to Loki
	loki, err := wasp.NewLokiClient(wasp.NewLokiConfig(userOverrides.LokiEndpoint, nil, nil, nil))
	require.NoError(t, err)
	defer loki.StopNow()

	gunMap := make(map[uint64]*DestinationGun)

	// Based on the config, initialize DestinationGun
	p := wasp.NewProfile()
	for selector, chain := range env.Chains {
		latesthdr, err := chain.Client.HeaderByNumber(ctx, nil)
		require.NoError(t, err)
		block := latesthdr.Number.Uint64()
		startBlocks[selector] = &block

		// Only create a destination gun if we have decided to send traffic to this chain
		for _, cs := range *userOverrides.EnabledDestionationChains {
			if cs == selector {
				gunMap[selector], err = NewDestinationGun(env.Logger, selector, *env, state.Chains[selector].Receiver.Address(), userOverrides, loki)
				if err != nil {
					lggr.Errorw("Failed to initialize DestinationGun for", "chainSelector", chain, "error", err)
					t.Fatal(err)
				}
			}
		}
	}

	for _, gun := range gunMap {
		p.Add(wasp.NewGenerator(&wasp.Config{
			T:           t,
			GenName:     "ccipLoad",
			LoadType:    wasp.RPS,
			CallTimeout: 5 * time.Second,
			Schedule:    wasp.Plain(1, time.Duration(*userOverrides.SecondsPerRequestPerDestination)*time.Second),
			// will need to be divided by number of chains
			// this schedule is per generator
			// in this example, it would be 1 request per 10seconds per generator (dest chain)
			// so if there are 3 generators, it would be 3 requests per 10 seconds over the network
			Gun:        gun,
			Labels:     CommonTestLabels,
			LokiConfig: wasp.NewLokiConfig(config.CCIP.Load.LokiEndpoint, nil, nil, nil),
			// use the same loki client using `NewLokiClient` with the same config for sending events
		}))
	}

	// find get fee revert
	csPair := ccipchangeset.SourceDestPair{
		SourceChainSelector: 12922642891491394802,
		DestChainSelector:   3379446385462418246,
	}
	res, err := state.Chains[csPair.SourceChainSelector].Router.IsChainSupported(nil, csPair.DestChainSelector)
	lggr.Infow("IsChainSupported", "res", res, "err", err)

	destChainConfig, err := state.Chains[csPair.SourceChainSelector].FeeQuoter.GetDestChainConfig(nil, csPair.DestChainSelector)
	lggr.Infow("GetDestChainConfig", "destChainConfig", destChainConfig, "err", err)

	// find the getFee revert
	_, err = p.Run(true)

	src, dst := env.Chains[csPair.SourceChainSelector], env.Chains[csPair.DestChainSelector]
	startblk := uint64(11654)

	seqNum := gunMap[csPair.DestChainSelector].seqNums[csPair].End.Load()
	_, err = ccipchangeset.ConfirmCommitWithExpectedSeqNumRange(t, src, dst, state.Chains[3379446385462418246].OffRamp, &startblk, cciptypes.SeqNumRange{
		cciptypes.SeqNum(seqNum - 1),
		cciptypes.SeqNum(seqNum - 1),
	}, false)

	ccipchangeset.ConfirmExecWithSeqNrsForAll(t, *env, state, map[ccipchangeset.SourceDestPair][]uint64{
		csPair: {seqNum - 1},
	}, startBlocks)

	// todo: create channels that watch for these events beforehand using WatchExecutionStateChanged and WatchCommitReportAccepted
	// rather than waiting for the generator to finish
	lokiLabels := map[string]string{}
	for chainSelector, startBlock := range startBlocks {
		wg.Add(1)
		go func(chainSelector uint64, startBlock *uint64) {
			defer wg.Done()
			lggr.Infow("Starting to query for events on ", "chainSelector", chainSelector, "startblock", startBlock)
			latesthdr, err := env.Chains[chainSelector].Client.HeaderByNumber(ctx, nil)
			require.NoError(t, err)
			lggr.Infow("Current block number", "chainSelector", chainSelector, "block", latesthdr.Number.Uint64())

			filterOpts := &bind.FilterOpts{
				Start:   *startBlock,
				End:     nil, // To the latest block
				Context: ctx,
			}

			offRamp := state.Chains[chainSelector].OffRamp
			// Filter CommitReportAccepted events
			commitIterator, err := offRamp.FilterCommitReportAccepted(filterOpts)
			require.NoError(t, err)

			fmt.Printf("Events on commitIterator %+v", commitIterator)

			for commitIterator.Next() {
				fmt.Printf("Events on commitIterator inside %+v", commitIterator)

				event := commitIterator.Event
				fmt.Printf("CommitReportAccepted event: %+v\n", event)

				blockNum := commitIterator.Event.Raw.BlockNumber
				header, err := env.Chains[chainSelector].Client.HeaderByNumber(ctx, big.NewInt(int64(blockNum)))
				require.NoError(t, err)
				timestamp := time.Unix(int64(header.Time), 0)

				for _, root := range event.MerkleRoots {
					lokiLabels, err = setLokiLabels(root.SourceChainSelector, chainSelector)
					require.NoError(t, err)

					for i := root.MinSeqNr; i <= root.MaxSeqNr; i++ {
						// todo: push loki calls to channel?
						fmt.Printf("Pushed loki for seqNumber %d\n", i)

						SendMetricsToLoki(lggr, loki, lokiLabels, &LokiMetric{
							EventType:      committed,
							Timestamp:      timestamp,
							SequenceNumber: i,
						})
					}
				}
			}

			// Filter ExecutionStateChanged events
			execIterator, err := state.Chains[chainSelector].OffRamp.FilterExecutionStateChanged(filterOpts, []uint64{chainSelector}, nil, nil)
			require.NoError(t, err)
			for execIterator.Next() {
				event := execIterator.Event
				fmt.Printf("ExecutionStateChanged event: %+v\n", event)

				blockNum := execIterator.Event.Raw.BlockNumber
				header, err := env.Chains[chainSelector].Client.HeaderByNumber(ctx, big.NewInt(int64(blockNum)))
				require.NoError(t, err)
				timestamp := time.Unix(int64(header.Time), 0)

				// todo: push loki calls to channel?
				lokiLabels, err = setLokiLabels(execIterator.Event.SourceChainSelector, chainSelector)
				require.NoError(t, err)

				fmt.Printf("Pushed loki exec for seqNumber %d\n", execIterator.Event.SequenceNumber)
				SendMetricsToLoki(lggr, loki, lokiLabels, &LokiMetric{
					EventType:      executed,
					Timestamp:      timestamp,
					GasUsed:        execIterator.Event.GasUsed.Uint64(),
					SequenceNumber: execIterator.Event.SequenceNumber,
				})

			}
		}(chainSelector, startBlock)
	}

	wg.Wait()

	// Stop the chains, cleanup the environment

	// crib.StopChains(env)
	// crib.StopNodes(env)
}
