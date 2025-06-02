package ccip

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

// Intention of this test is to ensure that the lane can be disabled and enabled correctly
// without disrupting the other lanes and in-flight requests are delivered.
func TestDisableLane(t *testing.T) {
	tenv, _, _ := testsetups.NewIntegrationEnvironment(t,
		testhelpers.WithNumOfChains(3),
		testhelpers.WithNumOfUsersPerChain(2),
	)

	e := tenv.Env
	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err)

	// add all lanes
	testhelpers.AddLanesForAll(t, &tenv, state)

	var (
		chains                 = e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilyEVM))
		chainA, chainB, chainC = chains[0], chains[1], chains[2]
		evmChains              = e.BlockChains.EVMChains()
		expectedSeqNumExec     = make(map[testhelpers.SourceDestPair][]uint64)
		startBlocks            = make(map[uint64]*uint64)
		pairs                  []testhelpers.SourceDestPair
		linkPrice              = deployment.E18Mult(100)
		wethPrice              = deployment.E18Mult(4000)
		noOfRequests           = 3
		sendmessage            = func(src, dest uint64, deployer *bind.TransactOpts) (*onramp.OnRampCCIPMessageSent, error) {
			return testhelpers.SendRequest(
				e,
				state,
				testhelpers.WithSender(deployer),
				testhelpers.WithSourceChain(src),
				testhelpers.WithDestChain(dest),
				testhelpers.WithTestRouter(false),
				testhelpers.WithEvm2AnyMessage(router.ClientEVM2AnyMessage{
					Receiver:     common.LeftPadBytes(state.MustGetEVMChainState(chainB).Receiver.Address().Bytes(), 32),
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
			latestHeader, err := evmChains[dest].Client.HeaderByNumber(testcontext.Get(t), nil)
			require.NoError(t, err)
			block := latestHeader.Number.Uint64()
			messageSentEvent, err := sendmessage(src, dest, evmChains[src].DeployerKey)
			require.NoError(t, err)
			expectedSeqNumExec[testhelpers.SourceDestPair{
				SourceChainSelector: src,
				DestChainSelector:   dest,
			}] = []uint64{messageSentEvent.SequenceNumber}
			startBlocks[dest] = &block
		}
	)

	// disable lane A -> B
	pairs = append(pairs, testhelpers.SourceDestPair{
		SourceChainSelector: chainA,
		DestChainSelector:   chainB,
	})
	testhelpers.RemoveLane(t, &tenv, chainA, chainB, false)
	// send a message to confirm it is reverted between A -> B
	assertSendRequestReverted(chainA, chainB, evmChains[chainA].Users[0])

	// send a message in other direction B -> A to confirm it is delivered
	assertRequestSent(chainB, chainA, evmChains[chainB].Users[0])
	testhelpers.ConfirmExecWithSeqNrsForAll(t, e, state, expectedSeqNumExec, startBlocks)

	// send a multiple message between A -> C and disable the lane while the requests are in-flight
	expectedSeqNumExec = make(map[testhelpers.SourceDestPair][]uint64)
	for range noOfRequests {
		assertRequestSent(chainA, chainC, evmChains[chainA].Users[1])
	}
	// disable lane A -> C while requests are getting sent in that lane
	pairs = append(pairs, testhelpers.SourceDestPair{
		SourceChainSelector: chainA,
		DestChainSelector:   chainC,
	})
	testhelpers.RemoveLane(t, &tenv, chainA, chainC, false)

	// confirm all in-flight messages are delivered in A -> C lane
	testhelpers.ConfirmExecWithSeqNrsForAll(t, e, state, expectedSeqNumExec, startBlocks)

	// now, as the lane is disabled, confirm that message sent in A -> C is reverted
	assertSendRequestReverted(chainA, chainC, evmChains[chainA].Users[0])

	// check getting token and gas price form fee quoter returns error when A -> C lane is disabled
	gp, err := state.MustGetEVMChainState(chainA).FeeQuoter.GetTokenAndGasPrices(&bind.CallOpts{
		Context: t.Context(),
	}, state.MustGetEVMChainState(chainC).Weth9.Address(), chainC)
	require.Error(t, err)
	require.Contains(t, err.Error(), "execution reverted")
	require.Nil(t, gp.GasPriceValue)
	require.Nil(t, gp.TokenPrice)

	// re-enable all the disabled lanes
	for _, pair := range pairs {
		testhelpers.AddLane(t, &tenv, pair.SourceChainSelector, pair.DestChainSelector, false,
			map[uint64]*big.Int{
				pair.DestChainSelector: testhelpers.DefaultGasPrice,
			},
			map[common.Address]*big.Int{
				state.MustGetEVMChainState(pair.SourceChainSelector).LinkToken.Address(): linkPrice,
				state.MustGetEVMChainState(pair.SourceChainSelector).Weth9.Address():     wethPrice,
			},
			v1_6.DefaultFeeQuoterDestChainConfig(true))
	}
	// send a message in all the lane including re-enabled lanes
	for _, pair := range pairs {
		assertRequestSent(pair.SourceChainSelector, pair.DestChainSelector, evmChains[pair.SourceChainSelector].Users[0])
	}
	// confirm all messages are delivered
	testhelpers.ConfirmExecWithSeqNrsForAll(t, e, state, expectedSeqNumExec, startBlocks)
}
