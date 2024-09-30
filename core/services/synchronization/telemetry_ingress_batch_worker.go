package synchronization

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	telemPb "github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
)

// telemetryIngressBatchWorker pushes telemetry in batches to the ingress server via wsrpc.
// A worker is created per ContractID.
type telemetryIngressBatchWorker struct {
	services.Service

	telemMaxBatchSize uint
	telemSendTimeout  time.Duration
	telemClient       telemPb.TelemClient
	chTelemetry       chan TelemPayload
	contractID        string
	telemType         TelemetryType
	logging           bool
	lggr              logger.Logger
	dropMessageCount  atomic.Uint32
}

// NewTelemetryIngressBatchWorker returns a worker for a given contractID that can send
// telemetry to the ingress server via WSRPC
func NewTelemetryIngressBatchWorker(
	telemMaxBatchSize uint,
	telemSendTimeout time.Duration,
	telemClient telemPb.TelemClient,
	chTelemetry chan TelemPayload,
	contractID string,
	telemType TelemetryType,
	lggr logger.Logger,
	logging bool,
) *telemetryIngressBatchWorker {
	return &telemetryIngressBatchWorker{
		telemSendTimeout:  telemSendTimeout,
		telemMaxBatchSize: telemMaxBatchSize,
		telemClient:       telemClient,
		chTelemetry:       chTelemetry,
		contractID:        contractID,
		telemType:         telemType,
		logging:           logging,
		lggr:              logger.Named(lggr, "TelemetryIngressBatchWorker"),
	}
}

// Send sends batched telemetry to the ingress server on an interval
func (tw *telemetryIngressBatchWorker) Send(ctx context.Context) {
	if len(tw.chTelemetry) == 0 {
		return
	}

	// Send batched telemetry to the ingress server, log any errors
	telemBatchReq := tw.BuildTelemBatchReq()
	ctx, cancel := context.WithTimeout(ctx, tw.telemSendTimeout)
	_, err := tw.telemClient.TelemBatch(ctx, telemBatchReq)
	cancel()

	if err != nil {
		tw.lggr.Warnf("Could not send telemetry: %v", err)
		return
	}
	if tw.logging {
		tw.lggr.Debugw("Successfully sent telemetry to ingress server", "contractID", telemBatchReq.ContractId, "telemType", telemBatchReq.TelemetryType, "telemetry", telemBatchReq.Telemetry)
	}
}

// logBufferFullWithExpBackoff logs messages at
// 1
// 2
// 4
// 8
// 16
// 32
// 64
// 100
// 200
// 300
// etc...
func (tw *telemetryIngressBatchWorker) logBufferFullWithExpBackoff(payload TelemPayload) {
	count := tw.dropMessageCount.Add(1)
	if count > 0 && (count%100 == 0 || count&(count-1) == 0) {
		tw.lggr.Warnw("telemetry ingress client buffer full, dropping message", "telemetry", payload.Telemetry, "droppedCount", count)
	}
}

// BuildTelemBatchReq reads telemetry off the worker channel and packages it into a batch request
func (tw *telemetryIngressBatchWorker) BuildTelemBatchReq() *telemPb.TelemBatchRequest {
	var telemBatch [][]byte

	// Read telemetry off the channel up to the max batch size
	for len(tw.chTelemetry) > 0 && len(telemBatch) < int(tw.telemMaxBatchSize) {
		telemPayload := <-tw.chTelemetry
		telemBatch = append(telemBatch, telemPayload.Telemetry)
	}

	return &telemPb.TelemBatchRequest{
		ContractId:    tw.contractID,
		TelemetryType: string(tw.telemType),
		Telemetry:     telemBatch,
		SentAt:        time.Now().UnixNano(),
	}
}
