package ccip

import (
	"encoding/hex"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	evm_fee_quoter "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/fee_quoter"
	module_fee_quoter "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain"

	suiBind "github.com/smartcontractkit/chainlink-sui/bindings/bind"
	suiutil "github.com/smartcontractkit/chainlink-sui/bindings/utils"
	sui_deployment "github.com/smartcontractkit/chainlink-sui/deployment"
	sui_cs "github.com/smartcontractkit/chainlink-sui/deployment/changesets"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	linkops "github.com/smartcontractkit/chainlink-sui/deployment/ops/link"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	mlt "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/messagelimitationstest"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/messagingtest"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

// sui2EvmMessagingFixtures is shared setup for Sui→EVM messaging tests split for CI isolation.
type sui2EvmMessagingFixtures struct {
	e                      testhelpers.DeployedEnv
	sourceChain, destChain uint64
	state                  stateview.CCIPOnChainState
	setup                  messagingtest.TestSetup
	suiLinkFeeToken        string
	standardMessage        []byte
	suiFQDestConfig        module_fee_quoter.DestChainConfig
	nonce                  *uint64
}

func prepareSui2EvmMessagingTest(t *testing.T) sui2EvmMessagingFixtures {
	t.Helper()
	ctx := testhelpers.Context(t)
	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithSuiChains(1),
	)

	evmChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyEVM))
	suiChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilySui))

	sourceChain := suiChainSelectors[0]
	destChain := evmChainSelectors[0]

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	suiState, err := sui_deployment.LoadOnchainStatesui(e.Env)
	require.NoError(t, err)

	t.Log("Source chain (Sui): ", sourceChain, "Dest chain (EVM): ", destChain)

	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.NoError(t, err)

	suiSenderAddr, err := e.Env.BlockChains.SuiChains()[sourceChain].Signer.GetAddress()
	require.NoError(t, err)

	normalizedAddr, err := suiutil.ConvertStringToAddressBytes(suiSenderAddr)
	require.NoError(t, err)

	suiSenderByte := normalizedAddr[:]

	_, output, err := commoncs.ApplyChangesets(t, e.Env, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.MintLinkToken{}, sui_cs.MintLinkTokenConfig{
			ChainSelector:  sourceChain,
			TokenPackageId: suiState[sourceChain].LinkTokenAddress,
			TreasuryCapId:  suiState[sourceChain].LinkTokenTreasuryCapId,
			Amount:         1000000000000, // 1000 Link with 1e9
		}),
	})
	require.NoError(t, err)

	rawOutput := output[0].Reports[0]
	outputMap, ok := rawOutput.Output.(sui_ops.OpTxResult[linkops.MintLinkTokenOutput])
	require.True(t, ok)

	sender := common.LeftPadBytes(suiSenderByte, 32)
	setup := messagingtest.NewTestSetupWithDeployedEnv(
		t,
		e,
		state,
		sourceChain,
		destChain,
		sender,
		false, // testRouter
	)
	suiLinkFeeToken := outputMap.Objects.MintedLinkTokenObjectId
	standardMessage := []byte("Hello EVM, from Sui!")

	suifeeQuoter, err := module_fee_quoter.NewFeeQuoter(suiState[sourceChain].CCIPAddress, e.Env.BlockChains.SuiChains()[sourceChain].Client)
	require.NoError(t, err)

	suiFQDestConfig, err := suifeeQuoter.DevInspect().GetDestChainConfig(ctx, &suiBind.CallOpts{
		Signer:           e.Env.BlockChains.SuiChains()[sourceChain].Signer,
		WaitForExecution: true,
	}, suiBind.Object{Id: suiState[sourceChain].CCIPObjectRef}, destChain)
	require.NoError(t, err, "Failed to get destination chain config")

	n := uint64(0)
	return sui2EvmMessagingFixtures{
		e: e, sourceChain: sourceChain, destChain: destChain,
		state: state, setup: setup,
		suiLinkFeeToken: suiLinkFeeToken, standardMessage: standardMessage,
		suiFQDestConfig: suiFQDestConfig, nonce: &n,
	}
}

