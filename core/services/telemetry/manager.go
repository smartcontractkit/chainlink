package telemetry

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	otelmetric "go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	common "github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
)

type Manager struct {
	services.Service
	eng *services.Engine

	bufferSize uint
	endpoints  []*telemetryEndpoint
	ks         keystore.CSA

	logging      bool
	maxBatchSize uint
	sendInterval time.Duration
	sendTimeout  time.Duration
	uniConn      bool
	useBatchSend bool

	chipService synchronization.ChipIngressService // nil when chip ingress is disabled

	chipIngressEndpoint string
	configGauge         otelmetric.Int64Gauge
}

type telemetryEndpoint struct {
	ChainID           string
	Network           string
	URL               *url.URL
	client            synchronization.TelemetryService
	PubKey            string
	chipIngressClient synchronization.ChipIngressService
}

// NewManager create a new telemetry manager that is responsible for configuring telemetry agents and generating the defined telemetry endpoints and monitoring endpoints
// chipService is the shared chip-ingress service to use when ChipIngressEnabled is true; nil when the flag is off.
func NewManager(cfg config.TelemetryIngress, csaKeyStore keystore.CSA, chipService synchronization.ChipIngressService, lggr common.Logger) *Manager {
	var chipIngressEndpoint string // "" when disabled
	if chipService != nil {
		lggr.Info("ChIP Ingress is enabled for telemetry")
		chipIngressEndpoint = cfg.ChipIngressEndpoint()
	}

	m := &Manager{
		bufferSize:          cfg.BufferSize(),
		ks:                  csaKeyStore,
		logging:             cfg.Logging(),
		maxBatchSize:        cfg.MaxBatchSize(),
		sendInterval:        cfg.SendInterval(),
		sendTimeout:         cfg.SendTimeout(),
		uniConn:             cfg.UniConn(),
		useBatchSend:        cfg.UseBatchSend(),
		chipService:         chipService,
		chipIngressEndpoint: chipIngressEndpoint,
	}
	m.Service, m.eng = services.Config{
		Name: "TelemetryManager",
		Start: func(ctx context.Context) error {
			if err := m.recordConfigMetric(ctx); err != nil {
				m.eng.Errorw("failed to record telemetry ingress config metric", "err", err)
			}
			m.eng.GoTick(services.TickerConfig{}.NewTicker(time.Hour), func(ctx context.Context) {
				if err := m.recordConfigMetric(ctx); err != nil {
					m.eng.Errorw("failed to record telemetry ingress config metric", "err", err)
				}
			})
			return nil
		},
		NewSubServices: func(lggr common.Logger) (subs []services.Service) {
			if m.chipService != nil {
				subs = append(subs, m.chipService)
			}
			for _, e := range cfg.Endpoints() {
				sub, err := m.newEndpoint(e, lggr, cfg)
				if err != nil {
					lggr.Error(err)
					continue
				}
				if sub != nil {
					subs = append(subs, sub)
				}
			}
			return
		},
	}.NewServiceEngine(lggr)

	return m
}

// GenMonitoringEndpoint creates a new monitoring endpoints based on the existing available endpoints defined in the core config TOML, if no endpoint for the network and chainID exists, a NOOP agent will be used and the telemetry will not be sent
func (m *Manager) GenMonitoringEndpoint(network string, chainID string, contractID string, telemType synchronization.TelemetryType) commontypes.MonitoringEndpoint {
	e, found := m.getEndpoint(network, chainID)

	if !found {
		m.eng.Warnf("no telemetry endpoint found for network %q chainID %q, telemetry %q for contractID %q will NOT be sent", network, chainID, telemType, contractID)
		return &NoopAgent{}
	}

	if m.chipService != nil {
		lggr := m.eng.Named("chip-ingress")
		adapter, err := NewChipIngressAgent(e.chipIngressClient, network, chainID, contractID, telemType, lggr)
		if err != nil {
			m.eng.Errorw("failed to create ChIP ingress agent, falling back to noop", "error", err, "network", network, "chainID", chainID, "contractID", contractID, "telemType", telemType)
			return &NoopAgent{}
		}
		return adapter
	}

	if m.useBatchSend {
		return NewTypedIngressAgentBatch(e.client, network, chainID, contractID, telemType)
	}

	return NewTypedIngressAgent(e.client, network, chainID, contractID, telemType)
}

