package ccipsolana

import (
	"bytes"
	cryptorand "crypto/rand"
	"math/rand"
	"testing"

	agbinary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/config"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/ccip"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

func TestMessageHasher_Any2Solana(t *testing.T) {
	any2AnyMsg, any2SolanaMsg := createAny2SolanaMessages(t)
	msgHasher := NewMessageHasherV1(logger.Test(t))
	actualHash, err := msgHasher.Hash(testutils.Context(t), any2AnyMsg)
	require.NoError(t, err)
	expectedHash, err := ccip.HashEvmToSolanaMessage(any2SolanaMsg, any2AnyMsg.Header.OnRamp)
	require.NoError(t, err)
	require.Equal(t, actualHash[:32], expectedHash)
}

func createAny2SolanaMessages(t *testing.T) (cciptypes.Message, ccip_router.Any2SolanaRampMessage) {
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

	extraArgs := ccip_router.SolanaExtraArgs{
		ComputeUnits: 1000,
		Accounts: []ccip_router.SolanaAccountMeta{
			{Pubkey: config.CcipReceiverProgram},
			{Pubkey: config.ReceiverTargetAccountPDA, IsWritable: true},
			{Pubkey: solana.SystemProgramID, IsWritable: false},
		},
	}
	var buf bytes.Buffer
	encoder := agbinary.NewBorshEncoder(&buf)
	extraArgs.MarshalWithEncoder(encoder)
	require.NoError(t, err)

	any2SolanaMsg := ccip_router.Any2SolanaRampMessage{
		Header: ccip_router.RampMessageHeader{
			MessageId:           messageID,
			SourceChainSelector: sourceChain,
			DestChainSelector:   destChain,
			SequenceNumber:      seqNum,
			Nonce:               nonce,
		},
		Sender:       sender,
		Receiver:     receiver,
		Data:         messageData,
		TokenAmounts: nil,
		ExtraArgs:    extraArgs,
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
		TokenAmounts:   nil,
		FeeToken:       []byte{},
		FeeTokenAmount: cciptypes.NewBigIntFromInt64(0),
		ExtraArgs:      buf.Bytes(),
	}
	return any2AnyMsg, any2SolanaMsg
}

func abiEncodedAddress(t *testing.T) []byte {
	addr := utils.RandomAddress()
	encoded, err := utils.ABIEncode(`[{"type": "address"}]`, addr)
	require.NoError(t, err)
	return encoded
}
