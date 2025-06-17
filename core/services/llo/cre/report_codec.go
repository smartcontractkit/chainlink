package cre

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"

	commonds "github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	datastreamsllo "github.com/smartcontractkit/chainlink-data-streams/llo"
)

var _ datastreamsllo.ReportCodec = ReportCodecCapabilityTrigger{}

type ReportCodecCapabilityTrigger struct {
	lggr  logger.Logger
	donID uint32
}

func NewReportCodecCapabilityTrigger(lggr logger.Logger, donID uint32) ReportCodecCapabilityTrigger {
	return ReportCodecCapabilityTrigger{lggr, donID}
}

// Encode a report into a capability trigger report
// the returned byte slice is the marshaled protobuf of [capabilitiespb.OCRTriggerReport]
func (r ReportCodecCapabilityTrigger) Encode(report datastreamsllo.Report, cd llotypes.ChannelDefinition) ([]byte, error) {
	if len(cd.Streams) != len(report.Values) {
		// Invariant violation
		return nil, fmt.Errorf("capability trigger expected %d streams, got %d", len(cd.Streams), len(report.Values))
	}
	if report.Specimen {
		// Not supported for now
		return nil, errors.New("capability trigger encoder does not currently support specimen reports")
	}
	payload := make([]*commonds.LLOStreamDecimal, len(report.Values))
	for i, stream := range report.Values {
		var d []byte
		switch stream.(type) {
		case nil:
			// Missing observations are nil
		case *datastreamsllo.Decimal:
			var err error
			d, err = stream.MarshalBinary()
			if err != nil {
				return nil, fmt.Errorf("failed to marshal decimal: %w", err)
			}
		default:
			return nil, fmt.Errorf("only decimal StreamValues are supported, got: %T", stream)
		}
		payload[i] = &commonds.LLOStreamDecimal{
			StreamID: cd.Streams[i].StreamID,
			Decimal:  d,
		}
	}
	ste := commonds.LLOStreamsTriggerEvent{
		Payload:                         payload,
		ObservationTimestampNanoseconds: report.ObservationTimestampNanoseconds,
	}
	outputs, err := values.WrapMap(ste)
	if err != nil {
		return nil, fmt.Errorf("failed to wrap map: %w", err)
	}
	p := &capabilitiespb.OCRTriggerReport{
		EventID:   r.EventID(report),
		Timestamp: report.ObservationTimestampNanoseconds,
		Outputs:   values.ProtoMap(outputs),
	}

	b, err := proto.MarshalOptions{Deterministic: true}.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal capability trigger report: %w", err)
	}
	return b, nil
}

func (r ReportCodecCapabilityTrigger) Verify(cd llotypes.ChannelDefinition) error {
	if len(cd.Opts) > 0 {
		return errors.New("capability trigger does not support channel definitions with options")
	}
	return nil
}

// EventID is expected to uniquely identify a (don, round)
func (r ReportCodecCapabilityTrigger) EventID(report datastreamsllo.Report) string {
	return fmt.Sprintf("streams_%d_%d", r.donID, report.ObservationTimestampNanoseconds)
}
