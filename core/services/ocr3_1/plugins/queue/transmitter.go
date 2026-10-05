package queue

import (
	"context"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2/triggers"
)

type contractTransmitter struct {
	lggr        logger.Logger
	fromAccount ocrtypes.Account
	eventSink   triggers.EventSink
}

func NewTransmitter(lggr logger.Logger, fromAccount ocrtypes.Account, eventSink triggers.EventSink) ocr3types.ContractTransmitter[[]byte] {
	return &contractTransmitter{
		lggr:        lggr,
		fromAccount: fromAccount,
		eventSink:   eventSink,
	}
}

func (c *contractTransmitter) Transmit(ctx context.Context, digest ocrtypes.ConfigDigest, seqNum uint64, rpi ocr3types.ReportWithInfo[[]byte], signatures []ocrtypes.AttributedOnchainSignature) error {
	var events []triggers.CoordinatedEvent
	//TODO deserialize list of trigger events from rpi
	for _, e := range events {
		if c.eventSink == nil {
			c.lggr.Warnw("Unable to execute trigger - no event sink configured", "event", e)
			continue
		}
		err := c.eventSink.ExecuteTrigger(ctx, e)
		if err != nil {
			c.lggr.Errorw("Failed to execute trigger", "event", e, "err", err)
		}
	}

	return nil
}

func (c *contractTransmitter) FromAccount(ctx context.Context) (ocrtypes.Account, error) {
	return c.fromAccount, nil
}
