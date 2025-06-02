package ccipnoop

import (
	"context"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

// NoopRMNCrypto is the RMNCrypto implementation.
type NoopRMNCrypto struct{}

func (r *NoopRMNCrypto) VerifyReportSignatures(
	_ context.Context,
	_ []cciptypes.RMNECDSASignature,
	_ cciptypes.RMNReport,
	_ []cciptypes.UnknownAddress,
) error {
	return nil
}