func Test_CCIP_Messaging_Sui2EVM_Success(t *testing.T) {
	fx := prepareSui2EvmMessagingTest(t)
	var out messagingtest.TestCaseOutput

	waitForSuiRPCSync(t, fx.e.Env.BlockChains.SuiChains()[fx.sourceChain])

	t.Run("Message to EVM", func(t *testing.T) {
		out = messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:              fx.setup,
				Nonce:                  fx.nonce,
				ValidationType:         messagingtest.ValidationTypeExec,
				Receiver:               fx.state.Chains[fx.destChain].Receiver.Address().Bytes(),
				ExtraArgs:              nil,
				Replayed:               true,
				FeeToken:               fx.suiLinkFeeToken,
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
			},
		)
	})

	waitForSuiRPCSync(t, fx.e.Env.BlockChains.SuiChains()[fx.sourceChain])

	// t.Run("Max Data Bytes - Should Succeed", func(t *testing.T) {
	// 	ctx := testhelpers.Context(t)
	// 	latestHead, err := testhelpers.LatestBlock(ctx, fx.e.Env, fx.destChain)
	// 	require.NoError(t, err)
	// 	message := []byte(strings.Repeat("0", int(16000)))
	// 	messagingtest.Run(t,
	// 		messagingtest.TestCase{
	// 			TestSetup:              fx.setup,
	// 			ValidationType:         messagingtest.ValidationTypeExec,
	// 			FeeToken:               fx.suiLinkFeeToken,
	// 			Receiver:               fx.state.Chains[fx.destChain].Receiver.Address().Bytes(),
	// 			MsgData:                message,
	// 			ExtraArgs:              testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(300000), false),
	// 			ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
	// 			ExtraAssertions: []func(t *testing.T){
	// 				func(t *testing.T) {
	// 					assertEvmMessageReceived(testhelpers.Context(t), t, fx.state, fx.destChain, latestHead, message)
	// 				},
	// 			},
	// 		},
	// 	)
	// })

	// waitForSuiRPCSync(t, fx.e.Env.BlockChains.SuiChains()[fx.sourceChain])

	// t.Run("Max Gas Limit - Should Succeed", func(t *testing.T) {
	// 	ctx := testhelpers.Context(t)
	// 	latestHead, err := testhelpers.LatestBlock(ctx, fx.e.Env, fx.destChain)
	// 	require.NoError(t, err)
	// 	messagingtest.Run(t,
	// 		messagingtest.TestCase{
	// 			TestSetup:              fx.setup,
	// 			ValidationType:         messagingtest.ValidationTypeExec,
	// 			FeeToken:               fx.suiLinkFeeToken,
	// 			Receiver:               fx.state.Chains[fx.destChain].Receiver.Address().Bytes(),
	// 			MsgData:                fx.standardMessage,
	// 			ExtraArgs:              testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(int64(fx.suiFQDestConfig.MaxPerMsgGasLimit)), false),
	// 			ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
	// 			ExtraAssertions: []func(t *testing.T){
	// 				func(t *testing.T) {
	// 					assertEvmMessageReceived(testhelpers.Context(t), t, fx.state, fx.destChain, latestHead, fx.standardMessage)
	// 				},
	// 			},
	// 		},
	// 	)
	// })

	t.Logf("out: %v\n", out)
}

// sui2EvmRevertHarness is shared mlt wiring for Sui→EVM source-revert cases (split across two CI jobs).
type sui2EvmRevertHarness struct {
	fx         sui2EvmMessagingFixtures
	mltSetup   mlt.TestSetup
	invalidSel mlt.TestSetup
}

func newSui2EvmRevertHarness(t *testing.T) sui2EvmRevertHarness {
	t.Helper()
	fx := prepareSui2EvmMessagingTest(t)
	mltSetup := mlt.NewTestSetup(
		t,
		fx.state,
		fx.sourceChain,
		fx.destChain,
		common.HexToAddress(fx.suiLinkFeeToken),
		fx.suiFQDestConfig,
		false, // testRouter
		true,  // validateResp
		mlt.WithDeployedEnv(fx.e),
	)
	invalidSel := mlt.NewTestSetup(
		t,
		fx.state,
		fx.sourceChain,
		fx.destChain,
		common.HexToAddress("0x0"),
		fx.suiFQDestConfig,
		false, // testRouter
		true,  // validateResp
		mlt.WithDeployedEnv(fx.e),
	)
	return sui2EvmRevertHarness{fx: fx, mltSetup: mltSetup, invalidSel: invalidSel}
}

