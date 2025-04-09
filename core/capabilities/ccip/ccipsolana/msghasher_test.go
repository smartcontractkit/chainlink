package ccipsolana

import (
	cryptorand "crypto/rand"
	"math/big"
	"math/rand"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common/mocks"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/config"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/ccip"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

var ExtraDataCodec = ccipcommon.NewExtraDataCodec(ccipcommon.NewExtraDataCodecParams(ccipevm.ExtraDataDecoder{}, ExtraDataDecoder{}))

func TestMessageHasher_EVM2SVM(t *testing.T) {
	any2AnyMsg, any2SolanaMsg, msgAccounts := createEVM2SolanaMessages(t)
	msgHasher := NewMessageHasherV1(logger.Test(t), ExtraDataCodec)
	actualHash, err := msgHasher.Hash(testutils.Context(t), any2AnyMsg)
	require.NoError(t, err)
	expectedHash, err := ccip.HashAnyToSVMMessage(any2SolanaMsg, any2AnyMsg.Header.OnRamp, msgAccounts)
	require.NoError(t, err)
	require.Equal(t, expectedHash, actualHash[:32])
}

func TestMessageHasher_InvalidReceiver(t *testing.T) {
	any2AnyMsg, _, _ := createEVM2SolanaMessages(t)

	// Set receiver to a []byte of 2 length
	any2AnyMsg.Receiver = []byte{0, 0}
	mockExtraDataCodec := &mocks.ExtraDataCodec{}
	mockExtraDataCodec.On("DecodeTokenAmountDestExecData", mock.Anything, mock.Anything).Return(map[string]any{
		"destGasAmount": uint32(10),
	}, nil)
	mockExtraDataCodec.On("DecodeExtraArgs", mock.Anything, mock.Anything).Return(map[string]any{
		"ComputeUnits":            uint32(1000),
		"AccountIsWritableBitmap": uint64(10),
		"Accounts": [][32]byte{
			[32]byte(config.CcipLogicReceiver.Bytes()),
			[32]byte(config.ReceiverTargetAccountPDA.Bytes()),
			[32]byte(solana.SystemProgramID.Bytes()),
		},
	}, nil)
	msgHasher := NewMessageHasherV1(logger.Test(t), mockExtraDataCodec)
	_, err := msgHasher.Hash(testutils.Context(t), any2AnyMsg)
	require.Error(t, err)
}

func TestMessageHasher_InvalidDestinationTokenAddress(t *testing.T) {
	any2AnyMsg, _, _ := createEVM2SolanaMessages(t)

	// Set DestTokenAddress to a []byte of 2 length
	any2AnyMsg.TokenAmounts[0].DestTokenAddress = []byte{0, 0}
	mockExtraDataCodec := &mocks.ExtraDataCodec{}
	mockExtraDataCodec.On("DecodeTokenAmountDestExecData", mock.Anything, mock.Anything).Return(map[string]any{
		"destGasAmount": uint32(10),
	}, nil)
	mockExtraDataCodec.On("DecodeExtraArgs", mock.Anything, mock.Anything).Return(map[string]any{
		"ComputeUnits":            uint32(1000),
		"AccountIsWritableBitmap": uint64(10),
		"Accounts": [][32]byte{
			[32]byte(config.CcipLogicReceiver.Bytes()),
			[32]byte(config.ReceiverTargetAccountPDA.Bytes()),
			[32]byte(solana.SystemProgramID.Bytes()),
		},
	}, nil)
	msgHasher := NewMessageHasherV1(logger.Test(t), mockExtraDataCodec)
	_, err := msgHasher.Hash(testutils.Context(t), any2AnyMsg)
	require.Error(t, err)
}

func createEVM2SolanaMessages(t *testing.T) (cciptypes.Message, ccip_offramp.Any2SVMRampMessage, []solana.PublicKey) {
	messageID := utils.RandomBytes32()
	sourceChain := uint64(5009297550715157269) // evm mainnet
	seqNum := rand.Uint64()
	nonce := rand.Uint64()
	destChain := rand.Uint64()

	messageData := make([]byte, rand.Intn(2048))
	_, err := cryptorand.Read(messageData)
	require.NoError(t, err)

	sender := abiEncodedAddress(t)
	receiver := solana.MustPublicKeyFromBase58("DS2tt4BX7YwCw7yrDNwbAdnYrxjeCPeGJbHmZEYC8RTb")
	tokenReceiver := solana.MustPublicKeyFromBase58("42Gia5bGsh8R2S44e37t9fsucap1qsgjr6GjBmWotgdF")
	extraArgs := ccip_offramp.Any2SVMRampExtraArgs{
		ComputeUnits:     uint32(10000),
		IsWritableBitmap: uint64(4),
	}
	abiEncodedExtraArgs := []byte{31, 59, 58, 186, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 39, 16, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 44, 230, 105, 156, 244, 184, 196, 235, 30, 58, 209, 82, 8, 202, 25, 73, 167, 169, 34, 150, 141, 129, 169, 150, 219, 160, 186, 44, 72, 156, 50, 170, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 160, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 44, 230, 105, 156, 244, 184, 196, 235, 30, 58, 209, 82, 8, 202, 25, 73, 167, 169, 34, 150, 141, 129, 169, 150, 219, 160, 186, 44, 72, 156, 50, 170}
	tokenAmount := cciptypes.NewBigInt(big.NewInt(rand.Int63()))
	destGasAmount, err := abiEncodeUint32(10)
	require.NoError(t, err)

	ccipTokenAmounts := make([]cciptypes.RampTokenAmount, 5)
	for z := 0; z < 5; z++ {
		ccipTokenAmounts[z] = cciptypes.RampTokenAmount{
			SourcePoolAddress: cciptypes.UnknownAddress("DS2tt4BX7YwCw7yrDNwbAdnYrxjeCPeGJbHmZEYC8RTb"),
			DestTokenAddress:  receiver.Bytes(),
			Amount:            tokenAmount,
			DestExecData:      destGasAmount,
		}
	}

	solTokenAmounts := make([]ccip_offramp.Any2SVMTokenTransfer, 5)
	for z := 0; z < 5; z++ {
		solTokenAmounts[z] = ccip_offramp.Any2SVMTokenTransfer{
			SourcePoolAddress: cciptypes.UnknownAddress("DS2tt4BX7YwCw7yrDNwbAdnYrxjeCPeGJbHmZEYC8RTb"),
			DestTokenAddress:  receiver,
			Amount:            ccip_offramp.CrossChainAmount{LeBytes: [32]uint8(encodeBigIntToFixedLengthLE(tokenAmount.Int, 32))},
			DestGasAmount:     uint32(10),
		}
	}

	any2SolanaMsg := ccip_offramp.Any2SVMRampMessage{
		Header: ccip_offramp.RampMessageHeader{
			MessageId:           messageID,
			SourceChainSelector: sourceChain,
			DestChainSelector:   destChain,
			SequenceNumber:      seqNum,
			Nonce:               nonce,
		},
		Sender:        sender,
		TokenReceiver: tokenReceiver,
		Data:          messageData,
		TokenAmounts:  solTokenAmounts,
		ExtraArgs:     extraArgs,
	}
	any2AnyMsg := cciptypes.Message{
		Header: cciptypes.RampMessageHeader{
			MessageID:           messageID,
			SourceChainSelector: cciptypes.ChainSelector(sourceChain),
			DestChainSelector:   cciptypes.ChainSelector(destChain),
			SequenceNumber:      cciptypes.SeqNum(seqNum),
			Nonce:               nonce,
			OnRamp:              abiEncodedAddress(t),
		},
		Sender:         sender,
		Receiver:       receiver.Bytes(),
		Data:           messageData,
		TokenAmounts:   ccipTokenAmounts,
		FeeToken:       []byte{},
		FeeTokenAmount: cciptypes.NewBigIntFromInt64(0),
		ExtraArgs:      abiEncodedExtraArgs,
	}

	msgAccounts := []solana.PublicKey{
		receiver,
		solana.MustPublicKeyFromBase58("42Gia5bGsh8R2S44e37t9fsucap1qsgjr6GjBmWotgdF"),
	}
	return any2AnyMsg, any2SolanaMsg, msgAccounts
}

func abiEncodedAddress(t *testing.T) []byte {
	addr := utils.RandomAddress()
	encoded, err := utils.ABIEncode(`[{"type": "address"}]`, addr)
	require.NoError(t, err)
	return encoded
}

func abiEncodeUint32(data uint32) ([]byte, error) {
	return utils.ABIEncode(`[{ "type": "uint32" }]`, data)
}
