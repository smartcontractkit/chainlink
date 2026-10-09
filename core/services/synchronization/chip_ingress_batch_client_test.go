package synchronization

import (
	"context"
	"maps"
	"slices"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/chipingress"
	chipingressmocks "github.com/smartcontractkit/chainlink-common/pkg/chipingress/mocks"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func newTestPayload(contractID string) TelemPayload {
	return TelemPayload{
		Telemetry:     []byte("telemetry-payload"),
		TelemType:     OCR2Median,
		ContractID:    contractID,
		ChainSelector: 12345,
		Network:       "EVM",
	}
}

// publishedBatches collects PublishBatch requests; the batch client publishes
// from its own goroutines, so access is guarded.
type publishedBatches struct {
	mu      sync.Mutex
	batches [][]*chipingress.CloudEventPb
}

func (p *publishedBatches) add(events []*chipingress.CloudEventPb) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.batches = append(p.batches, events)
}

func (p *publishedBatches) snapshot() [][]*chipingress.CloudEventPb {
	p.mu.Lock()
	defer p.mu.Unlock()
	return slices.Clone(p.batches)
}

type publishFunc func(*chipingress.CloudEventBatch) (*chipingress.PublishResponse, error)

// publishOK acknowledges every event; the batch client expects one positional result per event.
func publishOK(b *chipingress.CloudEventBatch) (*chipingress.PublishResponse, error) {
	results := make([]*chipingress.PublishResult, len(b.Events))
	for i := range results {
		results[i] = &chipingress.PublishResult{}
	}
	return &chipingress.PublishResponse{Results: results}, nil
}

// startBatchClient starts a batch client over a mock chip client whose PublishBatch
// records each request and then answers with publish.
func startBatchClient(t *testing.T, publish publishFunc) (ChipIngressService, *publishedBatches) {
	t.Helper()
	chipClient := chipingressmocks.NewClient(t)
	chipClient.On("Ping", mock.Anything, mock.Anything, mock.Anything).Return(&chipingress.PingResponse{}, nil).Maybe()
	chipClient.On("Close").Return(nil).Maybe() // batch.Client.Stop() closes the underlying client
	published := &publishedBatches{}
	chipClient.EXPECT().PublishBatch(mock.Anything, mock.Anything, mock.Anything).
		RunAndReturn(func(_ context.Context, b *chipingress.CloudEventBatch, _ ...grpc.CallOption) (*chipingress.PublishResponse, error) {
			published.add(b.Events)
			return publish(b)
		}).Maybe()

	client := NewTestChipIngressBatchClient(t, chipClient, false, time.Nanosecond)
	servicetest.Run(t, client)
	return client, published
}

// awaitFirstEvent waits for exactly one published batch holding one event and decodes it.
func awaitFirstEvent(t *testing.T, published *publishedBatches) chipingress.CloudEvent {
	t.Helper()
	require.EventuallyWithT(t, func(c *assert.CollectT) {
		assert.Len(c, published.snapshot(), 1)
	}, testutils.WaitTimeout(t), 10*time.Millisecond)
	batches := published.snapshot()
	require.Len(t, batches[0], 1)
	event, err := chipingress.ProtoToEvent(batches[0][0])
	require.NoError(t, err)
	return event
}

func counterValue(vec *prometheus.CounterVec, telemType TelemetryType) float64 {
	return testutil.ToFloat64(vec.WithLabelValues(chipIngressEndpointLabel, string(telemType)))
}