func Test_CCIP_Messaging_Sui2EVM_Revert_Part1(t *testing.T) {
	h := newSui2EvmRevertHarness(t)
	fx := h.fx

	t.Run("Max Data Bytes + 1 - Should Fail", func(t *testing.T) {
		waitForSuiRPCSync(t, fx.e.Env.BlockChains.SuiChains()[fx.sourceChain])
		message := []byte(strings.Repeat("0", int(fx.suiFQDestConfig.MaxDataBytes)+1))
		mlt.Run(mlt.TestCase{
			TestSetup: h.mltSetup,
			Name:      "Max Data Bytes + 1 - Should Fail",
			Msg: testhelpers.SuiSendRequest{
				Receiver:  fx.state.Chains[fx.destChain].Receiver.Address().Bytes(),
				Data:      message,
				FeeToken:  fx.suiLinkFeeToken,
				ExtraArgs: nil,
			},
			ExpRevert: true,
		})
	})

	t.Run("Max Data Bytes + 1 to EOA - Should Fail", func(t *testing.T) {
		waitForSuiRPCSync(t, fx.e.Env.BlockChains.SuiChains()[fx.sourceChain])
		message := []byte(strings.Repeat("0", int(fx.suiFQDestConfig.MaxDataBytes)+1))
		mlt.Run(mlt.TestCase{
			TestSetup: h.mltSetup,
			Name:      "Max Data Bytes + 1 to EOA - Should Fail",
			Msg: testhelpers.SuiSendRequest{
				Receiver:  fx.state.Chains[fx.destChain].Receiver.Address().Bytes(), // Sending to EOA
				Data:      message,
				FeeToken:  fx.suiLinkFeeToken,
				ExtraArgs: nil,
			},
			ExpRevert: true,
		})
	})

	t.Run("Max Gas Limit + 1 - Should Fail", func(t *testing.T) {
		waitForSuiRPCSync(t, fx.e.Env.BlockChains.SuiChains()[fx.sourceChain])
		message := fx.standardMessage
		mlt.Run(mlt.TestCase{
			TestSetup: h.mltSetup,
			Name:      "Max Gas Limit + 1 - Should Fail",
			Msg: testhelpers.SuiSendRequest{
				Receiver:  fx.state.Chains[fx.destChain].Receiver.Address().Bytes(),
				Data:      message,
				FeeToken:  fx.suiLinkFeeToken,
				ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(int64(fx.suiFQDestConfig.MaxPerMsgGasLimit)+1), false),
			},
			ExpRevert: true,
		})
	})
}

func Test_CCIP_Messaging_Sui2EVM_Revert_Part2(t *testing.T) {
	h := newSui2EvmRevertHarness(t)
	fx := h.fx

	t.Run("Missing ExtraArgs - Should Fail", func(t *testing.T) {
		waitForSuiRPCSync(t, fx.e.Env.BlockChains.SuiChains()[fx.sourceChain])
		message := fx.standardMessage
		mlt.Run(mlt.TestCase{
			TestSetup: h.mltSetup,
			Name:      "Missing ExtraArgs - Should Fail",
			Msg: testhelpers.SuiSendRequest{
				Receiver:  fx.state.Chains[fx.destChain].Receiver.Address().Bytes(),
				Data:      message,
				FeeToken:  fx.suiLinkFeeToken,
				ExtraArgs: []byte{},
			},
			ExpRevert: true,
		})
	})

	t.Run("Send message to invalid receiver - Should Fail", func(t *testing.T) {
		waitForSuiRPCSync(t, fx.e.Env.BlockChains.SuiChains()[fx.sourceChain])
		message := fx.standardMessage
		mlt.Run(mlt.TestCase{
			TestSetup: h.mltSetup,
			Name:      "Send message to invalid receiver - Should Fail",
			Msg: testhelpers.SuiSendRequest{
				Receiver:  []byte("0x0000"),
				Data:      message,
				FeeToken:  fx.suiLinkFeeToken,
				ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(300000), false),
			},
			ExpRevert: true,
		})
	})

	t.Run("Send message to invalid chain selector - Should Fail", func(t *testing.T) {
		waitForSuiRPCSync(t, fx.e.Env.BlockChains.SuiChains()[fx.sourceChain])
		message := []byte("Hello Sui, from EVM!")
		mlt.Run(mlt.TestCase{
			TestSetup: h.invalidSel,
			Name:      "Send message to invalid chain selector - Should Fail",
			Msg: testhelpers.SuiSendRequest{
				Receiver:  fx.state.Chains[fx.destChain].Receiver.Address().Bytes(),
				Data:      message,
				FeeToken:  fx.suiLinkFeeToken,
				ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(300000), false),
			},
			ExpRevert: true,
		})
	})
}

