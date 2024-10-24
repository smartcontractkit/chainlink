package ccipevm

import (
	"context"
	cryptorand "crypto/rand"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/message_hasher"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

// NOTE: these test cases are only EVM <-> EVM.
// Update these cases once we have non-EVM examples.
func TestMessageHasher_EVM2EVM(t *testing.T) {
	ctx := testutils.Context(t)
	d := testSetup(t)

	testCases := []evmExtraArgs{
		{version: "v1", gasLimit: big.NewInt(rand.Int63())},
		{version: "v2", gasLimit: big.NewInt(rand.Int63()), allowOOO: false},
		{version: "v2", gasLimit: big.NewInt(rand.Int63()), allowOOO: true},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("tc_%d", i), func(tt *testing.T) {
			testHasherEVM2EVM(ctx, tt, d, tc)
		})
	}
}

func testHasherEVM2EVM(ctx context.Context, t *testing.T, d *testSetupData, evmExtraArgs evmExtraArgs) {
	ccipMsg := createEVM2EVMMessage(t, d.contract, evmExtraArgs)

	var tokenAmounts []message_hasher.InternalAny2EVMTokenTransfer
	for _, rta := range ccipMsg.TokenAmounts {
		destGasAmount, err := abiDecodeUint32(rta.DestExecData)
		require.NoError(t, err)

		tokenAmounts = append(tokenAmounts, message_hasher.InternalAny2EVMTokenTransfer{
			SourcePoolAddress: rta.SourcePoolAddress,
			DestTokenAddress:  common.BytesToAddress(rta.DestTokenAddress),
			ExtraData:         rta.ExtraData[:],
			Amount:            rta.Amount.Int,
			DestGasAmount:     destGasAmount,
		})
	}
	evmMsg := message_hasher.InternalAny2EVMRampMessage{
		Header: message_hasher.InternalRampMessageHeader{
			MessageId:           ccipMsg.Header.MessageID,
			SourceChainSelector: uint64(ccipMsg.Header.SourceChainSelector),
			DestChainSelector:   uint64(ccipMsg.Header.DestChainSelector),
			SequenceNumber:      uint64(ccipMsg.Header.SequenceNumber),
			Nonce:               ccipMsg.Header.Nonce,
		},
		Sender:       ccipMsg.Sender,
		Receiver:     common.BytesToAddress(ccipMsg.Receiver),
		GasLimit:     evmExtraArgs.gasLimit,
		Data:         ccipMsg.Data,
		TokenAmounts: tokenAmounts,
	}

	expectedHash, err := d.contract.Hash(&bind.CallOpts{Context: ctx}, evmMsg, ccipMsg.Header.OnRamp)
	require.NoError(t, err)

	evmMsgHasher := NewMessageHasherV1(logger.Test(t))
	actualHash, err := evmMsgHasher.Hash(ctx, ccipMsg)
	require.NoError(t, err)

	require.Equal(t, fmt.Sprintf("%x", expectedHash), strings.TrimPrefix(actualHash.String(), "0x"))
}

type evmExtraArgs struct {
	version  string
	gasLimit *big.Int
	allowOOO bool
}

func createEVM2EVMMessage(t *testing.T, messageHasher *message_hasher.MessageHasher, evmExtraArgs evmExtraArgs) cciptypes.Message {
	messageID := utils.RandomBytes32()

	sourceTokenData := make([]byte, rand.Intn(2048))
	_, err := cryptorand.Read(sourceTokenData)
	require.NoError(t, err)

	sourceChain := rand.Uint64()
	seqNum := rand.Uint64()
	nonce := rand.Uint64()
	destChain := rand.Uint64()

	var extraArgsBytes []byte
	if evmExtraArgs.version == "v1" {
		extraArgsBytes, err = messageHasher.EncodeEVMExtraArgsV1(nil, message_hasher.ClientEVMExtraArgsV1{
			GasLimit: evmExtraArgs.gasLimit,
		})
		require.NoError(t, err)
	} else if evmExtraArgs.version == "v2" {
		extraArgsBytes, err = messageHasher.EncodeEVMExtraArgsV2(nil, message_hasher.ClientEVMExtraArgsV2{
			GasLimit:                 evmExtraArgs.gasLimit,
			AllowOutOfOrderExecution: evmExtraArgs.allowOOO,
		})
		require.NoError(t, err)
	} else {
		require.FailNowf(t, "unknown extra args version", "version: %s", evmExtraArgs.version)
	}

	messageData := make([]byte, rand.Intn(2048))
	_, err = cryptorand.Read(messageData)
	require.NoError(t, err)

	numTokens := rand.Intn(10)
	var sourceTokenDatas [][]byte
	for i := 0; i < numTokens; i++ {
		sourceTokenDatas = append(sourceTokenDatas, sourceTokenData)
	}

	var tokenAmounts []cciptypes.RampTokenAmount
	for i := 0; i < len(sourceTokenDatas); i++ {
		extraData := utils.RandomBytes32()
		encodedDestExecData, err := utils.ABIEncode(`[{ "type": "uint32" }]`, rand.Uint32())
		require.NoError(t, err)
		tokenAmounts = append(tokenAmounts, cciptypes.RampTokenAmount{
			SourcePoolAddress: abiEncodedAddress(t),
			DestTokenAddress:  abiEncodedAddress(t),
			ExtraData:         extraData[:],
			Amount:            cciptypes.NewBigInt(big.NewInt(0).SetUint64(rand.Uint64())),
			DestExecData:      encodedDestExecData,
		})
	}

	return cciptypes.Message{
		Header: cciptypes.RampMessageHeader{
			MessageID:           messageID,
			SourceChainSelector: cciptypes.ChainSelector(sourceChain),
			DestChainSelector:   cciptypes.ChainSelector(destChain),
			SequenceNumber:      cciptypes.SeqNum(seqNum),
			Nonce:               nonce,
			OnRamp:              abiEncodedAddress(t),
		},
		Sender:         abiEncodedAddress(t),
		Receiver:       abiEncodedAddress(t),
		Data:           messageData,
		TokenAmounts:   tokenAmounts,
		FeeToken:       abiEncodedAddress(t),
		FeeTokenAmount: cciptypes.NewBigInt(big.NewInt(0).SetUint64(rand.Uint64())),
		ExtraArgs:      extraArgsBytes,
	}
}

