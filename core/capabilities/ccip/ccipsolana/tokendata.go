package ccipsolana

import (
	"bytes"
	"context"
	"fmt"

	bin "github.com/gagliardetto/binary"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/latest/cctp_token_pool"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

type SolanaTokenDataEncoder struct{}

func NewSolanaTokenDataEncoder() SolanaTokenDataEncoder {
	return SolanaTokenDataEncoder{}
}

func (e SolanaTokenDataEncoder) EncodeUSDC(_ context.Context, message, attestation ccipocr3.Bytes) (ccipocr3.Bytes, error) {
	messageAndAttestation := cctp_token_pool.MessageAndAttestation{
		Message:     cctp_token_pool.CctpMessage{Data: message},
		Attestation: attestation,
	}
	buf := new(bytes.Buffer)
	err := bin.NewBorshEncoder(buf).Encode(messageAndAttestation)
	if err != nil {
		return nil, fmt.Errorf("failed to borsh encode USDC message and attestation: %w", err)
	}
	return buf.Bytes(), nil
}
