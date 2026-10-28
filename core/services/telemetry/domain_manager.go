package telemetry

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	chainselector "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/libocr/commontypes"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelmetric "go.opentelemetry.io/otel/metric"

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

// telemetryRoute is where one telemetry type goes: the Beholder domain whose emitter carries it,
// and the CloudEvent source and type it is sent with.
type telemetryRoute struct {
	domain beholder.Domain
	source string
	entity string
}

const (
	ocr2Entity = "offchainreporting2.TelemetryWrapper"
	ocr3Entity = "offchainreporting3.TelemetryWrapper"
)

// telemetryRoutes maps every supported telemetry type to its route. Types not listed (functions,
// llo-report, threshold, s4, rebalancer, CRE plugin types, ...) are not supported and are dropped.
var telemetryRoutes = map[synchronization.TelemetryType]telemetryRoute{
	synchronization.OCR:                      {beholder.DomainDataFeeds, "data-feeds.telemetry.ocr", "offchainreporting.TelemetryWrapper"},
	synchronization.OCR2Median:               {beholder.DomainDataFeeds, "data-feeds.telemetry.ocr2-median", ocr2Entity},
	synchronization.OCR3DataFeeds:            {beholder.DomainDataFeeds, "data-feeds.telemetry.ocr3-data-feeds", ocr3Entity},
	synchronization.EnhancedEA:               {beholder.DomainDataFeeds, "data-feeds.telemetry.enhanced-ea", "telem.EnhancedEA"},
	synchronization.OCR3Mercury:              {beholder.DomainDataStreams, "data-streams.telemetry.ocr3-mercury", ocr3Entity},
	synchronization.EnhancedEAMercury:        {beholder.DomainDataStreams, "data-streams.telemetry.enhanced-ea-mercury", "telem.EnhancedEAMercury"},
	synchronization.LLOObservation:           {beholder.DomainDataStreams, "data-streams.telemetry.llo-observation", "telem.LLOObservationTelemetry"},
	synchronization.LLOOutcome:               {beholder.DomainDataStreams, "data-streams.telemetry.llo-outcome", "telem.LLOOutcomeTelemetry"},
	synchronization.LLOAttributedObservation: {beholder.DomainDataStreams, "data-streams.telemetry.llo-attributed-observation", "telem.LLOAttributedObservationTelemetry"},
	synchronization.PipelineBridge:           {beholder.DomainDataStreams, "data-streams.telemetry.pipeline-bridge", "telem.LLOBridgeTelemetry"},
	synchronization.OCR2CCIPCommit:           {beholder.DomainCCIP, "ccip.telemetry.ocr2", ocr2Entity},
	synchronization.OCR2CCIPExec:             {beholder.DomainCCIP, "ccip.telemetry.ocr2", ocr2Entity},
	synchronization.OCR3CCIPCommit:           {beholder.DomainCCIP, "ccip.telemetry.ocr3", ocr3Entity},
	synchronization.OCR3CCIPExec:             {beholder.DomainCCIP, "ccip.telemetry.ocr3", ocr3Entity},
	// TODO(INFOPLAT-19709): head reports are emitted for every chain on every node; confirm ccip
	// is their domain.
	synchronization.HeadReport: {beholder.DomainCCIP, "ccip.telemetry.head-report", "telem.HeadReportRequest"},
}

// Metric names and attribute values. The manager's metrics are OTel (Beholder) metrics, so they
// reach us from every node rather than only from NOPs' own Prometheus.
const (
	metricMessagesSent        = "telemetry.messages_sent"
	metricMessagesSendErrors  = "telemetry.messages_send_errors"
	metricMessagesDropped     = "telemetry.messages_dropped"
	metricMessagesUnsupported = "telemetry.messages_unsupported"

	unsupportedReasonTelemetryType = "telemetry_type"
	unsupportedReasonChain         = "chain"
)

var _ MonitoringEndpointGenerator = (*DomainManager)(nil)

// DomainManager is a MonitoringEndpointGenerator that sends telemetry through Beholder's
// per-domain emitters. Each telemetry type maps to a domain; endpoints bind their domain's
// emitter when created. Network and chain only describe the event: nothing is configured or
// routed per network.
type DomainManager struct {
	services.Service
	eng *services.Engine

	emitters     map[beholder.Domain]beholder.DomainEmitter
	csaPubKeyHex string
	metrics      domainManagerMetrics
	now          func() time.Time
}

