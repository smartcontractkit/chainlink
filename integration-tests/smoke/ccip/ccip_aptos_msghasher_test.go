package ccip

import (
	"context"
	"math/big"
	"testing"

	"github.com/aptos-labs/aptos-go-sdk"
	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"

	ccipocr3common "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipaptos"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	aptos_call_opts "github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	aptos_ccip_offramp "github.com/smartcontractkit/chainlink-aptos/bindings/ccip_offramp/offramp"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	aptosstate "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/aptos"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

// Test_CCIP_AptosMessageHasher_OnChainVerification compares off-chain aptos msghasher.go implementation
// with on-chain Aptos Move offramp::calculate_message_hash()
func Test_CCIP_AptosMessageHasher_OnChainVerification(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := testhelpers.Context(t)

	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithAptosChains(1),
	)

	// Deploy CCIP contracts and load state
	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	// Get chain selectors
	evmChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyEVM))
	aptosChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyAptos))

	sourceChain := evmChainSelectors[0] // EVM source
	destChain := aptosChainSelectors[0] // Aptos destination

	// Setup off-chain message hasher
	extraDataCodec := ccipcommon.ExtraDataCodec(map[string]ccipcommon.SourceChainExtraDataCodec{
		chain_selectors.FamilyAptos: ccipaptos.ExtraDataDecoder{},
		chain_selectors.FamilyEVM:   ccipevm.ExtraDataDecoder{},
	})
	msgHasher := ccipaptos.NewMessageHasherV1(lggr, extraDataCodec)

	// Get deployed contract addresses
	ccipChainState := state.AptosChains[destChain]

	t.Run("EVM_to_Aptos_BasicMessage", func(t *testing.T) {
		msg := createBasicEVMToAptosMessage(t, sourceChain, destChain)
		verifyHashMatches(ctx, t, msgHasher, ccipChainState, msg, e)
	})

	t.Run("EVM_to_Aptos_WithTokens", func(t *testing.T) {
		msg := createEVMToAptosMessageWithTokens(t, sourceChain, destChain)
		verifyHashMatches(ctx, t, msgHasher, ccipChainState, msg, e)
	})

	t.Run("EVM_to_Aptos_EmptyData", func(t *testing.T) {
		msg := createEVMToAptosMessageWithEmptyData(t, sourceChain, destChain)
		verifyHashMatches(ctx, t, msgHasher, ccipChainState, msg, e)
	})

	t.Run("EVM_to_Aptos_LargeData_3KB", func(t *testing.T) {
		msg := createEVMToAptosMessageWithLargeData(t, sourceChain, destChain, 3000)
		verifyHashMatches(ctx, t, msgHasher, ccipChainState, msg, e)
	})
}

func verifyHashMatches(
	ctx context.Context,
	t *testing.T,
	msgHasher ccipocr3common.MessageHasher,
	ccipChainState aptosstate.CCIPChainState,
	msg ccipocr3common.Message,
	e testhelpers.DeployedEnv,
) {
	// Compute off-chain hash using Go implementation
	offChainHash, err := msgHasher.Hash(ctx, msg)
	require.NoError(t, err, "Off-chain hash computation failed")

	// Compute on-chain hash using Aptos Move contract
	onChainHash := computeOnChainHash(t, ccipChainState, msg, e)

	require.Equal(t, onChainHash[:], offChainHash[:],
		"On-chain and off-chain hash mismatch! \n"+
			"On-chain:  %s\n"+
			"Off-chain: %s\n"+
			"Message: %+v",
		hexutil.Encode(onChainHash[:]),
		hexutil.Encode(offChainHash[:]),
		msg)

	t.Logf("✓ Hash verification passed")
	t.Logf("  Onchain Hash: %s", hexutil.Encode(onChainHash[:]))
	t.Logf("  Offchain Hash: %s", hexutil.Encode(offChainHash[:]))
}

func computeOnChainHash(
	t *testing.T,
	ccipChainState aptosstate.CCIPChainState,
	msg ccipocr3common.Message,
	e testhelpers.DeployedEnv,
) [32]byte {
	destChain := uint64(msg.Header.DestChainSelector)

	aptosChain, exists := e.Env.BlockChains.AptosChains()[destChain]
	require.True(t, exists, "Aptos chain not found in dest (%d)", destChain)

	aptosClient := aptosChain.Client
	ccipAddr := ccipChainState.CCIPAddress
	offramp := aptos_ccip_offramp.NewOfframp(ccipAddr, aptosClient)
	gasLimit := parseGasLimitFromExtraArgs(msg.ExtraArgs)

	sourcePoolAddresses := make([][]byte, len(msg.TokenAmounts))
	destTokenAddresses := make([]aptos.AccountAddress, len(msg.TokenAmounts))
	destGasAmounts := make([]uint32, len(msg.TokenAmounts))
	extraDatas := make([][]byte, len(msg.TokenAmounts))
	amounts := make([]*big.Int, len(msg.TokenAmounts))

	for i, token := range msg.TokenAmounts {
		sourcePoolAddresses[i] = token.SourcePoolAddress
		var addr aptos.AccountAddress
		copy(addr[:], token.DestTokenAddress)
		destTokenAddresses[i] = addr
		destGasAmounts[i] = parseDestGasAmount(token.DestExecData)
		extraDatas[i] = token.ExtraData
		amounts[i] = token.Amount.Int
	}

	var receiver aptos.AccountAddress
	copy(receiver[:], msg.Receiver)

	result, err := offramp.CalculateMessageHash(
		&aptos_call_opts.CallOpts{},
		msg.Header.MessageID[:],
		uint64(msg.Header.SourceChainSelector),
		uint64(msg.Header.DestChainSelector),
		uint64(msg.Header.SequenceNumber),
		msg.Header.Nonce,
		msg.Sender,
		receiver,
		msg.Header.OnRamp,
		msg.Data,
		gasLimit,
		sourcePoolAddresses,
		destTokenAddresses,
		destGasAmounts,
		extraDatas,
		amounts,
	)
	require.NoError(t, err, "On chain offramp::calculate_message_hash() failed")

	var hash [32]byte
	copy(hash[:], result)
	return hash
}

