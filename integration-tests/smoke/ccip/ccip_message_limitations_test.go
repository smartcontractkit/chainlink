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

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	mlt "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/messagelimitationstest"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

func Test_CCIPMessageLimitations(t *testing.T) {
	ctx := testcontext.Get(t)
	callOpts := &bind.CallOpts{Context: ctx}

	testEnv, _, _ := testsetups.NewIntegrationEnvironment(t)
	chains := maps.Keys(testEnv.Env.BlockChains.EVMChains())

	onChainState, err := stateview.LoadOnchainState(testEnv.Env)
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

	chain0DestConfig, err := onChainState.MustGetEVMChainState(chains[0]).FeeQuoter.GetDestChainConfig(callOpts, chains[1])
	require.NoError(t, err)
	t.Logf("0->1 destination config: %+v", chain0DestConfig)

	testSetup := mlt.NewTestSetup(
		t,
		onChainState,
		chains[0],
		chains[1],
		srcToken.Address(),
		chain0DestConfig,
		false, // testRouter
		true,  // validateResp
		mlt.WithDeployedEnv(testEnv),
	)

	tcs := []mlt.TestCase{
		{
			TestSetup: testSetup,
			Name:      "hit limit on data",
			Msg: router.ClientEVM2AnyMessage{
				Receiver: common.LeftPadBytes(onChainState.MustGetEVMChainState(testSetup.DestChain).Receiver.Address().Bytes(), 32),
				Data:     []byte(strings.Repeat("0", int(testSetup.SrcFeeQuoterDestChainConfig.MaxDataBytes))),
				FeeToken: common.HexToAddress("0x0"),
			},
		},
		{
			TestSetup: testSetup,
			Name:      "hit limit on tokens",
			Msg: router.ClientEVM2AnyMessage{
				Receiver: common.LeftPadBytes(onChainState.MustGetEVMChainState(testSetup.DestChain).Receiver.Address().Bytes(), 32),
				TokenAmounts: slices.Repeat([]router.ClientEVMTokenAmount{
					{Token: testSetup.SrcToken, Amount: big.NewInt(1)},
				}, int(testSetup.SrcFeeQuoterDestChainConfig.MaxNumberOfTokensPerMsg)),
				FeeToken: common.HexToAddress("0x0"),
			},
		},
		{
			TestSetup: testSetup,
			Name:      "hit limit on gas limit",
			Msg: router.ClientEVM2AnyMessage{
				Receiver:  common.LeftPadBytes(onChainState.MustGetEVMChainState(testSetup.DestChain).Receiver.Address().Bytes(), 32),
				Data:      []byte(strings.Repeat("0", int(testSetup.SrcFeeQuoterDestChainConfig.MaxDataBytes))),
				FeeToken:  common.HexToAddress("0x0"),
				ExtraArgs: testhelpers.MakeEVMExtraArgsV2(uint64(testSetup.SrcFeeQuoterDestChainConfig.MaxPerMsgGasLimit), true),
			},
		},
		//{ // TODO: exec plugin never executed this message. CCIP-4471
		//	name:      "hit limit on maxDataBytes, tokens, gasLimit should succeed",
		//	fromChain: chains[0],
		//	toChain:   chains[1],
		//	msg: router.ClientEVM2AnyMessage{
		//		Receiver: common.LeftPadBytes(onChainState.MustGetEVMChainState(chains[1]].Receiver)Address().Bytes(), 32),
		//		Data:     []byte(strings.Repeat("0", int(chain0DestConfig.MaxDataBytes))),
		//		TokenAmounts: slices.Repeat([]router.ClientEVMTokenAmount{
		//			{Token: srcToken.Address(), Amount: big.NewInt(1)},
		//		}, int(chain0DestConfig.MaxNumberOfTokensPerMsg)),
		//		FeeToken:  common.HexToAddress("0x0"),
		//		ExtraArgs: changeset.MakeEVMExtraArgsV2(uint64(chain0DestConfig.MaxPerMsgGasLimit), true),
		//	},
		//},
		{
			TestSetup: testSetup,
			Name:      "exceeding maxDataBytes",
			Msg: router.ClientEVM2AnyMessage{
				Receiver:     common.LeftPadBytes(onChainState.MustGetEVMChainState(testSetup.DestChain).Receiver.Address().Bytes(), 32),
				Data:         []byte(strings.Repeat("0", int(testSetup.SrcFeeQuoterDestChainConfig.MaxDataBytes)+1)),
				TokenAmounts: []router.ClientEVMTokenAmount{},
				FeeToken:     common.HexToAddress("0x0"),
				ExtraArgs:    nil,
			},
			ExpRevert: true,
		},
		{
			TestSetup: testSetup,
			Name:      "exceeding number of tokens",
			Msg: router.ClientEVM2AnyMessage{
				Receiver: common.LeftPadBytes(onChainState.MustGetEVMChainState(testSetup.DestChain).Receiver.Address().Bytes(), 32),
				Data:     []byte("abc"),
				TokenAmounts: slices.Repeat([]router.ClientEVMTokenAmount{
					{Token: testSetup.SrcToken, Amount: big.NewInt(1)},
				}, int(testSetup.SrcFeeQuoterDestChainConfig.MaxNumberOfTokensPerMsg)+1),
				FeeToken:  common.HexToAddress("0x0"),
				ExtraArgs: nil,
			},
			ExpRevert: true,
		},
		{
			TestSetup: testSetup,
			Name:      "exceeding gas limit",
			Msg: router.ClientEVM2AnyMessage{
				Receiver:     common.LeftPadBytes(onChainState.MustGetEVMChainState(testSetup.DestChain).Receiver.Address().Bytes(), 32),
				Data:         []byte("abc"),
				TokenAmounts: []router.ClientEVMTokenAmount{},
				FeeToken:     common.HexToAddress("0x0"),
				ExtraArgs:    testhelpers.MakeEVMExtraArgsV2(uint64(testSetup.SrcFeeQuoterDestChainConfig.MaxPerMsgGasLimit)+1, true),
			},
			ExpRevert: true,
		},
	}

	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	// Send a message from each chain to every other chain.
	expectedSeqNum := make(map[testhelpers.SourceDestPair]uint64)
	expectedSeqNumExec := make(map[testhelpers.SourceDestPair][]uint64)
	for _, tc := range tcs {
		startBlocks[tc.DestChain] = nil

		tco := mlt.Run(tc)

		if tco.MsgSentEvent != nil {
			expectedSeqNum[testhelpers.SourceDestPair{
				SourceChainSelector: tc.SrcChain,
				DestChainSelector:   tc.DestChain,
			}] = tco.MsgSentEvent.SequenceNumber

			expectedSeqNumExec[testhelpers.SourceDestPair{
				SourceChainSelector: tc.SrcChain,
				DestChainSelector:   tc.DestChain,
			}] = []uint64{tco.MsgSentEvent.SequenceNumber}
		}
	}

	// Wait for all commit reports to land.
	testhelpers.ConfirmCommitForAllWithExpectedSeqNums(t, testEnv.Env, onChainState,
		testhelpers.ToSeqRangeMap(expectedSeqNum), startBlocks)
	// Wait for all exec reports to land
	testhelpers.ConfirmExecWithSeqNrsForAll(t, testEnv.Env, onChainState, expectedSeqNumExec, startBlocks)
}
