package v2

import (
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"
)

// stubTimeProvider simulates a TimeProvider whose DON-time fetch blocks on
// consensus by advancing a fake clock for the duration of the call.
type stubTimeProvider struct {
	clock    *clockwork.FakeClock
	donBlock time.Duration
	donErr   error
}

func (s *stubTimeProvider) GetNodeTime() time.Time { return s.clock.Now() }

func (s *stubTimeProvider) GetDONTime() (time.Time, error) {
	s.clock.Advance(s.donBlock) // simulate blocking on the DON-time consensus round
	return s.clock.Now(), s.donErr
}

func TestSuspensionTracker_SumsSequentialIntervals(t *testing.T) {
	t.Parallel()

	var tr suspensionTracker
	tr.add(100 * time.Millisecond)
	tr.add(0)                     // ignored
	tr.add(-5 * time.Millisecond) // ignored
	tr.add(250 * time.Millisecond)
	require.Equal(t, 350*time.Millisecond, tr.total())
}

func TestSuspensionTracker_NilSafe(t *testing.T) {
	t.Parallel()

	var tr *suspensionTracker
	tr.add(time.Second) // must not panic
	require.Equal(t, time.Duration(0), tr.total())
}

func TestMeasuredTimeProvider_RecordsDONTimeWait(t *testing.T) {
	t.Parallel()

	clock := clockwork.NewFakeClock()
	tracker := &suspensionTracker{}
	stub := &stubTimeProvider{clock: clock, donBlock: 10 * time.Second}
	mtp := newMeasuredTimeProvider(stub, clock, tracker)

	_, err := mtp.GetDONTime()
	require.NoError(t, err)
	require.Equal(t, 10*time.Second, tracker.total(), "DON-time wait should be recorded as guest suspension")

	// A second blocking call accumulates on top of the first.
	stub.donBlock = 2 * time.Second
	_, err = mtp.GetDONTime()
	require.NoError(t, err)
	require.Equal(t, 12*time.Second, tracker.total())
}

func TestMeasuredTimeProvider_NodeTimeNotMeasured(t *testing.T) {
	t.Parallel()

	clock := clockwork.NewFakeClock()
	tracker := &suspensionTracker{}
	stub := &stubTimeProvider{clock: clock}
	mtp := newMeasuredTimeProvider(stub, clock, tracker)

	_ = mtp.GetNodeTime() // local, instant; not a guest suspension
	require.Equal(t, time.Duration(0), tracker.total())
}
