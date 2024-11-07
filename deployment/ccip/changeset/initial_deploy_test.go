package changeset

import (
	"testing"

	"github.com/smartcontractkit/chainlink/deployment"

	ccdeploy "github.com/smartcontractkit/chainlink/deployment/ccip"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/stretchr/testify/require"
)

func TestInitialDeploy(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := ccdeploy.Context(t)
	tenv := ccdeploy.NewMemoryEnvironment(t, lggr, 3, 4)
	e := tenv.Env

	state, err := ccdeploy.LoadOnchainState(tenv.Env)
	require.NoError(t, err)
	require.NotNil(t, state.Chains[tenv.HomeChainSel].LinkToken)

	output, err := InitialDeploy(tenv.Env, ccdeploy.DeployCCIPContractConfig{
		HomeChainSel:   tenv.HomeChainSel,
		FeedChainSel:   tenv.FeedChainSel,
		ChainsToDeploy: tenv.Env.AllChainSelectors(),
		TokenConfig:    ccdeploy.NewTestTokenConfig(state.Chains[tenv.FeedChainSel].USDFeeds),
		MCMSConfig:     ccdeploy.NewTestMCMSConfig(t, e),
		OCRSecrets:     deployment.XXXGenerateTestOCRSecrets(),
	})
	require.NoError(t, err)
	// Get new state after migration.
	require.NoError(t, tenv.Env.ExistingAddresses.Merge(output.AddressBook))
	state, err = ccdeploy.LoadOnchainState(e)
	require.NoError(t, err)

	// Ensure capreg logs are up to date.
	ccdeploy.ReplayLogs(t, e.Offchain, tenv.ReplayBlocks)

	// Apply the jobs.
	for nodeID, jobs := range output.JobSpecs {
		for _, job := range jobs {
			// Note these auto-accept
			_, err := e.Offchain.ProposeJob(ctx,
				&jobv1.ProposeJobRequest{
					NodeId: nodeID,
					Spec:   job,
				})
			require.NoError(t, err)
		}
	}

	// Add all lanes
	require.NoError(t, ccdeploy.AddLanesForAll(e, state))
	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	// Send a message from each chain to every other chain.
	expectedSeqNum := make(map[uint64]uint64)

	for src := range e.Chains {
		for dest, destChain := range e.Chains {
			if src == dest {
				continue
			}
			latesthdr, err := destChain.Client.HeaderByNumber(testcontext.Get(t), nil)
			require.NoError(t, err)
			block := latesthdr.Number.Uint64()
			startBlocks[dest] = &block
			seqNum := ccdeploy.TestSendRequest(t, e, state, src, dest, false, nil)
			expectedSeqNum[dest] = seqNum
		}
	}

	// Wait for all commit reports to land.
	ccdeploy.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

	// Confirm token and gas prices are updated
	ccdeploy.ConfirmTokenPriceUpdatedForAll(t, e, state, startBlocks)
	ccdeploy.ConfirmGasPriceUpdatedForAll(t, e, state, startBlocks)

	// Wait for all exec reports to land
	ccdeploy.ConfirmExecWithSeqNrForAll(t, e, state, expectedSeqNum, startBlocks)
}