type domainManagerMetrics struct {
	sent        otelmetric.Int64Counter
	sendErrors  otelmetric.Int64Counter
	dropped     otelmetric.Int64Counter
	unsupported otelmetric.Int64Counter
}

// NewDomainManager resolves an emitter for every domain a telemetry route uses; a missing one is
// an error. csaPubKeyHex is the node's CSA public key, which legacy-mode chip ingress binds to
// the authenticated key.
func NewDomainManager(src DomainEmitterSource, csaPubKeyHex string, lggr logger.Logger) (*DomainManager, error) {
	if csaPubKeyHex == "" {
		return nil, errors.New("telemetry: CSA public key is required")
	}
	emitters := make(map[beholder.Domain]beholder.DomainEmitter)
	for _, route := range telemetryRoutes {
		if _, ok := emitters[route.domain]; ok {
			continue
		}
		em, err := src.ForDomain(route.domain)
		if err != nil {
			return nil, fmt.Errorf("telemetry: no Beholder emitter for domain %q: %w", route.domain, err)
		}
		emitters[route.domain] = em
	}
	metrics, err := newDomainManagerMetrics()
	if err != nil {
		return nil, err
	}

	m := &DomainManager{
		emitters:     emitters,
		csaPubKeyHex: csaPubKeyHex,
		metrics:      metrics,
		now:          time.Now,
	}
	m.Service, m.eng = services.Config{Name: "DomainManager"}.NewServiceEngine(lggr)
	return m, nil
}

