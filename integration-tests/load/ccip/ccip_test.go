package ccip

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"github.com/smartcontractkit/chainlink/deployment"
	ccipdeployment "github.com/smartcontractkit/chainlink/deployment/ccip"
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
	wg          sync.WaitGroup
	abPath      = "/Users/austin.wang/ccip-core/repos/chainlink/integration-tests/load/ccip/testfiles/ccip-v2-scripts-address-book.json"
	nodeIdsPath = "/Users/austin.wang/ccip-core/repos/chainlink/integration-tests/load/ccip/testfiles/ccip-v2-scripts-node-details.json"
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
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig([]string{"Load"}, tc.CCIP)
	if err != nil {
		t.Fatal(err)
	}
	l.Info().Interface("loadedTestConfig", config).Msg("Loaded Test Config")

	l.Info().Msg("Starting ccip load test")
	l.Info().
		Int("Number of Nodes", *(config.CCIP.Load.NoOfNodes)).
		Interface("config", config.CCIP.Load).
		Msg("Test Config")

	var env = generateEnvironment()

	var env = deployment.Environment{}
	// output, err := actions.SetupCCIPHomeChain(l, sethClient, config.CCIP, workerNodes)
	// require.NoError(t, err)
	// env, err = crib.DeployPrerequisites(output, config.CCIP)
	// merge addressbooks
	// env, err := crib.DeployCCIPContracts(output, config.CCIP)

	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	state, err := ccipdeployment.LoadOnchainState(env)
	require.NoError(t, err)

	cfgl := config.Logging.Loki

	// Parse all events from the simulated chains, send to Loki
	loki, err := wasp.NewLokiClient(wasp.NewLokiConfig(cfgl.Endpoint, cfgl.TenantId, cfgl.BasicAuth, cfgl.BearerToken))
	if err != nil {
		l.Error().Err(err).Msg("Failed to create Loki client")
		return
	}
	defer loki.StopNow()

	// Based on the config, initiate a DestinationGun
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
			Gun:        NewDestinationGun(l, selector, env, state.Chains[selector].Receiver.Address(), loki),
			Labels:     CommonTestLabels,
			LokiConfig: wasp.NewLokiConfig(cfgl.Endpoint, cfgl.TenantId, cfgl.BasicAuth, cfgl.BearerToken),
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
				block, err := env.Chains[chainSelector].Client.BlockByNumber(ctx, big.NewInt(int64(blockNum)))
				require.NoError(t, err)
				timestamp := time.Unix(int64(block.Time()), 0)

				for _, root := range event.MerkleRoots {
					lokiLabels = setLokiLabels(root.SourceChainSelector, chainSelector)

					for i := root.MinSeqNr; i <= root.MaxSeqNr; i++ {
						// todo: push loki calls to channel?
						SendMetricsToLoki(l, loki, lokiLabels, &LokiMetric{
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
				block, err := env.Chains[chainSelector].Client.BlockByNumber(ctx, big.NewInt(int64(blockNum)))
				require.NoError(t, err)
				timestamp := time.Unix(int64(block.Time()), 0)

				// todo: push loki calls to channel?
				lokiLabels = setLokiLabels(execIterator.Event.SourceChainSelector, chainSelector)
				SendMetricsToLoki(l, loki, lokiLabels, &LokiMetric{
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

func generateEnvironment() {
	ab := readAddressBook(abPath)
	nIds := readNodeIds(nodeIdsPath)

}
