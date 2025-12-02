package telem

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-data-streams/llo"

	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
)

func TestFingerprint(t *testing.T) {
	t.Parallel()

	ot := time.Now()
	bytes32 := sha256.Sum256([]byte(ot.String()))
	configDigest := bytes32[:]
	donID := uint32(2)
	streamID := uint32(123)
	channelID := uint32(234)

	tests := []struct {
		name        string
		msg         proto.Message
		typ         synchronization.TelemetryType
		fingerprint string
		ts          int32
		err         error
	}{
		{
			name: "successful observation",
			msg: &LLOObservationTelemetry{
				DonId:                donID,
				StreamId:             streamID,
				ConfigDigest:         configDigest,
				ObservationTimestamp: ot.UnixNano(),
			},
			typ:         synchronization.LLOObservation,
			fingerprint: fmt.Sprintf("%d-%d-%x", donID, streamID, configDigest),
			ts:          int32(ot.Unix()), //nolint:gosec // G115
			err:         nil,
		},
		{
			name: "successful outcome",
			msg: &llo.LLOOutcomeTelemetry{
				DonId:                           donID,
				ConfigDigest:                    configDigest,
				ObservationTimestampNanoseconds: uint64(ot.UnixNano()), //nolint:gosec // G115
			},
			typ:         synchronization.LLOOutcome,
			fingerprint: fmt.Sprintf("%d-%x", donID, configDigest),
			ts:          int32(ot.Unix()), //nolint:gosec // G115
			err:         nil,
		},
		{
			name: "successful report",
			msg: &llo.LLOReportTelemetry{
				DonId:                           donID,
				ChannelId:                       channelID,
				ConfigDigest:                    configDigest,
				ObservationTimestampNanoseconds: uint64(ot.UnixNano()), //nolint:gosec // G115
			},
			typ:         synchronization.LLOReport,
			fingerprint: fmt.Sprintf("%d-%d-%x", donID, channelID, configDigest),
			ts:          int32(ot.Unix()), //nolint:gosec // G115
			err:         nil,
		},
		{
			name: "successful bridge",
			msg: &LLOBridgeTelemetry{
				DonId:                donID,
				StreamId:             &streamID,
				SpecId:               345,
				BridgeAdapterName:    "bridge-adapter",
				ConfigDigest:         configDigest,
				ObservationTimestamp: ot.UnixNano(),
			},
			typ:         synchronization.PipelineBridge,
			fingerprint: fmt.Sprintf("%d-%d-%d-%s-%x", donID, streamID, 345, "bridge-adapter", configDigest),
			ts:          int32(ot.Unix()), //nolint:gosec // G115
			err:         nil,
		},
		{
			name: "unsupported telemetry type",
			typ:  synchronization.HeadReport,
			err:  errUnsupportedTelemetryType,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fp, ts, err := fingerprint(test.typ, test.msg)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.fingerprint, fp)
				assert.Equal(t, test.ts, ts)
			}
		})
	}
}

// TestSample ensures Sample() correctly chooses which messages to send.
func TestSample(t *testing.T) {
	t.Parallel()

	lggr := logger.TestSugared(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	samplr := newSampler(lggr, true)
	samplr.StartPruningLoop(ctx, &sync.WaitGroup{})

	t0 := time.Now()
	msg0 := &llo.LLOOutcomeTelemetry{
		DonId:                           2,
		ConfigDigest:                    []byte("digest"),
		ObservationTimestampNanoseconds: uint64(t0.UnixNano()), //nolint:gosec // G115
	}
	msg1 := &llo.LLOOutcomeTelemetry{
		DonId:                           2,
		ConfigDigest:                    []byte("digest"),
		ObservationTimestampNanoseconds: uint64(t0.Add(50 * time.Millisecond).UnixNano()), //nolint:gosec // G115
	}

	// Evaluate two messages from the same source with timestamp difference under our desired fidelity of 1s.
	// The first one should be sent, the second one should be omitted.
	shouldSend := samplr.Sample(synchronization.LLOOutcome, msg0)
	assert.True(t, shouldSend)
	shouldSend = samplr.Sample(synchronization.LLOOutcome, msg1)
	assert.False(t, shouldSend)
}

// TestPruningLoop ensures the pruning loop works as expected.
func TestPruningLoop(t *testing.T) {
	t.Parallel()

	lggr := logger.TestSugared(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	samplr := newSampler(lggr, true)
	// We need a prune time of at least one second in order to detect outdated entries.
	// The entries have second-based timestamps and cannot distinguish shorter intervals.
	samplr.prunePeriod = time.Second
	samplr.StartPruningLoop(ctx, &sync.WaitGroup{})

	msg := &llo.LLOOutcomeTelemetry{
		DonId:                           2,
		ConfigDigest:                    []byte("digest"),
		ObservationTimestampNanoseconds: uint64(time.Now().UnixNano()), //nolint:gosec // G115
	}
	fp, ots, err := fingerprint(synchronization.LLOOutcome, msg)
	require.NoError(t, err)

	hasSampleFn := func(s *sampler, fp string, ts int32) bool {
		s.samplesMu.Lock()
		defer s.samplesMu.Unlock()
		flag := s.samples[fp][ts] // we're OK with this panicking
		return flag != nil
	}

	msg2 := &llo.LLOOutcomeTelemetry{
		DonId:                           2,
		ConfigDigest:                    []byte("digest"),
		ObservationTimestampNanoseconds: uint64(time.Now().Add(10 * time.Second).UnixNano()), //nolint:gosec // G115
	}
	samplr.Sample(synchronization.LLOOutcome, msg2)
	fp2, ots2, _ := fingerprint(synchronization.LLOOutcome, msg2)
	require.True(t, hasSampleFn(samplr, fp2, ots2))

	require.False(t, hasSampleFn(samplr, fp, ots))

	// Sample a message.
	_ = samplr.Sample(synchronization.LLOOutcome, msg)
	// Ensure the message is among the samples.
	require.True(t, hasSampleFn(samplr, fp, ots))
	// Wait for the entry to expire.
	time.Sleep(3 * samplr.prunePeriod)
	// Ensure the message has been removed from the samples.
	require.False(t, hasSampleFn(samplr, fp, ots))
}

// TestPruningLoop_Exits ensures the pruning loop exits when its context is cancelled and doesn't leak goroutines.
func TestPruningLoop_Exits(t *testing.T) {
	t.Parallel()

	lggr := logger.TestSugared(t)
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	// Start the sampler and its loop. This increments the waitgroup's counter.
	samplr := newSampler(lggr, true)
	samplr.StartPruningLoop(ctx, wg)

	// Create a channel which will be closed when the waitgroup is done, i.e. when the loop is closed.
	ch := make(chan struct{})
	go func() {
		wg.Wait()
		close(ch)
	}()
	// A helper function to check if the channel is closed.
	doneFn := func(timeout time.Duration) bool {
		select {
		case <-ch:
			return true
		case <-time.After(timeout):
			return false
		}
	}

	// Assert the channel is not closed after 100ms.
	assert.False(t, doneFn(100*time.Millisecond))

	// Cancel the context and assert that the channel is closed within 100ms.
	cancel()
	assert.True(t, doneFn(100*time.Millisecond))
}
