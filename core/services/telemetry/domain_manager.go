package telemetry

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelmetric "go.opentelemetry.io/otel/metric"

	chainselector "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/chipingress"
	"github.com/smartcontractkit/chainlink-common/pkg/chipingress/batch"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
)

// DomainEmitterSource provides Beholder's per-domain emitters. *beholder.Client implements it.
type DomainEmitterSource interface {
	ForDomain(domain beholder.Domain) (beholder.DomainEmitter, error)
}

// supportedTelemetryDomains lists the telemetry types DomainManager sends and the Beholder domain
// whose emitter carries each: the types ingested in production. Their CloudEvent source and type
// come from synchronization.TelemetryTypeToDomainAndEntity. Every other type is dropped and counted.
// The TelemetryIngress.BeholderEnabled docs (core/config/docs/core.toml) list these types: keep them in sync.
// A package-level map: countUnsupported reads it per message and must not allocate.
var supportedTelemetryDomains = map[synchronization.TelemetryType]beholder.Domain{
	synchronization.OCR:               beholder.DomainDataFeeds,
	synchronization.OCR2Median:        beholder.DomainDataFeeds,
	synchronization.OCR3DataFeeds:     beholder.DomainDataFeeds,
	synchronization.EnhancedEA:        beholder.DomainDataFeeds,
	synchronization.OCR3Mercury:       beholder.DomainDataStreams,
	synchronization.OCR3_1Mercury:     beholder.DomainDataStreams,
	synchronization.EnhancedEAMercury: beholder.DomainDataStreams,
	synchronization.LLOObservation:    beholder.DomainDataStreams,
	synchronization.LLOOutcome:        beholder.DomainDataStreams,
	synchronization.PipelineBridge:    beholder.DomainDataStreams,
	synchronization.OCR2CCIPCommit:    beholder.DomainCCIP,
	synchronization.OCR2CCIPExec:      beholder.DomainCCIP,
	synchronization.OCR3CCIPCommit:    beholder.DomainCCIP,
	synchronization.OCR3CCIPExec:      beholder.DomainCCIP,
	// TODO(INFOPLAT-19709): head reports are emitted for every chain on every node; confirm
	// ccip is their domain.
	synchronization.HeadReport: beholder.DomainCCIP,
}

// Metric names and attribute values. The manager's metrics are OTel (Beholder) metrics, so they
// reach us from every node. The counters carry the route's domain and telemetry_type, the target
// chain's chain_selector, and error_code or reason where they apply.
//
// It also increments the existing Prometheus telemetry_client_messages_{sent,send_errors,dropped}
// counters for node-local dashboards, with two differences from the WSRPC path: they count
// messages, where the WSRPC batch worker counts batches, and their endpoint label is
// "beholder/<domain>" rather than the ingress URL. Unsupported types and chains are counted only
// in OTel (telemetry_messages_unsupported_total).
const (
	metricMessagesSent        = "telemetry_messages_sent_total"
	metricMessagesSendErrors  = "telemetry_messages_send_errors_total"
	metricMessagesDropped     = "telemetry_messages_dropped_total"
	metricMessagesUnsupported = "telemetry_messages_unsupported_total"

	unsupportedReasonTelemetryType = "telemetry_type"
	unsupportedReasonChain         = "chain"

	// errorCodeEmitRejected is the error_code of events the emitter refused outright (not
	// started, invalid, or outside its domain); other codes come from batch.ErrorCodeFor.
	errorCodeEmitRejected = "emit_rejected"
	errorCodeEventBuild   = "event_build"
)

var _ Manager = (*DomainManager)(nil)

// DomainManager is the Manager that sends telemetry through Beholder's per-domain
// emitters. Each supported telemetry type maps to a domain; endpoints bind their route
// when created, so sending does no lookups. Network and chain only describe the event:
// nothing is configured or routed per network.
type DomainManager struct {
	services.Service
	eng *services.Engine

	// ctx is cancelled when the manager closes. SendLog has no context (libocr's
	// MonitoringEndpoint), so sends and counters use this one.
	ctx    context.Context //nolint:containedctx // libocr's MonitoringEndpoint.SendLog takes no ctx
	cancel context.CancelFunc

	routes          map[synchronization.TelemetryType]*telemetryRoute
	csaPubKeyHex    string
	meterProvider   otelmetric.MeterProvider
	metrics         domainManagerMetrics
	now             func() time.Time
	reported        sync.Map // keys of unsupported types, chains and refusals already logged
	reportedKeys    atomic.Int64
	unsupportedOpts sync.Map // unsupportedKey -> otelmetric.AddOption
	metricOpts      sync.Map // metricKey -> otelmetric.AddOption
}

