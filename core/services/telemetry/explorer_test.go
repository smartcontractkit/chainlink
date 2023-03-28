package telemetry_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
)

func TestExplorerAgent(t *testing.T) {
	explorerClient := mocks.NewExplorerClient(t)
	explorerAgent := telemetry.NewExplorerAgent(explorerClient)
	monitoringEndpoint := explorerAgent.GenMonitoringEndpoint("0xa", synchronization.OCR)

	// Handle the Send call and store the logs
	var sentLog []byte
	explorerClient.On("Send", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("[]uint8"), synchronization.ExplorerBinaryMessage).Return().Run(func(args mock.Arguments) {
		sentLog = args[1].([]byte)
	})

	// Send the log to the monitoring endpoint
	log := []byte("test log")
	monitoringEndpoint.SendLog(log)

	// Logs should be sent to the mock as they were passed in
	assert.Equal(t, log, sentLog)
}
