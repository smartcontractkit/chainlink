package ccip

import (
	"math/big"
	"slices"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_2_0/router"
)

func Test_CCIPMessageLimitations(t *testing.T) {
	ctx := testcontext.Get(t)
	callOpts := &bind.CallOpts{Context: ctx}

	testEnv, _, _ := testsetups.NewIntegrationEnvironment(t)
	chains := maps.Keys(testEnv.Env.Chains)

	onChainState, err := changeset.LoadOnchainState(testEnv.Env)
	require.NoError(t, err)

	testhelpers.AddLanesForAll(t, &testEnv, onChainState)

	srcToken, _ := setupTokens(
		t,
		onChainState,
		testEnv,
		chains[0],
		chains[1],
		deployment.E18Mult(10_000),
		deployment.E18Mult(10_000),
	)

	chain0DestConfig, err := onChainState.Chains[chains[0]].FeeQuoter.GetDestChainConfig(callOpts, chains[1])
	require.NoError(t, err)
	t.Logf("0->1 destination config: %+v", chain0DestConfig)

	testMsgs := []struct {
		name      string
		fromChain uint64
		toChain   uint64
		msg       router.ClientEVM2AnyMessage
		expRevert bool
	}{
		{
			name:      "hit limit on data",
			fromChain: chains[0],
			toChain:   chains[1],
			msg: router.ClientEVM2AnyMessage{
				Receiver: common.LeftPadBytes(onChainState.Chains[chains[1]].Receiver.Address().Bytes(), 32),
				Data:     []byte(strings.Repeat("0", int(chain0DestConfig.MaxDataBytes))),
				FeeToken: common.HexToAddress("0x0"),
			},
		},
		{
			name:      "hit limit on tokens",
			fromChain: chains[0],
			toChain:   chains[1],
			msg: router.ClientEVM2AnyMessage{
				Receiver: common.LeftPadBytes(onChainState.Chains[chains[1]].Receiver.Address().Bytes(), 32),
				TokenAmounts: slices.Repeat([]router.ClientEVMTokenAmount{
					{Token: srcToken.Address(), Amount: big.NewInt(1)},
				}, int(chain0DestConfig.MaxNumberOfTokensPerMsg)),
				FeeToken: common.HexToAddress("0x0"),
			},
		},
		{
			name:      "hit limit on gas limit",
			fromChain: chains[0],
			toChain:   chains[1],
			msg: router.ClientEVM2AnyMessage{
				Receiver:  common.LeftPadBytes(onChainState.Chains[chains[1]].Receiver.Address().Bytes(), 32),
				Data:      []byte(strings.Repeat("0", int(chain0DestConfig.MaxDataBytes))),
				FeeToken:  common.HexToAddress("0x0"),
				ExtraArgs: testhelpers.MakeEVMExtraArgsV2(uint64(chain0DestConfig.MaxPerMsgGasLimit), true),
			},
		},
		//{ // TODO: exec plugin never executed this message. CCIP-4471
		//	name:      "hit limit on maxDataBytes, tokens, gasLimit should succeed",
		//	fromChain: chains[0],
		//	toChain:   chains[1],
		//	msg: router.ClientEVM2AnyMessage{
		//		Receiver: common.LeftPadBytes(onChainState.Chains[chains[1]].Receiver.Address().Bytes(), 32),
		//		Data:     []byte(strings.Repeat("0", int(chain0DestConfig.MaxDataBytes))),
		//		TokenAmounts: slices.Repeat([]router.ClientEVMTokenAmount{
		//			{Token: srcToken.Address(), Amount: big.NewInt(1)},
		//		}, int(chain0DestConfig.MaxNumberOfTokensPerMsg)),
		//		FeeToken:  common.HexToAddress("0x0"),
		//		ExtraArgs: changeset.MakeEVMExtraArgsV2(uint64(chain0DestConfig.MaxPerMsgGasLimit), true),
		//	},
		//},
		{
			name:      "exceeding maxDataBytes",
			fromChain: chains[0],
			toChain:   chains[1],
			msg: router.ClientEVM2AnyMessage{
				Receiver:     common.LeftPadBytes(onChainState.Chains[chains[1]].Receiver.Address().Bytes(), 32),
				Data:         []byte(strings.Repeat("0", int(chain0DestConfig.MaxDataBytes)+1)),
				TokenAmounts: []router.ClientEVMTokenAmount{},
				FeeToken:     common.HexToAddress("0x0"),
				ExtraArgs:    nil,
			},
			expRevert: true,
		},
		{
			name:      "exceeding number of tokens",
			fromChain: chains[0],
			toChain:   chains[1],
			msg: router.ClientEVM2AnyMessage{
				Receiver: common.LeftPadBytes(onChainState.Chains[chains[1]].Receiver.Address().Bytes(), 32),
				Data:     []byte("abc"),
				TokenAmounts: slices.Repeat([]router.ClientEVMTokenAmount{
					{Token: srcToken.Address(), Amount: big.NewInt(1)},
				}, int(chain0DestConfig.MaxNumberOfTokensPerMsg)+1),
				FeeToken:  common.HexToAddress("0x0"),
				ExtraArgs: nil,
			},
			expRevert: true,
		},
		{
			name:      "exceeding gas limit",
			fromChain: chains[0],
			toChain:   chains[1],
			msg: router.ClientEVM2AnyMessage{
				Receiver:     common.LeftPadBytes(onChainState.Chains[chains[1]].Receiver.Address().Bytes(), 32),
				Data:         []byte("abc"),
				TokenAmounts: []router.ClientEVMTokenAmount{},
				FeeToken:     common.HexToAddress("0x0"),
				ExtraArgs:    testhelpers.MakeEVMExtraArgsV2(uint64(chain0DestConfig.MaxPerMsgGasLimit)+1, true),
			},
			expRevert: true,
		},
	}

	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	// Send a message from each chain to every other chain.
	expectedSeqNum := make(map[testhelpers.SourceDestPair]uint64)
	expectedSeqNumExec := make(map[testhelpers.SourceDestPair][]uint64)
	for _, msg := range testMsgs {
		t.Logf("Sending msg: %s", msg.name)
		require.NotEqual(t, msg.fromChain, msg.toChain, "fromChain and toChain cannot be the same")
		startBlocks[msg.toChain] = nil
		msgSentEvent, err := testhelpers.DoSendRequest(
			t, testEnv.Env, onChainState,
			testhelpers.WithSourceChain(msg.fromChain),
			testhelpers.WithDestChain(msg.toChain),
			testhelpers.WithTestRouter(false),
			testhelpers.WithEvm2AnyMessage(msg.msg))

		if msg.expRevert {
			t.Logf("Message reverted as expected")
			require.Error(t, err)
			require.Contains(t, err.Error(), "execution reverted")
			continue
		}
		require.NoError(t, err)

		t.Logf("Message not reverted as expected")

		expectedSeqNum[testhelpers.SourceDestPair{
			SourceChainSelector: msg.fromChain,
			DestChainSelector:   msg.toChain,
		}] = msgSentEvent.SequenceNumber

		expectedSeqNumExec[testhelpers.SourceDestPair{
			SourceChainSelector: msg.fromChain,
			DestChainSelector:   msg.toChain,
		}] = []uint64{msgSentEvent.SequenceNumber}
	}

	// Wait for all commit reports to land.
	testhelpers.ConfirmCommitForAllWithExpectedSeqNums(t, testEnv.Env, onChainState, expectedSeqNum, startBlocks)
	// Wait for all exec reports to land
	testhelpers.ConfirmExecWithSeqNrsForAll(t, testEnv.Env, onChainState, expectedSeqNumExec, startBlocks)
}
