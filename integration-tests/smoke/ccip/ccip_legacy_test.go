package smoke

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
)

// This test does not run in CI, it is only written as an example of how to write a test for the legacy CCIP
func TestE2ELegacy(t *testing.T) {
	e, _ := testsetups.NewIntegrationEnvironment(t, changeset.WithLegacyDeployment())
	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	allChains := e.Env.AllChainSelectors()
	require.Len(t, allChains, 2)
	src, dest := allChains[0], allChains[1]
	srcChain := e.Env.Chains[src]
	destChain := e.Env.Chains[dest]
	pairs := []changeset.SourceDestPair{
		{SourceChainSelector: src, DestChainSelector: dest},
	}
	e.Env = v1_5.AddLanes(t, e.Env, state, pairs)
	// reload state after adding lanes
	state, err = changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	sentEvent, err := v1_5.SendRequest(t, e.Env, state,
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
	v1_5.WaitForCommit(t, srcChain, destChain, state.Chains[dest].CommitStore[src], sentEvent.Message.SequenceNumber)
	v1_5.WaitForExecute(t, srcChain, destChain, state.Chains[dest].EVM2EVMOffRamp[src], []uint64{sentEvent.Message.SequenceNumber}, destStartBlock.Number.Uint64())
}
