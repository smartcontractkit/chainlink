package ccip

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
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

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test_CCIP_Messaging_Sui2EVM(t *testing.T) {
	ctx := testhelpers.Context(t)
	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithSuiChains(1),
	)

	evmChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyEVM))
	suiChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilySui))

	fmt.Println("EVM: ", evmChainSelectors[0])
	fmt.Println("Sui: ", suiChainSelectors[0])

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

	// SUI FeeToken
	// mint link token to use as feeToken
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

	var (
		nonce  uint64
		sender = common.LeftPadBytes(suiSenderByte, 32)
		out    messagingtest.TestCaseOutput
		setup  = messagingtest.NewTestSetupWithDeployedEnv(
			t,
			e,
			state,
			sourceChain,
			destChain,
			sender,
			false, // testRouter
		)

		suiLinkFeeToken = outputMap.Objects.MintedLinkTokenObjectId
		standardMessage = []byte("Hello EVM, from Sui!")
	)

	suifeeQuoter, err := module_fee_quoter.NewFeeQuoter(suiState[sourceChain].CCIPAddress, e.Env.BlockChains.SuiChains()[sourceChain].Client)
	require.NoError(t, err)

	suiFeeQuoterDestChainConfig, err := suifeeQuoter.DevInspect().GetDestChainConfig(ctx, &suiBind.CallOpts{
		Signer:           e.Env.BlockChains.SuiChains()[sourceChain].Signer,
		WaitForExecution: true,
	}, suiBind.Object{Id: suiState[sourceChain].CCIPObjectRef}, destChain)
	require.NoError(t, err, "Failed to get destination chain config")

	t.Run("Message to EVM", func(t *testing.T) {
		require.NoError(t, err)
		out = messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:              setup,
				Nonce:                  &nonce,
				ValidationType:         messagingtest.ValidationTypeExec,
				Receiver:               state.Chains[destChain].Receiver.Address().Bytes(),
				ExtraArgs:              nil,
				Replayed:               true,
				FeeToken:               suiLinkFeeToken,
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
			},
		)
	})

	// For testing messages that revert on source
	mltTestSetup := mlt.NewTestSetup(
		t,
		state,
		sourceChain,
		destChain,
		common.HexToAddress(suiLinkFeeToken),
		suiFeeQuoterDestChainConfig,
		false, // testRouter
		true,  // validateResp
		mlt.WithDeployedEnv(e),
	)

	invalidDestChainSelectorTestSetup := mlt.NewTestSetup(
		t,
		state,
		sourceChain,
		destChain,
		common.HexToAddress("0x0"),
		suiFeeQuoterDestChainConfig,
		false, // testRouter
		true,  // validateResp
		mlt.WithDeployedEnv(e),
	)

	t.Run("Max Data Bytes - Should Succeed", func(t *testing.T) {
		latestHead, err := testhelpers.LatestBlock(ctx, e.Env, destChain)
		require.NoError(t, err)
		message := []byte(strings.Repeat("0", int(16000)))
		messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:      setup,
				ValidationType: messagingtest.ValidationTypeExec,
				FeeToken:       suiLinkFeeToken,
				Receiver:       state.Chains[destChain].Receiver.Address().Bytes(),
				MsgData:        message,
				// Just ensuring enough gas is provided to execute the message, doesn't matter if it's way too much
				ExtraArgs:              testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(300000), false),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) { assertEvmMessageReceived(ctx, t, state, destChain, latestHead, message) },
				},
			},
		)
	})

	t.Run("Max Gas Limit - Should Succeed", func(t *testing.T) {
		latestHead, err := testhelpers.LatestBlock(ctx, e.Env, destChain)
		require.NoError(t, err)
		messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:              setup,
				ValidationType:         messagingtest.ValidationTypeExec,
				FeeToken:               suiLinkFeeToken,
				Receiver:               state.Chains[destChain].Receiver.Address().Bytes(),
				MsgData:                standardMessage,
				ExtraArgs:              testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(int64(suiFeeQuoterDestChainConfig.MaxPerMsgGasLimit)), false),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) { assertEvmMessageReceived(ctx, t, state, destChain, latestHead, standardMessage) },
				},
			},
		)
	})

	t.Run("Max Data Bytes + 1 - Should Fail", func(t *testing.T) {
		message := []byte(strings.Repeat("0", int(suiFeeQuoterDestChainConfig.MaxDataBytes)+1))
		mlt.Run(mlt.TestCase{
			TestSetup: mltTestSetup,
			Name:      "Max Data Bytes + 1 - Should Fail",
			Msg: testhelpers.SuiSendRequest{
				Receiver:  state.Chains[destChain].Receiver.Address().Bytes(),
				Data:      message,
				FeeToken:  suiLinkFeeToken,
				ExtraArgs: nil,
			},
			ExpRevert: true,
		})
	})

	t.Run("Max Data Bytes + 1 to EOA - Should Fail", func(t *testing.T) {
		message := []byte(strings.Repeat("0", int(suiFeeQuoterDestChainConfig.MaxDataBytes)+1))
		mlt.Run(mlt.TestCase{
			TestSetup: mltTestSetup,
			Name:      "Max Data Bytes + 1 to EOA - Should Fail",
			Msg: testhelpers.SuiSendRequest{
				Receiver:  state.Chains[destChain].Receiver.Address().Bytes(), // Sending to EOA
				Data:      message,
				FeeToken:  suiLinkFeeToken,
				ExtraArgs: nil,
			},
			ExpRevert: true,
		})
	})

	t.Run("Max Gas Limit + 1 - Should Fail", func(t *testing.T) {
		message := standardMessage
		mlt.Run(mlt.TestCase{
			TestSetup: mltTestSetup,
			Name:      "Max Gas Limit + 1 - Should Fail",
			Msg: testhelpers.SuiSendRequest{
				Receiver:  state.Chains[destChain].Receiver.Address().Bytes(),
				Data:      message,
				FeeToken:  suiLinkFeeToken,
				ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(int64(suiFeeQuoterDestChainConfig.MaxPerMsgGasLimit)+1), false),
			},
			ExpRevert: true,
		})
	})

	t.Run("Missing ExtraArgs - Should Fail", func(t *testing.T) {
		message := standardMessage
		mlt.Run(mlt.TestCase{
			TestSetup: mltTestSetup,
			Name:      "Missing ExtraArgs - Should Fail",
			Msg: testhelpers.SuiSendRequest{
				Receiver:  state.Chains[destChain].Receiver.Address().Bytes(),
				Data:      message,
				FeeToken:  suiLinkFeeToken,
				ExtraArgs: []byte{},
			},
			ExpRevert: true,
		})
	})

	t.Run("Send message to invalid receiver - Should Fail", func(t *testing.T) {
		message := standardMessage
		mlt.Run(mlt.TestCase{
			TestSetup: mltTestSetup,
			Name:      "Send message to invalid receiver - Should Fail",
			Msg: testhelpers.SuiSendRequest{
				Receiver:  []byte("0x0000"),
				Data:      message,
				FeeToken:  suiLinkFeeToken,
				ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(300000), false),
			},
			ExpRevert: true,
		})
	})

	t.Run("Send message to invalid chain selector - Should Fail", func(t *testing.T) {
		message := []byte("Hello Sui, from EVM!")
		mlt.Run(mlt.TestCase{
			TestSetup: invalidDestChainSelectorTestSetup,
			Name:      "Send message to invalid chain selector - Should Fail",
			Msg: testhelpers.SuiSendRequest{
				Receiver:  state.Chains[destChain].Receiver.Address().Bytes(),
				Data:      message,
				FeeToken:  suiLinkFeeToken,
				ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(300000), false),
			},
			ExpRevert: true,
		})
	})

	fmt.Printf("out: %v\n", out)
}

