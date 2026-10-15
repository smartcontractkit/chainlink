package synchronization

import (
	"context"
	"fmt"
	"strconv"
	"sync"
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
	// unhealthyAfterPingFailures marks the service unhealthy only after about a minute of
	// failed pings, so a brief telemetry-server blip doesn't flip the node's /health.
	unhealthyAfterPingFailures = 12
	chipIngressUnreachable     = "chip-ingress unreachable"

	// logEvery limits each kind of repeated warning (per telemetry type) to one line per interval.
	logEvery = 30 * time.Second
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
	logs         logLimiter
	pingInterval time.Duration // healthPingInterval; shortened by tests
}

// ChipIngressBatchConfig sizes the shared chip-ingress batch client. Values come from
// the [TelemetryIngress] ChipIngress* settings.
type ChipIngressBatchConfig struct {
	BufferSize         uint          // messages buffered across all jobs before new ones are dropped
	MaxBatchSize       uint          // messages per PublishBatch request
	MaxConcurrentSends int           // PublishBatch requests in flight
	SendInterval       time.Duration // max wait before flushing an incomplete batch
	SendTimeout        time.Duration // per-request PublishBatch timeout
	MaxGRPCRequestSize int           // max serialized PublishBatch request size; larger batches are split
	DrainTimeout       time.Duration // max time Close waits to publish queued messages and run their callbacks
	Logging            bool          // debug-log every successful send
}

// NewChipIngressBatchClient wraps chipClient in the chainlink-common batch client.
// The batch client owns chipClient: Close publishes queued messages, then closes the connection.
func NewChipIngressBatchClient(chipClient chipingress.Client, csaPubKeyHex string, cfg ChipIngressBatchConfig, lggr logger.Logger) (ChipIngressService, error) {
	batchClient, err := batch.NewBatchClient(chipClient,
		batch.WithMessageBuffer(int(cfg.BufferSize)),
		batch.WithBatchSize(int(cfg.MaxBatchSize)),
		batch.WithMaxConcurrentSends(cfg.MaxConcurrentSends),
		batch.WithBatchInterval(cfg.SendInterval),
		batch.WithMaxPublishTimeout(cfg.SendTimeout),
		batch.WithMaxGRPCRequestSize(cfg.MaxGRPCRequestSize),
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
		logs:         logLimiter{every: logEvery},
		pingInterval: healthPingInterval,
	}
	c.Service, c.eng = services.Config{
		Name:  "ChipIngressBatchClient",
		Start: c.start,
		Close: c.close,
	}.NewServiceEngine(lggr)

	return c, nil
}

// start starts the batch client and a periodic health ping. The ping drives
// telemetry_client_connection_status{endpoint="chip-ingress"} on every probe, and
// marks the service unhealthy only after a sustained outage.
func (c *chipIngressBatchClient) start(ctx context.Context) error {
	c.batchClient.Start(ctx)
	failures := 0 // only touched by the ticker goroutine
	c.eng.GoTick(services.TickerConfig{}.NewTicker(c.pingInterval), func(ctx context.Context) {
		pingCtx, cancel := context.WithTimeout(ctx, healthPingTimeout)
		defer cancel()
		if _, err := c.chipClient.Ping(pingCtx, &chipingresspb.EmptyRequest{}); err != nil {
			TelemetryClientConnectionStatus.WithLabelValues(chipIngressEndpointLabel).Set(0)
			if failures++; failures >= unhealthyAfterPingFailures {
				c.eng.SetHealthCond(chipIngressUnreachable, err)
			}
			return
		}
		TelemetryClientConnectionStatus.WithLabelValues(chipIngressEndpointLabel).Set(1)
		failures = 0
		c.eng.ClearHealthCond(chipIngressUnreachable)
	})
	return nil
}

// close publishes everything still queued, waiting up to the batch client's shutdown
// timeout for sends and callbacks, then closes the chip-ingress client. Batches still
// draining after the timeout fail against the closed client and count as send errors.
func (c *chipIngressBatchClient) close() error {
	c.batchClient.Stop()
	return nil
}

