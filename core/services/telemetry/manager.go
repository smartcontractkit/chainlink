package telemetry

import (
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
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

	logging                     bool
	maxBatchSize                uint
	sendInterval                time.Duration
	sendTimeout                 time.Duration
	uniConn                     bool
	useBatchSend                bool
	MonitoringEndpointGenerator MonitoringEndpointGenerator
}

type telemetryEndpoint struct {
	ChainID string
	Network string
	URL     *url.URL
	client  synchronization.TelemetryService
	PubKey  string
}

// NewManager create a new telemetry manager that is responsible for configuring telemetry agents and generating the defined telemetry endpoints and monitoring endpoints
func NewManager(cfg config.TelemetryIngress, csaKeyStore keystore.CSA, lggr logger.Logger) *Manager {
	m := &Manager{
		bufferSize:   cfg.BufferSize(),
		ks:           csaKeyStore,
		logging:      cfg.Logging(),
		maxBatchSize: cfg.MaxBatchSize(),
		sendInterval: cfg.SendInterval(),
		sendTimeout:  cfg.SendTimeout(),
		uniConn:      cfg.UniConn(),
		useBatchSend: cfg.UseBatchSend(),
	}
	m.Service, m.eng = services.Config{
		Name: "TelemetryManager",
		NewSubServices: func(lggr common.Logger) (subs []services.Service) {
			for _, e := range cfg.Endpoints() {
				if sub, err := m.newEndpoint(e, lggr, cfg); err != nil {
					lggr.Error(err)
				} else {
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

	if m.useBatchSend {
		return NewIngressAgentBatch(e.client, network, chainID, contractID, telemType)
	}

	return NewIngressAgent(e.client, network, chainID, contractID, telemType)
}

func (m *Manager) newEndpoint(e config.TelemetryIngressEndpoint, lggr logger.Logger, cfg config.TelemetryIngress) (services.Service, error) {
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

	lggr = logger.Sugared(lggr).Named(e.Network()).Named(e.ChainID())
	var tClient synchronization.TelemetryService
	if m.useBatchSend {
		tClient = synchronization.NewTelemetryIngressBatchClient(e.URL(), e.ServerPubKey(), m.ks, cfg.Logging(), lggr, cfg.BufferSize(), cfg.MaxBatchSize(), cfg.SendInterval(), cfg.SendTimeout(), cfg.UniConn())
	} else {
		tClient = synchronization.NewTelemetryIngressClient(e.URL(), e.ServerPubKey(), m.ks, cfg.Logging(), lggr, cfg.BufferSize())
	}

	te := telemetryEndpoint{
		Network: strings.ToUpper(e.Network()),
		ChainID: strings.ToUpper(e.ChainID()),
		URL:     e.URL(),
		PubKey:  e.ServerPubKey(),
		client:  tClient,
	}

	m.endpoints = append(m.endpoints, &te)
	return te.client, nil
}

func (m *Manager) getEndpoint(network string, chainID string) (*telemetryEndpoint, bool) {
	for _, e := range m.endpoints {
		if e.Network == strings.ToUpper(network) && e.ChainID == strings.ToUpper(chainID) {
			return e, true
		}
	}
	return nil, false
}
