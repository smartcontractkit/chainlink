package telemetry

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"sync"
	"testing"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/chipingress"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
)

const testCSAPubKeyHex = "755294c8a1f0"

// fakeDomainEmitter records EmitEvent calls and reports callbackErr as each event's fate.
type fakeDomainEmitter struct {
	mu          sync.Mutex
	events      []chipingress.CloudEvent
	emitErr     error // returned by EmitEvent (callback not invoked)
	callbackErr error // reported to the callback
}

func (f *fakeDomainEmitter) Emit(context.Context, []byte, ...any) error { return nil }
func (f *fakeDomainEmitter) Close() error                               { return nil }

func (f *fakeDomainEmitter) EmitEvent(_ context.Context, event chipingress.CloudEvent, callback func(error)) error {
	if f.emitErr != nil {
		return f.emitErr
	}
	f.mu.Lock()
	f.events = append(f.events, event)
	f.mu.Unlock()
	if callback != nil {
		callback(f.callbackErr)
	}
	return nil
}

func (f *fakeDomainEmitter) sent() []chipingress.CloudEvent {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]chipingress.CloudEvent(nil), f.events...)
}

// fakeEmitterSource stands in for *beholder.Client.
type fakeEmitterSource map[beholder.Domain]*fakeDomainEmitter

func (s fakeEmitterSource) ForDomain(d beholder.Domain) (beholder.DomainEmitter, error) {
	em, ok := s[d]
	if !ok {
		return nil, errors.New("no emitter for domain " + string(d))
	}
	return em, nil
}

func newFakeEmitterSource() fakeEmitterSource {
	return fakeEmitterSource{
		beholder.DomainDataFeeds:   {},
		beholder.DomainDataStreams: {},
		beholder.DomainCCIP:        {},
	}
}

// installTestMeter routes otel.Meter to a manual reader for the test's duration.
func installTestMeter(t *testing.T) *sdkmetric.ManualReader {
	t.Helper()
	reader := sdkmetric.NewManualReader()
	prev := otel.GetMeterProvider()
	otel.SetMeterProvider(sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader)))
	t.Cleanup(func() { otel.SetMeterProvider(prev) })
	return reader
}

// counterValue sums the int64 counter name over data points whose attributes include want.
func counterValue(t *testing.T, reader *sdkmetric.ManualReader, name string, want ...attribute.KeyValue) int64 {
	t.Helper()
	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(context.Background(), &rm))
	var total int64
	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			if m.Name != name {
				continue
			}
			sum, ok := m.Data.(metricdata.Sum[int64])
			require.True(t, ok, "%s is an int64 counter", name)
			for _, dp := range sum.DataPoints {
				matches := true
				for _, kv := range want {
					if v, ok := dp.Attributes.Value(kv.Key); !ok || v != kv.Value {
						matches = false
						break
					}
				}
				if matches {
					total += dp.Value
				}
			}
		}
	}
	return total
}

func newTestDomainManager(t *testing.T, src fakeEmitterSource) *DomainManager {
	t.Helper()
	m, err := NewDomainManager(src, testCSAPubKeyHex, logger.Test(t))
	require.NoError(t, err)
	return m
}

var ethMainnetChainID = strconv.FormatUint(chainselectors.ETHEREUM_MAINNET.EvmChainID, 10)

func TestNewDomainManager_RequiresEveryDomain(t *testing.T) {
	src := newFakeEmitterSource()
	delete(src, beholder.DomainCCIP)

	_, err := NewDomainManager(src, testCSAPubKeyHex, logger.Test(t))
	require.Error(t, err, "a missing domain emitter fails construction")
}

func TestNewDomainManager_RequiresCSAPublicKey(t *testing.T) {
	_, err := NewDomainManager(newFakeEmitterSource(), "", logger.Test(t))
	require.Error(t, err)
}

