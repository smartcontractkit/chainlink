package log

import (
	"testing"

	"github.com/test-go/testify/assert"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
)

var _ Listener = testListener{}

type testListener struct {
	jobID int32
}

func (tl testListener) JobID() int32        { return tl.jobID }
func (tl testListener) HandleLog(Broadcast) { panic("not implemented") }

func newTestListener(t *testing.T, jobID int32) testListener {
	return testListener{jobID}
}

func TestRegistrationsUnit_PanicsOnDoubleSubscribe(t *testing.T) {
	l := newTestListener(t, 1)

	contractAddr := testutils.NewAddress()
	opts := ListenerOpts{Contract: contractAddr}
	sub := &subscriber{l, opts}

	r := newRegistrations(logger.TestLogger(t), *testutils.FixtureChainID)
	r.addSubscriber(sub)

	l2 := newTestListener(t, 1)
	opts2 := ListenerOpts{Contract: contractAddr}
	sub2 := &subscriber{l2, opts2}

	// Different subscriber same job ID is ok
	r.addSubscriber(sub2)

	// Adding same subscriber twice is not ok
	assert.Panics(t, func() {
		r.addSubscriber(sub2)
	}, "expected adding same subscription twice to panic")

	r.removeSubscriber(sub)

	// Removing subscriber twice also panics
	assert.Panics(t, func() {
		r.removeSubscriber(sub)
	}, "expected removing a subscriber twice to panic")

	// Now we can add it again
	r.addSubscriber(sub)
}
