package v2

import (
	"sync/atomic"
	"time"

	"github.com/jonboulle/clockwork"
)

// suspensionTracker accumulates the wall-clock time during which the WASM guest
// is suspended inside a blocking DON-time host call (GetDONTime).
//
// Why this exists:
// Compute (RESOURCE_TYPE_COMPUTE) is metered as the end-to-end wall-clock around
// Module.Execute. That wall-clock also includes time the guest spends blocked in
// GetDONTime, which can stall for a full DON-time consensus round. That wait is
// not compute: it is non-deterministic across nodes (a node whose request lands
// in the current round returns in microseconds; a node whose request slips to
// the next round blocks for ~a round), and it has produced >100x divergence in
// per-node metered compute - and therefore consumed credits - for
// otherwise-identical executions on the same DON.
//
// Why summing is correct:
// The WASM guest is single-threaded, so its GetDONTime calls happen one at a
// time and never overlap with guest execution. Summing the individual wait
// intervals therefore yields the total non-compute DON-time wait, and
//
//	computeTime = wallClock - totalDONTimeWait
//
// Scope: this only covers the DON-time wait, which is the one blocking host call
// the engine observes (via TimeProvider). It is NOT a general "guest suspended"
// accountant - see the package note in the PR for the capability/secrets-await
// gap, which is only measurable inside chainlink-common's host runtime.
type suspensionTracker struct {
	totalNanos atomic.Int64
}

// add records a single guest-suspension interval. Safe for concurrent use and
// nil-safe so callers never have to guard the pointer.
func (s *suspensionTracker) add(d time.Duration) {
	if s == nil || d <= 0 {
		return
	}
	s.totalNanos.Add(int64(d))
}

// total returns the cumulative guest-suspended duration recorded so far.
func (s *suspensionTracker) total() time.Duration {
	if s == nil {
		return 0
	}
	return time.Duration(s.totalNanos.Load())
}

// measuredTimeProvider decorates a TimeProvider, recording the wall-clock spent
// in each blocking time call into the supplied suspension tracker.
//
// GetDONTime can block on DON-time consensus; that wait must not be billed as
// compute. GetNodeTime is served locally and is effectively instant, so it is
// inherited unchanged from the embedded TimeProvider (no measurement needed).
type measuredTimeProvider struct {
	TimeProvider
	clock   clockwork.Clock
	tracker *suspensionTracker
}

func newMeasuredTimeProvider(tp TimeProvider, clock clockwork.Clock, tracker *suspensionTracker) measuredTimeProvider {
	return measuredTimeProvider{TimeProvider: tp, clock: clock, tracker: tracker}
}

func (m measuredTimeProvider) GetDONTime() (time.Time, error) {
	start := m.clock.Now()
	t, err := m.TimeProvider.GetDONTime()
	m.tracker.add(m.clock.Now().Sub(start))
	return t, err
}
