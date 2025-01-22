package ccip

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
)

// Intention of this test is to ensure that the lane can be disabled and enabled correctly without disrupting the other lanes.
func TestDisableLane(t *testing.T) {
	tenv, _, _ := testsetups.NewIntegrationEnvironment(t,
		testhelpers.WithNumOfChains(3),
	)

	e := tenv.Env
	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	// Add all lanes
	testhelpers.AddLanesForAll(t, &tenv, state)

	chains := e.AllChainSelectors()
	chainA, chainB, chainC := chains[0], chains[1], chains[2]

	assertSendRequestReverted := func(src, dest uint64) {
		_, err = testhelpers.DoSendRequest(
			t,
			e,
			state,
			testhelpers.WithSender(e.Chains[src].DeployerKey),
			testhelpers.WithSourceChain(src),
			testhelpers.WithDestChain(dest),
			testhelpers.WithTestRouter(false),
			testhelpers.WithEvm2AnyMessage(router.ClientEVM2AnyMessage{
				Receiver:     common.LeftPadBytes(state.Chains[chainB].Receiver.Address().Bytes(), 32),
				Data:         []byte("hello"),
				TokenAmounts: nil,
				FeeToken:     common.HexToAddress("0x0"),
				ExtraArgs:    nil,
			}))
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted")
	}

	expectedSeqNumExec := make(map[testhelpers.SourceDestPair][]uint64)
	startBlocks := make(map[uint64]*uint64)
	assertRequestSent := func(src, dest uint64) {
		latestHeader, err := e.Chains[dest].Client.HeaderByNumber(testcontext.Get(t), nil)
		require.NoError(t, err)
		block := latestHeader.Number.Uint64()
		msgSentEvent := testhelpers.TestSendRequest(t, e, state, src, dest, false, router.ClientEVM2AnyMessage{
			Receiver:     common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32),
			Data:         []byte("hello"),
			TokenAmounts: nil,
			FeeToken:     common.HexToAddress("0x0"),
			ExtraArgs:    nil,
		})
		expectedSeqNumExec[testhelpers.SourceDestPair{
			SourceChainSelector: src,
			DestChainSelector:   dest,
		}] = []uint64{msgSentEvent.SequenceNumber}
		startBlocks[dest] = &block
	}

	// Disable specific lane A -> B
	testhelpers.UpdateLane(t, &tenv, chainA, chainB, false, false)
	// send a transaction to confirm it is reverted between chainA and chainB
	assertSendRequestReverted(chainA, chainB)
	assertRequestSent(chainB, chainA)
	// Disable B -> A
	testhelpers.UpdateLane(t, &tenv, chainB, chainA, false, false)
	assertSendRequestReverted(chainB, chainA)

	// send transactions in other lanes and ensure they are delivered
	assertRequestSent(chainA, chainC)
	assertRequestSent(chainC, chainA)
	assertRequestSent(chainB, chainC)
	assertRequestSent(chainC, chainB)

	// re-enable the lanes A -> B and B -> A
	testhelpers.UpdateLane(t, &tenv, chainA, chainB, false, true)
	testhelpers.UpdateLane(t, &tenv, chainB, chainA, false, true)
	// send a transaction in all the lane including re-enabled lanes
	pairs := testsetups.GetSourceDestPairs([]uint64{chainA, chainB, chainC})
	for _, pair := range pairs {
		assertRequestSent(pair.SourceChainSelector, pair.DestChainSelector)
	}
	// Confirm all exec reports
	testhelpers.ConfirmExecWithSeqNrsForAll(t, e, state, expectedSeqNumExec, startBlocks)
}
