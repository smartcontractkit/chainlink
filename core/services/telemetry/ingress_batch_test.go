package telemetry_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
)

func TestIngressAgentBatch(t *testing.T) {
	telemetryBatchClient := mocks.NewTelemetryService(t)
	ingressAgentBatch := telemetry.NewIngressAgentWrapper(telemetryBatchClient)
	monitoringEndpoint := ingressAgentBatch.GenMonitoringEndpoint("test-network", "test-chainID", "0xa", synchronization.OCR)

	// Handle the Send call and store the telem
	var telemPayload synchronization.TelemPayload
	telemetryBatchClient.On("Send", mock.Anything, mock.AnythingOfType("[]uint8"), mock.AnythingOfType("string"), mock.AnythingOfType("TelemetryType")).Return().Run(func(args mock.Arguments) {
		telemPayload = synchronization.TelemPayload{
			Telemetry:  args[1].([]byte),
			ContractID: args[2].(string),
			TelemType:  args[3].(synchronization.TelemetryType),
		}
	})

	// Send the log to the monitoring endpoint
	log := []byte("test log")
	monitoringEndpoint.SendLog(log)

	// Telemetry should be sent to the mock as expected
	assert.Equal(t, log, telemPayload.Telemetry)
	assert.Equal(t, synchronization.OCR, telemPayload.TelemType)
	assert.Equal(t, "0xa", telemPayload.ContractID)
}
