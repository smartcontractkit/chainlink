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
	// chipIngressEndpointLabel is the "endpoint" label on telemetry_client_* metrics.
	chipIngressEndpointLabel = "chip-ingress"
	// legacyTelemetryBatchClientName is the client_name label on chip_ingress.batch.* metrics.
	legacyTelemetryBatchClientName = "legacy_telemetry"

	healthPingInterval = 5 * time.Second
	healthPingTimeout  = 2 * time.Second
)

// chipIngressBatchClient sends the node's legacy OCR telemetry to chip-ingress,
// replacing the WSRPC path to OTI (OCR Telemetry Ingress). One instance is
// shared by every job. Each message becomes a CloudEvent whose extensions
// chip-ingress turns into the same Kafka key and headers OTI used to produce.
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

// ChipIngressBatchConfig sizes the shared chip-ingress batch client. Values come from
// [TelemetryIngress] (ChipIngress* fields plus SendInterval/SendTimeout).
type ChipIngressBatchConfig struct {
	BufferSize         uint          // messages buffered across all jobs before new ones are dropped
	MaxBatchSize       uint          // messages per PublishBatch request
	MaxConcurrentSends int           // PublishBatch requests in flight
	SendInterval       time.Duration // max wait before flushing an incomplete batch
	SendTimeout        time.Duration // per-request PublishBatch timeout
	DrainTimeout       time.Duration // max time to flush the buffer on Close
	Logging            bool          // debug-log every successful send
}

// NewChipIngressBatchClient wraps chipClient in the chainlink-common batch client.
// The batch client owns chipClient: Close flushes the buffer, then closes the connection.
func NewChipIngressBatchClient(chipClient chipingress.Client, csaPubKeyHex string, cfg ChipIngressBatchConfig, lggr logger.Logger) (ChipIngressService, error) {
	batchClient, err := batch.NewBatchClient(chipClient,
		batch.WithMessageBuffer(int(cfg.BufferSize)),
		batch.WithBatchSize(int(cfg.MaxBatchSize)),
		batch.WithMaxConcurrentSends(cfg.MaxConcurrentSends),
		batch.WithBatchInterval(cfg.SendInterval),
		batch.WithMaxPublishTimeout(cfg.SendTimeout),
		batch.WithShutdownTimeout(cfg.DrainTimeout),
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
		logging:      cfg.Logging,
	}
	c.Service, c.eng = services.Config{
		Name:  "ChipIngressBatchClient",
		Start: c.start,
		Close: c.close,
	}.NewServiceEngine(lggr)

	return c, nil
}

// start starts the batch client and a periodic health ping that drives
// telemetry_client_connection_status{endpoint="chip-ingress"}.
func (c *chipIngressBatchClient) start(ctx context.Context) error {
	c.batchClient.Start(ctx)
	c.eng.GoTick(services.TickerConfig{}.NewTicker(healthPingInterval), func(ctx context.Context) {
		pingCtx, cancel := context.WithTimeout(ctx, healthPingTimeout)
		defer cancel()
		connected := float64(1)
		if _, err := c.chipClient.Ping(pingCtx, &chipingresspb.EmptyRequest{}); err != nil {
			connected = 0
			c.eng.EmitHealthErr(err)
		}
		TelemetryClientConnectionStatus.WithLabelValues(chipIngressEndpointLabel).Set(connected)
	})
	return nil
}

// close flushes queued messages (up to the batch client's shutdown timeout),
// then closes the underlying chip-ingress client.
func (c *chipIngressBatchClient) close() error {
	c.batchClient.Stop()
	return nil
}

// Send builds the CloudEvent and queues it. It never blocks the caller (a libocr
// goroutine): if the event can't be built or the buffer is full, the message is
// dropped and counted.
func (c *chipIngressBatchClient) Send(_ context.Context, payload TelemPayload) {
	event, err := c.newEvent(payload, time.Now())
	if err == nil {
		err = c.batchClient.QueueMessage(event, c.onPublished(payload))
	}
	if err != nil {
		c.drop(payload, err)
		return
	}
	c.dropCount.Store(0)
}

