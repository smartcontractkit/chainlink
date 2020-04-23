package fluxmonitor

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/eth/contracts"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func (fm *concreteFluxMonitor) MockLogBroadcaster() *mockLogBroadcaster {
	mock := mockLogBroadcaster{}
	fm.logBroadcaster = &mock
	return &mock
}

type mockLogBroadcaster struct {
	Started bool
	utils.DependentAwaiter
}

func (mlb *mockLogBroadcaster) Start() {
	mlb.Started = true
}
func (mlb *mockLogBroadcaster) Register(common.Address, eth.LogListener) bool {
	return false
}
func (mlb *mockLogBroadcaster) Unregister(common.Address, eth.LogListener) {}
func (mlb *mockLogBroadcaster) Stop()                                      {}

type MockableLogBroadcaster interface {
	MockLogBroadcaster() *mockLogBroadcaster
}

func TestFluxMonitor_PollingDeviationChecker_HandlesNilLogs(t *testing.T) {
	p := &PollingDeviationChecker{}
	var logNewRound *contracts.LogNewRound
	assert.NotPanics(t, func() {
		p.respondToLog(logNewRound)
	})
	var logAnswerUpdated *contracts.LogAnswerUpdated
	assert.NotPanics(t, func() {
		p.respondToLog(logAnswerUpdated)
	})
	var randomType interface{}
	assert.NotPanics(t, func() {
		p.respondToLog(randomType)
	})
}
