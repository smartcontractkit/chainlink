package smoke

import (
	"testing"

	"github.com/stretchr/testify/require"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"

	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	ccdeploy "github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/changeset"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestInitialDeployOnLocal(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := ccdeploy.Context(t)
	tenv, _, _ := ccdeploy.NewLocalDevEnvironment(t, lggr)
	e := tenv.Env

	state, err := ccdeploy.LoadOnchainState(tenv.Env, tenv.Ab)
	require.NoError(t, err)

	feeds := state.Chains[tenv.FeedChainSel].USDFeeds
	tokenConfig := ccdeploy.NewTokenConfig()
	tokenConfig.UpsertTokenInfo(ccdeploy.LinkSymbol,
		pluginconfig.TokenInfo{
			AggregatorAddress: feeds[ccdeploy.LinkSymbol].Address().String(),
			Decimals:          ccdeploy.LinkDecimals,
			DeviationPPB:      cciptypes.NewBigIntFromInt64(1e9),
		},
	)
	// Apply migration
	output, err := changeset.InitialDeployChangeSet(tenv.Ab, tenv.Env, ccdeploy.DeployCCIPContractConfig{
		HomeChainSel:       tenv.HomeChainSel,
		FeedChainSel:       tenv.FeedChainSel,
		ChainsToDeploy:     tenv.Env.AllChainSelectors(),
		TokenConfig:        tokenConfig,
		MCMSConfig:         ccdeploy.NewTestMCMSConfig(t, e),
		CapabilityRegistry: state.Chains[tenv.HomeChainSel].CapabilityRegistry.Address(),
		FeeTokenContracts:  tenv.FeeTokenContracts,
		OCRSecrets:         deployment.XXXGenerateTestOCRSecrets(),
	})
	require.NoError(t, err)
	// Get new state after migration.
	state, err = ccdeploy.LoadOnchainState(e, tenv.Ab)
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
			seqNum := ccdeploy.SendRequest(t, e, state, src, dest, false)
			expectedSeqNum[dest] = seqNum
		}
	}

	// Wait for all commit reports to land.
	ccdeploy.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

	// After commit is reported on all chains, token prices should be updated in FeeQuoter.
	for dest := range e.Chains {
		linkAddress := state.Chains[dest].LinkToken.Address()
		feeQuoter := state.Chains[dest].FeeQuoter
		timestampedPrice, err := feeQuoter.GetTokenPrice(nil, linkAddress)
		require.NoError(t, err)
		require.Equal(t, ccdeploy.MockLinkPrice, timestampedPrice.Value)
	}

	// Wait for all exec reports to land
	ccdeploy.ConfirmExecWithSeqNrForAll(t, e, state, expectedSeqNum, startBlocks)

	// TODO: Apply the proposal.
}