// evm2SuiMessagingFixtures is shared setup for EVM→Sui messaging tests split for CI isolation.
type evm2SuiMessagingFixtures struct {
	e                      testhelpers.DeployedEnv
	sourceChain, destChain uint64
	state                  stateview.CCIPOnChainState
	setup                  messagingtest.TestSetup
	receiverByte           []byte
	receiverObjectIDs      [][32]byte
	srcFQDestConfig        evm_fee_quoter.FeeQuoterDestChainConfig
	nativeFeeToken         string
}

func prepareEVM2SuiMessagingTest(t *testing.T) evm2SuiMessagingFixtures {
	t.Helper()
	ctx := testcontext.Get(t)
	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithSuiChains(1),
	)

	evmChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyEVM))
	suiChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilySui))

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	sourceChain := evmChainSelectors[0]
	destChain := suiChainSelectors[0]

	t.Log("Source chain (EVM): ", sourceChain, "Dest chain (Sui): ", destChain)

	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.NoError(t, err)

	sender := common.LeftPadBytes(e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey.From.Bytes(), 32)
	setup := messagingtest.NewTestSetupWithDeployedEnv(
		t,
		e,
		state,
		sourceChain,
		destChain,
		sender,
		false, // test router
	)

	_, output, err := commoncs.ApplyChangesets(t, e.Env, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.DeployDummyReceiver{}, sui_cs.DeployDummyReceiverConfig{
			SuiChainSelector: destChain,
			McmsOwner:        "0x1",
		}),
	})
	require.NoError(t, err)

	rawOutput := output[0].Reports[0]

	outputMap, ok := rawOutput.Output.(sui_ops.OpTxResult[ccipops.DeployDummyReceiverObjects])
	require.True(t, ok)

	id := strings.TrimPrefix(outputMap.PackageId, "0x")
	receiverByteDecoded, err := hex.DecodeString(id)
	require.NoError(t, err)

	_, _, err = commoncs.ApplyChangesets(t, e.Env, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.RegisterDummyReceiver{}, sui_cs.RegisterDummyReceiverConfig{
			SuiChainSelector:       destChain,
			OwnerCapObjectId:       outputMap.Objects.OwnerCapObjectId,
			CCIPObjectRefObjectId:  state.SuiChains[destChain].CCIPObjectRef,
			DummyReceiverPackageId: outputMap.PackageId,
		}),
	})
	require.NoError(t, err)

	receiverByte := receiverByteDecoded

	var clockObj [32]byte
	copy(clockObj[:], hexutil.MustDecode(
		"0x0000000000000000000000000000000000000000000000000000000000000006",
	))

	var stateObj [32]byte
	copy(stateObj[:], hexutil.MustDecode(
		outputMap.Objects.CCIPReceiverStateObjectId,
	))

	receiverObjectIDs := [][32]byte{clockObj, stateObj}

	srcFQDestConfig, err := state.Chains[sourceChain].FeeQuoter.GetDestChainConfig(&bind.CallOpts{Context: ctx}, destChain)
	require.NoError(t, err, "Failed to get destination chain config")

	return evm2SuiMessagingFixtures{
		e: e, sourceChain: sourceChain, destChain: destChain,
		state: state, setup: setup,
		receiverByte: receiverByte, receiverObjectIDs: receiverObjectIDs,
		srcFQDestConfig: srcFQDestConfig, nativeFeeToken: "0x0",
	}
}