func newDomainManagerMetrics() (m domainManagerMetrics, err error) {
	meter := otel.Meter("github.com/smartcontractkit/chainlink/v2/core/services/telemetry")
	if m.sent, err = meter.Int64Counter(metricMessagesSent,
		otelmetric.WithDescription("Telemetry messages accepted by chip ingress")); err != nil {
		return m, err
	}
	if m.sendErrors, err = meter.Int64Counter(metricMessagesSendErrors,
		otelmetric.WithDescription("Telemetry messages that failed to publish")); err != nil {
		return m, err
	}
	if m.dropped, err = meter.Int64Counter(metricMessagesDropped,
		otelmetric.WithDescription("Telemetry messages dropped before publishing")); err != nil {
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
	route, ok := telemetryRoutes[telemType]
	if !ok {
		m.eng.Warnw("Unsupported telemetry type: its telemetry will be dropped", "telemType", telemType, "network", network, "chainID", chainID, "contractID", contractID)
		return m.unsupportedAgent(telemType, unsupportedReasonTelemetryType)
	}
	selector, err := chainSelector(network, chainID)
	if err != nil {
		m.eng.Warnw("Unsupported chain: its telemetry will be dropped", "err", err, "telemType", telemType, "network", network, "chainID", chainID, "contractID", contractID)
		return m.unsupportedAgent(telemType, unsupportedReasonChain)
	}
	return &domainAgent{
		m:          m,
		route:      route,
		emitter:    m.emitters[route.domain],
		telemType:  telemType,
		network:    network,
		contractID: contractID,
		selector:   selector,
		attrs:      routeAttrs(route, telemType),
	}
}

// GenMultitypeMonitoringEndpoint returns an endpoint whose messages carry their own type.
func (m *DomainManager) GenMultitypeMonitoringEndpoint(network string, chainID string, contractID string) MultitypeMonitoringEndpoint {
	selector, err := chainSelector(network, chainID)
	if err != nil {
		m.eng.Warnw("Unsupported chain: its telemetry will be dropped", "err", err, "network", network, "chainID", chainID, "contractID", contractID)
		return m.unsupportedAgent("", unsupportedReasonChain)
	}
	return &multitypeDomainAgent{m: m, network: network, contractID: contractID, selector: selector}
}

func chainSelector(network, chainID string) (uint64, error) {
	details, err := chainselector.GetChainDetailsByChainIDAndFamily(chainID, strings.ToLower(network))
	if err != nil {
		return 0, err
	}
	return details.ChainSelector, nil
}

func routeAttrs(route telemetryRoute, telemType synchronization.TelemetryType) attribute.Set {
	return attribute.NewSet(
		attribute.String("domain", string(route.domain)),
		attribute.String("telemetry_type", string(telemType)),
	)
}

// send builds the event and hands it to emitter. Fates are counted, never logged per message.
func (m *DomainManager) send(emitter beholder.DomainEmitter, route telemetryRoute, telemType synchronization.TelemetryType,
	network, contractID string, selector uint64, attrs attribute.Set, telemetry []byte) {
	// SendLog has no context (libocr's MonitoringEndpoint); the emitter only queues.
	ctx := context.Background()
	ev, err := m.newEvent(route, telemType, network, contractID, selector, telemetry)
	if err != nil {
		m.metrics.dropped.Add(ctx, 1, otelmetric.WithAttributeSet(attrs))
		return
	}
	err = emitter.EmitEvent(ctx, ev, func(sendErr error) {
		if sendErr == nil {
			m.metrics.sent.Add(ctx, 1, otelmetric.WithAttributeSet(attrs))
			return
		}
		errAttrs := append(attrs.ToSlice(), attribute.String("error_code", batch.ErrorCodeFor(sendErr)))
		m.metrics.sendErrors.Add(ctx, 1, otelmetric.WithAttributes(errAttrs...))
	})
	if err != nil {
		m.metrics.dropped.Add(ctx, 1, otelmetric.WithAttributeSet(attrs))
	}
}

// newEvent builds the CloudEvent legacy-mode chip ingress expects. Each extension becomes a
// "ce_<name>" Kafka header; the comments name the header OTI used for the same value. Chip
// ingress derives nodeoperatorname/nodename (from the CSA key's NOP info) and receivedat itself.
func (m *DomainManager) newEvent(route telemetryRoute, telemType synchronization.TelemetryType,
	network, contractID string, selector uint64, telemetry []byte) (chipingress.CloudEvent, error) {
	now := m.now()
	event, err := chipingress.NewEvent(route.source, route.entity, telemetry, map[string]any{"time": now.UTC()})
	if err != nil {
		return event, fmt.Errorf("failed creating CloudEvent: %w", err)
	}
	event.SetExtension("legacytelemetry", "true")                         // marks node-sent legacy telemetry
	event.SetExtension("telemetrytype", string(telemType))                // "telemetry-type"
	event.SetExtension("chainselector", strconv.FormatUint(selector, 10)) // chain-selectors ID of the job's chain
	event.SetExtension("networkname", network)                            // e.g. "EVM"
	event.SetExtension("contractid", contractID)                          // "contract-address"
	event.SetExtension("csapublickey", m.csaPubKeyHex)                    // "node-operator-csa-public-key"
	event.SetExtension("partitionkey", contractID+"-"+m.csaPubKeyHex)     // Kafka record key
	event.SetExtension("sentat", strconv.FormatInt(now.UnixNano(), 10))   // "sent-at", unix ns
	return event, nil
}

func (m *DomainManager) unsupportedAgent(telemType synchronization.TelemetryType, reason string) *unsupportedAgent {
	return &unsupportedAgent{m: m, telemType: telemType, reason: reason}
}

// domainAgent sends one telemetry type through its domain's emitter.
type domainAgent struct {
	m          *DomainManager
	route      telemetryRoute
	emitter    beholder.DomainEmitter
	telemType  synchronization.TelemetryType
	network    string
	contractID string
	selector   uint64
	attrs      attribute.Set
}

func (a *domainAgent) SendLog(telemetry []byte) {
	a.m.send(a.emitter, a.route, a.telemType, a.network, a.contractID, a.selector, a.attrs, telemetry)
}

// multitypeDomainAgent routes each message by the type it carries.
type multitypeDomainAgent struct {
	m          *DomainManager
	network    string
	contractID string
	selector   uint64
}

func (a *multitypeDomainAgent) SendTypedLog(telemType synchronization.TelemetryType, telemetry []byte) {
	route, ok := telemetryRoutes[telemType]
	if !ok {
		a.m.unsupportedAgent(telemType, unsupportedReasonTelemetryType).SendLog(telemetry)
		return
	}
	a.m.send(a.m.emitters[route.domain], route, telemType, a.network, a.contractID, a.selector, routeAttrs(route, telemType), telemetry)
}

// unsupportedAgent drops everything it is given and counts it.
type unsupportedAgent struct {
	m         *DomainManager
	telemType synchronization.TelemetryType
	reason    string
}

func (a *unsupportedAgent) SendLog([]byte) {
	a.m.metrics.unsupported.Add(context.Background(), 1, otelmetric.WithAttributes(
		attribute.String("telemetry_type", string(a.telemType)),
		attribute.String("reason", a.reason),
	))
}

func (a *unsupportedAgent) SendTypedLog(telemType synchronization.TelemetryType, telemetry []byte) {
	if a.telemType == "" {
		a.m.unsupportedAgent(telemType, a.reason).SendLog(telemetry)
		return
	}
	a.SendLog(telemetry)
}
