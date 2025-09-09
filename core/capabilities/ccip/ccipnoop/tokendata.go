package ccipnoop

import (
	"context"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

type tokenDataEncoder struct{}

func NewTokenDataEncoder() ccipocr3.TokenDataEncoder {
	return tokenDataEncoder{}
}

func (e tokenDataEncoder) EncodeUSDC(_ context.Context, message ccipocr3.Bytes, attestation ccipocr3.Bytes) (ccipocr3.Bytes, error) {
	return []byte{}, nil
}
