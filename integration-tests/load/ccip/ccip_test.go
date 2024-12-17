package ccip

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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
	wg sync.WaitGroup
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

	cribEnv := crib.NewDevspaceEnvFromStateDir(CRIB_DIRECTORY)

	cribDeployOutput := cribEnv.GetConfig()
	env, err := crib.NewDeployEnvironmentFromCribOutput(lggr, cribDeployOutput)
	require.NoError(t, err)
	require.NotNil(t, env)

	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	state, err := ccipchangeset.LoadOnchainState(*env)
	require.NoError(t, err)

	// Parse all events from the simulated chains, send to Loki
	loki, err := wasp.NewLokiClient(wasp.NewLokiConfig(config.CCIP.Load.LokiEndpoint, nil, nil, nil))
	require.NoError(t, err)
	defer loki.StopNow()

	// Based on the config, initialize DestinationGun
	p := wasp.NewProfile()
	for selector, chain := range env.Chains {
		latesthdr, err := chain.Client.HeaderByNumber(ctx, nil)
		require.NoError(t, err)
		block := latesthdr.Number.Uint64()
		startBlocks[selector] = &block

		p.Add(wasp.NewGenerator(&wasp.Config{
			T:           t,
			GenName:     "ccipLoad",
			LoadType:    wasp.RPS,
			CallTimeout: 20 * time.Second,
			Schedule:    wasp.Plain(1, 20*time.Second),
			// will need to be divided by number of chains
			// this schedule is per generator
			// in this example, it would be 1 request per 10seconds per generator (dest chain)
			// so if there are 3 generators, it would be 3 requests per 10 seconds over the network
			Gun:        NewDestinationGun(env.Logger, selector, *env, state.Chains[selector].Receiver.Address(), loki),
			Labels:     CommonTestLabels,
			LokiConfig: wasp.NewLokiConfig(config.CCIP.Load.LokiEndpoint, nil, nil, nil),
			// use the same loki client using `NewLokiClient` with the same config for sending events
		}))
	}

	_, err = p.Run(true)

	lokiLabels := map[string]string{}
	for chainSelector, startBlock := range startBlocks {
		wg.Add(1)
		go func(chainSelector uint64, startBlock *uint64) {
			defer wg.Done()
			filterOpts := &bind.FilterOpts{
				Start:   *startBlock,
				End:     nil, // To the latest block
				Context: ctx,
			}

			offRamp := state.Chains[chainSelector].OffRamp
			// Filter CommitReportAccepted events
			commitIterator, err := offRamp.FilterCommitReportAccepted(filterOpts)
			require.NoError(t, err)

			for commitIterator.Next() {
				event := commitIterator.Event
				fmt.Printf("CommitReportAccepted event: %+v\n", event)

				blockNum := commitIterator.Event.Raw.BlockNumber
				header, err := env.Chains[chainSelector].Client.HeaderByNumber(ctx, big.NewInt(int64(blockNum)))
				require.NoError(t, err)
				timestamp := time.Unix(int64(header.Time), 0)

				for _, root := range event.MerkleRoots {
					lokiLabels = setLokiLabels(root.SourceChainSelector, chainSelector)

					for i := root.MinSeqNr; i <= root.MaxSeqNr; i++ {
						// todo: push loki calls to channel?
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
				lokiLabels = setLokiLabels(execIterator.Event.SourceChainSelector, chainSelector)
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
