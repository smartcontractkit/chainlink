package ccipaptos

import (
	"context"
	"errors"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

type AptosTokenDataEncoder struct{}

func NewAptosTokenDataEncoder() AptosTokenDataEncoder {
	return AptosTokenDataEncoder{}
}

func (e AptosTokenDataEncoder) EncodeUSDC(_ context.Context, message, attestation ccipocr3.Bytes) (ccipocr3.Bytes, error) {
	return nil, errors.New("not implemented")
}
