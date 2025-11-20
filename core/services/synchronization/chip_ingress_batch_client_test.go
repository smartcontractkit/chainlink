package synchronization_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink-common/pkg/chipingress"
	chipingressmocks "github.com/smartcontractkit/chainlink-common/pkg/chipingress/mocks"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
)

func TestChipIngressBatchClient_HappyPath(t *testing.T) {
	g := gomega.NewWithT(t)

	// Create mocks
	chipClient := chipingressmocks.NewClient(t)
	chipClient.On("Ping", mock.Anything, mock.Anything, mock.Anything).Return(&chipingress.PingResponse{}, nil).Maybe()

	// Wire up the chip ingress client with instant interval for testing
	sendInterval := time.Nanosecond
	chipIngressClient := synchronization.NewTestChipIngressBatchClient(t, chipClient, false, sendInterval)
	servicetest.Run(t, chipIngressClient)

	// Create telemetry payloads for different contracts
	telemPayload1 := synchronization.TelemPayload{
		Telemetry:     []byte("Mock telem 1"),
		ContractID:    "0x1",
		TelemType:     synchronization.OCR,
		ChainSelector: 12345,
		Domain:        "data-feeds",
		Entity:        "ocr.v1.telemetry",
	}
	telemPayload3 := synchronization.TelemPayload{
		Telemetry:     []byte("Mock telem 3"),
		ContractID:    "0x3",
		TelemType:     synchronization.OCR2Functions,
		ChainSelector: 67890,
		Domain:        "functions",
		Entity:        "ocr.v2.functions.telemetry",
	}

	// Assert telemetry payloads for each contract are correctly sent to chip ingress
	var contractCounter1 atomic.Uint32
	var contractCounter3 atomic.Uint32
	chipClient.On("PublishBatch", mock.Anything, mock.Anything, mock.Anything).Return(&chipingress.PublishResponse{}, nil).Run(func(args mock.Arguments) {
		batch := args.Get(1).(*chipingress.CloudEventBatch)

		for _, protoEvent := range batch.Events {
			event, err := chipingress.ProtoToEvent(protoEvent)
			assert.NoError(t, err)

			attrs := protoEvent.GetAttributes()
			contractID := attrs["contractid"].GetCeString()

			if contractID == "0x1" {
				contractCounter1.Add(1)
				assert.Equal(t, telemPayload1.Telemetry, event.Data())
				assert.Equal(t, string(synchronization.OCR), attrs["telemtype"].GetCeString())
			}
			if contractID == "0x3" {
				contractCounter3.Add(1)
				assert.Equal(t, telemPayload3.Telemetry, event.Data())
				assert.Equal(t, string(synchronization.OCR2Functions), attrs["telemtype"].GetCeString())
			}
		}
	})

	// Send telemetry
	testCtx := testutils.Context(t)
	chipIngressClient.Send(testCtx, telemPayload1.Telemetry, telemPayload1.ContractID, telemPayload1.TelemType, telemPayload1.ChainSelector, telemPayload1.Domain, telemPayload1.Entity)
	chipIngressClient.Send(testCtx, telemPayload3.Telemetry, telemPayload3.ContractID, telemPayload3.TelemType, telemPayload3.ChainSelector, telemPayload3.Domain, telemPayload3.Entity)
	time.Sleep(sendInterval * 2)
	chipIngressClient.Send(testCtx, telemPayload1.Telemetry, telemPayload1.ContractID, telemPayload1.TelemType, telemPayload1.ChainSelector, telemPayload1.Domain, telemPayload1.Entity)
	chipIngressClient.Send(testCtx, telemPayload1.Telemetry, telemPayload1.ContractID, telemPayload1.TelemType, telemPayload1.ChainSelector, telemPayload1.Domain, telemPayload1.Entity)

	// Wait for the telemetry to be handled
	g.Eventually(func() []uint32 {
		return []uint32{contractCounter1.Load(), contractCounter3.Load()}
	}).Should(gomega.Equal([]uint32{3, 1}))
}

func TestChipIngressBatchClient_MultipleBatches(t *testing.T) {
	g := gomega.NewWithT(t)

	chipClient := chipingressmocks.NewClient(t)
	chipClient.On("Ping", mock.Anything, mock.Anything, mock.Anything).Return(&chipingress.PingResponse{}, nil).Maybe()

	sendInterval := time.Millisecond
	chipIngressClient := synchronization.NewTestChipIngressBatchClient(t, chipClient, false, sendInterval)
	servicetest.Run(t, chipIngressClient)

	telemPayload := synchronization.TelemPayload{
		Telemetry:     []byte("Test telemetry"),
		ContractID:    "0xabc",
		TelemType:     synchronization.OCR2Median,
		ChainSelector: 12345,
		Domain:        "data-feeds",
		Entity:        "ocr.v2.median.telemetry",
	}

	var batchCount atomic.Uint32
	chipClient.On("PublishBatch", mock.Anything, mock.Anything, mock.Anything).Return(&chipingress.PublishResponse{}, nil).Run(func(args mock.Arguments) {
		batchCount.Add(1)
	})

	testCtx := testutils.Context(t)
	// Send multiple messages to trigger multiple batches
	for i := 0; i < 10; i++ {
		chipIngressClient.Send(testCtx, telemPayload.Telemetry, telemPayload.ContractID, telemPayload.TelemType, telemPayload.ChainSelector, telemPayload.Domain, telemPayload.Entity)
		if i%3 == 0 {
			time.Sleep(sendInterval * 2) // Allow batch to be sent
		}
	}

	// Wait for batches to be sent
	g.Eventually(func() uint32 {
		return batchCount.Load()
	}, 200*time.Millisecond).Should(gomega.BeNumerically(">=", 2))
}

