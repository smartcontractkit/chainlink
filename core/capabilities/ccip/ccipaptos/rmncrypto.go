package ccipaptos

import (
	"context"
	"errors"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

// AptosRMNCrypto is the RMNCrypto implementation for Aptos chains.
type AptosRMNCrypto struct{}

func (r *AptosRMNCrypto) VerifyReportSignatures(
	_ context.Context,
	_ []cciptypes.RMNECDSASignature,
	_ cciptypes.RMNReport,
	_ []cciptypes.UnknownAddress,
) error {
	return errors.New("not implemented")
}
