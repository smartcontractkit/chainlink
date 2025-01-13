package ccipsolana

import (
	"context"
	"fmt"

	agbinary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/ccip"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

// MessageHasherV1 implements the MessageHasher interface.
// Compatible with:
// - "OnRamp 1.6.0-dev"
type MessageHasherV1 struct {
	lggr logger.Logger
}

func NewMessageHasherV1(lggr logger.Logger) *MessageHasherV1 {
	return &MessageHasherV1{
		lggr: lggr,
	}
}

// Hash implements the MessageHasher interface.
func (h *MessageHasherV1) Hash(_ context.Context, msg cciptypes.Message) (cciptypes.Bytes32, error) {
	h.lggr.Debugw("hashing message", "msg", msg)

	anyToSolanaMessage := ccip_router.Any2SolanaRampMessage{}
	anyToSolanaMessage.Header = ccip_router.RampMessageHeader{
		SourceChainSelector: uint64(msg.Header.SourceChainSelector),
		DestChainSelector:   uint64(msg.Header.DestChainSelector),
		SequenceNumber:      uint64(msg.Header.SequenceNumber),
		MessageId:           msg.Header.MessageID,
		Nonce:               msg.Header.Nonce,
	}
	anyToSolanaMessage.Receiver = solana.PublicKeyFromBytes(msg.Receiver)
	anyToSolanaMessage.Sender = msg.Sender
	anyToSolanaMessage.Data = msg.Data
	for _, ta := range msg.TokenAmounts {
		anyToSolanaMessage.TokenAmounts = append(anyToSolanaMessage.TokenAmounts, ccip_router.Any2SolanaTokenTransfer{
			SourcePoolAddress: ta.SourcePoolAddress,
			DestTokenAddress:  solana.PublicKeyFromBytes(ta.DestTokenAddress),
			ExtraData:         ta.ExtraData,
			Amount:            tokens.ToLittleEndianU256(ta.Amount.Int.Uint64()),
		})
	}
	decoder := agbinary.NewBorshDecoder(msg.ExtraArgs)
	err := anyToSolanaMessage.ExtraArgs.UnmarshalWithDecoder(decoder)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to decode ExtraArgs: %w", err)
	}

	hash, err := ccip.HashEvmToSolanaMessage(anyToSolanaMessage, msg.Header.OnRamp)

	return [32]byte(hash), err
}

// Interface compliance check
var _ cciptypes.MessageHasher = (*MessageHasherV1)(nil)
