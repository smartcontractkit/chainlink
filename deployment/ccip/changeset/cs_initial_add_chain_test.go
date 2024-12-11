package changeset

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/stretchr/testify/require"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
)

func TestInitialAddChainAppliedTwice(t *testing.T) {
	// This already applies the initial add chain changeset.
	e := NewMemoryEnvironment(t)

	state, err := LoadOnchainState(e.Env)
	require.NoError(t, err)

	// now try to apply it again for the second time
	// Build the per chain config.
	allChains := e.Env.AllChainSelectors()
	tokenConfig := NewTestTokenConfig(state.Chains[e.FeedChainSel].USDFeeds)
	chainConfigs := make(map[uint64]CCIPOCRParams)
	timelockContractsPerChain := make(map[uint64]*commonchangeset.TimelockExecutionContracts)

	for _, chain := range allChains {
		timelockContractsPerChain[chain] = &commonchangeset.TimelockExecutionContracts{
			Timelock:  state.Chains[chain].Timelock,
			CallProxy: state.Chains[chain].CallProxy,
		}
		tokenInfo := tokenConfig.GetTokenInfo(e.Env.Logger, state.Chains[chain].LinkToken, state.Chains[chain].Weth9)
		ocrParams := DefaultOCRParams(e.FeedChainSel, tokenInfo, []pluginconfig.TokenDataObserverConfig{})
		chainConfigs[chain] = ocrParams
	}
	e.Env, err = commonchangeset.ApplyChangesets(t, e.Env, timelockContractsPerChain, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(ConfigureNewChains),
			Config: NewChainsConfig{
				HomeChainSel:       e.HomeChainSel,
				FeedChainSel:       e.FeedChainSel,
				ChainConfigByChain: chainConfigs,
			},
		},
	})
	require.NoError(t, err)
	// send requests
	chain1, chain2 := allChains[0], allChains[1]

	_, err = AddLanes(e.Env, AddLanesConfig{
		LaneConfigs: []LaneConfig{
			{
				SourceSelector:        chain1,
				DestSelector:          chain2,
				InitialPricesBySource: DefaultInitialPrices,
				FeeQuoterDestChain:    DefaultFeeQuoterDestChainConfig(),
				TestRouter:            true,
			},
		},
	})
	require.NoError(t, err)
	ReplayLogs(t, e.Env.Offchain, e.ReplayBlocks)
	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	// Send a message from each chain to every other chain.
	expectedSeqNumExec := make(map[SourceDestPair][]uint64)
	expectedSeqNum := make(map[SourceDestPair]uint64)
	latesthdr, err := e.Env.Chains[chain2].Client.HeaderByNumber(testcontext.Get(t), nil)
	require.NoError(t, err)
	block := latesthdr.Number.Uint64()
	startBlocks[chain2] = &block
	msgSentEvent := TestSendRequest(t, e.Env, state, chain1, chain2, true, router.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(state.Chains[chain2].Receiver.Address().Bytes(), 32),
		Data:         []byte("hello"),
		TokenAmounts: nil,
		FeeToken:     common.HexToAddress("0x0"),
		ExtraArgs:    nil,
	})

	expectedSeqNum[SourceDestPair{
		SourceChainSelector: chain1,
		DestChainSelector:   chain2,
	}] = msgSentEvent.SequenceNumber
	expectedSeqNumExec[SourceDestPair{
		SourceChainSelector: chain1,
		DestChainSelector:   chain2,
	}] = []uint64{msgSentEvent.SequenceNumber}
	ConfirmCommitForAllWithExpectedSeqNums(t, e.Env, state, expectedSeqNum, startBlocks)
	ConfirmExecWithSeqNrsForAll(t, e.Env, state, expectedSeqNumExec, startBlocks)
}
