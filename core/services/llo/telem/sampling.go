package telem

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-data-streams/llo"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
)

const (
	defaultPrunePeriod = 10 * time.Second

	samplerDelimiter = "-"
)

var (
	errUnsupportedTelemetryType = errors.New("unsupported telemetry type")
)

// sampler keeps track of what kind of telemetry has already been sent to the collection point and decides whether the
// next telemetry package will be sent or dropped.
type sampler struct {
	// samples keeps track of the telemetry samples we've already sent (or approved for sending).
	// The format is `map[fingerprint][observation timestamp in seconds]any`. We intentionally use int32 because it's
	// enough for seconds but not enough for nanos, so we can't mix them up.
	samples   map[string]map[int32]any
	samplesMu sync.Mutex

	prunePeriod time.Duration // exists, so we can test pruning
	lggr        logger.Logger
}

func newSampler(lgger logger.SugaredLogger) *sampler {
	return &sampler{
		samples:     make(map[string]map[int32]any),
		prunePeriod: defaultPrunePeriod,
		lggr:        lgger,
	}
}

// Sample is the method which decides whether we're going to send the data downstream or not.
func (s *sampler) Sample(typ synchronization.TelemetryType, msg proto.Message) bool {
	fp, ots, err := fingerprint(typ, msg)
	if err != nil {
		if err != errUnsupportedTelemetryType {
			s.lggr.Warnw("Couldn't determine fingerprint", "type", typ, "err", err)
		}
		return true
	}

	s.samplesMu.Lock()
	defer s.samplesMu.Unlock()
	// Do we have any records for this fingerprint?
	if _, ok := s.samples[fp]; !ok {
		s.samples[fp] = make(map[int32]any)
	}
	// Do we already have a record for this fingerprint and this second?
	if _, ok := s.samples[fp][ots]; !ok {
		s.samples[fp][ots] = struct{}{}
		return true
	}
	// We already have a record, and we don't need to send another one.
	return false
}

// StartPruningLoop starts a regular check routine which removes old entries from the sampling records.
//
// This method is non-blocking. It starts a goroutine and returns.
func (s *sampler) StartPruningLoop(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		t := time.NewTicker(s.prunePeriod)
		defer t.Stop()

		for {
			select {
			case <-t.C:
				s.pruneStorage()
			case <-ctx.Done():
				return
			}
		}
	}()
}

// pruneStorage removes all records which are older than a predefined period (s.prunePeriod).
func (s *sampler) pruneStorage() {
	s.samplesMu.Lock()
	defer s.samplesMu.Unlock()

	cutoff := int32(time.Now().Add(-s.prunePeriod).Unix())
	for _, ots := range s.samples {
		for ts := range ots {
			if ts < cutoff {
				delete(ots, ts)
			}
		}
	}
}

// fingerprint combines unique characteristics of each supported telemetry report type and constructs a string
// fingerprint of it. It returns the fingerprint, together with an observation timestamp in seconds.
// TODO improve encoding efficiency by switching from string to hex([]byte) or string([]byte).
func fingerprint(typ synchronization.TelemetryType, msg proto.Message) (string, int32, error) {
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
		return strings.Join(traits, samplerDelimiter), nanosToSec(m.ObservationTimestamp), nil
	case synchronization.LLOOutcome:
		m, ok := msg.(*llo.LLOOutcomeTelemetry)
		if !ok || m == nil {
			return "", 0, errors.New("invalid telemetry type, expected LLOOutcomeTelemetry")
		}
		traits := []string{
			fmt.Sprint(m.DonId),
			fmt.Sprint(hex.EncodeToString(m.ConfigDigest)),
		}
		return strings.Join(traits, samplerDelimiter), nanosToSec(int64(m.ObservationTimestampNanoseconds)), nil
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
		return strings.Join(traits, samplerDelimiter), nanosToSec(int64(m.ObservationTimestampNanoseconds)), nil
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
		return strings.Join(traits, samplerDelimiter), nanosToSec(m.ObservationTimestamp), nil
	default:
		return "", 0, errUnsupportedTelemetryType
	}
}

func nanosToSec(n int64) int32 {
	return int32(n / int64(time.Second))
}
