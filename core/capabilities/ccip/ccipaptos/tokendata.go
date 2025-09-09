package ccipaptos

import (
	"context"
	"errors"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

type AptosTokenDataEncoder struct{}

func NewAptosTokenDataEncoder() AptosTokenDataEncoder {
	return AptosTokenDataEncoder{}
}

func (e AptosTokenDataEncoder) EncodeUSDC(_ context.Context, message cciptypes.Bytes, attestation cciptypes.Bytes) (cciptypes.Bytes, error) {
	return nil, errors.New("not implemented")
}