func Test_CCIP_Messaging_EVM2Sui_Success(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/CCIP-11130")

	fx := prepareEVM2SuiMessagingTest(t)
	var nonce uint64

	waitForSuiRPCSync(t, fx.e.Env.BlockChains.SuiChains()[fx.destChain])

	t.Run("Message to Sui", func(t *testing.T) {
		message := []byte("Hello Sui, from EVM!")
		messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:              fx.setup,
				Nonce:                  &nonce,
				ValidationType:         messagingtest.ValidationTypeExec,
				Receiver:               fx.receiverByte,
				MsgData:                message,
				ExtraArgs:              testhelpers.MakeSuiExtraArgs(1000000, true, fx.receiverObjectIDs, [32]byte{}),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
			},
		)
	})

	waitForSuiRPCSync(t, fx.e.Env.BlockChains.SuiChains()[fx.destChain])
}

// evm2SuiRevertHarness is shared mlt wiring for EVM→Sui source-revert cases (split across two CI jobs).
type evm2SuiRevertHarness struct {
	fx         evm2SuiMessagingFixtures
	mltSetup   mlt.TestSetup
	invalidSel mlt.TestSetup
}

func newEVM2SuiRevertHarness(t *testing.T) evm2SuiRevertHarness {
	t.Helper()
	fx := prepareEVM2SuiMessagingTest(t)
	waitForSuiRPCSync(t, fx.e.Env.BlockChains.SuiChains()[fx.destChain])
	mltSetup := mlt.NewTestSetup(
		t,
		fx.state,
		fx.sourceChain,
		fx.destChain,
		common.HexToAddress("0x0"),
		fx.srcFQDestConfig,
		false, // testRouter
		true,  // validateResp
		mlt.WithDeployedEnv(fx.e),
	)
	invalidSel := mlt.NewTestSetup(
		t,
		fx.state,
		fx.sourceChain,
		fx.destChain,
		common.HexToAddress("0x0"),
		fx.srcFQDestConfig,
		false, // testRouter
		true,  // validateResp
		mlt.WithDeployedEnv(fx.e),
	)
	return evm2SuiRevertHarness{fx: fx, mltSetup: mltSetup, invalidSel: invalidSel}
}

func Test_CCIP_Messaging_EVM2Sui_Revert_Part1(t *testing.T) {
	h := newEVM2SuiRevertHarness(t)
	fx := h.fx

	t.Run("Max Data Bytes + 1 - Should Fail", func(t *testing.T) {
		message := []byte(strings.Repeat("0", int(fx.srcFQDestConfig.MaxDataBytes)+1))
		mlt.Run(mlt.TestCase{
			TestSetup: h.mltSetup,
			Name:      "Max Data Bytes + 1 - Should Fail",
			Msg: router.ClientEVM2AnyMessage{
				Receiver:  fx.receiverByte,
				Data:      message,
				FeeToken:  common.HexToAddress(fx.nativeFeeToken),
				ExtraArgs: testhelpers.MakeSuiExtraArgs(uint64(fx.srcFQDestConfig.MaxPerMsgGasLimit+1), true, [][32]byte{}, [32]byte{}),
			},
			ExpRevert: true,
		})
	})

	t.Run("Max Data Bytes + 1 to EOA - Should Fail", func(t *testing.T) {
		message := []byte(strings.Repeat("0", int(fx.srcFQDestConfig.MaxDataBytes)+1))
		mlt.Run(mlt.TestCase{
			TestSetup: h.mltSetup,
			Name:      "Max Data Bytes + 1 to EOA - Should Fail",
			Msg: router.ClientEVM2AnyMessage{
				Receiver:  fx.receiverByte, // Sending to EOA
				Data:      message,
				FeeToken:  common.HexToAddress(fx.nativeFeeToken),
				ExtraArgs: testhelpers.MakeSuiExtraArgs(uint64(fx.srcFQDestConfig.MaxPerMsgGasLimit)+1, true, fx.receiverObjectIDs, [32]byte{}),
			},
			ExpRevert: true,
		})
	})

	t.Run("Missing ExtraArgs - Should Fail", func(t *testing.T) {
		message := []byte("Hello Sui, from EVM!")
		mlt.Run(mlt.TestCase{
			TestSetup: h.mltSetup,
			Name:      "Missing ExtraArgs - Should Fail",
			Msg: router.ClientEVM2AnyMessage{
				Receiver:  fx.receiverByte,
				Data:      message,
				FeeToken:  common.HexToAddress(fx.nativeFeeToken),
				ExtraArgs: []byte{},
			},
			ExpRevert: true,
		})
	})
}