// telemetryRoute is where one telemetry type goes, with its metric handles precomputed.
type telemetryRoute struct {
	telemType synchronization.TelemetryType
	domain    beholder.Domain
	source    string
	entity    string
	emitter   beholder.DomainEmitter

	attrs attribute.Set

	promSent, promSendErrors, promDropped prometheus.Counter
}

type domainManagerMetrics struct {
	sent        otelmetric.Int64Counter
	sendErrors  otelmetric.Int64Counter
	dropped     otelmetric.Int64Counter
	unsupported otelmetric.Int64Counter
}

// DomainManagerOpt configures NewDomainManager.
type DomainManagerOpt func(*DomainManager)

// WithMeterProvider sets the OTel meter provider for the manager's metrics. The default is the
// global provider, which is Beholder's once the Beholder client is set globally.
func WithMeterProvider(mp otelmetric.MeterProvider) DomainManagerOpt {
	return func(m *DomainManager) { m.meterProvider = mp }
}

// NewDomainManager resolves an emitter for every domain the supported telemetry types use; a
// missing one is an error. csaPubKeyHex is the node's CSA public key, which legacy-mode chip
// ingress binds to the authenticated key.
func NewDomainManager(src DomainEmitterSource, csaPubKeyHex string, lggr logger.Logger, opts ...DomainManagerOpt) (*DomainManager, error) {
	if src == nil {
		return nil, errors.New("telemetry: a Beholder domain emitter source is required")
	}
	if csaPubKeyHex == "" {
		return nil, errors.New("telemetry: CSA public key is required")
	}
	m := &DomainManager{
		csaPubKeyHex:  csaPubKeyHex,
		meterProvider: otel.GetMeterProvider(),
		now:           time.Now,
	}
	for _, opt := range opts {
		opt(m)
	}
	metrics, err := newDomainManagerMetrics(m.meterProvider)
	if err != nil {
		return nil, err
	}
	m.metrics = metrics

	// Resolve domains in a stable order so a missing one is reported the same way every time.
	var domains []beholder.Domain
	for _, d := range supportedTelemetryDomains {
		if !slices.Contains(domains, d) {
			domains = append(domains, d)
		}
	}
	slices.Sort(domains)
	emitters := make(map[beholder.Domain]beholder.DomainEmitter, len(domains))
	for _, d := range domains {
		em, err := src.ForDomain(d)
		if err != nil {
			return nil, fmt.Errorf("telemetry: no Beholder emitter for domain %q: %w", d, err)
		}
		emitters[d] = em
	}

	m.routes = make(map[synchronization.TelemetryType]*telemetryRoute, len(supportedTelemetryDomains))
	for telemType, d := range supportedTelemetryDomains {
		source, entity, err := synchronization.TelemetryTypeToDomainAndEntity(telemType)
		if err != nil {
			return nil, fmt.Errorf("telemetry: no CloudEvent source for %q: %w", telemType, err)
		}
		attrs := attribute.NewSet(
			attribute.String("domain", string(d)),
			attribute.String("telemetry_type", string(telemType)),
		)
		endpoint := "beholder/" + string(d)
		m.routes[telemType] = &telemetryRoute{
			telemType:      telemType,
			domain:         d,
			source:         source,
			entity:         entity,
			emitter:        emitters[d],
			attrs:          attrs,
			promSent:       synchronization.TelemetryClientMessagesSent.WithLabelValues(endpoint, string(telemType)),
			promSendErrors: synchronization.TelemetryClientMessagesSendErrors.WithLabelValues(endpoint, string(telemType)),
			promDropped:    synchronization.TelemetryClientMessagesDropped.WithLabelValues(endpoint, string(telemType)),
		}
	}

	m.ctx, m.cancel = context.WithCancel(context.Background())
	m.Service, m.eng = services.Config{
		Name:  "DomainManager",
		Close: func() error { m.cancel(); return nil },
	}.NewServiceEngine(lggr)
	return m, nil
}

