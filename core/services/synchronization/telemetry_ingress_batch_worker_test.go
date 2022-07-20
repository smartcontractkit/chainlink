package synchronization_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/core/services/synchronization/mocks"
	"github.com/stretchr/testify/assert"
)

func TestTelemetryIngressWorker_BuildTelemBatchReq(t *testing.T) {
	telemPayload := synchronization.TelemPayload{
		Ctx:        context.Background(),
		Telemetry:  []byte("Mock telemetry"),
		ContractID: "0xa",
	}

	maxTelemBatchSize := 3
	chTelemetry := make(chan synchronization.TelemPayload, 10)
	worker := synchronization.NewTelemetryIngressBatchWorker(
		uint(maxTelemBatchSize),
		time.Millisecond*1,
		time.Second,
		new(mocks.TelemClient),
		&sync.WaitGroup{},
		make(chan struct{}),
		chTelemetry,
		"0xa",
		logger.TestLogger(t),
		false,
	)

	chTelemetry <- telemPayload
	chTelemetry <- telemPayload
	chTelemetry <- telemPayload
	chTelemetry <- telemPayload
	chTelemetry <- telemPayload

	// Batch request should not exceed the max batch size
	batchReq1 := worker.BuildTelemBatchReq()
	assert.Equal(t, batchReq1.ContractId, "0xa")
	assert.Len(t, batchReq1.Telemetry, maxTelemBatchSize)
	assert.Len(t, chTelemetry, 2)

	// Remainder of telemetry should be batched on next call
	batchReq2 := worker.BuildTelemBatchReq()
	assert.Equal(t, batchReq2.ContractId, "0xa")
	assert.Len(t, batchReq2.Telemetry, 2)
	assert.Len(t, chTelemetry, 0)
}
