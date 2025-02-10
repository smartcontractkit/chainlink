package ccipsolana

import (
	"context"
	"errors"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

// SolanaRMNCrypto is the RMNCrypto implementation for Solana chains.
type SolanaRMNCrypto struct{}

func (r *SolanaRMNCrypto) VerifyReportSignatures(
	_ context.Context,
	_ []cciptypes.RMNECDSASignature,
	_ cciptypes.RMNReport,
	_ []cciptypes.UnknownAddress,
) error {
	return errors.New("not implemented")
}
