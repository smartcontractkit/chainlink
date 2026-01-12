package synchronization_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

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
		Network:       "ethereum-mainnet",
	}
	telemPayload3 := synchronization.TelemPayload{
		Telemetry:     []byte("Mock telem 3"),
		ContractID:    "0x3",
		TelemType:     synchronization.OCR2Functions,
		ChainSelector: 67890,
		Domain:        "functions",
		Entity:        "ocr.v2.functions.telemetry",
		Network:       "polygon-mainnet",
	}

	// Assert telemetry payloads for each contract are correctly sent to chip ingress
	var contractCounter1 atomic.Uint32
	var contractCounter3 atomic.Uint32
	chipClient.On("PublishBatch", mock.Anything, mock.Anything, mock.Anything).Return(&chipingress.PublishResponse{}, nil).Run(func(args mock.Arguments) {
		batch := args.Get(1).(*chipingress.CloudEventBatch)

		for _, protoEvent := range batch.Events {
			event, err := chipingress.ProtoToEvent(protoEvent)
			require.NoError(t, err)

			attrs := protoEvent.GetAttributes()
			contractID := attrs["contractid"].GetCeString()

			if contractID == "0x1" {
				contractCounter1.Add(1)
				assert.Equal(t, telemPayload1.Telemetry, event.Data())
				assert.Equal(t, string(synchronization.OCR), attrs["telemetrytype"].GetCeString())
				assert.Equal(t, telemPayload1.Network, attrs["networkname"].GetCeString())
				assert.Empty(t, attrs["nodeoperatorname"].GetCeString())
				assert.Empty(t, attrs["nodename"].GetCeString())
				assert.Empty(t, attrs["nodecsapublickey"].GetCeString())
			}
			if contractID == "0x3" {
				contractCounter3.Add(1)
				assert.Equal(t, telemPayload3.Telemetry, event.Data())
				assert.Equal(t, string(synchronization.OCR2Functions), attrs["telemetrytype"].GetCeString())
				assert.Equal(t, telemPayload3.Network, attrs["networkname"].GetCeString())
				assert.Empty(t, attrs["nodeoperatorname"].GetCeString())
				assert.Empty(t, attrs["nodename"].GetCeString())
				assert.Empty(t, attrs["nodecsapublickey"].GetCeString())
			}
		}
	})

	// Send telemetry
	testCtx := testutils.Context(t)
	chipIngressClient.Send(testCtx, telemPayload1)
	chipIngressClient.Send(testCtx, telemPayload3)
	time.Sleep(sendInterval * 2)
	chipIngressClient.Send(testCtx, telemPayload1)
	chipIngressClient.Send(testCtx, telemPayload1)

	// Wait for the telemetry to be handled
	g.Eventually(func() []uint32 {
		return []uint32{contractCounter1.Load(), contractCounter3.Load()}
	}, testutils.WaitTimeout(t)).Should(gomega.Equal([]uint32{3, 1}))
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
		Network:       "ethereum-mainnet",
	}

	var batchCount atomic.Uint32
	chipClient.On("PublishBatch", mock.Anything, mock.Anything, mock.Anything).Return(&chipingress.PublishResponse{}, nil).Run(func(args mock.Arguments) {
		batchCount.Add(1)
	})

	testCtx := testutils.Context(t)
	// Send multiple messages to trigger multiple batches
	for i := 0; i < 10; i++ {
		chipIngressClient.Send(testCtx, telemPayload)
		if i%3 == 0 {
			time.Sleep(sendInterval * 2) // Allow batch to be sent
		}
	}

	// Wait for batches to be sent
	g.Eventually(batchCount.Load, testutils.WaitTimeout(t)).Should(gomega.BeNumerically(">=", 2))
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
		Network:       "ethereum-mainnet",
	}

	payloadOCR2 := synchronization.TelemPayload{
		Telemetry:     []byte("OCR2 telemetry"),
		ContractID:    "0x123",
		TelemType:     synchronization.OCR2Median,
		ChainSelector: 1,
		Domain:        "data-feeds",
		Entity:        "ocr.v2.median.telemetry",
		Network:       "ethereum-mainnet",
	}

	var ocrCount, ocr2Count atomic.Uint32
	chipClient.On("PublishBatch", mock.Anything, mock.Anything, mock.Anything).Return(&chipingress.PublishResponse{}, nil).Run(func(args mock.Arguments) {
		batch := args.Get(1).(*chipingress.CloudEventBatch)
		for _, protoEvent := range batch.Events {
			attrs := protoEvent.GetAttributes()
			telemType := attrs["telemetrytype"].GetCeString()
			if telemType == string(synchronization.OCR) {
				ocrCount.Add(1)
			} else if telemType == string(synchronization.OCR2Median) {
				ocr2Count.Add(1)
			}
		}
	})

	testCtx := testutils.Context(t)
	chipIngressClient.Send(testCtx, payloadOCR)
	chipIngressClient.Send(testCtx, payloadOCR2)

	g.Eventually(func() []uint32 {
		return []uint32{ocrCount.Load(), ocr2Count.Load()}
	}, testutils.WaitTimeout(t)).Should(gomega.Equal([]uint32{1, 1}))
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
		Network:       "ethereum-mainnet",
	}

	// Create a cancelled context
	ctx, cancel := context.WithCancel(testutils.Context(t))
	cancel()

	// Should not panic or block when context is cancelled
	chipIngressClient.Send(ctx, telemPayload)
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
		Network:       "ethereum-mainnet",
	}

	var messageCount atomic.Uint32
	chipClient.On("PublishBatch", mock.Anything, mock.Anything, mock.Anything).Return(&chipingress.PublishResponse{}, nil).Run(func(args mock.Arguments) {
		batch := args.Get(1).(*chipingress.CloudEventBatch)
		// #nosec G115 -- len() returns non-negative int, safe to convert to uint32
		messageCount.Add(uint32(len(batch.Events)))
	})

	testCtx := testutils.Context(t)
	// Send multiple messages with same contract and type - should reuse worker
	for i := 0; i < 5; i++ {
		chipIngressClient.Send(testCtx, telemPayload)
	}

	// Wait for all messages to be sent
	g.Eventually(messageCount.Load, testutils.WaitTimeout(t)).Should(gomega.Equal(uint32(5)))
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
		Network:       "ethereum-mainnet",
	}

	var capturedChainSelector atomic.Value
	chipClient.On("PublishBatch", mock.Anything, mock.Anything, mock.Anything).Return(&chipingress.PublishResponse{}, nil).Run(func(args mock.Arguments) {
		batch := args.Get(1).(*chipingress.CloudEventBatch)
		if len(batch.Events) > 0 {
			attrs := batch.Events[0].GetAttributes()
			capturedChainSelector.Store(attrs["chainselector"].GetCeString())
		}
	})

	testCtx := testutils.Context(t)
	chipIngressClient.Send(testCtx, telemPayload)

	g.Eventually(func() string {
		v := capturedChainSelector.Load()
		if v == nil {
			return ""
		}
		return v.(string)
	}, testutils.WaitTimeout(t)).Should(gomega.Equal("123456789"))
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
	g.Eventually(pingCallCount.Load, testutils.WaitTimeout(t), 100*time.Millisecond).Should(gomega.BeNumerically(">=", 2))
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
	g.Eventually(pingCallCount.Load, testutils.WaitTimeout(t), 100*time.Millisecond).Should(gomega.BeNumerically(">=", 2))
}
