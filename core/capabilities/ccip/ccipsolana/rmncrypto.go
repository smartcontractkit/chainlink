package ccipsolana

import (
	"context"
	"errors"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// SolanaRMNCrypto is the RMNCrypto implementation for Solana chains.
type SolanaRMNCrypto struct {
	lggr logger.Logger
}

func (r *SolanaRMNCrypto) VerifyReportSignatures(
	_ context.Context,
	_ []cciptypes.RMNECDSASignature,
	_ cciptypes.RMNReport,
	_ []cciptypes.UnknownAddress,
) error {
	return errors.New("not implemented")
}
