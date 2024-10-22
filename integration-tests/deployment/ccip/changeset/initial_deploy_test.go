package changeset

import (
	"testing"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	ccdeploy "github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestInitialDeploy(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := ccdeploy.Context(t)
	tenv := ccdeploy.NewMemoryEnvironment(t, lggr, 3)
	e := tenv.Env

	state, err := ccdeploy.LoadOnchainState(tenv.Env, tenv.Ab)
	require.NoError(t, err)
	require.NotNil(t, state.Chains[tenv.HomeChainSel].LinkToken)

	feeds := state.Chains[tenv.FeedChainSel].USDFeeds
	tokenConfig := ccdeploy.NewTokenConfig()
	tokenConfig.UpsertTokenInfo(ccdeploy.LinkSymbol,
		pluginconfig.TokenInfo{
			AggregatorAddress: feeds[ccdeploy.LinkSymbol].Address().String(),
			Decimals:          ccdeploy.LinkDecimals,
			DeviationPPB:      cciptypes.NewBigIntFromInt64(1e9),
		},
	)
	tokenConfig.UpsertTokenInfo(ccdeploy.WethSymbol,
		pluginconfig.TokenInfo{
			AggregatorAddress: feeds[ccdeploy.WethSymbol].Address().String(),
			Decimals:          ccdeploy.WethDecimals,
			DeviationPPB:      cciptypes.NewBigIntFromInt64(1e9),
		},
	)

	output, err := InitialDeployChangeSet(tenv.Ab, tenv.Env, ccdeploy.DeployCCIPContractConfig{
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

	// Initial state for tokens and gas prices
	initialGasUpdates := getInitialGasUpdates(t, e, state)
	initialTokenUpdates := getInitialTokenUpdates(t, e, state)

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

	// Token and Gas prices should be updated in FeeQuoter
	assertUpdatedGas(t, e, state, initialGasUpdates)
	assertUpdatedTokens(t, e, state, initialTokenUpdates)

	// Wait for all exec reports to land
	ccdeploy.ConfirmExecWithSeqNrForAll(t, e, state, expectedSeqNum, startBlocks)
}

func getInitialGasUpdates(
	t *testing.T,
	e deployment.Environment,
	state ccdeploy.CCIPOnChainState,
) map[uint64]map[uint64]fee_quoter.InternalTimestampedPackedUint224 {
	lggr := logger.TestLogger(t)
	srcToDestGasPriceTimestamps := make(map[uint64]map[uint64]fee_quoter.InternalTimestampedPackedUint224)
	for src := range e.Chains {
		feeQuoter := state.Chains[src].FeeQuoter
		for dest := range e.Chains {
			if src == dest {
				continue
			}
			gasUpdate, err := feeQuoter.GetDestinationChainGasPrice(nil, dest)
			require.NoError(t, err)
			require.NotNil(t, gasUpdate)
			require.Equal(t, ccdeploy.InitialGasPrice, gasUpdate.Value)
			lggr.Infow("Gas price",
				"src", src,
				"dest", dest,
				"gasUpdate", gasUpdate)
			if srcToDestGasPriceTimestamps[src] == nil {
				srcToDestGasPriceTimestamps[src] = make(map[uint64]fee_quoter.InternalTimestampedPackedUint224)
			}
			srcToDestGasPriceTimestamps[src][dest] = gasUpdate
		}
	}
	return srcToDestGasPriceTimestamps
}

func assertUpdatedGas(
	t *testing.T,
	e deployment.Environment,
	state ccdeploy.CCIPOnChainState,
	initialUpdates map[uint64]map[uint64]fee_quoter.InternalTimestampedPackedUint224,
) {
	lggr := logger.TestLogger(t)
	for src := range e.Chains {
		feeQuoter := state.Chains[src].FeeQuoter
		for dest := range e.Chains {
			if src == dest {
				continue
			}
			gasUpdate, err := feeQuoter.GetDestinationChainGasPrice(nil, dest)
			require.NoError(t, err)
			require.NotNil(t, gasUpdate)
			// Different value
			require.NotEqual(t, initialUpdates[src][dest].Value, gasUpdate.Value)
			// Newer timestamp
			require.True(t, initialUpdates[src][dest].Timestamp < gasUpdate.Timestamp)
			lggr.Infow("Gas price",
				"src", src,
				"dest", dest,
				"gasUpdate", gasUpdate)
		}
	}

}

func getInitialTokenUpdates(
	t *testing.T,
	e deployment.Environment,
	state ccdeploy.CCIPOnChainState,
) map[uint64]fee_quoter.InternalTimestampedPackedUint224 {
	lggr := logger.TestLogger(t)
	srcToDestTokenPriceTimestamps := make(map[uint64]fee_quoter.InternalTimestampedPackedUint224)
	for chain := range e.Chains {
		feeQuoter := state.Chains[chain].FeeQuoter
		linkAddress := state.Chains[chain].LinkToken.Address()
		linkUpdate, err := feeQuoter.GetTokenPrice(nil, linkAddress)
		require.NoError(t, err)
		require.NotNil(t, linkUpdate)
		require.Equal(t, ccdeploy.InitialLinkPrice, linkUpdate.Value)
		lggr.Infow("LinkPrice",
			"chain", chain,
			"LinkUpdate", linkUpdate)
		srcToDestTokenPriceTimestamps[chain] = linkUpdate
	}
	return srcToDestTokenPriceTimestamps
}

func assertUpdatedTokens(
	t *testing.T,
	e deployment.Environment,
	state ccdeploy.CCIPOnChainState,
	initialUpdates map[uint64]fee_quoter.InternalTimestampedPackedUint224,
) {
	lggr := logger.TestLogger(t)
	for chain := range e.Chains {
		feeQuoter := state.Chains[chain].FeeQuoter
		linkAddress := state.Chains[chain].LinkToken.Address()
		tokenUpdate, err := feeQuoter.GetTokenPrice(nil, linkAddress)
		require.NoError(t, err)
		require.NotNil(t, tokenUpdate)
		// Different value
		require.NotEqual(t, initialUpdates[chain].Value, tokenUpdate.Value)
		// Newer timestamp
		require.True(t, initialUpdates[chain].Timestamp < tokenUpdate.Timestamp)
		lggr.Infow("LinkPrice",
			"chain", chain,
			"LinkUpdate", tokenUpdate)
	}
}
