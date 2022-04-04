package dydx

import (
	"context"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

var _ types.ContractTransmitter = (*ContractTracker)(nil)

// Transmit uploads the result to dydx API endpoint, by making an HTTP call to
// the dydx ExternalAdapter. The EA does the actual signing and uploading work.
func (c *ContractTracker) Transmit(
	ctx context.Context,
	reportCtx types.ReportContext,
	report types.Report,
	sigs []types.AttributedOnchainSignature,
) error {

	// TODO: Implement the call to the dydx EA which uploads the result to the API endpoint

	c.ansLock.RLock()
	defer c.ansLock.RUnlock()
	c.answer = Answer{
		Data:      nil,
		Timestamp: time.Now(),
		epoch:     reportCtx.Epoch,
		round:     reportCtx.Round,
	}
	return nil
}

// Returns the latest epoch from the last stored transmission.
func (c *ContractTracker) LatestConfigDigestAndEpoch(
	ctx context.Context,
) (
	configDigest types.ConfigDigest,
	epoch uint32,
	err error,
) {
	digester, err := c.digester.configDigest()
	c.ansLock.RLock()
	defer c.ansLock.RUnlock()
	return digester, c.answer.epoch, err
}

// TODO: Check if returning an item from StaticTransmitters value is good enough
func (c *ContractTracker) FromAccount() types.Account {
	return StaticTransmitters[0]
}
