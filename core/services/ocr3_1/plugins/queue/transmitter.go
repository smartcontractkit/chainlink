package queue

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
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
	var report *pb.ConsensusQueueReport
	err := proto.Unmarshal(rpi.Report, report)
	if err != nil {
		return fmt.Errorf("failed to unmarshal report: %w", err)
	}
	// Execute each event
	for _, message := range report.Events {
		var event triggers.CoordinatedEvent
		if err := event.FromProto(message); err != nil {
			c.lggr.Errorw("Failed to convert trigger event proto", "err", err)
			continue
		}
		if c.eventSink == nil {
			c.lggr.Warnw("Unable to execute trigger - no event sink configured", "event", event)
			continue
		}
		err := c.eventSink.ExecuteTrigger(ctx, event)
		if err != nil {
			c.lggr.Errorw("Failed to execute trigger", "event", event, "err", err)
		}
	}

	return nil
}

func (c *contractTransmitter) FromAccount(ctx context.Context) (ocrtypes.Account, error) {
	return c.fromAccount, nil
}