func Test_CCIP_Messaging_EVM2Sui(t *testing.T) {
	lggr := logger.TestLogger(t)
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

	lggr.Debug("Source chain (EVM): ", sourceChain, "Dest chain (Sui): ", destChain)

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
		nativeFeeToken = "0x0"
	)

	// Deploy SUI Receiver
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

	// register the receiver
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

	srcFeeQuoterDestChainConfig, err := state.Chains[sourceChain].FeeQuoter.GetDestChainConfig(&bind.CallOpts{Context: ctx}, destChain)
	require.NoError(t, err, "Failed to get destination chain config")

	t.Run("Message to Sui", func(t *testing.T) {
		// ccipChainState := state.SuiChains[destChain]
		message := []byte("Hello Sui, from EVM!")
		messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:              setup,
				Nonce:                  &nonce,
				ValidationType:         messagingtest.ValidationTypeExec,
				Receiver:               receiverByte,
				MsgData:                message,
				ExtraArgs:              testhelpers.MakeSuiExtraArgs(1000000, true, receiverObjectIDs, [32]byte{}),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
			},
		)
	})

	// TODO: consider using this for single commit with multiple report
	// tcs := []testhelpers.TestTransferRequest{
	// 	{
	// 		Name:           "Message to Sui (valid receiver)",
	// 		SourceChain:    sourceChain,
	// 		DestChain:      destChain,
	// 		Receiver:       receiverByte,
	// 		Data:           []byte("Hello Sui, from EVM!"),
	// 		ExtraArgs:      testhelpers.MakeSuiExtraArgs(1000000, true, receiverObjectIDs, [32]byte{}),
	// 		ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
	// 	},
	// 	{
	// 		Name:           "Message to Sui (zero receiver)",
	// 		SourceChain:    sourceChain,
	// 		DestChain:      destChain,
	// 		Receiver:       []byte{},
	// 		Data:           []byte("Hello Sui, from EVM!"),
	// 		ExtraArgs:      testhelpers.MakeSuiExtraArgs(0, true, [][32]byte{}, [32]byte{}),
	// 		ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
	// 	},
	// }

	// SUI MaxDataBytes won't exactly be srcFeeQuoterDestChainConfig.MaxDataBytes because we add following additional overhead;
	//  suiExpandedDataLength +=
	// ((receiverObjectIdsLength + Client.SUI_MESSAGING_ACCOUNTS_OVERHEAD) * Client.SUI_ACCOUNT_BYTE_SIZE);
	// t.Run("Message to Sui with valid receiver with data bytes = max data bytes allowed", func(t *testing.T) {
	// 	message := []byte(strings.Repeat("0", int(16000)))
	// 	messagingtest.Run(t,
	// 		messagingtest.TestCase{
	// 			TestSetup:              setup,
	// 			Nonce:                  &nonce,
	// 			ValidationType:         messagingtest.ValidationTypeExec,
	// 			Receiver:               receiverByte,
	// 			MsgData:                message,
	// 			ExtraArgs:              testhelpers.MakeSuiExtraArgs(3000000, true, receiverObjectIDs, [32]byte{}),
	// 			ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
	// 		},
	// 	)
	// })

	// REVERT CASES

	// For testing messages that revert on source
	mltTestSetup := mlt.NewTestSetup(
		t,
		state,
		sourceChain,
		destChain,
		common.HexToAddress("0x0"),
		srcFeeQuoterDestChainConfig,
		false, // testRouter
		true,  // validateResp
		mlt.WithDeployedEnv(e),
	)

	invalidDestChainSelectorTestSetup := mlt.NewTestSetup(
		t,
		state,
		sourceChain,
		destChain,
		common.HexToAddress("0x0"),
		srcFeeQuoterDestChainConfig,
		false, // testRouter
		true,  // validateResp
		mlt.WithDeployedEnv(e),
	)

	t.Run("Max Data Bytes + 1 - Should Fail", func(t *testing.T) {
		message := []byte(strings.Repeat("0", int(srcFeeQuoterDestChainConfig.MaxDataBytes)+1))
		mlt.Run(mlt.TestCase{
			TestSetup: mltTestSetup,
			Name:      "Max Data Bytes + 1 - Should Fail",
			Msg: router.ClientEVM2AnyMessage{
				Receiver:  receiverByte,
				Data:      message,
				FeeToken:  common.HexToAddress(nativeFeeToken),
				ExtraArgs: testhelpers.MakeSuiExtraArgs(uint64(srcFeeQuoterDestChainConfig.MaxPerMsgGasLimit+1), true, [][32]byte{}, [32]byte{}),
			},
			ExpRevert: true,
		})
	})

	t.Run("Max Data Bytes + 1 to EOA - Should Fail", func(t *testing.T) {
		message := []byte(strings.Repeat("0", int(srcFeeQuoterDestChainConfig.MaxDataBytes)+1))
		mlt.Run(mlt.TestCase{
			TestSetup: mltTestSetup,
			Name:      "Max Data Bytes + 1 to EOA - Should Fail",
			Msg: router.ClientEVM2AnyMessage{
				Receiver:  receiverByte, // Sending to EOA
				Data:      message,
				FeeToken:  common.HexToAddress(nativeFeeToken),
				ExtraArgs: testhelpers.MakeSuiExtraArgs(uint64(srcFeeQuoterDestChainConfig.MaxPerMsgGasLimit)+1, true, receiverObjectIDs, [32]byte{}),
			},
			ExpRevert: true,
		})
	})

	t.Run("Missing ExtraArgs - Should Fail", func(t *testing.T) {
		message := []byte("Hello Sui, from EVM!")
		mlt.Run(mlt.TestCase{
			TestSetup: mltTestSetup,
			Name:      "Missing ExtraArgs - Should Fail",
			Msg: router.ClientEVM2AnyMessage{
				Receiver:  receiverByte,
				Data:      message,
				FeeToken:  common.HexToAddress(nativeFeeToken),
				ExtraArgs: []byte{},
			},
			ExpRevert: true,
		})
	})

	t.Run("OutOfOrder Execution False - Should Fail", func(t *testing.T) {
		message := []byte("Hello Sui, from EVM!")
		mlt.Run(mlt.TestCase{
			TestSetup: mltTestSetup,
			Name:      "OutOfOrder Execution False - Should Fail",
			Msg: router.ClientEVM2AnyMessage{
				Receiver:  receiverByte,
				Data:      message,
				FeeToken:  common.HexToAddress(nativeFeeToken),
				ExtraArgs: testhelpers.MakeSuiExtraArgs(100000, false, [][32]byte{}, [32]byte{}),
			},
			ExpRevert: true,
		})
	})

	t.Run("Send message to invalid receiver - Should Fail", func(t *testing.T) {
		message := []byte("Hello Sui, from EVM!")
		mlt.Run(mlt.TestCase{
			TestSetup: mltTestSetup,
			Name:      "Send message to invalid receiver - Should Fail",
			Msg: router.ClientEVM2AnyMessage{
				Receiver:  []byte("0x000"),
				Data:      message,
				FeeToken:  common.HexToAddress(nativeFeeToken),
				ExtraArgs: testhelpers.MakeSuiExtraArgs(100000, false, [][32]byte{}, [32]byte{}),
			},
			ExpRevert: true,
		})
	})

	t.Run("Send message to invalid chain selector - Should Fail", func(t *testing.T) {
		message := []byte("Hello Sui, from EVM!")
		mlt.Run(mlt.TestCase{
			TestSetup: invalidDestChainSelectorTestSetup,
			Name:      "Send message to invalid chain selector - Should Fail",
			Msg: router.ClientEVM2AnyMessage{
				Receiver:  receiverByte,
				Data:      message,
				FeeToken:  common.HexToAddress(nativeFeeToken),
				ExtraArgs: testhelpers.MakeSuiExtraArgs(100000, false, [][32]byte{}, [32]byte{}),
			},
			ExpRevert: true,
		})
	})
}

func Test_CCIP_EVM2Sui_ZeroReceiver(t *testing.T) {
	lggr := logger.TestLogger(t)
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

	lggr.Debug("Source chain (EVM): ", sourceChain, "Dest chain (Sui): ", destChain)

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
