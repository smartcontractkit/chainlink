package cre

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/shopspring/decimal"

	commonds "github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	datastreamsllo "github.com/smartcontractkit/chainlink-data-streams/llo"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"

	"google.golang.org/protobuf/proto"
)

var _ datastreamsllo.ReportCodec = ReportCodecCapabilityTrigger{}

type ReportCodecCapabilityTrigger struct {
	lggr  logger.Logger
	donID uint32
}

func NewReportCodecCapabilityTrigger(lggr logger.Logger, donID uint32) ReportCodecCapabilityTrigger {
	return ReportCodecCapabilityTrigger{lggr, donID}
}

type ReportCodecCapabilityTriggerMultiplier struct {
	Multiplier decimal.Decimal   `json:"multiplier"`
	StreamID   llotypes.StreamID `json:"streamID"`
}

// Opts format remains unchanged
type ReportCodecCapabilityTriggerOpts struct {
	// EXAMPLE
	//
	// [{streamID: 1000000001, "multiplier":"10000"}, ...]
	//
	// The total number of streams must be n, where n is the number of
	// top-level elements in this ReportCodecCapabilityTriggerMultipliers array
	Multipliers []ReportCodecCapabilityTriggerMultiplier `json:"multipliers"`
}

func (r *ReportCodecCapabilityTriggerOpts) Decode(opts []byte) error {
	if len(opts) == 0 {
		return nil
	}
	decoder := json.NewDecoder(bytes.NewReader(opts))
	decoder.DisallowUnknownFields() // Error on unrecognized fields
	return decoder.Decode(r)
}

func (r *ReportCodecCapabilityTriggerOpts) Encode() ([]byte, error) {
	return json.Marshal(r)
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

	// NOTE: It seems suboptimal to have to parse the opts on every encode but
	// not sure how to avoid it. Should be negligible performance hit as long
	// as Opts is small.
	opts := ReportCodecCapabilityTriggerOpts{}
	if err := (&opts).Decode(cd.Opts); err != nil {
		return nil, fmt.Errorf("failed to decode opts; got: '%s'; %w", cd.Opts, err)
	}

	payload := make([]*commonds.LLOStreamDecimal, len(report.Values))
	for i, stream := range report.Values {
		var d []byte
		switch v := stream.(type) {
		case nil:
			// Missing observations are nil
		case *datastreamsllo.Decimal:
			multipliedStreamValue := v.Decimal()

			if len(opts.Multipliers) != 0 {
				multipliedStreamValue = multipliedStreamValue.Mul(opts.Multipliers[i].Multiplier)
			}

			var err error
			d, err = multipliedStreamValue.MarshalBinary()
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
	opts := new(ReportCodecCapabilityTriggerOpts)
	if err := opts.Decode(cd.Opts); err != nil {
		return fmt.Errorf("invalid Opts, got: %q; %w", cd.Opts, err)
	}
	if opts != nil && opts.Multipliers != nil {
		if len(opts.Multipliers) != len(cd.Streams) {
			return fmt.Errorf("multipliers length %d != StreamValues length %d", len(opts.Multipliers), len(cd.Streams))
		}

		for i, stream := range cd.Streams {
			if opts.Multipliers[i].StreamID != stream.StreamID {
				return fmt.Errorf("LLO StreamID %d mismatched with Multiplier StreamID %d", stream.StreamID, opts.Multipliers[i].StreamID)
			}
			if !(opts.Multipliers[i].Multiplier.IsInteger()) {
				return fmt.Errorf("multiplier for StreamID %d must be an integer", opts.Multipliers[i].StreamID)
			}
			if opts.Multipliers[i].Multiplier.IsZero() {
				return fmt.Errorf("multiplier for StreamID %d can't be zero", opts.Multipliers[i].StreamID)
			}
			if opts.Multipliers[i].Multiplier.IsNegative() {
				return fmt.Errorf("multiplier for StreamID %d can't be negative", opts.Multipliers[i].StreamID)
			}
		}
	}
	return nil
}

// EventID is expected to uniquely identify a (don, round)
func (r ReportCodecCapabilityTrigger) EventID(report datastreamsllo.Report) string {
	return fmt.Sprintf("streams_%d_%d", r.donID, report.ObservationTimestampNanoseconds)
}
