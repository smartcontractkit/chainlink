package telemetry

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/chipingress"
	"github.com/smartcontractkit/chainlink-common/pkg/chipingress/batch"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
)

const testCSAPubKeyHex = "755294c8a1f0"

// fakeDomainEmitter records EmitEvent calls. Like the real emitter, it can report each event's
// fate from another goroutine (async).
type fakeDomainEmitter struct {
	async       bool
	emitErr     error // returned by EmitEvent; the callback is then not invoked
	callbackErr error // reported to the callback as the event's fate

	mu     sync.Mutex
	events []chipingress.CloudEvent
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
	if callback == nil {
		return nil
	}
	if f.async {
		go callback(f.callbackErr)
	} else {
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

// newTestDomainManager returns a manager whose OTel metrics go to its own reader, so tests can
// run in parallel.
func newTestDomainManager(t *testing.T, src DomainEmitterSource) (*DomainManager, *sdkmetric.ManualReader) {
	t.Helper()
	reader := sdkmetric.NewManualReader()
	m, err := NewDomainManager(src, testCSAPubKeyHex, logger.Test(t),
		WithMeterProvider(sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))))
	require.NoError(t, err)
	return m, reader
}

// counterValue sums the int64 counter name over data points whose attributes include want.
func counterValue(t *testing.T, reader *sdkmetric.ManualReader, name string, want ...attribute.KeyValue) int64 {
	t.Helper()
	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &rm))
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

var ethMainnetChainID = strconv.FormatUint(chainselectors.ETHEREUM_MAINNET.EvmChainID, 10)

func TestNewDomainManager_RequiresEveryDomain(t *testing.T) {
	t.Parallel()
	src := newFakeEmitterSource()
	delete(src, beholder.DomainCCIP)
	delete(src, beholder.DomainDataStreams)

	_, err := NewDomainManager(src, testCSAPubKeyHex, logger.Test(t))
	require.Error(t, err)
	assert.Contains(t, err.Error(), `"ccip"`, "domains are resolved in a stable order, so the same one is reported every time")
}

func TestNewDomainManager_RequiresEmitterSourceAndCSAPublicKey(t *testing.T) {
	t.Parallel()
	_, err := NewDomainManager(nil, testCSAPubKeyHex, logger.Test(t))
	require.Error(t, err)
	_, err = NewDomainManager(newFakeEmitterSource(), "", logger.Test(t))
	require.Error(t, err)
}

func TestDomainManager_RoutesEachTypeToItsDomain(t *testing.T) {
	t.Parallel()
	const ocr2, ocr3 = "offchainreporting2.TelemetryWrapper", "offchainreporting3.TelemetryWrapper"
	for _, tc := range []struct {
		telemType synchronization.TelemetryType
		domain    beholder.Domain
		source    string
		entity    string
	}{
		{synchronization.OCR, beholder.DomainDataFeeds, "data-feeds.telemetry.ocr", "offchainreporting.TelemetryWrapper"},
		{synchronization.OCR2Median, beholder.DomainDataFeeds, "data-feeds.telemetry.ocr2-median", ocr2},
		{synchronization.OCR3DataFeeds, beholder.DomainDataFeeds, "data-feeds.telemetry.ocr3-data-feeds", ocr3},
		{synchronization.EnhancedEA, beholder.DomainDataFeeds, "data-feeds.telemetry.enhanced-ea", "telem.EnhancedEA"},
		{synchronization.OCR3Mercury, beholder.DomainDataStreams, "data-streams.telemetry.ocr3-mercury", ocr3},
		{synchronization.EnhancedEAMercury, beholder.DomainDataStreams, "data-streams.telemetry.enhanced-ea-mercury", "telem.EnhancedEAMercury"},
		{synchronization.LLOObservation, beholder.DomainDataStreams, "data-streams.telemetry.llo-observation", "telem.LLOObservationTelemetry"},
		{synchronization.LLOOutcome, beholder.DomainDataStreams, "data-streams.telemetry.llo-outcome", "telem.LLOOutcomeTelemetry"},
		{synchronization.PipelineBridge, beholder.DomainDataStreams, "data-streams.telemetry.pipeline-bridge", "telem.LLOBridgeTelemetry"},
		{synchronization.OCR2CCIPCommit, beholder.DomainCCIP, "ccip.telemetry.ocr2", ocr2},
		{synchronization.OCR2CCIPExec, beholder.DomainCCIP, "ccip.telemetry.ocr2", ocr2},
		{synchronization.OCR3CCIPCommit, beholder.DomainCCIP, "ccip.telemetry.ocr3", ocr3},
		{synchronization.OCR3CCIPExec, beholder.DomainCCIP, "ccip.telemetry.ocr3", ocr3},
		{synchronization.HeadReport, beholder.DomainCCIP, "ccip.telemetry.head-report", "telem.HeadReportRequest"},
	} {
		t.Run(string(tc.telemType), func(t *testing.T) {
			t.Parallel()
			src := newFakeEmitterSource()
			m, _ := newTestDomainManager(t, src)

			m.GenMonitoringEndpoint("EVM", ethMainnetChainID, "0xabc", tc.telemType).SendLog([]byte("payload"))

			for d, em := range src {
				if d != tc.domain {
					assert.Empty(t, em.sent(), "nothing on %s", d)
					continue
				}
				require.Len(t, em.sent(), 1, "sent on %s", d)
				assert.Equal(t, tc.source, em.sent()[0].Source())
				assert.Equal(t, tc.entity, em.sent()[0].Type())
			}
		})
	}
}

