package changeset

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/stretchr/testify/require"
)

func TestInitialDeploy(t *testing.T) {
	lggr := logger.TestLogger(t)
	tenv := NewMemoryEnvironmentWithJobsAndContracts(t, lggr, memory.MemoryEnvironmentConfig{
		Chains:             3,
		Nodes:              4,
		Bootstraps:         1,
		NumOfUsersPerChain: 1,
	}, nil)
	e := tenv.Env
	state, err := LoadOnchainState(e)
	require.NoError(t, err)
	// Add all lanes
	require.NoError(t, AddLanesForAll(e, state))
	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	// Send a message from each chain to every other chain.
	expectedSeqNum := make(map[SourceDestPair]uint64)
	expectedSeqNumExec := make(map[SourceDestPair][]uint64)

	for src := range e.Chains {
		for dest, destChain := range e.Chains {
			if src == dest {
				continue
			}
			latesthdr, err := destChain.Client.HeaderByNumber(testcontext.Get(t), nil)
			require.NoError(t, err)
			block := latesthdr.Number.Uint64()
			startBlocks[dest] = &block
			require.GreaterOrEqual(t, len(tenv.Users[src]), 1)
			msgSentEvent, err := DoSendRequest(t, e, state,
				WithTestRouter(false),
				WithSourceChain(src),
				WithDestChain(dest),
				WithEvm2AnyMessage(router.ClientEVM2AnyMessage{
					Receiver:     common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32),
					Data:         []byte("hello"),
					TokenAmounts: nil,
					FeeToken:     common.HexToAddress("0x0"),
					ExtraArgs:    nil,
				}),
				WithSender(tenv.Users[src][0]),
			)
			require.NoError(t, err)
			expectedSeqNum[SourceDestPair{
				SourceChainSelector: src,
				DestChainSelector:   dest,
			}] = msgSentEvent.SequenceNumber
			expectedSeqNumExec[SourceDestPair{
				SourceChainSelector: src,
				DestChainSelector:   dest,
			}] = []uint64{msgSentEvent.SequenceNumber}
		}
	}

	// Wait for all commit reports to land.
	ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

	// Confirm token and gas prices are updated
	ConfirmTokenPriceUpdatedForAll(t, e, state, startBlocks,
		DefaultInitialPrices.LinkPrice, DefaultInitialPrices.WethPrice)
	// TODO: Fix gas prices?
	//ConfirmGasPriceUpdatedForAll(t, e, state, startBlocks)
	//
	//// Wait for all exec reports to land
	ConfirmExecWithSeqNrsForAll(t, e, state, expectedSeqNumExec, startBlocks)
}
