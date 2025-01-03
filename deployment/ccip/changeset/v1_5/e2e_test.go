package v1_5

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
)

// This test only works if the destination chain id is 1337
// Otherwise it shows error for offchain and onchain config digest mismatch
func TestE2ELegacy(t *testing.T) {
	e := changeset.NewMemoryEnvironment(
		t,
		changeset.WithLegacyDeployment(),
		changeset.WithChains(3),
		changeset.WithChainIds([]uint64{chainselectors.GETH_TESTNET.EvmChainID}))
	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	allChains := e.Env.AllChainSelectorsExcluding([]uint64{chainselectors.GETH_TESTNET.Selector})
	require.Contains(t, e.Env.AllChainSelectors(), chainselectors.GETH_TESTNET.Selector)
	require.Len(t, allChains, 2)
	src, dest := allChains[1], chainselectors.GETH_TESTNET.Selector
	srcChain := e.Env.Chains[src]
	destChain := e.Env.Chains[dest]
	pairs := []changeset.SourceDestPair{
		{SourceChainSelector: src, DestChainSelector: dest},
	}
	e.Env = AddLanes(t, e.Env, state, pairs)
	// reload state after adding lanes
	state, err = changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	sentEvent, err := SendRequest(t, e.Env, state,
		changeset.WithSourceChain(src),
		changeset.WithDestChain(dest),
		changeset.WithTestRouter(false),
		changeset.WithEvm2AnyMessage(router.ClientEVM2AnyMessage{
			Receiver:     common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32),
			Data:         []byte("hello"),
			TokenAmounts: nil,
			FeeToken:     common.HexToAddress("0x0"),
			ExtraArgs:    nil,
		}),
	)
	require.NoError(t, err)
	require.NotNil(t, sentEvent)
	destStartBlock, err := destChain.Client.HeaderByNumber(context.Background(), nil)
	require.NoError(t, err)
	WaitForCommit(t, srcChain, destChain, state.Chains[dest].CommitStore[src], sentEvent.Message.SequenceNumber)
	WaitForExecute(t, srcChain, destChain, state.Chains[dest].EVM2EVMOffRamp[src], []uint64{sentEvent.Message.SequenceNumber}, destStartBlock.Number.Uint64())
}
