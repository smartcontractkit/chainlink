package ccipsui

import (
	"context"
	"errors"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

type SuiTokenDataEncoder struct{}

func NewSuiTokenDataEncoder() SuiTokenDataEncoder {
	return SuiTokenDataEncoder{}
}

func (e SuiTokenDataEncoder) EncodeUSDC(_ context.Context, message cciptypes.Bytes, attestation cciptypes.Bytes) (cciptypes.Bytes, error) {
	return nil, errors.New("not implemented")
}
