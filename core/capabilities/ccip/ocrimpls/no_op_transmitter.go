package ocrimpls

import (
	"context"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

var _ ocr3types.ContractTransmitter[[]byte] = (*noOpTransmitter)(nil)

// NewNoOpTransmitter creates a new no-op transmitter. It is intended to be used in
// role DONs where the node that is participating cannot transmit to the destination chain.
func NewNoOpTransmitter(lggr logger.Logger) *noOpTransmitter {
	return &noOpTransmitter{
		lggr: lggr.Named("NoOpTransmitter"),
	}
}

// noOpTransmitter is an implementation of ocr3types.ContractTransmitter[[]byte]
// that does nothing. It is intended to be used in role DONs where the
// node that is participating cannot transmit to the destination chain.
type noOpTransmitter struct {
	lggr logger.Logger
}

// FromAccount implements ocr3types.ContractTransmitter.
func (n *noOpTransmitter) FromAccount(context.Context) (types.Account, error) {
	n.lggr.Errorw("no-op transmitter FromAccount() called, it shouldn't be!")
	// Return nil because even if we incorrectly call this, the transmission
	// schedule should eventually go to another node that can transmit and hopefully succeeds.
	// If we return an error it'll look like there's really something wrong, when in
	// fact it's just a no-op.
	return types.Account(""), nil
}

// Transmit implements ocr3types.ContractTransmitter.
func (n *noOpTransmitter) Transmit(context.Context, types.ConfigDigest, uint64, ocr3types.ReportWithInfo[[]byte], []types.AttributedOnchainSignature) error {
	n.lggr.Errorw("no-op transmitter Transmit() called, it shouldn't be!")
	// Return nil because even if we incorrectly call this, the transmission
	// schedule should eventually go to another node that can transmit and hopefully succeeds.
	// If we return an error it'll look like there's really something wrong, when in
	// fact it's just a no-op.
	return nil
}