func TestDomainManager_DropsTypesItDoesNotSupport(t *testing.T) {
	t.Parallel()
	src := newFakeEmitterSource()
	m, reader := newTestDomainManager(t, src)

	unsupported := []synchronization.TelemetryType{
		synchronization.OCR2Functions,
		synchronization.FunctionsRequests,
		synchronization.LLOAttributedObservation,
		synchronization.LLOReport,
		synchronization.OCR2Threshold,
		synchronization.OCR2S4,
		synchronization.OCR3Rebalancer,
		synchronization.OCR3CCIPBootstrap,
		"vault-plugin",
	}
	for _, tt := range unsupported {
		ep := m.GenMonitoringEndpoint("EVM", ethMainnetChainID, "0xabc", tt)
		require.NotNil(t, ep)
		ep.SendLog([]byte("p"))
		ep.SendLog([]byte("p"))
	}

	for d, em := range src {
		assert.Empty(t, em.sent(), "nothing sent on %s", d)
	}
	assert.Equal(t, int64(2), counterValue(t, reader, metricMessagesUnsupported,
		attribute.String("telemetry_type", string(synchronization.LLOAttributedObservation)),
		attribute.String("reason", unsupportedReasonTelemetryType)), "counted per message")
	assert.Equal(t, int64(2*len(unsupported)), counterValue(t, reader, metricMessagesUnsupported))
}

func TestDomainManager_EventCarriesLegacyTelemetryExtensions(t *testing.T) {
	t.Parallel()
	src := newFakeEmitterSource()
	m, _ := newTestDomainManager(t, src)
	now := time.Date(2026, 10, 28, 12, 0, 0, 0, time.UTC)
	m.now = func() time.Time { return now }

	m.GenMonitoringEndpoint("EVM", ethMainnetChainID, "0xabc", synchronization.OCR2Median).SendLog([]byte("payload"))

	events := src[beholder.DomainDataFeeds].sent()
	require.Len(t, events, 1)
	ev := events[0]

	require.NoError(t, ev.Validate())
	assert.Equal(t, []byte("payload"), ev.Data())
	assert.True(t, ev.Time().Equal(now))

	ext := ev.Extensions()
	assert.Equal(t, "true", ext["legacytelemetry"])
	assert.Equal(t, "ocr2-median", ext["telemetrytype"])
	assert.Equal(t, strconv.FormatUint(chainselectors.ETHEREUM_MAINNET.Selector, 10), ext["chainselector"])
	assert.Equal(t, "evm", ext["networkname"], "network is normalized to lower case")
	assert.Equal(t, "0xabc", ext["contractid"])
	assert.Equal(t, testCSAPubKeyHex, ext["csapublickey"])
	assert.Equal(t, "0xabc-"+testCSAPubKeyHex, ext["partitionkey"])
	assert.Equal(t, strconv.FormatInt(now.UnixNano(), 10), ext["sentat"], "sentat is unix nanoseconds, as the node sends OTI today")
	recorded, ok := ext["recordedtime"]
	require.True(t, ok)
	assert.Contains(t, recorded.(interface{ String() string }).String(), "2026-10-28T12:00:00", "recordedtime uses the manager's clock")
}