func TestChipIngressBatchClient_DifferentTelemetryTypes(t *testing.T) {
	g := gomega.NewWithT(t)

	chipClient := chipingressmocks.NewClient(t)
	chipClient.On("Ping", mock.Anything, mock.Anything, mock.Anything).Return(&chipingress.PingResponse{}, nil).Maybe()

	sendInterval := time.Nanosecond
	chipIngressClient := synchronization.NewTestChipIngressBatchClient(t, chipClient, false, sendInterval)
	servicetest.Run(t, chipIngressClient)

	// Create payloads with different telemetry types but same contract
	payloadOCR := synchronization.TelemPayload{
		Telemetry:     []byte("OCR telemetry"),
		ContractID:    "0x123",
		TelemType:     synchronization.OCR,
		ChainSelector: 1,
		Domain:        "data-feeds",
		Entity:        "ocr.v1.telemetry",
	}

	payloadOCR2 := synchronization.TelemPayload{
		Telemetry:     []byte("OCR2 telemetry"),
		ContractID:    "0x123",
		TelemType:     synchronization.OCR2Median,
		ChainSelector: 1,
		Domain:        "data-feeds",
		Entity:        "ocr.v2.median.telemetry",
	}

	var ocrCount, ocr2Count atomic.Uint32
	chipClient.On("PublishBatch", mock.Anything, mock.Anything, mock.Anything).Return(&chipingress.PublishResponse{}, nil).Run(func(args mock.Arguments) {
		batch := args.Get(1).(*chipingress.CloudEventBatch)
		for _, protoEvent := range batch.Events {
			attrs := protoEvent.GetAttributes()
			telemType := attrs["telemtype"].GetCeString()
			if telemType == string(synchronization.OCR) {
				ocrCount.Add(1)
			} else if telemType == string(synchronization.OCR2Median) {
				ocr2Count.Add(1)
			}
		}
	})

	testCtx := testutils.Context(t)
	chipIngressClient.Send(testCtx, payloadOCR.Telemetry, payloadOCR.ContractID, payloadOCR.TelemType, payloadOCR.ChainSelector, payloadOCR.Domain, payloadOCR.Entity)
	chipIngressClient.Send(testCtx, payloadOCR2.Telemetry, payloadOCR2.ContractID, payloadOCR2.TelemType, payloadOCR2.ChainSelector, payloadOCR2.Domain, payloadOCR2.Entity)

	g.Eventually(func() []uint32 {
		return []uint32{ocrCount.Load(), ocr2Count.Load()}
	}).Should(gomega.Equal([]uint32{1, 1}))
}

func TestChipIngressBatchClient_ContextCancellation(t *testing.T) {
	chipClient := chipingressmocks.NewClient(t)
	chipClient.On("Ping", mock.Anything, mock.Anything, mock.Anything).Return(&chipingress.PingResponse{}, nil).Maybe()
	chipClient.On("PublishBatch", mock.Anything, mock.Anything, mock.Anything).Return(&chipingress.PublishResponse{}, nil).Maybe()

	sendInterval := time.Nanosecond
	chipIngressClient := synchronization.NewTestChipIngressBatchClient(t, chipClient, false, sendInterval)
	servicetest.Run(t, chipIngressClient)

	telemPayload := synchronization.TelemPayload{
		Telemetry:     []byte("Test telemetry"),
		ContractID:    "0xdef",
		TelemType:     synchronization.OCR2Functions,
		ChainSelector: 67890,
		Domain:        "functions",
		Entity:        "ocr.v2.functions.telemetry",
	}

	// Create a cancelled context
	ctx, cancel := context.WithCancel(testutils.Context(t))
	cancel()

	// Should not panic or block when context is cancelled
	chipIngressClient.Send(ctx, telemPayload.Telemetry, telemPayload.ContractID, telemPayload.TelemType, telemPayload.ChainSelector, telemPayload.Domain, telemPayload.Entity)
}

