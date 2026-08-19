package ccipnoop

import (
	"context"

	ccipocr3common "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

type tokenDataEncoder struct{}

func NewTokenDataEncoder() ccipocr3common.TokenDataEncoder {
	return tokenDataEncoder{}
}

func (e tokenDataEncoder) EncodeUSDC(_ context.Context, message ccipocr3common.Bytes, attestation ccipocr3common.Bytes) (ccipocr3common.Bytes, error) {
	return []byte{}, nil
}
