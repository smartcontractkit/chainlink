package synchronization

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/chipingress"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

const (
	chipIngress = "chip-ingress"
)

// chipIngressBatchWorker mirrors telemetryIngressBatchWorker but targets the ChIP ingress client.
// A worker is created per (contractID, telemetry type) pair.
type chipIngressBatchWorker struct {
	services.Service

	maxBatchSize     uint
	sendTimeout      time.Duration
	chipClient       chipingress.Client
	chTelemetry      chan TelemPayload
	contractID       string
	telemType        TelemetryType
	logging          bool
	lggr             logger.Logger
	dropMessageCount atomic.Uint32
}

// NewChipIngressBatchWorker returns a worker for a given contractID that can send
// telemetry to the chip ingress server via PublishBatch.
func NewChipIngressBatchWorker(
	maxBatchSize uint,
	sendTimeout time.Duration,
	chipClient chipingress.Client,
	chTelemetry chan TelemPayload,
	contractID string,
	telemType TelemetryType,
	lggr logger.Logger,
	logging bool,
) *chipIngressBatchWorker {
	return &chipIngressBatchWorker{
		maxBatchSize: maxBatchSize,
		sendTimeout:  sendTimeout,
		chipClient:   chipClient,
		chTelemetry:  chTelemetry,
		contractID:   contractID,
		telemType:    telemType,
		logging:      logging,
		lggr:         logger.Named(lggr, "ChipIngressBatchWorker"),
	}
}

// Send sends batched telemetry to the chip ingress server on an interval.
func (cw *chipIngressBatchWorker) Send(ctx context.Context) {
	if len(cw.chTelemetry) == 0 {
		return
	}

	batch := cw.BuildCloudEventBatch()
	if batch == nil || len(batch.Events) == 0 {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, cw.sendTimeout)
	defer cancel()

	_, err := cw.chipClient.PublishBatch(ctx, batch)
	if err != nil {
		cw.lggr.Warnf("Could not send telemetry via ChIP ingress: %v", err)
		TelemetryClientMessagesSendErrors.WithLabelValues(chipIngress, string(cw.telemType)).Inc()
		return
	}

	TelemetryClientMessagesSent.WithLabelValues(chipIngress, string(cw.telemType)).Inc()
	if cw.logging {
		cw.lggr.Debugw("Successfully sent telemetry to ChIP ingress", "contractID", cw.contractID, "telemType", cw.telemType, "batchSize", len(batch.Events))
	}
}

// logBufferFullWithExpBackoff logs messages at 1,2,4,8,16,32,64,100,200,300,... when the buffer is full.
func (cw *chipIngressBatchWorker) logBufferFullWithExpBackoff(payload TelemPayload) {
	count := cw.dropMessageCount.Add(1)
	TelemetryClientMessagesDropped.WithLabelValues(chipIngress, string(cw.telemType)).Inc()

	if count > 0 && (count%100 == 0 || count&(count-1) == 0) {
		cw.lggr.Warnw("chip ingress client buffer full, dropping message", "telemetry", payload.Telemetry, "droppedCount", count)
	}
}

// BuildCloudEventBatch reads telemetry off the worker channel and packages it into a CloudEvent batch.
func (cw *chipIngressBatchWorker) BuildCloudEventBatch() *chipingress.CloudEventBatch {
	var events []chipingress.CloudEvent

	// #nosec G115 -- maxBatchSize is uint, safe to convert to int for comparison
	for len(cw.chTelemetry) > 0 && len(events) < int(cw.maxBatchSize) {
		payload := <-cw.chTelemetry
		event, err := cw.payloadToEvent(payload)
		if err != nil {
			cw.lggr.Warnw("failed to build CloudEvent for ChIP ingress", "error", err, "contractID", payload.ContractID, "telemType", payload.TelemType)
			TelemetryClientMessagesDropped.WithLabelValues(chipIngress, string(cw.telemType)).Inc()
			continue
		}
		events = append(events, event)
	}

	if len(events) == 0 {
		return nil
	}

	batch, err := chipingress.EventsToBatch(events)
	if err != nil {
		cw.lggr.Warnw("failed to convert CloudEvents to batch", "error", err)
		return nil
	}

	return batch
}

func (cw *chipIngressBatchWorker) payloadToEvent(payload TelemPayload) (chipingress.CloudEvent, error) {
	domain := payload.Domain
	entity := payload.Entity
	if domain == "" || entity == "" {
		var err error
		domain, entity, err = TelemetryTypeToDomainAndEntity(payload.TelemType)
		if err != nil {
			return chipingress.CloudEvent{}, fmt.Errorf("missing domain/entity for telemType %s: %w", payload.TelemType, err)
		}
	}

	now := time.Now().UTC()

	attrs := map[string]any{
		"time":            now,
		"datacontenttype": "application/protobuf",
	}

	event, err := chipingress.NewEvent(domain, entity, payload.Telemetry, attrs)
	if err != nil {
		return chipingress.CloudEvent{}, fmt.Errorf("failed creating CloudEvent: %w", err)
	}

	event.SetExtension("telemetrytype", string(payload.TelemType))
	event.SetExtension("chainselector", strconv.FormatUint(payload.ChainSelector, 10))
	event.SetExtension("contractid", payload.ContractID)
	event.SetExtension("networkname", payload.Network)
	event.SetExtension("sentat", now.Format(time.RFC3339))
	// These attributes need to be set in chip-ingress server side
	event.SetExtension("nodeoperatorname", "")
	event.SetExtension("nodename", "")
	event.SetExtension("nodecsapublickey", "")

	return event, nil
}