func abiEncodedAddress(t *testing.T) []byte {
	addr := utils.RandomAddress()
	encoded, err := utils.ABIEncode(`[{"type": "address"}]`, addr)
	require.NoError(t, err)
	return encoded
}

type testSetupData struct {
	contractAddr common.Address
	contract     *message_hasher.MessageHasher
	sb           *backends.SimulatedBackend
	auth         *bind.TransactOpts
}

func testSetup(t *testing.T) *testSetupData {
	transactor := testutils.MustNewSimTransactor(t)
	simulatedBackend := backends.NewSimulatedBackend(core.GenesisAlloc{
		transactor.From: {Balance: assets.Ether(1000).ToInt()},
	}, 30e6)

	// Deploy the contract
	address, _, _, err := message_hasher.DeployMessageHasher(transactor, simulatedBackend)
	require.NoError(t, err)
	simulatedBackend.Commit()

	// Setup contract client
	contract, err := message_hasher.NewMessageHasher(address, simulatedBackend)
	require.NoError(t, err)

	return &testSetupData{
		contractAddr: address,
		contract:     contract,
		sb:           simulatedBackend,
		auth:         transactor,
	}
}

func TestMessagerHasher_againstRmnSharedVector(t *testing.T) {
	const (
		messageID           = "c6f553ab71282f01324bbdbcc82e22a7e66efbcd108881ecc4cdbd728aed9b1e"
		onRampAddress       = "0000000000000000000000007a2088a1bfc9d81c55368ae168c2c02570cb814f"
		dataField           = "68656c6c6f"
		receiverAddress     = "677df0cb865368207999f2862ece576dc56d8df6"
		extraArgs           = "181dcf100000000000000000000000000000000000000000000000000000000000030d400000000000000000000000000000000000000000000000000000000000000000"
		senderAddress       = "f39fd6e51aad88f6f4ce6ab8827279cfffb92266"
		feeToken            = "9fe46736679d2d9a65f0992f2272de9f3c7fa6e0"
		sourceChainSelector = 3379446385462418246
		destChainSelector   = 12922642891491394802
		expectedMsgHash     = "0x1c61fef7a3dd153943419c1101031316ed7b7a3d75913c34cbe8628033f5924f"
	)

	h := NewMessageHasherV1(logger.Test(t))
	msgH, err := h.Hash(context.Background(), cciptypes.Message{
		Header: cciptypes.RampMessageHeader{
			MessageID:           cciptypes.Bytes32(common.Hex2Bytes(messageID)),
			SourceChainSelector: sourceChainSelector,
			DestChainSelector:   destChainSelector,
			SequenceNumber:      1,
			Nonce:               1,
			MsgHash:             cciptypes.Bytes32{},
			OnRamp:              common.HexToAddress(onRampAddress).Bytes(),
		},
		Sender:       common.HexToAddress(senderAddress).Bytes(),
		Data:         common.Hex2Bytes(dataField),
		Receiver:     common.Hex2Bytes(receiverAddress),
		ExtraArgs:    common.Hex2Bytes(extraArgs),
		FeeToken:     common.HexToAddress(feeToken).Bytes(),
		TokenAmounts: []cciptypes.RampTokenAmount{},
	})
	require.NoError(t, err)
	require.Equal(t, expectedMsgHash, msgH.String())
}