func TestDomainManager_RoutesEachTypeToItsDomain(t *testing.T) {
	for _, tc := range []struct {
		telemType synchronization.TelemetryType
		domain    beholder.Domain
		source    string
	}{
		{synchronization.OCR, beholder.DomainDataFeeds, "data-feeds.telemetry.ocr"},
		{synchronization.OCR2Median, beholder.DomainDataFeeds, "data-feeds.telemetry.ocr2-median"},
		{synchronization.OCR3DataFeeds, beholder.DomainDataFeeds, "data-feeds.telemetry.ocr3-data-feeds"},
		{synchronization.EnhancedEA, beholder.DomainDataFeeds, "data-feeds.telemetry.enhanced-ea"},
		{synchronization.OCR3Mercury, beholder.DomainDataStreams, "data-streams.telemetry.ocr3-mercury"},
		{synchronization.EnhancedEAMercury, beholder.DomainDataStreams, "data-streams.telemetry.enhanced-ea-mercury"},
		{synchronization.LLOObservation, beholder.DomainDataStreams, "data-streams.telemetry.llo-observation"},
		{synchronization.LLOOutcome, beholder.DomainDataStreams, "data-streams.telemetry.llo-outcome"},
		{synchronization.LLOAttributedObservation, beholder.DomainDataStreams, "data-streams.telemetry.llo-attributed-observation"},
		{synchronization.PipelineBridge, beholder.DomainDataStreams, "data-streams.telemetry.pipeline-bridge"},
		{synchronization.OCR2CCIPCommit, beholder.DomainCCIP, "ccip.telemetry.ocr2"},
		{synchronization.OCR2CCIPExec, beholder.DomainCCIP, "ccip.telemetry.ocr2"},
		{synchronization.OCR3CCIPCommit, beholder.DomainCCIP, "ccip.telemetry.ocr3"},
		{synchronization.OCR3CCIPExec, beholder.DomainCCIP, "ccip.telemetry.ocr3"},
		{synchronization.HeadReport, beholder.DomainCCIP, "ccip.telemetry.head-report"},
	} {
		t.Run(string(tc.telemType), func(t *testing.T) {
			src := newFakeEmitterSource()
			m := newTestDomainManager(t, src)

			m.GenMonitoringEndpoint("EVM", ethMainnetChainID, "0xabc", tc.telemType).SendLog([]byte("payload"))

			for d, em := range src {
				if d == tc.domain {
					require.Len(t, em.sent(), 1, "sent on %s", d)
					assert.Equal(t, tc.source, em.sent()[0].Source())
				} else {
					assert.Empty(t, em.sent(), "nothing on %s", d)
				}
			}
		})
	}
}

func TestDomainManager_EventCarriesLegacyTelemetryExtensions(t *testing.T) {
	src := newFakeEmitterSource()
	m := newTestDomainManager(t, src)

	m.GenMonitoringEndpoint("EVM", ethMainnetChainID, "0xabc", synchronization.OCR2Median).SendLog([]byte("payload"))

	events := src[beholder.DomainDataFeeds].sent()
	require.Len(t, events, 1)
	ev := events[0]

	require.NoError(t, ev.Validate())
	assert.Equal(t, "offchainreporting2.TelemetryWrapper", ev.Type())
	assert.Equal(t, []byte("payload"), ev.Data())
	assert.False(t, ev.Time().IsZero())

	ext := ev.Extensions()
	assert.Equal(t, "true", ext["legacytelemetry"])
	assert.Equal(t, "ocr2-median", ext["telemetrytype"])
	assert.Equal(t, strconv.FormatUint(chainselectors.ETHEREUM_MAINNET.Selector, 10), ext["chainselector"])
	assert.Equal(t, "EVM", ext["networkname"])
	assert.Equal(t, "0xabc", ext["contractid"])
	assert.Equal(t, testCSAPubKeyHex, ext["csapublickey"])
	assert.Equal(t, "0xabc-"+testCSAPubKeyHex, ext["partitionkey"])
	sentAt, ok := ext["sentat"].(string)
	require.True(t, ok)
	_, err := strconv.ParseInt(sentAt, 10, 64)
	assert.NoError(t, err, "sentat is unix nanoseconds")
}

func TestDomainManager_SendsForAnyNetworkWithoutPerNetworkConfig(t *testing.T) {
	src := newFakeEmitterSource()
	m := newTestDomainManager(t, src)

	for _, chain := range []chainselectors.Chain{chainselectors.ETHEREUM_MAINNET, chainselectors.POLYGON_MAINNET} {
		m.GenMonitoringEndpoint("EVM", strconv.FormatUint(chain.EvmChainID, 10), "0xabc", synchronization.OCR2Median).SendLog([]byte("p"))
	}
	assert.Len(t, src[beholder.DomainDataFeeds].sent(), 2)
}

