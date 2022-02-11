package ocrcommon

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
)

type testService struct {
}

func (ts *testService) Start() error {
	return nil
}

func (ts *testService) Close() error {
	return nil
}

func TestDependentService(t *testing.T) {
	ready := make(chan struct{})
	lggr := logger.TestLogger(t)
	ocrLggr := logger.NewOCRWrapper(lggr, false, func(string) {})
	ds := NewDependentOCRService(ready, &testService{}, ocrLggr)
	done := make(chan struct{})
	go func() {
		defer close(done)
		ds.run()
	}()
	ready <- struct{}{}
	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Errorf("service did not shutdown in time")
	}
}
