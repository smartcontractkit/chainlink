package synchronization

import (
	"context"
	"testing"
	"time"

	cepb "github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/chipingress"
	"github.com/smartcontractkit/chainlink-common/pkg/chipingress/pb"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type noopChipIngressPublisher struct{}

func (noopChipIngressPublisher) Publish(ctx context.Context, event *cepb.CloudEvent, opts ...grpc.CallOption) (*pb.PublishResponse, error) {
	return &pb.PublishResponse{}, nil
}

func (noopChipIngressPublisher) PublishBatch(ctx context.Context, batch *pb.CloudEventBatch, opts ...grpc.CallOption) (*pb.PublishResponse, error) {
	return &pb.PublishResponse{}, nil
}

type partialChipIngressPublisher struct {
	resp *pb.PublishResponse
}

func (p partialChipIngressPublisher) Publish(ctx context.Context, event *cepb.CloudEvent, opts ...grpc.CallOption) (*pb.PublishResponse, error) {
	return &pb.PublishResponse{}, nil
}

func (p partialChipIngressPublisher) PublishBatch(ctx context.Context, batch *pb.CloudEventBatch, opts ...grpc.CallOption) (*pb.PublishResponse, error) {
	return p.resp, nil
}

func (p partialChipIngressPublisher) Ping(ctx context.Context, req *pb.EmptyRequest, opts ...grpc.CallOption) (*pb.PingResponse, error) {
	return &pb.PingResponse{}, nil
}

func (p partialChipIngressPublisher) StreamEvents(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[pb.StreamEventsRequest, pb.StreamEventsResponse], error) {
	return nil, nil
}

func (p partialChipIngressPublisher) RegisterSchema(ctx context.Context, req *pb.RegisterSchemaRequest, opts ...grpc.CallOption) (*pb.RegisterSchemaResponse, error) {
	return &pb.RegisterSchemaResponse{}, nil
}

func (p partialChipIngressPublisher) Close() error {
	return nil
}

func (p partialChipIngressPublisher) RegisterSchemas(ctx context.Context, schemas ...*pb.Schema) (map[string]int, error) {
	return nil, nil
}

func (noopChipIngressPublisher) Ping(ctx context.Context, req *pb.EmptyRequest, opts ...grpc.CallOption) (*pb.PingResponse, error) {
	return &pb.PingResponse{}, nil
}

func (noopChipIngressPublisher) StreamEvents(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[pb.StreamEventsRequest, pb.StreamEventsResponse], error) {
	return nil, nil
}

func (noopChipIngressPublisher) RegisterSchema(ctx context.Context, req *pb.RegisterSchemaRequest, opts ...grpc.CallOption) (*pb.RegisterSchemaResponse, error) {
	return &pb.RegisterSchemaResponse{}, nil
}

func (noopChipIngressPublisher) Close() error {
	return nil
}

func (noopChipIngressPublisher) RegisterSchemas(ctx context.Context, schemas ...*pb.Schema) (map[string]int, error) {
	return nil, nil
}

func TestChipIngressBatchWorker_Send_PartialDelivery(t *testing.T) {
	t.Parallel()
	// Verify that Send groups per-event errors into a single WARN log line
	// (partial delivery, transactionEnabled=false).
	partialResp := &pb.PublishResponse{
		Results: []*pb.PublishResult{
			{
				EventId: "evt-1",
				Error: &pb.PublishError{
					ErrorCode: pb.PublishErrorCode_PUBLISH_ERROR_CODE_VALIDATION_FAILED,
					Reason:    "schema not found",
				},
			},
			{
				EventId: "evt-2",
				Error: &pb.PublishError{
					ErrorCode: pb.PublishErrorCode_PUBLISH_ERROR_CODE_VALIDATION_FAILED,
					Reason:    "schema not found",
				},
			},
		},
	}
	publisher := partialChipIngressPublisher{resp: partialResp}

	lggr, observed := logger.TestLoggerObserved(t, zap.WarnLevel)
	chTelemetry := make(chan TelemPayload, 5)
	worker := NewChipIngressBatchWorker(
		2,
		time.Second,
		publisher,
		chTelemetry,
		"0xabc",
		OCR,
		lggr,
		true,
	)

	chTelemetry <- TelemPayload{
		Telemetry:     []byte("payload1"),
		TelemType:     OCR,
		ContractID:    "0xabc",
		Domain:        "data-feeds",
		Entity:        "ocr.v1.telemetry",
		ChainSelector: 7700,
		Network:       "EVM",
		ChainID:       "1",
	}
	chTelemetry <- TelemPayload{
		Telemetry:     []byte("payload2"),
		TelemType:     OCR,
		ContractID:    "0xabc",
		Domain:        "data-feeds",
		Entity:        "ocr.v1.telemetry",
		ChainSelector: 7700,
		Network:       "EVM",
		ChainID:       "1",
	}

	require.NotPanics(t, func() { worker.Send(t.Context()) })
	assert.Empty(t, chTelemetry)

	// Two failed events must produce exactly one grouped WARN log line.
	logs := observed.FilterMessage("chip ingress partial delivery errors")
	require.Equal(t, 1, logs.Len(), "expected exactly one grouped log line for 2 failed events")
	assert.Equal(t, zap.WarnLevel, logs.All()[0].Level)
}

func TestChipIngressBatchWorker_Send_PartialDeliveryThrottled(t *testing.T) {
	t.Parallel()
	// A persistent partial-delivery condition (e.g. missing schema) must not log on every
	// Send call - it should follow the same backoff cadence as logBufferFullWithExpBackoff.
	singleFailureResp := &pb.PublishResponse{
		Results: []*pb.PublishResult{
			{
				EventId: "evt-1",
				Error: &pb.PublishError{
					ErrorCode: pb.PublishErrorCode_PUBLISH_ERROR_CODE_SCHEMA_MISSING,
					Reason:    "schema not found",
				},
			},
		},
	}
	publisher := partialChipIngressPublisher{resp: singleFailureResp}

	lggr, observed := logger.TestLoggerObserved(t, zap.WarnLevel)
	chTelemetry := make(chan TelemPayload, 1)
	worker := NewChipIngressBatchWorker(
		1,
		time.Second,
		publisher,
		chTelemetry,
		"0xabc",
		OCR,
		lggr,
		true,
	)

	payload := TelemPayload{
		Telemetry:     []byte("payload"),
		TelemType:     OCR,
		ContractID:    "0xabc",
		Domain:        "data-feeds",
		Entity:        "ocr.v1.telemetry",
		ChainSelector: 7700,
		Network:       "EVM",
		ChainID:       "1",
	}

	// Drop counts after each send: 1, 2, 3. Backoff logs at 1 and 2 (powers of two), not at 3.
	for range 3 {
		chTelemetry <- payload
		require.NotPanics(t, func() { worker.Send(t.Context()) })
	}

	logs := observed.FilterMessage("chip ingress partial delivery errors")
	assert.Equal(t, 2, logs.Len(), "expected the third consecutive failure to be throttled")
}

func TestChipIngressBatchWorker_BuildCloudEventBatch(t *testing.T) {
	maxBatchSize := 3
	chTelemetry := make(chan TelemPayload, 10)
	// #nosec G115 -- maxBatchSize is a small positive constant, safe to convert to uint
	worker := NewChipIngressBatchWorker(
		uint(maxBatchSize),
		time.Second,
		noopChipIngressPublisher{},
		chTelemetry,
		"0xabc",
		OCR,
		logger.TestLogger(t),
		false,
	)

	payload := TelemPayload{
		Telemetry:     []byte("payload-1"),
		TelemType:     OCR,
		ContractID:    "0xabc",
		Domain:        "data-feeds",
		Entity:        "ocr.v1.telemetry",
		ChainSelector: 7700,
		Network:       "EVM",
		ChainID:       "1",
	}

	// enqueue more payloads than maxBatchSize to ensure batching occurs
	for range 5 {
		chTelemetry <- payload
	}

	batch1 := worker.BuildCloudEventBatch()
	require.NotNil(t, batch1)
	require.Len(t, batch1.Events, maxBatchSize)
	assert.Len(t, chTelemetry, 5-maxBatchSize)

	evt, err := chipingress.ProtoToEvent(batch1.Events[0])
	require.NoError(t, err)
	assert.Equal(t, "data-feeds", evt.Source())
	assert.Equal(t, "ocr.v1.telemetry", evt.Type())
	assert.Equal(t, []byte("payload-1"), evt.Data())

	attrs := batch1.Events[0].GetAttributes()
	require.Contains(t, attrs, "telemetrytype")
	require.Contains(t, attrs, "chainselector")
	assert.Equal(t, string(OCR), attrs["telemetrytype"].GetCeString())
	assert.Equal(t, "7700", attrs["chainselector"].GetCeString())

	batch2 := worker.BuildCloudEventBatch()
	require.NotNil(t, batch2)
	require.Len(t, batch2.Events, 2)
	assert.Empty(t, chTelemetry)
}

func TestChipIngressBatchWorker_BuildCloudEventBatchUsesMapping(t *testing.T) {
	chTelemetry := make(chan TelemPayload, 1)
	worker := NewChipIngressBatchWorker(
		1,
		time.Second,
		noopChipIngressPublisher{},
		chTelemetry,
		"0xdef",
		OCR,
		logger.TestLogger(t),
		false,
	)

	chTelemetry <- TelemPayload{
		Telemetry:     []byte("payload"),
		TelemType:     OCR,
		ContractID:    "0xdef",
		ChainSelector: 9001,
		Network:       "EVM",
		ChainID:       "137",
	}

	batch := worker.BuildCloudEventBatch()
	require.NotNil(t, batch)
	require.Len(t, batch.Events, 1)

	evt, err := chipingress.ProtoToEvent(batch.Events[0])
	require.NoError(t, err)
	assert.Equal(t, "data-feeds.telemetry.ocr", evt.Source())
}