// Send builds the CloudEvent and queues it. It never blocks the caller (a libocr
// goroutine). Types chip-ingress has no mapping for are counted as unsupported;
// a full buffer or a shutting-down client drops the message.
func (c *chipIngressBatchClient) Send(_ context.Context, payload TelemPayload) {
	telemType := string(payload.TelemType)
	domain, entity, err := TelemetryTypeToDomainAndEntity(payload.TelemType)
	if err != nil {
		// Expected on LLO nodes (llo-report); OTI drops these types too.
		TelemetryClientMessagesUnsupported.WithLabelValues(chipIngressEndpointLabel, telemType).Inc()
		if ok, suppressed := c.logs.allow("unsupported:" + telemType); ok {
			c.eng.Debugw("Telemetry type not supported by ChIP ingress; not sending", "telemType", telemType, "suppressed", suppressed)
		}
		return
	}
	event, err := c.newEvent(payload, domain, entity, time.Now())
	if err == nil {
		err = c.batchClient.QueueMessage(event, c.onPublished(payload.TelemType, payload.ContractID))
	}
	if err != nil {
		TelemetryClientMessagesDropped.WithLabelValues(chipIngressEndpointLabel, telemType).Inc()
		if ok, suppressed := c.logs.allow("drop:" + telemType); ok {
			c.eng.Warnw("Dropping telemetry message for ChIP ingress", "error", err, "telemType", telemType, "suppressed", suppressed)
		}
	}
}

// onPublished returns the callback the batch client runs once the message's batch
// has been published, counting each message as sent or failed.
func (c *chipIngressBatchClient) onPublished(telemType TelemetryType, contractID string) func(error) {
	return func(err error) {
		if err != nil {
			TelemetryClientMessagesSendErrors.WithLabelValues(chipIngressEndpointLabel, string(telemType)).Inc()
			code := batch.ErrorCodeFor(err)
			if ok, suppressed := c.logs.allow("error:" + string(telemType) + ":" + code); ok {
				c.eng.Warnw("Could not send telemetry via ChIP ingress", "error", err, "errorCode", code, "telemType", telemType, "suppressed", suppressed)
			}
			return
		}
		TelemetryClientMessagesSent.WithLabelValues(chipIngressEndpointLabel, string(telemType)).Inc()
		if c.logging {
			c.eng.Debugw("Successfully sent telemetry to ChIP ingress", "contractID", contractID, "telemType", telemType)
		}
	}
}

// newEvent converts a payload into the CloudEvent chip-ingress expects for legacy
// telemetry. Each extension becomes a "ce_<name>" Kafka header; the comments name the
// header OTI used for the same value. A legacy-mode chip-ingress routes on telemetrytype
// and chainselector; source/type (domain/entity) only matter to a non-legacy deployment.
//
// The node does not send nodeoperatorname/nodename or receivedat: chip-ingress derives
// them from the CSA key's NOP info (WithNOPLookup) and its own clock. nodeoperatorkey is
// unknown to the node (spike INFOPLAT-19426).
func (c *chipIngressBatchClient) newEvent(payload TelemPayload, domain, entity string, now time.Time) (*chipingress.CloudEventPb, error) {
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

	return chipingress.EventToProto(event)
}

// logLimiter allows one log line per key per interval and counts what it suppressed,
// so repeated warnings stay readable however fast messages arrive.
type logLimiter struct {
	every time.Duration
	mu    sync.Mutex
	keys  map[string]*logLimiterKey
}

type logLimiterKey struct {
	last       time.Time
	suppressed int
}

func (l *logLimiter) allow(key string) (ok bool, suppressed int) {
	now := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.keys == nil {
		l.keys = make(map[string]*logLimiterKey)
	}
	k, seen := l.keys[key]
	if !seen {
		l.keys[key] = &logLimiterKey{last: now}
		return true, 0
	}
	if now.Sub(k.last) < l.every {
		k.suppressed++
		return false, 0
	}
	suppressed = k.suppressed
	k.last, k.suppressed = now, 0
	return true, suppressed
}
