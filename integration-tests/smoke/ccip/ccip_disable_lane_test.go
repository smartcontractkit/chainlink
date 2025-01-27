package ccip

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
)

// Intention of this test is to ensure that the lane can be disabled and enabled correctly without disrupting the other lanes.
func TestDisableLane(t *testing.T) {
	tenv, _, _ := testsetups.NewIntegrationEnvironment(t,
		testhelpers.WithNumOfChains(3),
		testhelpers.WithNumOfUsersPerChain(2),
	)

	e := tenv.Env
	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	// add all lanes
	testhelpers.AddLanesForAll(t, &tenv, state)

	var (
		chains                 = e.AllChainSelectors()
		chainA, chainB, chainC = chains[0], chains[1], chains[2]
		expectedSeqNumExec     = make(map[testhelpers.SourceDestPair][]uint64)
		startBlocks            = make(map[uint64]*uint64)
		pairs                  = testsetups.GetSourceDestPairs([]uint64{chainA, chainB, chainC})

		sendmessage = func(src, dest uint64, deployer *bind.TransactOpts) (*onramp.OnRampCCIPMessageSent, error) {
			return testhelpers.DoSendRequest(
				t,
				e,
				state,
				testhelpers.WithSender(deployer),
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
		}

		assertSendRequestReverted = func(src, dest uint64, deployer *bind.TransactOpts) {
			_, err = sendmessage(src, dest, deployer)
			require.Error(t, err)
			require.Contains(t, err.Error(), "execution reverted")
		}

		assertRequestSent = func(src, dest uint64, deployer *bind.TransactOpts) {
			latestHeader, err := e.Chains[dest].Client.HeaderByNumber(testcontext.Get(t), nil)
			require.NoError(t, err)
			block := latestHeader.Number.Uint64()
			messageSentEvent, err := sendmessage(src, dest, e.Chains[src].DeployerKey)
			require.NoError(t, err)
			expectedSeqNumExec[testhelpers.SourceDestPair{
				SourceChainSelector: src,
				DestChainSelector:   dest,
			}] = []uint64{messageSentEvent.SequenceNumber}
			startBlocks[dest] = &block
		}
	)

	// disable lane A -> B
	testhelpers.RemoveLane(t, &tenv, chainA, chainB, false)
	// send a message to confirm it is reverted between chainA and chainB
	assertSendRequestReverted(chainA, chainB, e.Chains[chainA].Users[0])
	// send a message between chainB and chainA to confirm it is not reverted
	assertRequestSent(chainB, chainA, e.Chains[chainB].Users[0])
	// disable lane B -> A
	testhelpers.RemoveLane(t, &tenv, chainB, chainA, false)
	assertSendRequestReverted(chainB, chainA, e.Chains[chainB].Users[0])

	// send message in other lanes and ensure they are delivered
	go func() {
		assertRequestSent(chainA, chainC, e.Chains[chainA].Users[1])
		assertRequestSent(chainC, chainA, e.Chains[chainC].Users[1])
		assertRequestSent(chainB, chainC, e.Chains[chainB].Users[1])
		assertRequestSent(chainC, chainB, e.Chains[chainC].Users[1])
	}()
	// disable lanes between A & C and C & B while requests are getting sent
	testhelpers.RemoveLane(t, &tenv, chainA, chainC, false)
	testhelpers.RemoveLane(t, &tenv, chainC, chainA, false)
	testhelpers.RemoveLane(t, &tenv, chainB, chainC, false)
	testhelpers.RemoveLane(t, &tenv, chainC, chainB, false)
	// check fee quoter returns error when the lane is disabled
	gp, err := state.Chains[chainA].FeeQuoter.GetTokenAndGasPrices(&bind.CallOpts{
		Context: tests.Context(t),
	}, state.Chains[chainB].Weth9.Address(), chainB)
	require.Error(t, err)
	require.Contains(t, err.Error(), "execution reverted")
	require.Nil(t, gp.GasPriceValue)
	require.Nil(t, gp.TokenPrice)
	// confirm that message sent in all lanes are reverted after disabling the lanes
	for _, pair := range pairs {
		assertSendRequestReverted(pair.SourceChainSelector, pair.DestChainSelector, e.Chains[pair.SourceChainSelector].Users[0])
	}
	// re-enable all the lanes
	var (
		linkPrice = deployment.E18Mult(100)
		wethPrice = deployment.E18Mult(4000)
	)
	for _, pair := range pairs {
		testhelpers.AddLane(t, &tenv, pair.SourceChainSelector, pair.DestChainSelector, false,
			map[uint64]*big.Int{
				pair.DestChainSelector: testhelpers.DefaultGasPrice,
			},
			map[common.Address]*big.Int{
				state.Chains[pair.SourceChainSelector].LinkToken.Address(): linkPrice,
				state.Chains[pair.SourceChainSelector].Weth9.Address():     wethPrice,
			},
			changeset.DefaultFeeQuoterDestChainConfig(true))
	}
	// send a message in all the lane including re-enabled lanes
	for _, pair := range pairs {
		assertRequestSent(pair.SourceChainSelector, pair.DestChainSelector, e.Chains[pair.SourceChainSelector].Users[0])
	}
	// confirm all messages are delivered
	testhelpers.ConfirmExecWithSeqNrsForAll(t, e, state, expectedSeqNumExec, startBlocks)
}
