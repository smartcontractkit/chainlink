package ccipnoop

import (
	"context"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

type NoopTokenDataEncoder struct{}

func NewNoopTokenDataEncoder() NoopTokenDataEncoder {
	return NoopTokenDataEncoder{}
}

func (e NoopTokenDataEncoder) EncodeUSDC(_ context.Context, message cciptypes.Bytes, attestation cciptypes.Bytes) (cciptypes.Bytes, error) {
	return []byte{}, nil
}
