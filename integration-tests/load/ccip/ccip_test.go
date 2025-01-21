package ccip

import (
	"sync"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"

	"github.com/stretchr/testify/require"

	"context"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	ccipchangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	crib "github.com/smartcontractkit/chainlink/deployment/environment/crib"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

var (
	CommonTestLabels = map[string]string{
		"branch": "ccip_load_crib",
		"commit": "ccip_load_crib",
	}
	wg sync.WaitGroup
	// todo: add multiple keys and rotate them when sending messages
	SIM_CHAIN_TEST_KEY = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
)

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

	// get user defined configurations
	config, err := tc.GetConfig([]string{"Load"}, tc.CCIP)
	require.NoError(t, err)
	lggr.Infof("loaded ccip test config: %+v", config.CCIP.Load)
	userOverrides := config.CCIP.Load

	// apply user defined test duration
	timeout, err := time.ParseDuration(*userOverrides.TestDuration)
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// generate environment from crib-produced files
	cribEnv := crib.NewDevspaceEnvFromStateDir(*userOverrides.CribEnvDirectory)
	cribDeployOutput, err := cribEnv.GetConfig(SIM_CHAIN_TEST_KEY)
	require.NoError(t, err)
	env, err := crib.NewDeployEnvironmentFromCribOutput(lggr, cribDeployOutput)
	require.NoError(t, err)
	require.NotNil(t, env)

	// initialize loki using endpoint from user defined env vars
	loki, err := wasp.NewLokiClient(wasp.NewLokiConfig(userOverrides.LokiEndpoint, nil, nil, nil))
	require.NoError(t, err)
	defer loki.StopNow()

	// Keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	state, err := ccipchangeset.LoadOnchainState(*env)
	require.NoError(t, err)

	errChan := make(chan error)
	finalSeqNrCommitChannels := make(map[uint64]chan finalSeqNrReport)
	finalSeqNrExecChannels := make(map[uint64]chan finalSeqNrReport)

	// gunMap holds a destinationGun for every enabled destination chain
	gunMap := make(map[uint64]*DestinationGun)
	p := wasp.NewProfile()
	// Only create a destination gun if we have decided to send traffic to this chain
	for _, cs := range *userOverrides.EnabledDestionationChains {
		latesthdr, err := env.Chains[cs].Client.HeaderByNumber(ctx, nil)
		require.NoError(t, err)
		block := latesthdr.Number.Uint64()
		startBlocks[cs] = &block

		gunMap[cs], err = NewDestinationGun(env.Logger, cs, *env, state.Chains[cs].Receiver.Address(), userOverrides, loki)
		if err != nil {
			lggr.Errorw("Failed to initialize DestinationGun for", "chainSelector", cs, "error", err)
			t.Fatal(err)
		}

		otherChains := env.AllChainSelectorsExcluding([]uint64{cs})
		finalSeqNrCommitChannels[cs] = make(chan finalSeqNrReport)
		finalSeqNrExecChannels[cs] = make(chan finalSeqNrReport)

		wg.Add(1)
		go subscribeDeferredCommitEvents(
			ctx,
			lggr,
			state.Chains[cs].OffRamp,
			otherChains,
			&block,
			loki,
			cs,
			env.Chains[cs].Client,
			finalSeqNrCommitChannels[cs],
			errChan,
			&wg)

		wg.Add(1)
		go subscribeExecutionEvents(
			ctx,
			lggr,
			state.Chains[cs].OffRamp,
			otherChains,
			&block,
			loki,
			cs,
			env.Chains[cs].Client,
			finalSeqNrExecChannels[cs],
			errChan,
			&wg)
	}

	loadDuration, err := time.ParseDuration(*userOverrides.LoadDuration)
	require.NoError(t, err)
	requestFrequency, err := time.ParseDuration(*userOverrides.RequestFrequency)
	require.NoError(t, err)

	for _, gun := range gunMap {
		p.Add(wasp.NewGenerator(&wasp.Config{
			T:           t,
			GenName:     "ccipLoad",
			LoadType:    wasp.RPS,
			CallTimeout: 5 * time.Second,
			// 1 request per second for n seconds
			Schedule: wasp.Plain(1, loadDuration),
			// limit requests to 1 per duration
			RateLimitUnitDuration: requestFrequency,
			// will need to be divided by number of chains
			// this schedule is per generator
			// in this example, it would be 1 request per 5seconds per generator (dest chain)
			// so if there are 3 generators, it would be 3 requests per 5 seconds over the network
			Gun:        gun,
			Labels:     CommonTestLabels,
			LokiConfig: wasp.NewLokiConfig(config.CCIP.Load.LokiEndpoint, nil, nil, nil),
			// use the same loki client using `NewLokiClient` with the same config for sending events
		}))
	}

	_, err = p.Run(true)

	// wait for offchain to complete handling load fully
	execExpectedSeqNums := make(map[testhelpers.SourceDestPair][]uint64)
	commitExepectedSeqNums := make(map[testhelpers.SourceDestPair]uint64)
	for destSel, gun := range gunMap {
		for csPair, seqNums := range gun.seqNums {
			lggr.Debugw("pushing finalized sequence numbers for ",
				"chainSelector", gun.chainSelector,
				"sourceChainSelector", csPair.SourceChainSelector,
				"seqNums", seqNums)
			finalSeqNrCommitChannels[csPair.DestChainSelector] <- finalSeqNrReport{
				sourceChainSelector: csPair.SourceChainSelector,
				expectedSeqNrRange: ccipocr3.SeqNumRange{
					ccipocr3.SeqNum(seqNums.Start.Load()), ccipocr3.SeqNum(seqNums.End.Load()),
				},
			}

			finalSeqNrExecChannels[csPair.DestChainSelector] <- finalSeqNrReport{
				sourceChainSelector: csPair.SourceChainSelector,
				expectedSeqNrRange: ccipocr3.SeqNumRange{
					ccipocr3.SeqNum(seqNums.Start.Load()), ccipocr3.SeqNum(seqNums.End.Load()),
				},
			}
		}

		defer close(finalSeqNrCommitChannels[destSel])
		defer close(finalSeqNrExecChannels[destSel])
	}

	testhelpers.ConfirmCommitForAllWithExpectedSeqNums(t, *env, state, commitExepectedSeqNums, startBlocks)
	testhelpers.ConfirmExecWithSeqNrsForAll(t, *env, state, execExpectedSeqNums, startBlocks)

	// wait for offchain to complete handling load fully
	//execExpectedSeqNums := make(map[ccipchangeset.SourceDestPair][]uint64)
	//commitExepectedSeqNums := make(map[ccipchangeset.SourceDestPair]uint64)
	//for _, gun := range gunMap {
	//	for csPair := range gun.seqNums {
	//		commitExepectedSeqNums[csPair] = gun.seqNums[csPair].End.Load()
	//		for i := gun.seqNums[csPair].Start.Load(); i <= gun.seqNums[csPair].End.Load(); i++ {
	//			execExpectedSeqNums[csPair] = append(execExpectedSeqNums[csPair], i)
	//		}
	//	}
	//}

	//ccipchangeset.ConfirmCommitForAllWithExpectedSeqNums(t, *env, state, commitExepectedSeqNums, startBlocks)
	//ccipchangeset.ConfirmExecWithSeqNrsForAll(t, *env, state, execExpectedSeqNums, startBlocks)

	// todo: create channels that subscribe to these events beforehand using WatchExecutionStateChanged and WatchCommitReportAccepted
	//lokiLabels := map[string]string{}
	//for chainSelector, startBlock := range startBlocks {
	//	filterOpts := &bind.FilterOpts{
	//		Start:   *startBlock,
	//		End:     nil, // To the latest block
	//		Context: ctx,
	//	}

	//wg.Add(1)
	//go func(chainSelector uint64, startBlock *uint64, filterOpts *bind.FilterOpts) {
	//	defer wg.Done()
	//	lggr.Infow("Starting to query for events on ", "chainSelector", chainSelector, "startblock", startBlock)
	//	latesthdr, err := env.Chains[chainSelector].Client.HeaderByNumber(ctx, nil)
	//	require.NoError(t, err)
	//	lggr.Infow("Current block number", "chainSelector", chainSelector, "block", latesthdr.Number.Uint64())
	//
	//	offRamp := state.Chains[chainSelector].OffRamp
	//	// Filter CommitReportAccepted events
	//	commitIterator, err := offRamp.FilterCommitReportAccepted(filterOpts)
	//	require.NoError(t, err)
	//
	//	for commitIterator.Next() {
	//		event := commitIterator.Event
	//
	//		blockNum := commitIterator.Event.Raw.BlockNumber
	//		header, err := env.Chains[chainSelector].Client.HeaderByNumber(ctx, big.NewInt(int64(blockNum)))
	//		require.NoError(t, err)
	//		timestamp := time.Unix(int64(header.Time), 0)
	//
	//		for _, root := range event.MerkleRoots {
	//			lokiLabels, err = setLokiLabels(root.SourceChainSelector, chainSelector)
	//			require.NoError(t, err)
	//
	//			for i := root.MinSeqNr; i <= root.MaxSeqNr; i++ {
	//				SendMetricsToLoki(lggr, loki, lokiLabels, &LokiMetric{
	//					EventType:      committed,
	//					Timestamp:      timestamp,
	//					SequenceNumber: i,
	//				})
	//				lggr.Infow("pushed loki commit event for ", "seqNumber", i, "src", root.SourceChainSelector, "dest", chainSelector)
	//
	//			}
	//		}
	//	}
	//}(chainSelector, startBlock, filterOpts)

	//for sourceCS := range env.Chains {
	//	wg.Add(1)
	//	go func(srcSelector uint64, startBlock *uint64, filterOpts *bind.FilterOpts) {
	//		defer wg.Done()
	//		csPair := ccipchangeset.SourceDestPair{
	//			SourceChainSelector: srcSelector,
	//			DestChainSelector:   chainSelector,
	//		}
	//		// Filter ExecutionStateChanged events
	//		execIterator, err := state.Chains[chainSelector].OffRamp.FilterExecutionStateChanged(filterOpts, []uint64{srcSelector}, execExpectedSeqNums[csPair], nil)
	//		require.NoError(t, err)
	//
	//		for execIterator.Next() {
	//			blockNum := execIterator.Event.Raw.BlockNumber
	//			header, err := env.Chains[chainSelector].Client.HeaderByNumber(ctx, big.NewInt(int64(blockNum)))
	//			require.NoError(t, err)
	//			timestamp := time.Unix(int64(header.Time), 0)
	//
	//			lokiLabels, err = setLokiLabels(execIterator.Event.SourceChainSelector, chainSelector)
	//			require.NoError(t, err)
	//
	//			SendMetricsToLoki(lggr, loki, lokiLabels, &LokiMetric{
	//				EventType:      executed,
	//				Timestamp:      timestamp,
	//				GasUsed:        execIterator.Event.GasUsed.Uint64(),
	//				SequenceNumber: execIterator.Event.SequenceNumber,
	//			})
	//			lggr.Infow("pushed loki exec event for ", "seqNumber", execIterator.Event.SequenceNumber, "src", execIterator.Event.SourceChainSelector, "dest", chainSelector)
	//		}
	//	}(sourceCS, startBlock, filterOpts)
	//}
	//}

	wg.Wait()
	lggr.Infow("finished wait group")
	for err := range errChan {
		require.NoError(t, err)
	}
}