func newDomainManagerMetrics(mp otelmetric.MeterProvider) (m domainManagerMetrics, err error) {
	meter := mp.Meter("github.com/smartcontractkit/chainlink/v2/core/services/telemetry")
	if m.sent, err = meter.Int64Counter(metricMessagesSent,
		otelmetric.WithDescription("Telemetry messages accepted by chip ingress")); err != nil {
		return m, err
	}
	if m.sendErrors, err = meter.Int64Counter(metricMessagesSendErrors,
		otelmetric.WithDescription("Telemetry messages that failed to publish")); err != nil {
		return m, err
	}
	if m.dropped, err = meter.Int64Counter(metricMessagesDropped,
		otelmetric.WithDescription("Telemetry messages dropped before publishing (buffer full, shutdown, refused)")); err != nil {
		return m, err
	}
	if m.unsupported, err = meter.Int64Counter(metricMessagesUnsupported,
		otelmetric.WithDescription("Telemetry messages dropped because their type or chain is not supported")); err != nil {
		return m, err
	}
	return m, nil
}

// GenMonitoringEndpoint returns an endpoint for one telemetry type. Unsupported types and chains
// get an endpoint that drops and counts what it is given.
func (m *DomainManager) GenMonitoringEndpoint(network string, chainID string, contractID string, telemType synchronization.TelemetryType) commontypes.MonitoringEndpoint {
	route, ok := m.routes[telemType]
	if !ok {
		m.warnUnsupportedType(telemType)
		return &unsupportedAgent{m: m, telemType: telemType, reason: unsupportedReasonTelemetryType}
	}
	selector, ok := m.chainSelector(network, chainID)
	if !ok {
		return &unsupportedAgent{m: m, telemType: telemType, reason: unsupportedReasonChain}
	}
	return &domainAgent{m: m, route: route, target: m.newEventTarget(network, contractID, selector)}
}

// GenMultitypeMonitoringEndpoint returns an endpoint whose messages carry their own type.
func (m *DomainManager) GenMultitypeMonitoringEndpoint(network string, chainID string, contractID string) MultitypeMonitoringEndpoint {
	selector, ok := m.chainSelector(network, chainID)
	if !ok {
		return &unsupportedAgent{m: m, reason: unsupportedReasonChain}
	}
	return &multitypeDomainAgent{m: m, target: m.newEventTarget(network, contractID, selector)}
}

// eventTarget holds an endpoint's per-event values, computed once when the endpoint is created.
type eventTarget struct {
	network      string // lower case, e.g. "evm"
	contractID   string
	selector     string // decimal chain selector
	partitionKey string // <contractID>-<csaPubKeyHex>, the Kafka record key, as OTI's
}

func (m *DomainManager) newEventTarget(network, contractID string, selector uint64) eventTarget {
	return eventTarget{
		network:      strings.ToLower(network),
		contractID:   contractID,
		selector:     strconv.FormatUint(selector, 10),
		partitionKey: contractID + "-" + m.csaPubKeyHex,
	}
}

// chainSelector resolves the job's chain to its chain-selectors selector, which legacy-mode chip
// ingress maps to the network's topics. A chain that does not resolve (e.g. one newer than the
// pinned chain-selectors) loses its telemetry, so it is logged at Error, once per chain.
func (m *DomainManager) chainSelector(network, chainID string) (uint64, bool) {
	details, err := chainselector.GetChainDetailsByChainIDAndFamily(chainID, strings.ToLower(network))
	if err != nil {
		m.reportOnce("chain:"+strings.ToLower(network)+"/"+chainID, func() {
			m.eng.Errorw("Chain has no chain selector: its telemetry will be dropped", "network", network, "chainID", chainID, "err", err)
		})
		return 0, false
	}
	return details.ChainSelector, true
}

// maxReportedKeys caps how many distinct one-time reports are remembered. The keys are
// sender-influenced (telemetry types, network/chainID pairs), so the cap bounds
// worst-case memory; past it, new keys are only counted — never logged or retained.
const maxReportedKeys = 25

func (m *DomainManager) reportOnce(key string, report func()) {
	if _, loaded := m.reported.LoadOrStore(key, struct{}{}); loaded {
		return
	}
	if m.reportedKeys.Add(1) > maxReportedKeys {
		m.reported.Delete(key)
		return
	}
	report()
}

