package ocrcommon

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/logger"
)

type testService struct {
	started, closed chan struct{}
}

func (ts *testService) Start() error {
	ts.started <- struct{}{}
	return nil
}

func (ts *testService) Close() error {
	ts.closed <- struct{}{}
	return nil
}

func TestDependentService(t *testing.T) {
	ready := make(chan struct{})
	lggr := logger.TestLogger(t)
	ocrLggr := logger.NewOCRWrapper(lggr, false, func(string) {})
	started := make(chan struct{}, 1)
	closed := make(chan struct{}, 1)
	ds := NewDependentOCRService(ready, &testService{started: started, closed: closed}, ocrLggr)
	done := make(chan struct{})
	go func() {
		defer close(done)
		// Should block until ready is called
		ds.run()
	}()

	// It should NOT start until we send ready
	select {
	case <-started:
		t.Errorf("service started before signal")
	case <-time.After(1 * time.Second):
	}

	// Send ready
	ready <- struct{}{}

	// It should start and return
	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Errorf("dependent service did not start")
	}

	// Closing should close the started service
	require.NoError(t, ds.Close())
	select {
	case <-closed:
	case <-time.After(1 * time.Second):
		t.Errorf("service did not shutdown in time")
	}
}