func TestDomainManager_ResolvesChainsInEveryFamily(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		network, chainID string
		selector         uint64
	}{
		{"EVM", ethMainnetChainID, chainselectors.ETHEREUM_MAINNET.Selector},
		{"evm", ethMainnetChainID, chainselectors.ETHEREUM_MAINNET.Selector},
		{"Solana", "5eykt4UsFv8P8NJdTREpY1vzqKqZKvdpKuc147dw2N9d", 124615329519749607}, // genesis hash: the relayer's chain ID
		{"StarkNet", "SN_MAIN", 511843109281680063},
		{"Aptos", "1", 4741433654826277614},
		{"Tron", "728126428", 1546563616611573945},
		{"TON", "-239", 16448340667252469081},
		{"Sui", "1", 17529533435026248318},
	} {
		t.Run(tc.network+"/"+tc.chainID, func(t *testing.T) {
			t.Parallel()
			src := newFakeEmitterSource()
			m, _ := newTestDomainManager(t, src)

			m.GenMonitoringEndpoint(tc.network, tc.chainID, "0xabc", synchronization.OCR2Median).SendLog([]byte("p"))

			events := src[beholder.DomainDataFeeds].sent()
			require.Len(t, events, 1)
			assert.Equal(t, strconv.FormatUint(tc.selector, 10), events[0].Extensions()["chainselector"])
			assert.Equal(t, strings.ToLower(tc.network), events[0].Extensions()["networkname"])
		})
	}
}

func TestDomainManager_UnresolvableChainIsDroppedAndCounted(t *testing.T) {
	t.Parallel()
	src := newFakeEmitterSource()
	m, reader := newTestDomainManager(t, src)

	m.GenMonitoringEndpoint("EVM", "999999999999", "0xabc", synchronization.OCR2Median).SendLog([]byte("p"))
	mt := m.GenMultitypeMonitoringEndpoint("EVM", "999999999999", "0xabc/1")
	mt.SendTypedLog(synchronization.LLOOutcome, []byte("p"))
	mt.SendTypedLog(synchronization.LLOObservation, []byte("p"))

	for d, em := range src {
		assert.Empty(t, em.sent(), "nothing sent on %s", d)
	}
	reason := attribute.String("reason", unsupportedReasonChain)
	assert.Equal(t, int64(3), counterValue(t, reader, metricMessagesUnsupported, reason))
	assert.Equal(t, int64(1), counterValue(t, reader, metricMessagesUnsupported, reason,
		attribute.String("telemetry_type", string(synchronization.LLOOutcome))), "multitype counts carry the message's type")
}

func TestDomainManager_MultitypeEndpointRoutesPerMessage(t *testing.T) {
	t.Parallel()
	src := newFakeEmitterSource()
	m, reader := newTestDomainManager(t, src)

	ep := m.GenMultitypeMonitoringEndpoint("EVM", ethMainnetChainID, "0xabc/1")
	ep.SendTypedLog(synchronization.LLOOutcome, []byte("p"))
	ep.SendTypedLog(synchronization.LLOReport, []byte("p")) // unsupported: dropped

	events := src[beholder.DomainDataStreams].sent()
	require.Len(t, events, 1)
	assert.Equal(t, "data-streams.telemetry.llo-outcome", events[0].Source())
	assert.Equal(t, "llo-outcome", events[0].Extensions()["telemetrytype"])
	assert.Equal(t, int64(1), counterValue(t, reader, metricMessagesUnsupported,
		attribute.String("telemetry_type", string(synchronization.LLOReport))))
}

func TestDomainManager_CountsEachEventsFate(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name        string
		emitErr     error
		callbackErr error
		metric      string
		errorCode   string
	}{
		{name: "sent", metric: metricMessagesSent},
		{name: "server unavailable", callbackErr: status.Error(codes.Unavailable, "down"), metric: metricMessagesSendErrors, errorCode: codes.Unavailable.String()},
		{name: "server without no-ack", callbackErr: status.Error(codes.Unimplemented, "no"), metric: metricMessagesSendErrors, errorCode: codes.Unimplemented.String()},
		{name: "buffer full", callbackErr: batch.ErrMessageBufferFull, metric: metricMessagesDropped, errorCode: batch.ErrorCodeFor(batch.ErrMessageBufferFull)},
		{name: "emitter shutting down", callbackErr: batch.ErrClientShutdown, metric: metricMessagesDropped, errorCode: batch.ErrorCodeFor(batch.ErrClientShutdown)},
		{name: "emitter refuses the event", emitErr: errors.New("not started"), metric: metricMessagesDropped, errorCode: errorCodeEmitRejected},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			src := newFakeEmitterSource()
			src[beholder.DomainDataFeeds].async = true
			src[beholder.DomainDataFeeds].emitErr = tc.emitErr
			src[beholder.DomainDataFeeds].callbackErr = tc.callbackErr
			m, reader := newTestDomainManager(t, src)

			m.GenMonitoringEndpoint("EVM", ethMainnetChainID, "0xabc", synchronization.OCR2Median).SendLog([]byte("p"))

			want := []attribute.KeyValue{
				attribute.String("domain", string(beholder.DomainDataFeeds)),
				attribute.String("telemetry_type", string(synchronization.OCR2Median)),
			}
			if tc.errorCode != "" {
				want = append(want, attribute.String("error_code", tc.errorCode))
			}
			require.Eventually(t, func() bool { return counterValue(t, reader, tc.metric, want...) == 1 },
				5*time.Second, 10*time.Millisecond)
			for _, other := range []string{metricMessagesSent, metricMessagesSendErrors, metricMessagesDropped} {
				if other != tc.metric {
					assert.Zero(t, counterValue(t, reader, other), "%s", other)
				}
			}
		})
	}
}