func (m *Manager) GenMultitypeMonitoringEndpoint(network string, chainID string, contractID string) MultitypeMonitoringEndpoint {
	e, found := m.getEndpoint(network, chainID)

	if !found {
		m.eng.Warnf("no telemetry endpoint found for network %q chainID %q, telemetry for contractID %q will NOT be sent", network, chainID, contractID)
		return &NoopAgent{}
	}

	if m.chipService != nil {
		lggr := m.eng.Named("chip-ingress-multipletype")
		adapter, err := NewChipIngressAgentMultitype(e.chipIngressClient, network, chainID, contractID, lggr)
		if err != nil {
			m.eng.Errorw("failed to create ChIP ingress multitype agent, falling back to noop", "error", err, "network", network, "chainID", chainID, "contractID", contractID)
			return &NoopAgent{}
		}
		return adapter
	}

	if m.useBatchSend {
		return NewMultiIngressAgentBatch(e.client, network, chainID, contractID)
	}

	return NewMultiIngressAgent(e.client, network, chainID, contractID)
}

func (m *Manager) newEndpoint(e config.TelemetryIngressEndpoint, lggr common.Logger, cfg config.TelemetryIngress) (services.Service, error) {
	if e.Network() == "" {
		return nil, errors.New("cannot add telemetry endpoint, network cannot be empty")
	}

	if e.ChainID() == "" {
		return nil, errors.New("cannot add telemetry endpoint, chainID cannot be empty")
	}

	if e.URL() == nil {
		return nil, errors.New("cannot add telemetry endpoint, URL cannot be empty")
	}

	if e.ServerPubKey() == "" {
		return nil, errors.New("cannot add telemetry endpoint, ServerPubKey cannot be empty")
	}

	if _, found := m.getEndpoint(e.Network(), e.ChainID()); found {
		return nil, errors.Errorf("cannot add telemetry endpoint for network %q and chainID %q, endpoint already exists", e.Network(), e.ChainID())
	}

	lggr = common.Sugared(lggr).Named(e.Network()).Named(e.ChainID())

	if m.chipService != nil {
		lggr.Infof("Using chip-ingress service for network %q chainID %q", e.Network(), e.ChainID())
		// When ChIP ingress is enabled, all endpoints share the one ChipIngressService
		// (injected by the application); entries only select which chains send telemetry.
		te := telemetryEndpoint{
			Network:           strings.ToUpper(e.Network()),
			ChainID:           strings.ToUpper(e.ChainID()),
			URL:               e.URL(),
			PubKey:            e.ServerPubKey(),
			chipIngressClient: m.chipService,
		}
		m.endpoints = append(m.endpoints, &te)
		return nil, nil
	}

	// Otherwise use the traditional telemetry service
	var tClient synchronization.TelemetryService
	if m.useBatchSend {
		tClient = synchronization.NewTelemetryIngressBatchClient(e.URL(), e.ServerPubKey(), m.ks, cfg.Logging(), lggr, cfg.BufferSize(), cfg.MaxBatchSize(), cfg.SendInterval(), cfg.SendTimeout(), cfg.UniConn())
	} else {
		tClient = synchronization.NewTelemetryIngressClient(e.URL(), e.ServerPubKey(), m.ks, lggr, cfg.BufferSize())
	}

	te := telemetryEndpoint{
		Network: strings.ToUpper(e.Network()),
		ChainID: strings.ToUpper(e.ChainID()),
		URL:     e.URL(),
		client:  tClient,
		PubKey:  e.ServerPubKey(),
	}

	m.endpoints = append(m.endpoints, &te)
	return te.client, nil
}

// recordConfigMetric records the telemetry ingress config info metric, mirroring
// beholder's ConfigRecorder "config.info" pattern.
func (m *Manager) recordConfigMetric(ctx context.Context) error {
	if m.configGauge == nil {
		gauge, err := beholder.GetMeter().Int64Gauge("telemetry_ingress.config.info",
			otelmetric.WithDescription("Telemetry ingress config info metric"),
			otelmetric.WithUnit("{info}"))
		if err != nil {
			return fmt.Errorf("failed to create telemetry ingress config metric: %w", err)
		}
		m.configGauge = gauge
	}
	m.configGauge.Record(ctx, 1, otelmetric.WithAttributes(
		attribute.Bool("chip_ingress_enabled", m.chipService != nil),
		// "" when disabled
		attribute.String("chip_ingress_endpoint", m.chipIngressEndpoint),
	))
	return nil
}

func (m *Manager) getEndpoint(network string, chainID string) (*telemetryEndpoint, bool) {
	for _, e := range m.endpoints {
		if e.Network == strings.ToUpper(network) && e.ChainID == strings.ToUpper(chainID) {
			return e, true
		}
	}
	return nil, false
}