func TestDomainManager_UnsupportedTypesAreDroppedAndCounted(t *testing.T) {
	reader := installTestMeter(t)
	src := newFakeEmitterSource()
	m := newTestDomainManager(t, src)

	unsupported := []synchronization.TelemetryType{
		synchronization.OCR2Functions,
		synchronization.FunctionsRequests,
		synchronization.OCR2Threshold,
		synchronization.OCR2S4,
		synchronization.OCR3Rebalancer,
		synchronization.LLOReport,
		"vault-plugin",
	}
	for _, tt := range unsupported {
		ep := m.GenMonitoringEndpoint("EVM", ethMainnetChainID, "0xabc", tt)
		require.NotNil(t, ep)
		ep.SendLog([]byte("p"))
	}

	for d, em := range src {
		assert.Empty(t, em.sent(), "nothing sent on %s", d)
	}
	assert.Equal(t, int64(1), counterValue(t, reader, metricMessagesUnsupported,
		attribute.String("telemetry_type", string(synchronization.OCR2Functions)),
		attribute.String("reason", unsupportedReasonTelemetryType)))
	assert.Equal(t, int64(len(unsupported)), counterValue(t, reader, metricMessagesUnsupported))
}

func TestDomainManager_UnresolvableChainIsDroppedAndCounted(t *testing.T) {
	reader := installTestMeter(t)
	src := newFakeEmitterSource()
	m := newTestDomainManager(t, src)

	m.GenMonitoringEndpoint("EVM", "999999999999", "0xabc", synchronization.OCR2Median).SendLog([]byte("p"))

	assert.Empty(t, src[beholder.DomainDataFeeds].sent())
	assert.Equal(t, int64(1), counterValue(t, reader, metricMessagesUnsupported,
		attribute.String("reason", unsupportedReasonChain)))
}

func TestDomainManager_MultitypeEndpointRoutesPerMessage(t *testing.T) {
	src := newFakeEmitterSource()
	m := newTestDomainManager(t, src)

	ep := m.GenMultitypeMonitoringEndpoint("EVM", ethMainnetChainID, "0xabc/1")
	ep.SendTypedLog(synchronization.LLOOutcome, []byte("p"))
	ep.SendTypedLog(synchronization.LLOReport, []byte("p")) // unmapped: dropped

	events := src[beholder.DomainDataStreams].sent()
	require.Len(t, events, 1)
	assert.Equal(t, "data-streams.telemetry.llo-outcome", events[0].Source())
	assert.Equal(t, "llo-outcome", events[0].Extensions()["telemetrytype"])
}

func TestDomainManager_MetricsFollowEachEventsFate(t *testing.T) {
	dataFeeds := attribute.String("domain", string(beholder.DomainDataFeeds))
	ocr2 := attribute.String("telemetry_type", string(synchronization.OCR2Median))

	t.Run("sent", func(t *testing.T) {
		reader := installTestMeter(t)
		m := newTestDomainManager(t, newFakeEmitterSource())
		m.GenMonitoringEndpoint("EVM", ethMainnetChainID, "0xabc", synchronization.OCR2Median).SendLog([]byte("p"))
		assert.Equal(t, int64(1), counterValue(t, reader, metricMessagesSent, dataFeeds, ocr2))
	})
	t.Run("send error from the per-message result", func(t *testing.T) {
		reader := installTestMeter(t)
		src := newFakeEmitterSource()
		src[beholder.DomainDataFeeds].callbackErr = errors.New("rejected")
		m := newTestDomainManager(t, src)
		m.GenMonitoringEndpoint("EVM", ethMainnetChainID, "0xabc", synchronization.OCR2Median).SendLog([]byte("p"))
		assert.Equal(t, int64(1), counterValue(t, reader, metricMessagesSendErrors, dataFeeds, ocr2))
		assert.Zero(t, counterValue(t, reader, metricMessagesSent))
	})
	t.Run("dropped when the emitter refuses the event", func(t *testing.T) {
		reader := installTestMeter(t)
		src := newFakeEmitterSource()
		src[beholder.DomainDataFeeds].emitErr = errors.New("stopped")
		m := newTestDomainManager(t, src)
		m.GenMonitoringEndpoint("EVM", ethMainnetChainID, "0xabc", synchronization.OCR2Median).SendLog([]byte("p"))
		assert.Equal(t, int64(1), counterValue(t, reader, metricMessagesDropped, dataFeeds, ocr2))
	})
}

// Beholder domain emitters reject events whose source is not under their domain, so every
// route's source must be the domain or a subdomain of it.
func TestTelemetryRoutes_SourcesSitUnderTheirDomain(t *testing.T) {
	for telemType, route := range telemetryRoutes {
		d := string(route.domain)
		assert.True(t, route.source == d || strings.HasPrefix(route.source, d+"."),
			"%s: source %q is not under domain %q", telemType, route.source, d)
	}
}