func createBasicEVMToAptosMessage(t *testing.T, sourceChain, destChain uint64) ccipocr3common.Message {
	messageIDBytes := hexutil.MustDecode("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	var messageID ccipocr3common.Bytes32
	copy(messageID[:], messageIDBytes)

	onRampBytes := hexutil.MustDecode("0x47a1f0a819457f01153f35c6b6b0d42e2e16e91e")
	senderBytes := hexutil.MustDecode("0xd87929a32cf0cbdc9e2d07ffc7c33344079de727")
	receiverBytes := hexutil.MustDecode("0xbd8a1fb0af25dc8700d2d302cfbae718c3b2c3c61cfe47f58a45b1126c006490")

	extraArgs := testhelpers.MakeEVMExtraArgsV2(500000, true)

	return ccipocr3common.Message{
		Header: ccipocr3common.RampMessageHeader{
			MessageID:           messageID,
			SourceChainSelector: ccipocr3common.ChainSelector(sourceChain),
			DestChainSelector:   ccipocr3common.ChainSelector(destChain),
			SequenceNumber:      ccipocr3common.SeqNum(42),
			Nonce:               123,
			OnRamp:              onRampBytes,
		},
		Sender:       senderBytes,
		Receiver:     receiverBytes,
		Data:         []byte("hello CCIPReceiver"),
		ExtraArgs:    extraArgs,
		TokenAmounts: []ccipocr3common.RampTokenAmount{},
	}
}

func createEVMToAptosMessageWithTokens(t *testing.T, sourceChain, destChain uint64) ccipocr3common.Message {
	msg := createBasicEVMToAptosMessage(t, sourceChain, destChain)

	srcPool1 := hexutil.MustDecode("0xabcdef1234567890abcdef1234567890abcdef12")
	destToken1Bytes := hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000005678")
	extraData1 := hexutil.MustDecode("0x00112233")

	srcPool2 := hexutil.MustDecode("0x123456789abcdef123456789abcdef123456789a")
	destToken2Bytes := hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000009abc")
	extraData2 := hexutil.MustDecode("0xffeeddcc")

	destExecData1, err := utils.ABIEncode(`[{"type":"uint32"}]`, uint32(10000))
	require.NoError(t, err)
	destExecData2, err := utils.ABIEncode(`[{"type":"uint32"}]`, uint32(20000))
	require.NoError(t, err)

	msg.TokenAmounts = []ccipocr3common.RampTokenAmount{
		{
			SourcePoolAddress: srcPool1,
			DestTokenAddress:  destToken1Bytes,
			ExtraData:         extraData1,
			Amount:            ccipocr3common.NewBigInt(big.NewInt(1000000)),
			DestExecData:      destExecData1,
		},
		{
			SourcePoolAddress: srcPool2,
			DestTokenAddress:  destToken2Bytes,
			ExtraData:         extraData2,
			Amount:            ccipocr3common.NewBigInt(big.NewInt(5000000)),
			DestExecData:      destExecData2,
		},
	}
	return msg
}

func createEVMToAptosMessageWithEmptyData(t *testing.T, sourceChain, destChain uint64) ccipocr3common.Message {
	msg := createBasicEVMToAptosMessage(t, sourceChain, destChain)
	msg.Data = []byte{}
	return msg
}

func createEVMToAptosMessageWithLargeData(t *testing.T, sourceChain, destChain uint64, size int) ccipocr3common.Message {
	msg := createBasicEVMToAptosMessage(t, sourceChain, destChain)
	msg.Data = make([]byte, size)
	for i := range msg.Data {
		msg.Data[i] = byte(i % 256)
	}
	return msg
}

func parseGasLimitFromExtraArgs(extraArgs []byte) *big.Int {
	evmDecoder := ccipevm.ExtraDataDecoder{}
	if decodedMap, err := evmDecoder.DecodeExtraArgsToMap(extraArgs); err == nil {
		if gasLimit, exists := decodedMap["gasLimit"]; exists {
			if gl, ok := gasLimit.(*big.Int); ok {
				return gl
			}
		}
	}
	return big.NewInt(200000)
}

func parseDestGasAmount(destExecData []byte) uint32 {
	evmDecoder := ccipevm.ExtraDataDecoder{}
	if decodedMap, err := evmDecoder.DecodeDestExecDataToMap(destExecData); err == nil {
		if destGasAmount, exists := decodedMap["destGasAmount"]; exists {
			if gasAmount, ok := destGasAmount.(uint32); ok {
				return gasAmount
			}
		}
	}
	return 50000
}