func TestChipIngressBatchClient_Event(t *testing.T) {
	t.Parallel()
	client, published := startBatchClient(t, publishOK)

	client.Send(t.Context(), newTestPayload("0xabc"))
	event := awaitFirstEvent(t, published)

	// Source/type come from the telemetry type and pick the Kafka topic.
	domain, entity, err := TelemetryTypeToDomainAndEntity(OCR2Median)
	require.NoError(t, err)
	assert.Equal(t, domain, event.Source())
	assert.Equal(t, entity, event.Type())

	// Exact extension set from the OTI-parity table plus the ones the
	// library stamps (seqnum by the batch client, recordedtime by NewEvent).
	exts := event.Extensions()
	assert.ElementsMatch(t, []string{
		"legacytelemetry", "telemetrytype", "chainselector", "networkname",
		"contractid", "csapublickey", "partitionkey", "sentat", "receivedat",
		"seqnum", "recordedtime",
	}, slices.Collect(maps.Keys(exts)))
	assert.Equal(t, "true", exts["legacytelemetry"])
	assert.Equal(t, string(OCR2Median), exts["telemetrytype"])
	assert.Equal(t, "12345", exts["chainselector"])
	assert.Equal(t, "EVM", exts["networkname"])
	assert.Equal(t, "0xabc", exts["contractid"])
	assert.Equal(t, "deadbeef", exts["csapublickey"])
	assert.Equal(t, "0xabc-deadbeef", exts["partitionkey"])

	// sentat is unix nanos; receivedat is unix millis of the same instant.
	sentAt, err := strconv.ParseInt(exts["sentat"].(string), 10, 64)
	require.NoError(t, err)
	receivedAt, err := strconv.ParseInt(exts["receivedat"].(string), 10, 64)
	require.NoError(t, err)
	assert.Equal(t, sentAt/1_000_000, receivedAt)
	assert.WithinDuration(t, time.Now(), time.Unix(0, sentAt), time.Minute)
}

func TestChipIngressBatchClient_UnmappedTypeDropped(t *testing.T) { //nolint:paralleltest // asserts deltas on package-global prometheus metrics
	client, published := startBatchClient(t, publishOK)

	payload := newTestPayload("0xabc")
	payload.TelemType = LLOReport // no chip domain/entity mapping
	before := counterValue(TelemetryClientMessagesDropped, payload.TelemType)
	// The event is rejected synchronously in Send, before it reaches the queue.
	client.Send(t.Context(), payload)

	assert.InDelta(t, 1, counterValue(TelemetryClientMessagesDropped, payload.TelemType)-before, 0.001)
	assert.Empty(t, published.snapshot(), "PublishBatch must never be called for unmapped types")
}

func TestChipIngressBatchClient_CountsPerMessage(t *testing.T) { //nolint:paralleltest // asserts deltas on package-global prometheus metrics
	for name, tc := range map[string]struct { //nolint:paralleltest // subtests share package-global prometheus metrics
		publish publishFunc
		counter *prometheus.CounterVec
	}{
		"sent":        {publishOK, TelemetryClientMessagesSent},
		"send errors": {func(*chipingress.CloudEventBatch) (*chipingress.PublishResponse, error) { return nil, assert.AnError }, TelemetryClientMessagesSendErrors},
	} {
		t.Run(name, func(t *testing.T) {
			client, _ := startBatchClient(t, tc.publish)
			const n = 3
			before := counterValue(tc.counter, OCR2Median)
			for range n {
				client.Send(t.Context(), newTestPayload("0xabc"))
			}
			require.EventuallyWithT(t, func(c *assert.CollectT) {
				assert.InDelta(c, float64(n), counterValue(tc.counter, OCR2Median)-before, 0.001, "must count messages, not batches")
			}, testutils.WaitTimeout(t), 10*time.Millisecond)
		})
	}
}

func TestChipIngressBatchClient_HealthMonitoring(t *testing.T) { //nolint:paralleltest // resets a package-global prometheus gauge
	// The status gauge is package-global; reset it so a value left by another test can't satisfy the check.
	TelemetryClientConnectionStatus.WithLabelValues(chipIngressEndpointLabel).Set(0)
	startBatchClient(t, publishOK)

	require.EventuallyWithT(t, func(c *assert.CollectT) {
		assert.InDelta(c, 1, testutil.ToFloat64(TelemetryClientConnectionStatus.WithLabelValues(chipIngressEndpointLabel)), 0.001)
	}, testutils.WaitTimeout(t), 10*time.Millisecond)
}
