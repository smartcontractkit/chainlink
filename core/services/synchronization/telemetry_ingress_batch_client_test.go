package synchronization_test

import (
	"context"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/csakey"
	ksmocks "github.com/smartcontractkit/chainlink/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/core/services/synchronization/mocks"
	telemPb "github.com/smartcontractkit/chainlink/core/services/synchronization/telem"
)

func TestTelemetryIngressBatchClient_HappyPath(t *testing.T) {
	g := gomega.NewWithT(t)

	// Create mocks
	telemClient := new(mocks.TelemClient)
	csaKeystore := new(ksmocks.CSA)

	// Set mock handlers for keystore
	key := cltest.DefaultCSAKey
	keyList := []csakey.KeyV2{key}
	csaKeystore.On("GetAll").Return(keyList, nil)

	// Wire up the telem ingress client
	url := &url.URL{}
	serverPubKeyHex := "33333333333"
	sendInterval := time.Millisecond * 5
	telemIngressClient := synchronization.NewTestTelemetryIngressBatchClient(t, url, serverPubKeyHex, csaKeystore, false, telemClient, sendInterval)
	require.NoError(t, telemIngressClient.Start())

	// Create telemetry payloads for different contracts
	telemPayload1 := synchronization.TelemPayload{
		Ctx:        context.Background(),
		Telemetry:  []byte("Mock telem 1"),
		ContractID: "0x1",
	}
	telemPayload2 := synchronization.TelemPayload{
		Ctx:        context.Background(),
		Telemetry:  []byte("Mock telem 2"),
		ContractID: "0x2",
	}
	telemPayload3 := synchronization.TelemPayload{
		Ctx:        context.Background(),
		Telemetry:  []byte("Mock telem 3"),
		ContractID: "0x3",
	}

	// Assert telemetry payloads for each contract are correctly sent to wsrpc
	var contractCounter1 atomic.Uint32
	var contractCounter2 atomic.Uint32
	var contractCounter3 atomic.Uint32
	telemClient.On("TelemBatch", mock.Anything, mock.Anything).Return(nil, nil).Run(func(args mock.Arguments) {
		telemBatchReq := args.Get(1).(*telemPb.TelemBatchRequest)

		if telemBatchReq.ContractId == "0x1" {
			for _, telem := range telemBatchReq.Telemetry {
				contractCounter1.Inc()
				assert.Equal(t, telemPayload1.Telemetry, telem)
			}
		}
		if telemBatchReq.ContractId == "0x2" {
			for _, telem := range telemBatchReq.Telemetry {
				contractCounter2.Inc()
				assert.Equal(t, telemPayload2.Telemetry, telem)
			}
		}
		if telemBatchReq.ContractId == "0x3" {
			for _, telem := range telemBatchReq.Telemetry {
				contractCounter3.Inc()
				assert.Equal(t, telemPayload3.Telemetry, telem)
			}
		}
	})

	// Send telemetry
	telemIngressClient.Send(telemPayload1)
	telemIngressClient.Send(telemPayload2)
	telemIngressClient.Send(telemPayload3)
	time.Sleep(sendInterval * 2)
	telemIngressClient.Send(telemPayload1)
	telemIngressClient.Send(telemPayload1)
	telemIngressClient.Send(telemPayload2)

	// Wait for the telemetry to be handled
	g.Eventually(contractCounter1.Load, "5s").Should(gomega.Equal(uint32(3)))
	g.Eventually(contractCounter2.Load, "5s").Should(gomega.Equal(uint32(2)))
	g.Eventually(contractCounter3.Load, "5s").Should(gomega.Equal(uint32(1)))

	// Client should shut down
	telemIngressClient.Close()
}

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
