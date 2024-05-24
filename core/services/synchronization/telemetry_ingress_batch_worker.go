package synchronization

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	telemPb "github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
)

// telemetryIngressBatchWorker pushes telemetry in batches to the ingress server via wsrpc.
// A worker is created per ContractID.
type telemetryIngressBatchWorker struct {
	services.Service

	telemMaxBatchSize uint
	telemSendInterval time.Duration
	telemSendTimeout  time.Duration
	telemClient       telemPb.TelemClient
	wgDone            *sync.WaitGroup
	chDone            services.StopChan
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
	telemSendInterval time.Duration,
	telemSendTimeout time.Duration,
	telemClient telemPb.TelemClient,
	wgDone *sync.WaitGroup,
	chDone chan struct{},
	chTelemetry chan TelemPayload,
	contractID string,
	telemType TelemetryType,
	globalLogger logger.Logger,
	logging bool,
) *telemetryIngressBatchWorker {
	return &telemetryIngressBatchWorker{
		telemSendInterval: telemSendInterval,
		telemSendTimeout:  telemSendTimeout,
		telemMaxBatchSize: telemMaxBatchSize,
		telemClient:       telemClient,
		wgDone:            wgDone,
		chDone:            chDone,
		chTelemetry:       chTelemetry,
		contractID:        contractID,
		telemType:         telemType,
		logging:           logging,
		lggr:              globalLogger.Named("TelemetryIngressBatchWorker"),
	}
}

// Start sends batched telemetry to the ingress server on an interval
func (tw *telemetryIngressBatchWorker) Start() {
	tw.wgDone.Add(1)
	sendTicker := time.NewTicker(tw.telemSendInterval)

	go func() {
		defer tw.wgDone.Done()

		for {
			select {
			case <-sendTicker.C:
				if len(tw.chTelemetry) == 0 {
					continue
				}

				// Send batched telemetry to the ingress server, log any errors
				telemBatchReq := tw.BuildTelemBatchReq()
				ctx, cancel := tw.chDone.CtxCancel(context.WithTimeout(context.Background(), tw.telemSendTimeout))
				_, err := tw.telemClient.TelemBatch(ctx, telemBatchReq)
				cancel()

				if err != nil {
					tw.lggr.Warnf("Could not send telemetry: %v", err)
					continue
				}
				if tw.logging {
					tw.lggr.Debugw("Successfully sent telemetry to ingress server", "contractID", telemBatchReq.ContractId, "telemType", telemBatchReq.TelemetryType, "telemetry", telemBatchReq.Telemetry)
				}
			case <-tw.chDone:
				return
			}
		}
	}()
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
