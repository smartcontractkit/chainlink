package ocrimpls

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var _ ocr3types.ContractTransmitter[[]byte] = (*noOpTransmitter)(nil)

const errMsg = "no-op transmitter %s called, it shouldn't be! Check the CCIPHome ChainConfig for this chain and ensure the 'readers' field is correctly set"

// NewNoOpTransmitter creates a new no-op transmitter. It is intended to be used in
// role DONs where the node that is participating cannot transmit to the destination chain.
func NewNoOpTransmitter(lggr logger.Logger, myP2PID string) *noOpTransmitter {
	return &noOpTransmitter{
		lggr:    lggr,
		myP2PID: myP2PID,
	}
}

// noOpTransmitter is an implementation of ocr3types.ContractTransmitter[[]byte]
// that does nothing. It is intended to be used in role DONs where the
// node that is participating cannot transmit to the destination chain.
type noOpTransmitter struct {
	lggr    logger.Logger
	myP2PID string
}

// FromAccount implements ocr3types.ContractTransmitter.
func (n *noOpTransmitter) FromAccount(context.Context) (types.Account, error) {
	n.lggr.Criticalw(fmt.Sprintf(errMsg, "FromAccount()"),
		"myP2PID", n.myP2PID,
	)
	// Return nil because even if we incorrectly call this, the transmission
	// schedule should eventually go to another node that can transmit and hopefully succeeds.
	// If we return an error it'll look like there's really something wrong, when in
	// fact it's just a no-op.
	return types.Account(""), nil
}

// Transmit implements ocr3types.ContractTransmitter.
func (n *noOpTransmitter) Transmit(_ context.Context, digest types.ConfigDigest, seqNr uint64, _ ocr3types.ReportWithInfo[[]byte], _ []types.AttributedOnchainSignature) error {
	n.lggr.Criticalw(fmt.Sprintf(errMsg, "Transmit()"),
		"myP2PID", n.myP2PID,
		"configDigest", digest.Hex(),
		"seqNr", seqNr,
	)
	// Return nil because even if we incorrectly call this, the transmission
	// schedule should eventually go to another node that can transmit and hopefully succeeds.
	// If we return an error it'll look like there's really something wrong, when in
	// fact it's just a no-op.
	return nil
}
