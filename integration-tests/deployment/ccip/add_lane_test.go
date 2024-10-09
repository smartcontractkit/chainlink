package ccipdeployment

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// TestAddLane covers the workflow of adding a lane
// between existing supported chains in CCIP.
func TestAddLane(t *testing.T) {
	e := NewMemoryEnvironmentWithJobs(t, logger.TestLogger(t), 2)
	// Here we have CR + nodes set up, but no CCIP contracts deployed.
	state, err := LoadOnchainState(e.Env, e.Ab)
	require.NoError(t, err)
	selectors := e.Env.AllChainSelectors()
	from, to := selectors[0], selectors[1]
	feeds := state.Chains[e.FeedChainSel].USDFeeds
	tokenConfig := NewTokenConfig()
	tokenConfig.UpsertTokenInfo(LinkSymbol,
		pluginconfig.TokenInfo{
			AggregatorAddress: feeds[LinkSymbol].Address().String(),
			Decimals:          LinkDecimals,
			DeviationPPB:      cciptypes.NewBigIntFromInt64(1e9),
		},
	)
	feeTokenContracts := make(map[uint64]FeeTokenContracts)
	for _, sel := range []uint64{from, to} {
		feeTokenContracts[sel] = e.FeeTokenContracts[sel]
	}
	// Set up CCIP contracts and a DON per chain.
	err = DeployCCIPContracts(e.Env, e.Ab, DeployCCIPContractConfig{
		HomeChainSel:       e.HomeChainSel,
		FeedChainSel:       e.FeedChainSel,
		TokenConfig:        tokenConfig,
		MCMSConfig:         NewTestMCMSConfig(t, e.Env),
		FeeTokenContracts:  feeTokenContracts,
		ChainsToDeploy:     []uint64{from, to},
		CapabilityRegistry: state.Chains[e.HomeChainSel].CapabilityRegistry.Address(),
	})
	require.NoError(t, err)

	// We expect no lanes available on any chain.
	state, err = LoadOnchainState(e.Env, e.Ab)
	require.NoError(t, err)
	for _, sel := range []uint64{from, to} {
		chain := state.Chains[sel]
		offRamps, err := chain.Router.GetOffRamps(nil)
		require.NoError(t, err)
		require.Len(t, offRamps, 0)
	}

	replayBlocks, err := LatestBlocksByChain(testcontext.Get(t), e.Env.Chains)
	require.NoError(t, err)
	// Add one lane and send traffic.
	require.NoError(t, AddLane(e.Env, state, from, to))
	for _, sel := range []uint64{from, to} {
		chain := state.Chains[sel]
		offRamps, err := chain.Router.GetOffRamps(nil)
		require.NoError(t, err)
		if sel == to {
			require.Len(t, offRamps, 1)
			srcCfg, err := chain.OffRamp.GetSourceChainConfig(nil, from)
			require.NoError(t, err)
			require.Equal(t, common.LeftPadBytes(state.Chains[from].OnRamp.Address().Bytes(), 32), srcCfg.OnRamp)
		} else {
			require.Len(t, offRamps, 0)
		}
	}

	time.Sleep(30 * time.Second)
	ReplayLogs(t, e.Env.Offchain, replayBlocks)

	latesthdr, err := e.Env.Chains[to].Client.HeaderByNumber(testcontext.Get(t), nil)
	require.NoError(t, err)
	startBlock := latesthdr.Number.Uint64()
	seqNum := SendRequest(t, e.Env, state, from, to, false)
	require.Equal(t, uint64(1), seqNum)
	require.NoError(t,
		ConfirmCommitWithExpectedSeqNumRange(
			t, e.Env.Chains[from], e.Env.Chains[to],
			state.Chains[to].OffRamp, &startBlock,
			cciptypes.SeqNumRange{
				cciptypes.SeqNum(seqNum),
				cciptypes.SeqNum(seqNum),
			}))
	require.NoError(t, ConfirmExecWithSeqNr(t, e.Env.Chains[from], e.Env.Chains[to], state.Chains[to].OffRamp, &startBlock, seqNum))

	// Add another lane
	require.NoError(t, AddLane(e.Env, state, to, from))
	// disable onRamp for previous lane from -> to
	tx, err := state.Chains[to].OffRamp.ApplySourceChainConfigUpdates(
		e.Env.Chains[to].DeployerKey, []offramp.OffRampSourceChainConfigArgs{
			{
				Router:              state.Chains[to].Router.Address(),
				SourceChainSelector: from,
				IsEnabled:           false,
				OnRamp:              common.LeftPadBytes(state.Chains[from].OnRamp.Address().Bytes(), 32),
			},
		})
	require.NoError(t, err)
	_, err = deployment.ConfirmIfNoError(e.Env.Chains[to], tx, err)
	require.NoError(t, err)

	// Send traffic on the first lane and it should fail
	latesthdr, err = e.Env.Chains[to].Client.HeaderByNumber(testcontext.Get(t), nil)
	require.NoError(t, err)
	startBlock = latesthdr.Number.Uint64()
	seqNum = SendRequest(t, e.Env, state, from, to, false)
	require.Equal(t, uint64(1), seqNum)
	require.Error(t, ConfirmExecWithSeqNr(t, e.Env.Chains[from], e.Env.Chains[to], state.Chains[to].OffRamp, &startBlock, seqNum))

	// Send traffic on the second lane and it should succeed
	latesthdr, err = e.Env.Chains[from].Client.HeaderByNumber(testcontext.Get(t), nil)
	require.NoError(t, err)
	startBlock = latesthdr.Number.Uint64()
	seqNum = SendRequest(t, e.Env, state, to, from, false)
	require.Equal(t, uint64(1), seqNum)
	require.Error(t, ConfirmExecWithSeqNr(t, e.Env.Chains[to], e.Env.Chains[from], state.Chains[from].OffRamp, &startBlock, seqNum))
}
