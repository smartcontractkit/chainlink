package synchronization

import (
	"context"
	"testing"
	"time"

	cepb "github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	for i := 0; i < 5; i++ {
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
