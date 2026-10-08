package synchronization

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"
	"time"

	chipingress "github.com/smartcontractkit/chainlink-common/pkg/chipingress"
	batch "github.com/smartcontractkit/chainlink-common/pkg/chipingress/batch"
	chipingresspb "github.com/smartcontractkit/chainlink-common/pkg/chipingress/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

const (
	chipIngress = "chip-ingress" // endpoint label on telemetry_client_* metrics
	// legacyTelemetryBatchClientName is the client_name label on chip_ingress.batch.* metrics.
	legacyTelemetryBatchClientName = "legacy_telemetry"
	// Shared-buffer sizing mirrors chainlink-common beholder.DefaultConfig() chip batch defaults.
	chipIngressMessageBufferSize  = 10_000
	chipIngressMaxBatchSize       = 1_000
	chipIngressMaxConcurrentSends = 10
)

type chipIngressBatchClient struct {
	services.Service
	eng *services.Engine

	chipClient   chipingress.Client // used only for health Ping; owned and closed by batchClient.Stop
	batchClient  *batch.Client
	csaPubKeyHex string
	logging      bool
	dropCount    atomic.Uint32
	errorCount   atomic.Uint32
}

// NewChipIngressBatchClient returns a client backed by the shared
// chainlink-common chip-ingress batch client that can send telemetry to the
// chip ingress server. The batch client owns chipClient: Stop flushes the
// buffer and then closes the underlying gRPC connection.
func NewChipIngressBatchClient(chipClient chipingress.Client, csaPubKeyHex string, logging bool, lggr logger.Logger, sendInterval, sendTimeout time.Duration) (ChipIngressService, error) {
	batchClient, err := batch.NewBatchClient(chipClient,
		batch.WithMessageBuffer(chipIngressMessageBufferSize),
		batch.WithBatchSize(chipIngressMaxBatchSize),
		batch.WithMaxConcurrentSends(chipIngressMaxConcurrentSends),
		batch.WithBatchInterval(sendInterval),
		batch.WithMaxPublishTimeout(sendTimeout),
		batch.WithClientName(legacyTelemetryBatchClientName),
		// Send builds a fresh event per message and never touches it after queueing,
		// so skip the defensive per-message proto.Clone.
		batch.WithEventClone(false),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create chip-ingress batch client: %w", err)
	}

	c := &chipIngressBatchClient{
		chipClient:   chipClient,
		batchClient:  batchClient,
		csaPubKeyHex: csaPubKeyHex,
		logging:      logging,
	}
	c.Service, c.eng = services.Config{
		Name:  "ChipIngressBatchClient",
		Start: c.start,
		Close: c.close,
	}.NewServiceEngine(lggr)

	return c, nil
}

// start starts the batch client and a 5s health ping that drives
// telemetry_client_connection_status{endpoint="chip-ingress"}.
func (cc *chipIngressBatchClient) start(ctx context.Context) error {
	cc.batchClient.Start(ctx)
	cc.eng.GoTick(services.TickerConfig{}.NewTicker(5*time.Second), func(ctx context.Context) {
		pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
		connected := float64(1)
		if _, err := cc.chipClient.Ping(pingCtx, &chipingresspb.EmptyRequest{}); err != nil {
			connected = 0
			cc.eng.EmitHealthErr(err)
		}
		TelemetryClientConnectionStatus.WithLabelValues(chipIngress).Set(connected)
	})
	return nil
}

// close stops the batch client: it flushes queued messages (up to the shutdown
// timeout) and then closes the underlying chip-ingress client.
func (cc *chipIngressBatchClient) close() error {
	cc.batchClient.Stop()
	return nil
}

// Send converts a telemetry payload into a CloudEvent and queues it on the
// shared batch client. Queueing is non-blocking: if the shared buffer is full
// the message is dropped and a warning is logged with exponential back-off.
func (cc *chipIngressBatchClient) Send(ctx context.Context, payload TelemPayload) {
	now := time.Now()
	ev, err := cc.payloadToEvent(payload, now)
	if err != nil {
		TelemetryClientMessagesDropped.WithLabelValues(chipIngress, string(payload.TelemType)).Inc()
		cc.eng.Warnw("failed to build CloudEvent for ChIP ingress", "error", err, "contractID", payload.ContractID, "telemType", payload.TelemType)
		return
	}

	err = cc.batchClient.QueueMessage(ev, func(sendErr error) {
		if sendErr == nil {
			TelemetryClientMessagesSent.WithLabelValues(chipIngress, string(payload.TelemType)).Inc()
			if cc.logging {
				cc.eng.Debugw("Successfully sent telemetry to ChIP ingress", "contractID", payload.ContractID, "telemType", payload.TelemType)
			}
			return
		}
		TelemetryClientMessagesSendErrors.WithLabelValues(chipIngress, string(payload.TelemType)).Inc()
		if count := cc.errorCount.Add(1); shouldLogCount(count) {
			cc.eng.Warnw("Could not send telemetry via ChIP ingress",
				"error", sendErr,
				"errorCode", batch.ErrorCodeFor(sendErr),
				"telemType", payload.TelemType,
				"errorCount", count)
		}
	})
	if err != nil {
		// batch.ErrMessageBufferFull or batch.ErrClientShutdown
		TelemetryClientMessagesDropped.WithLabelValues(chipIngress, string(payload.TelemType)).Inc()
		if count := cc.dropCount.Add(1); shouldLogCount(count) {
			cc.eng.Warnw("dropping telemetry message for ChIP ingress",
				"error", err,
				"contractID", payload.ContractID,
				"telemType", payload.TelemType,
				"droppedCount", count)
		}
		return
	}
	cc.dropCount.Store(0)
}

// payloadToEvent converts a telemetry payload into a CloudEvent with the
// OTI-parity extension set. sentat and receivedat are both stamped with now
// (enqueue time); chip-ingress legacy mode may overwrite receivedat with the
// true server-receive time later.
func (cc *chipIngressBatchClient) payloadToEvent(payload TelemPayload, now time.Time) (*chipingress.CloudEventPb, error) {
	domain, entity, err := TelemetryTypeToDomainAndEntity(payload.TelemType)
	if err != nil {
		return nil, err
	}
	event, err := chipingress.NewEvent(domain, entity, payload.Telemetry, map[string]any{"time": now.UTC()})
	if err != nil {
		return nil, fmt.Errorf("failed creating CloudEvent: %w", err)
	}

	event.SetExtension("legacytelemetry", "true")
	event.SetExtension("telemetrytype", string(payload.TelemType))
	event.SetExtension("chainselector", strconv.FormatUint(payload.ChainSelector, 10))
	event.SetExtension("networkname", payload.Network)
	event.SetExtension("contractid", payload.ContractID)
	event.SetExtension("csapublickey", cc.csaPubKeyHex)
	event.SetExtension("partitionkey", payload.ContractID+"-"+cc.csaPubKeyHex)
	event.SetExtension("sentat", strconv.FormatInt(now.UnixNano(), 10))
	event.SetExtension("receivedat", strconv.FormatInt(now.UnixMilli(), 10))
	// nodeoperatorname / nodename are stamped server-side by chip-ingress via
	// csa-auth (WithNOPLookup); the node must not send them, or chip-ingress
	// would emit duplicate conflicting headers. nodeoperatorkey is not known to
	// the node at all (spike INFOPLAT-19426).

	return chipingress.EventToProto(event)
}
