package autotelemetry21

import (
	"sync"
	"testing"
)

func TestNewAutomationCustomTelemetryService(t *testing.T) {
	// me := &MockMonitoringEndpoint{}
	// lggr := &MockLogger{}
	// blocksub := &MockBlockSubscriber{}
	// configTracker := &MockContractConfigTracker{}

	// service, err := NewAutomationCustomTelemetryService(me, lggr, blocksub, configTracker)
	// if err != nil {
	// 	t.Errorf("Expected no error, but got: %v", err)
	// }
}

type MockMonitoringEndpoint struct {
	sentLogs [][]byte
	lock     sync.RWMutex
}

func (me *MockMonitoringEndpoint) SendLog(log []byte) {
	me.lock.Lock()
	defer me.lock.Unlock()
	me.sentLogs = append(me.sentLogs, log)
}

func (me *MockMonitoringEndpoint) LogCount() int {
	me.lock.RLock()
	defer me.lock.RUnlock()
	return len(me.sentLogs)
}