func Test_CCIP_Messaging_EVM2Sui_Revert_Part2(t *testing.T) {
	h := newEVM2SuiRevertHarness(t)
	fx := h.fx

	t.Run("OutOfOrder Execution False - Should Fail", func(t *testing.T) {
		message := []byte("Hello Sui, from EVM!")
		mlt.Run(mlt.TestCase{
			TestSetup: h.mltSetup,
			Name:      "OutOfOrder Execution False - Should Fail",
			Msg: router.ClientEVM2AnyMessage{
				Receiver:  fx.receiverByte,
				Data:      message,
				FeeToken:  common.HexToAddress(fx.nativeFeeToken),
				ExtraArgs: testhelpers.MakeSuiExtraArgs(100000, false, [][32]byte{}, [32]byte{}),
			},
			ExpRevert: true,
		})
	})

	t.Run("Send message to invalid receiver - Should Fail", func(t *testing.T) {
		message := []byte("Hello Sui, from EVM!")
		mlt.Run(mlt.TestCase{
			TestSetup: h.mltSetup,
			Name:      "Send message to invalid receiver - Should Fail",
			Msg: router.ClientEVM2AnyMessage{
				Receiver:  []byte("0x000"),
				Data:      message,
				FeeToken:  common.HexToAddress(fx.nativeFeeToken),
				ExtraArgs: testhelpers.MakeSuiExtraArgs(100000, false, [][32]byte{}, [32]byte{}),
			},
			ExpRevert: true,
		})
	})

	t.Run("Send message to invalid chain selector - Should Fail", func(t *testing.T) {
		message := []byte("Hello Sui, from EVM!")
		mlt.Run(mlt.TestCase{
			TestSetup: h.invalidSel,
			Name:      "Send message to invalid chain selector - Should Fail",
			Msg: router.ClientEVM2AnyMessage{
				Receiver:  fx.receiverByte,
				Data:      message,
				FeeToken:  common.HexToAddress(fx.nativeFeeToken),
				ExtraArgs: testhelpers.MakeSuiExtraArgs(100000, false, [][32]byte{}, [32]byte{}),
			},
			ExpRevert: true,
		})
	})
}

func Test_CCIP_EVM2Sui_ZeroReceiver(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/CCIP-11130")
	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithSuiChains(1),
	)

	evmChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyEVM))
	suiChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilySui))

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	sourceChain := evmChainSelectors[0]
	destChain := suiChainSelectors[0]

	t.Log("Source chain (EVM): ", sourceChain, "Dest chain (Sui): ", destChain)

	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.NoError(t, err)

	var (
		nonce  uint64
		sender = common.LeftPadBytes(e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey.From.Bytes(), 32)
		setup  = messagingtest.NewTestSetupWithDeployedEnv(
			t,
			e,
			state,
			sourceChain,
			destChain,
			sender,
			false, // test router
		)
	)

	waitForSuiRPCSync(t, e.Env.BlockChains.SuiChains()[destChain])

	t.Run("Message to Sui with zero receiver", func(t *testing.T) {
		message := []byte("Hello Sui, from EVM!")
		messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:              setup,
				Nonce:                  &nonce,
				ValidationType:         messagingtest.ValidationTypeExec,
				Receiver:               []byte{},
				MsgData:                message,
				ExtraArgs:              testhelpers.MakeSuiExtraArgs(0, true, [][32]byte{}, [32]byte{}),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
			},
		)
	})
}