// send builds the event and hands it to the route's emitter. Fates are counted; publish failures
// are logged (rate-limited) by the emitter, and events that never reach it are logged here, once
// per route and error code, since they mean telemetry is being lost.
func (m *DomainManager) send(route *telemetryRoute, target eventTarget, telemetry []byte) {
	ctx := m.ctx
	ev, err := m.newEvent(route, target, telemetry)
	if err != nil {
		m.refused(ctx, route, target, errorCodeEventBuild, err)
		return
	}
	if err := route.emitter.EmitEvent(ctx, ev, func(sendErr error) { m.recordFate(ctx, route, target, sendErr) }); err != nil {
		m.refused(ctx, route, target, errorCodeEmitRejected, err)
	}
}

// refused counts an event that never reached the emitter's buffer and logs the first one per
// route and error code at Error.
func (m *DomainManager) refused(ctx context.Context, route *telemetryRoute, target eventTarget, code string, err error) {
	m.recordDropped(ctx, route, target, code)
	m.reportOnce("refused:"+string(route.telemType)+"/"+code, func() {
		m.eng.Errorw("Telemetry refused: it is being dropped (further refusals are only counted)",
			"telemType", route.telemType, "domain", route.domain, "errorCode", code, "err", err)
	})
}

// recordFate counts an event's fate as reported by the emitter: published, failed to publish,
// or dropped before publishing (full buffer or shutdown).
func (m *DomainManager) recordFate(ctx context.Context, route *telemetryRoute, target eventTarget, err error) {
	switch {
	case err == nil:
		m.metrics.sent.Add(ctx, 1, m.metricOpt(route, target, ""))
		route.promSent.Inc()
	case errors.Is(err, batch.ErrMessageBufferFull), errors.Is(err, batch.ErrClientShutdown):
		m.recordDropped(ctx, route, target, batch.ErrorCodeFor(err))
	default:
		m.metrics.sendErrors.Add(ctx, 1, m.metricOpt(route, target, batch.ErrorCodeFor(err)))
		route.promSendErrors.Inc()
	}
}

func (m *DomainManager) recordDropped(ctx context.Context, route *telemetryRoute, target eventTarget, code string) {
	m.metrics.dropped.Add(ctx, 1, m.metricOpt(route, target, code))
	route.promDropped.Inc()
}

// metricKey identifies a cached metric option: the target chain's selector, the route's
// telemetry type and the error code ("" for the sent counter).
type metricKey struct {
	selector  string
	telemType synchronization.TelemetryType
	errorCode string
}

// metricOpt returns the route's attributes plus the target chain's chain_selector and, when code
// is non-empty, the error code. Cached per metricKey: a bounded set.
func (m *DomainManager) metricOpt(route *telemetryRoute, target eventTarget, code string) otelmetric.AddOption {
	key := metricKey{selector: target.selector, telemType: route.telemType, errorCode: code}
	if v, ok := m.metricOpts.Load(key); ok {
		return v.(otelmetric.AddOption)
	}
	kvs := append(route.attrs.ToSlice(), attribute.String("chain_selector", target.selector))
	if code != "" {
		kvs = append(kvs, attribute.String("error_code", code))
	}
	v, _ := m.metricOpts.LoadOrStore(key, otelmetric.WithAttributes(kvs...))
	return v.(otelmetric.AddOption)
}

// unsupportedTypeLabel clamps an unsupported telemetry type for metric labels and log
// dedup keys: the generic telemetry adapter accepts arbitrary strings, so a raw value
// would grow the metric series and the option cache without bound. Types with a
// CloudEvent mapping (synchronization.TelemetryTypeToDomainAndEntity is the registry)
// keep their label; anything else is "unknown" — its raw value is logged once by
// warnUnsupportedType.
func unsupportedTypeLabel(telemType synchronization.TelemetryType) string {
	if _, _, err := synchronization.TelemetryTypeToDomainAndEntity(telemType); err == nil {
		return string(telemType)
	}
	return "unknown"
}

// unsupportedKey identifies a cached unsupported-metric option: the clamped telemetry
// type label and the drop reason. Bounded by construction — labels come from the
// CloudEvent registry plus "unknown", reasons from a fixed set.
type unsupportedKey struct {
	telemType string
	reason    string
}

