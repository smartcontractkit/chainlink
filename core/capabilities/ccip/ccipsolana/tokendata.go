package ccipsolana

import (
	"context"
	"errors"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

type SolanaTokenDataEncoder struct{}

func NewSolanaTokenDataEncoder() SolanaTokenDataEncoder {
	return SolanaTokenDataEncoder{}
}

func (e SolanaTokenDataEncoder) EncodeUSDC(_ context.Context, message cciptypes.Bytes, attestation cciptypes.Bytes) (cciptypes.Bytes, error) {
	return nil, errors.New("not implemented")
}
