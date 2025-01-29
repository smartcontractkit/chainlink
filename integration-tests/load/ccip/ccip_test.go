package ccip

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"context"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
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
)

// todo: add multiple keys and rotate them when sending messages
// this key only works on simulated geth chains in crib
const simChainTestKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

// step 1: setup
// Parse the test config, initialize CRIB with configurations defined
// step 2: subscribe
// Create event subscribers on the offramp
// step 3: load
// Use wasp to initiate load
// step 4: teardown
// Stop the chains, cleanup the environment
func TestCCIPLoad_RPS(t *testing.T) {
	// comment out when executing the test
	t.Skip("Skipping test as this test should not be auto triggered")
	lggr := logger.Test(t)

	// get user defined configurations
	config, err := tc.GetConfig([]string{"Load"}, tc.CCIP)
	require.NoError(t, err)
	lggr.Infof("loaded ccip test config: %+v", config.CCIP.Load)
	userOverrides := config.CCIP.Load
	userOverrides.Validate(t)

	// apply user defined test duration
	ctx, cancel := context.WithTimeout(context.Background(), userOverrides.GetTestDuration())
	defer cancel()

	// generate environment from crib-produced files
	cribEnv := crib.NewDevspaceEnvFromStateDir(*userOverrides.CribEnvDirectory)
	cribDeployOutput, err := cribEnv.GetConfig(simChainTestKey)
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

		wg.Add(2)
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

	requestFrequency, err := time.ParseDuration(*userOverrides.RequestFrequency)
	require.NoError(t, err)

	for _, gun := range gunMap {
		p.Add(wasp.NewGenerator(&wasp.Config{
			T:           t,
			GenName:     "ccipLoad",
			LoadType:    wasp.RPS,
			CallTimeout: userOverrides.GetLoadDuration(),
			// 1 request per second for n seconds
			Schedule: wasp.Plain(1, userOverrides.GetLoadDuration()),
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
	require.NoError(t, err)

	for _, gun := range gunMap {
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
	}

	wg.Wait()
	lggr.Infow("finished wait group")
	for err := range errChan {
		require.NoError(t, err)
	}
}
