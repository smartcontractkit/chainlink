package changeset

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
)

func TestAddLanesWithTestRouter(t *testing.T) {
	t.Parallel()
	e, _ := NewMemoryEnvironment(t)
	// Here we have CR + nodes set up, but no CCIP contracts deployed.
	state, err := LoadOnchainState(e.Env)
	require.NoError(t, err)

	selectors := e.Env.AllChainSelectors()
	chain1, chain2 := selectors[0], selectors[1]
	AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, chain1, chain2, true)
	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	// Send a message from each chain to every other chain.
	expectedSeqNumExec := make(map[SourceDestPair][]uint64)
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
	expectedSeqNumExec[SourceDestPair{
		SourceChainSelector: chain1,
		DestChainSelector:   chain2,
	}] = []uint64{msgSentEvent.SequenceNumber}
	ConfirmExecWithSeqNrsForAll(t, e.Env, state, expectedSeqNumExec, startBlocks)
}
