package functions

import (
	"context"

	"github.com/pkg/errors"
)

// Simple wrapper around a channel to transmit offchain reports between
// OCR plugin and Gateway connector
type OffchainTransmitter interface {
	TransmitReport(ctx context.Context, report *OffchainResponse) error
	ReportChannel() chan *OffchainResponse
}

type offchainTransmitter struct {
	reportCh chan *OffchainResponse
}

func NewOffchainTransmitter(chanSize uint32) OffchainTransmitter {
	return &offchainTransmitter{
		reportCh: make(chan *OffchainResponse, chanSize),
	}
}

func (t *offchainTransmitter) TransmitReport(ctx context.Context, report *OffchainResponse) error {
	select {
	case t.reportCh <- report:
		return nil
	case <-ctx.Done():
		return errors.New("context cancelled")
	}
}

func (t *offchainTransmitter) ReportChannel() chan *OffchainResponse {
	return t.reportCh
}