// onPublished returns the callback the batch client runs once the message's batch
// has been published, counting each message as sent or failed.
func (c *chipIngressBatchClient) onPublished(payload TelemPayload) func(error) {
	return func(err error) {
		if err != nil {
			TelemetryClientMessagesSendErrors.WithLabelValues(chipIngressEndpointLabel, string(payload.TelemType)).Inc()
			if count := c.errorCount.Add(1); shouldLogCount(count) {
				c.eng.Warnw("Could not send telemetry via ChIP ingress",
					"error", err,
					"errorCode", batch.ErrorCodeFor(err),
					"telemType", payload.TelemType,
					"errorCount", count)
			}
			return
		}
		TelemetryClientMessagesSent.WithLabelValues(chipIngressEndpointLabel, string(payload.TelemType)).Inc()
		if c.logging {
			c.eng.Debugw("Successfully sent telemetry to ChIP ingress", "contractID", payload.ContractID, "telemType", payload.TelemType)
		}
	}
}

// drop counts a message that never reached the queue: its telemetry type has no
// chip-ingress mapping, the buffer is full, or the client is shutting down.
func (c *chipIngressBatchClient) drop(payload TelemPayload, err error) {
	TelemetryClientMessagesDropped.WithLabelValues(chipIngressEndpointLabel, string(payload.TelemType)).Inc()
	if count := c.dropCount.Add(1); shouldLogCount(count) {
		c.eng.Warnw("Dropping telemetry message for ChIP ingress",
			"error", err,
			"contractID", payload.ContractID,
			"telemType", payload.TelemType,
			"droppedCount", count)
	}
}

// newEvent converts a payload into the CloudEvent chip-ingress expects for legacy
// telemetry. The event source/type select the Kafka topic. Each extension becomes a
// "ce_<name>" Kafka header; the comments name the header OTI used for the same value.
//
// The node does not send nodeoperatorname/nodename: chip-ingress looks them up from
// the CSA key (WithNOPLookup), and sending them too would duplicate those headers.
// nodeoperatorkey is unknown to the node (spike INFOPLAT-19426).
func (c *chipIngressBatchClient) newEvent(payload TelemPayload, now time.Time) (*chipingress.CloudEventPb, error) {
	domain, entity, err := TelemetryTypeToDomainAndEntity(payload.TelemType)
	if err != nil {
		return nil, err
	}
	event, err := chipingress.NewEvent(domain, entity, payload.Telemetry, map[string]any{"time": now.UTC()})
	if err != nil {
		return nil, fmt.Errorf("failed creating CloudEvent: %w", err)
	}

	event.SetExtension("legacytelemetry", "true")                                      // marks node-sent legacy telemetry
	event.SetExtension("telemetrytype", string(payload.TelemType))                     // "telemetry-type"
	event.SetExtension("chainselector", strconv.FormatUint(payload.ChainSelector, 10)) // chain-selectors ID of the job's chain
	event.SetExtension("networkname", payload.Network)                                 // endpoint network, e.g. "EVM"
	event.SetExtension("contractid", payload.ContractID)                               // "contract-address"
	event.SetExtension("csapublickey", c.csaPubKeyHex)                                 // "node-operator-csa-public-key"
	event.SetExtension("partitionkey", payload.ContractID+"-"+c.csaPubKeyHex)          // Kafka record key
	event.SetExtension("sentat", strconv.FormatInt(now.UnixNano(), 10))                // "sent-at", unix ns
	event.SetExtension("receivedat", strconv.FormatInt(now.UnixMilli(), 10))           // "received-at", unix ms; chip-ingress may overwrite with its receive time

	return chipingress.EventToProto(event)
}
