package telem

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-data-streams/llo"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
)

const samplerDelimiter = "-"

// sampler keeps track of what kind of telemetry has already been sent to the collection point and decides whether the
// next telemetry package will be sent or dropped.
type sampler struct {
	// samples keeps track of the telemetry samples we've already sent (or approved for sending).
	// The format is `map[fingerprint][observation timestamp in seconds]any`
	samples map[string]map[int64]any
	lggr    logger.Logger
}

func newSampler(lgger logger.SugaredLogger) *sampler {
	return &sampler{
		samples: make(map[string]map[int64]any),
		lggr:    lgger,
	}
}

// TODO implement a loop which evicts data older than 10 seconds. start this on telemetry start

// TODO implement config option. should allow enabling/disabling and changing the sampling frequency.

// Sample is the method which decides whether we're going to send the data downstream or not.
func (s *sampler) Sample(typ synchronization.TelemetryType, msg proto.Message) bool {

	fp, ots, err := fingerprint(typ, msg)
	if err != nil {
		s.lggr.Warnw("Couldn't determine fingerprint", "type", typ, "err", err)
		// default to sending data
		return true
	}
	// Do we have any records for this fingerprint?
	if _, ok := s.samples[fp]; !ok {
		s.samples[fp] = make(map[int64]any)
	}
	// Do we already have a record for this fingerprint and this second?
	if _, ok := s.samples[fp][nanosToSec(ots)]; !ok {
		s.samples[fp][nanosToSec(ots)] = struct{}{}
		return true
	}
	// We already have a record, and we don't need to send another one.
	return false
}

// fingerprint combines unique characteristics of each supported telemetry report type and constructs a string
// fingerprint of it. It returns the fingerprint, together with a nanosecond observation timestamp.
// TODO improve encoding efficiency by switching from string to hex([]byte) or string([]byte).
func fingerprint(typ synchronization.TelemetryType, msg proto.Message) (string, int64, error) {
	switch typ {
	case synchronization.LLOObservation:
		m, ok := msg.(*LLOObservationTelemetry)
		if !ok || m == nil {
			return "", 0, errors.New("invalid telemetry type, expected LLOObservation")
		}
		traits := []string{
			fmt.Sprint(m.DonId),
			fmt.Sprint(m.GetStreamId()),
			fmt.Sprint(hex.EncodeToString(m.ConfigDigest)),
		}
		return strings.Join(traits, samplerDelimiter), m.ObservationTimestamp, nil

	case synchronization.LLOOutcome:
		m, ok := msg.(*llo.LLOOutcomeTelemetry)
		if !ok || m == nil {
			return "", 0, errors.New("invalid telemetry type, expected LLOOutcomeTelemetry")
		}
		traits := []string{
			fmt.Sprint(m.DonId),
			fmt.Sprint(hex.EncodeToString(m.ConfigDigest)),
		}
		return strings.Join(traits, samplerDelimiter), int64(m.ObservationTimestampNanoseconds), nil
	case synchronization.LLOReport:
		m, ok := msg.(*llo.LLOReportTelemetry)
		if !ok || m == nil {
			return "", 0, errors.New("invalid telemetry type, expected LLOReportTelemetry")
		}
		traits := []string{
			fmt.Sprint(m.DonId),
			fmt.Sprint(m.ChannelId),
			fmt.Sprint(hex.EncodeToString(m.ConfigDigest)),
		}
		return strings.Join(traits, samplerDelimiter), int64(m.ObservationTimestampNanoseconds), nil
	case synchronization.PipelineBridge:
		m, ok := msg.(*LLOBridgeTelemetry)
		if !ok || m == nil {
			return "", 0, errors.New("invalid telemetry type, expected LLOBridgeTelemetry")
		}
		traits := []string{
			fmt.Sprint(m.DonId),
			fmt.Sprint(m.GetStreamId()),
			fmt.Sprint(m.SpecId),
			fmt.Sprint(m.BridgeAdapterName),
			fmt.Sprint(hex.EncodeToString(m.ConfigDigest)),
		}
		return strings.Join(traits, samplerDelimiter), m.ObservationTimestamp, nil
	default:
		// TODO We should probably return a sentinel error here and always send these downstream.
		// 	This is not really an error, it's just a "we got something we don't intend to sample" situation.
		return "", 0, fmt.Errorf("unsupported telemetry type: %s", typ)
	}
}

func nanosToSec(n int64) int64 {
	return n / int64(time.Second)
}
