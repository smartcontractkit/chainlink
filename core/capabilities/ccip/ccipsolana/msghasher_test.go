package ccipsolana

import (
	"bytes"
	cryptorand "crypto/rand"
	"math/big"
	"math/rand"
	"testing"

	agbinary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common/mocks"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/config"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/ccip"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-integrations/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

func TestMessageHasher_Any2Solana(t *testing.T) {
	any2AnyMsg, any2SolanaMsg, msgAccounts := createAny2SolanaMessages(t)
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
	actualHash, err := msgHasher.Hash(testutils.Context(t), any2AnyMsg)
	require.NoError(t, err)
	expectedHash, err := ccip.HashAnyToSVMMessage(any2SolanaMsg, any2AnyMsg.Header.OnRamp, msgAccounts)
	require.NoError(t, err)
	require.Equal(t, expectedHash, actualHash[:32])
}

func createAny2SolanaMessages(t *testing.T) (cciptypes.Message, ccip_offramp.Any2SVMRampMessage, []solana.PublicKey) {
	messageID := utils.RandomBytes32()

	sourceChain := rand.Uint64()
	seqNum := rand.Uint64()
	nonce := rand.Uint64()
	destChain := rand.Uint64()

	messageData := make([]byte, rand.Intn(2048))
	_, err := cryptorand.Read(messageData)
	require.NoError(t, err)

	sender := abiEncodedAddress(t)
	receiver := solana.MustPublicKeyFromBase58("DS2tt4BX7YwCw7yrDNwbAdnYrxjeCPeGJbHmZEYC8RTb")
	computeUnit := uint32(1000)
	bitmap := uint64(10)

	extraArgs := ccip_offramp.Any2SVMRampExtraArgs{
		ComputeUnits:     computeUnit,
		IsWritableBitmap: bitmap,
	}
	var buf bytes.Buffer
	encoder := agbinary.NewBorshEncoder(&buf)
	err = extraArgs.MarshalWithEncoder(encoder)
	require.NoError(t, err)
	tokenAmount := cciptypes.NewBigInt(big.NewInt(rand.Int63()))

	ccipTokenAmounts := make([]cciptypes.RampTokenAmount, 5)
	for z := 0; z < 5; z++ {
		ccipTokenAmounts[z] = cciptypes.RampTokenAmount{
			SourcePoolAddress: cciptypes.UnknownAddress("DS2tt4BX7YwCw7yrDNwbAdnYrxjeCPeGJbHmZEYC8RTb"),
			DestTokenAddress:  receiver.Bytes(),
			Amount:            tokenAmount,
			DestExecDataDecoded: map[string]any{
				"destGasAmount": uint32(10),
			},
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
		TokenReceiver: receiver,
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
		ExtraArgs:      buf.Bytes(),
		ExtraArgsDecoded: map[string]any{
			"ComputeUnits":            computeUnit,
			"AccountIsWritableBitmap": bitmap,
			"Accounts": [][32]byte{
				[32]byte(config.CcipLogicReceiver.Bytes()),
				[32]byte(config.ReceiverTargetAccountPDA.Bytes()),
				[32]byte(solana.SystemProgramID.Bytes()),
			},
		},
	}

	msgAccounts := []solana.PublicKey{
		config.CcipLogicReceiver,
		config.ReceiverTargetAccountPDA,
		solana.SystemProgramID,
	}
	return any2AnyMsg, any2SolanaMsg, msgAccounts
}

func abiEncodedAddress(t *testing.T) []byte {
	addr := utils.RandomAddress()
	encoded, err := utils.ABIEncode(`[{"type": "address"}]`, addr)
	require.NoError(t, err)
	return encoded
}