// countUnsupported counts a dropped message whose type or chain is not supported. The
// telemetry_type label is clamped (see unsupportedTypeLabel) and the domain is always
// recorded ("unknown" when the type maps to none), so domain-filtered loss dashboards
// see every drop. The option is cached per (label, reason): hot path.
func (m *DomainManager) countUnsupported(telemType synchronization.TelemetryType, reason string) {
	label := unsupportedTypeLabel(telemType)
	key := unsupportedKey{telemType: label, reason: reason}
	v, ok := m.unsupportedOpts.Load(key)
	if !ok {
		d, mapped := supportedTelemetryDomains[telemType]
		if !mapped {
			d = "unknown"
		}
		attrs := []attribute.KeyValue{
			attribute.String("telemetry_type", label),
			attribute.String("reason", reason),
			attribute.String("domain", string(d)),
		}
		v, _ = m.unsupportedOpts.LoadOrStore(key, otelmetric.WithAttributes(attrs...))
	}
	m.metrics.unsupported.Add(m.ctx, 1, v.(otelmetric.AddOption))
}

func (m *DomainManager) warnUnsupportedType(telemType synchronization.TelemetryType) {
	m.reportOnce("type:"+string(telemType), func() {
		m.eng.Warnw("Telemetry type is not supported: its telemetry will be dropped", "telemType", telemType)
	})
}

// newEvent builds the CloudEvent legacy-mode chip ingress expects. Each extension becomes a
// "ce_<name>" Kafka header; the comments name the header OTI used for the same value. Chip
// ingress derives nodeoperatorname/nodename (from the CSA key's NOP info) and receivedat itself.
func (m *DomainManager) newEvent(route *telemetryRoute, target eventTarget, telemetry []byte) (chipingress.CloudEvent, error) {
	now := m.now()
	event, err := chipingress.NewEvent(route.source, route.entity, telemetry, map[string]any{
		"time":         now.UTC(),
		"recordedtime": now.UTC(),
	})
	if err != nil {
		return event, fmt.Errorf("failed creating CloudEvent: %w", err)
	}
	event.SetExtension("legacytelemetry", "true")                       // marks node-sent legacy telemetry
	event.SetExtension("telemetrytype", string(route.telemType))        // "telemetry-type"
	event.SetExtension("chainselector", target.selector)                // chain-selectors ID of the job's chain
	event.SetExtension("networkname", target.network)                   // lower case, e.g. "evm"
	event.SetExtension("contractid", target.contractID)                 // "contract-address"
	event.SetExtension("csapublickey", m.csaPubKeyHex)                  // "node-operator-csa-public-key"
	event.SetExtension("partitionkey", target.partitionKey)             // Kafka record key, as OTI's
	event.SetExtension("sentat", strconv.FormatInt(now.UnixNano(), 10)) // "sent-at", unix ns as over WSRPC
	return event, nil
}

// domainAgent sends one telemetry type through its domain's emitter.
type domainAgent struct {
	m      *DomainManager
	route  *telemetryRoute
	target eventTarget
}

func (a *domainAgent) SendLog(telemetry []byte) {
	a.m.send(a.route, a.target, telemetry)
}

// multitypeDomainAgent routes each message by the type it carries.
type multitypeDomainAgent struct {
	m      *DomainManager
	target eventTarget
}

func (a *multitypeDomainAgent) SendTypedLog(telemType synchronization.TelemetryType, telemetry []byte) {
	route, ok := a.m.routes[telemType]
	if !ok {
		a.m.warnUnsupportedType(telemType)
		a.m.countUnsupported(telemType, unsupportedReasonTelemetryType)
		return
	}
	a.m.send(route, a.target, telemetry)
}

// unsupportedAgent drops everything it is given and counts it. telemType is empty for a
// multitype endpoint, whose messages carry their own type.
type unsupportedAgent struct {
	m         *DomainManager
	telemType synchronization.TelemetryType
	reason    string
}

func (a *unsupportedAgent) SendLog([]byte) {
	a.m.countUnsupported(a.telemType, a.reason)
}

func (a *unsupportedAgent) SendTypedLog(telemType synchronization.TelemetryType, _ []byte) {
	if a.telemType != "" {
		telemType = a.telemType
	}
	a.m.countUnsupported(telemType, a.reason)
}
