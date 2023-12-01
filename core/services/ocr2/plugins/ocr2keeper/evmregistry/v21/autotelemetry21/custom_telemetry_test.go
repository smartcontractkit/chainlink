package autotelemetry21

import (
	"sync"
	"testing"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"

	headtracker "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	evm "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21"
)

// const historySize = 4
// const blockSize = int64(4)
const finality = uint32(4)

func TestNewAutomationCustomTelemetryService(t *testing.T) {
	me := &MockMonitoringEndpoint{}
	lggr := logger.TestLogger(t)
	var hb headtracker.HeadBroadcaster
	var lp logpoller.LogPoller

	bs := evm.NewBlockSubscriber(hb, lp, finality, lggr)
	// configTracker := &MockContractConfigTracker{}
	var configTracker types.ContractConfigTracker

	service, err := NewAutomationCustomTelemetryService(me, lggr, bs, configTracker)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	service.monitoringEndpoint.SendLog([]byte("test"))
	assert.Equal(t, me.LogCount(), 1)
	service.monitoringEndpoint.SendLog([]byte("test2"))
	assert.Equal(t, me.LogCount(), 2)
	service.Close()
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
