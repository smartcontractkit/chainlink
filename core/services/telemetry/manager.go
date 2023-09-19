package telemetry

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/commontypes"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

//// Client encapsulates all the functionality needed to
//// send telemetry to the ingress server using wsrpc
//type Client interface {
//	services.ServiceCtx
//	Send(context.Context, synchronization.TelemPayload)
//}

type Manager struct {
	utils.StartStopOnce
	bufferSize                  uint
	endpoints                   []*telemetryEndpoint
	ks                          keystore.CSA
	lggr                        logger.Logger
	logging                     bool
	maxBatchSize                uint
	sendInterval                time.Duration
	sendTimeout                 time.Duration
	uniConn                     bool
	useBatchSend                bool
	MonitoringEndpointGenerator MonitoringEndpointGenerator
}

type telemetryEndpoint struct {
	utils.StartStopOnce
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
		endpoints:    nil,
		ks:           csaKeyStore,
		lggr:         lggr.Named("TelemetryManager"),
		logging:      cfg.Logging(),
		maxBatchSize: cfg.MaxBatchSize(),
		sendInterval: cfg.SendInterval(),
		sendTimeout:  cfg.SendTimeout(),
		uniConn:      cfg.UniConn(),
		useBatchSend: cfg.UseBatchSend(),
	}
	for _, e := range cfg.Endpoints() {
		if err := m.addEndpoint(e); err != nil {
			m.lggr.Error(err.Error())
		}

	}
	return m
}

func (m *Manager) Start(ctx context.Context) error {
	return m.StartOnce("TelemetryManager", func() error {
		var err error
		for _, e := range m.endpoints {
			err = multierr.Append(err, e.client.Start(ctx))
		}
		return err
	})
}
func (m *Manager) Close() error {
	return m.StopOnce("TelemetryManager", func() error {
		var err error
		for _, e := range m.endpoints {
			err = multierr.Append(err, e.client.Close())
		}
		return err
	})
}

func (m *Manager) Name() string {
	return m.lggr.Name()
}

func (m *Manager) HealthReport() map[string]error {
	hr := make(map[string]error)
	hr[m.lggr.Name()] = m.Healthy()
	for _, e := range m.endpoints {
		name := fmt.Sprintf("%s.%s.%s", m.lggr.Name(), e.Network, e.ChainID)
		hr[name] = e.StartStopOnce.Healthy()
	}
	return hr
}

// GenMonitoringEndpoint creates a new monitoring endpoints based on the existing available endpoints defined in the core config TOML, if no endpoint for the network and chainID exists, a NOOP agent will be used and the telemetry will not be sent
func (m *Manager) GenMonitoringEndpoint(contractID string, telemType synchronization.TelemetryType, network string, chainID string) commontypes.MonitoringEndpoint {

	e, found := m.getEndpoint(network, chainID)

	if !found {
		m.lggr.Warnf("no telemetry endpoint found for network %q chainID %q, telemetry %q for contactID %q will NOT be sent", network, chainID, telemType, contractID)
		return &NoopAgent{}
	}

	if m.useBatchSend {
		return NewIngressAgentBatch(e.client, contractID, telemType, network, chainID)
	}

	return NewIngressAgent(e.client, contractID, telemType, network, chainID)

}

func (m *Manager) addEndpoint(e config.TelemetryIngressEndpoint) error {
	if e.Network() == "" {
		return errors.New("cannot add telemetry endpoint, network cannot be empty")
	}

	if e.ChainID() == "" {
		return errors.New("cannot add telemetry endpoint, chainID cannot be empty")
	}

	if e.URL() == nil {
		return errors.New("cannot add telemetry endpoint, URL cannot be empty")
	}

	if e.ServerPubKey() == "" {
		return errors.New("cannot add telemetry endpoint, ServerPubKey cannot be empty")
	}

	if _, found := m.getEndpoint(e.Network(), e.ChainID()); found {
		return errors.Errorf("cannot add telemetry endpoint for network %q and chainID %q, endpoint already exists", e.Network(), e.ChainID())
	}

	var tClient synchronization.TelemetryService
	if m.useBatchSend {
		tClient = synchronization.NewTelemetryIngressBatchClient(e.URL(), e.ServerPubKey(), m.ks, m.logging, m.lggr, m.bufferSize, m.maxBatchSize, m.sendInterval, m.sendTimeout, m.uniConn)
	} else {
		tClient = synchronization.NewTelemetryIngressClient(e.URL(), e.ServerPubKey(), m.ks, m.logging, m.lggr, m.bufferSize)
	}

	te := telemetryEndpoint{
		Network: strings.ToUpper(e.Network()),
		ChainID: strings.ToUpper(e.ChainID()),
		URL:     e.URL(),
		PubKey:  e.ServerPubKey(),
		client:  tClient,
	}

	m.endpoints = append(m.endpoints, &te)
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