func TestDomainManager_ConcurrentSends(t *testing.T) {
	t.Parallel()
	src := newFakeEmitterSource()
	for _, em := range src {
		em.async = true
	}
	m, reader := newTestDomainManager(t, src)

	const workers, perWorker = 16, 50
	var wg sync.WaitGroup
	for i := range workers {
		wg.Go(func() {
			ep := m.GenMonitoringEndpoint("EVM", ethMainnetChainID, "0x"+strconv.Itoa(i), synchronization.OCR2Median)
			mt := m.GenMultitypeMonitoringEndpoint("EVM", ethMainnetChainID, "0x"+strconv.Itoa(i)+"/1")
			for range perWorker {
				ep.SendLog([]byte("p"))
				mt.SendTypedLog(synchronization.LLOOutcome, []byte("p"))
			}
		})
	}
	wg.Wait()

	require.Eventually(t, func() bool { return counterValue(t, reader, metricMessagesSent) == 2*workers*perWorker },
		5*time.Second, 10*time.Millisecond)
	assert.Len(t, src[beholder.DomainDataFeeds].sent(), workers*perWorker)
	assert.Len(t, src[beholder.DomainDataStreams].sent(), workers*perWorker)
}

// The Prometheus counters are process-global, so this test does not run in parallel: sequential
// tests finish before parallel ones resume, which keeps the deltas below exact.
//
//nolint:paralleltest // asserts deltas on process-global Prometheus counters
func TestDomainManager_AlsoCountsInPrometheus(t *testing.T) {
	const endpoint = "beholder/data-feeds"
	tt := string(synchronization.OCR2Median)
	sent := testutil.ToFloat64(synchronization.TelemetryClientMessagesSent.WithLabelValues(endpoint, tt))
	sendErrs := testutil.ToFloat64(synchronization.TelemetryClientMessagesSendErrors.WithLabelValues(endpoint, tt))
	dropped := testutil.ToFloat64(synchronization.TelemetryClientMessagesDropped.WithLabelValues(endpoint, tt))

	for _, callbackErr := range []error{nil, status.Error(codes.Unavailable, "down"), batch.ErrMessageBufferFull} {
		src := newFakeEmitterSource()
		src[beholder.DomainDataFeeds].callbackErr = callbackErr
		m, _ := newTestDomainManager(t, src)
		m.GenMonitoringEndpoint("EVM", ethMainnetChainID, "0xabc", synchronization.OCR2Median).SendLog([]byte("p"))
	}

	assert.Equal(t, sent+1, testutil.ToFloat64(synchronization.TelemetryClientMessagesSent.WithLabelValues(endpoint, tt)))
	assert.Equal(t, sendErrs+1, testutil.ToFloat64(synchronization.TelemetryClientMessagesSendErrors.WithLabelValues(endpoint, tt)))
	assert.Equal(t, dropped+1, testutil.ToFloat64(synchronization.TelemetryClientMessagesDropped.WithLabelValues(endpoint, tt)))
}

// Beholder domain emitters reject events whose source is not under their domain, so every
// supported type's source must be the domain or a subdomain of it.
func TestSupportedTelemetry_SourcesSitUnderTheirDomain(t *testing.T) {
	t.Parallel()
	for telemType, d := range supportedTelemetryDomains() {
		source, _, err := synchronization.TelemetryTypeToDomainAndEntity(telemType)
		require.NoError(t, err, "%s has a CloudEvent source", telemType)
		assert.True(t, source == string(d) || strings.HasPrefix(source, string(d)+"."),
			"%s: source %q is not under domain %q", telemType, source, d)
	}
}