func TestChipIngressBatchClient_WorkerReuse(t *testing.T) {
	g := gomega.NewWithT(t)

	chipClient := chipingressmocks.NewClient(t)
	chipClient.On("Ping", mock.Anything, mock.Anything, mock.Anything).Return(&chipingress.PingResponse{}, nil).Maybe()

	sendInterval := time.Nanosecond
	chipIngressClient := synchronization.NewTestChipIngressBatchClient(t, chipClient, false, sendInterval)
	servicetest.Run(t, chipIngressClient)

	telemPayload := synchronization.TelemPayload{
		Telemetry:     []byte("Test telemetry"),
		ContractID:    "0xabc",
		TelemType:     synchronization.OCR2Automation,
		ChainSelector: 99999,
		Domain:        "automation",
		Entity:        "ocr.v2.automation.telemetry",
	}

	var messageCount atomic.Uint32
	chipClient.On("PublishBatch", mock.Anything, mock.Anything, mock.Anything).Return(&chipingress.PublishResponse{}, nil).Run(func(args mock.Arguments) {
		batch := args.Get(1).(*chipingress.CloudEventBatch)
		messageCount.Add(uint32(len(batch.Events)))
	})

	testCtx := testutils.Context(t)
	// Send multiple messages with same contract and type - should reuse worker
	for i := 0; i < 5; i++ {
		chipIngressClient.Send(testCtx, telemPayload.Telemetry, telemPayload.ContractID, telemPayload.TelemType, telemPayload.ChainSelector, telemPayload.Domain, telemPayload.Entity)
	}

	// Wait for all messages to be sent
	g.Eventually(func() uint32 {
		return messageCount.Load()
	}).Should(gomega.Equal(uint32(5)))
}

func TestChipIngressBatchClient_ChainSelectorInAttributes(t *testing.T) {
	g := gomega.NewWithT(t)

	chipClient := chipingressmocks.NewClient(t)
	chipClient.On("Ping", mock.Anything, mock.Anything, mock.Anything).Return(&chipingress.PingResponse{}, nil).Maybe()

	sendInterval := time.Nanosecond
	chipIngressClient := synchronization.NewTestChipIngressBatchClient(t, chipClient, false, sendInterval)
	servicetest.Run(t, chipIngressClient)

	expectedChainSelector := uint64(123456789)
	telemPayload := synchronization.TelemPayload{
		Telemetry:     []byte("Test telemetry"),
		ContractID:    "0xtest",
		TelemType:     synchronization.OCR3CCIPCommit,
		ChainSelector: expectedChainSelector,
		Domain:        "ccip",
		Entity:        "ocr.v3.ccip.commit.telemetry",
	}

	var capturedChainSelector string
	chipClient.On("PublishBatch", mock.Anything, mock.Anything, mock.Anything).Return(&chipingress.PublishResponse{}, nil).Run(func(args mock.Arguments) {
		batch := args.Get(1).(*chipingress.CloudEventBatch)
		if len(batch.Events) > 0 {
			attrs := batch.Events[0].GetAttributes()
			capturedChainSelector = attrs["chainselector"].GetCeString()
		}
	})

	testCtx := testutils.Context(t)
	chipIngressClient.Send(testCtx, telemPayload.Telemetry, telemPayload.ContractID, telemPayload.TelemType, telemPayload.ChainSelector, telemPayload.Domain, telemPayload.Entity)

	g.Eventually(func() string {
		return capturedChainSelector
	}).Should(gomega.Equal("123456789"))
}

func TestChipIngressBatchClient_HealthMonitoring(t *testing.T) {
	g := gomega.NewWithT(t)

	chipClient := chipingressmocks.NewClient(t)

	sendInterval := time.Nanosecond
	chipIngressClient := synchronization.NewTestChipIngressBatchClient(t, chipClient, false, sendInterval)

	// Mock Ping to succeed initially
	var pingCallCount atomic.Uint32
	chipClient.On("Ping", mock.Anything, mock.Anything, mock.Anything).Return(&chipingress.PingResponse{}, nil).Run(func(args mock.Arguments) {
		pingCallCount.Add(1)
	})

	// Start the client which should start health monitoring
	servicetest.Run(t, chipIngressClient)

	// Wait for at least 2 ping calls to verify health monitoring is running
	g.Eventually(func() uint32 {
		return pingCallCount.Load()
	}, 15*time.Second, 100*time.Millisecond).Should(gomega.BeNumerically(">=", 2))
}

func TestChipIngressBatchClient_HealthMonitoring_PingFailure(t *testing.T) {
	g := gomega.NewWithT(t)

	chipClient := chipingressmocks.NewClient(t)

	sendInterval := time.Nanosecond
	chipIngressClient := synchronization.NewTestChipIngressBatchClient(t, chipClient, false, sendInterval)

	// Mock Ping to fail
	var pingCallCount atomic.Uint32
	chipClient.On("Ping", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Run(func(args mock.Arguments) {
		pingCallCount.Add(1)
	})

	// Start the client which should start health monitoring
	servicetest.Run(t, chipIngressClient)

	// Wait for at least 2 ping calls to verify health monitoring continues despite failures
	g.Eventually(func() uint32 {
		return pingCallCount.Load()
	}, 15*time.Second, 100*time.Millisecond).Should(gomega.BeNumerically(">=", 2))
}
