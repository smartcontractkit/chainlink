package chainlink

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"go.opentelemetry.io/otel/metric"

	commonservices "github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/timeutil"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/static"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
)

type Heartbeat struct {
	commonservices.Service
	eng *commonservices.Engine

	opts    HeartbeatConfig
	beat    time.Duration
	emitter custmsg.MessageEmitter // Add emitter field
	meter   metric.Meter           // Add meter field
}

type HeartbeatConfig struct {
	Beat         time.Duration
	Lggr         logger.Logger
	P2P          string
	AppID        string
	CSAPublicKey string
	EVMNodeURLs  map[string]string // Maps WSURL_chainID/HTTPURL_chainID to comma-separated domains
}

func NewHeartbeatConfig(cfg ApplicationOpts) HeartbeatConfig {
	csaKey := ""
	csaKeys, err := cfg.KeyStore.CSA().GetAll()
	if err != nil {
		cfg.Logger.Errorw("failed to get CSA keys", "err", err)
	}

	if len(csaKeys) > 0 {
		csaKey = csaKeys[0].PublicKeyString()
	} else {
		cfg.Logger.Warn("no CSA key found for heartbeat")
	}

	evmNodeURLs := make(map[string]string)
	if cfg.Config.EVMEnabled() {
		for _, evmCfg := range cfg.Config.EVMConfigs() {
			if evmCfg.ChainID == nil || !evmCfg.IsEnabled() {
				continue
			}
			chainID := evmCfg.ChainID.String()
			var wsDomains, httpDomains []string
			for _, node := range evmCfg.Nodes {
				if node.WSURL != nil {
					if domain := ExtractDomainFromURLString(node.WSURL.String()); domain != "" {
						wsDomains = append(wsDomains, domain)
					}
				}
				if node.HTTPURL != nil {
					if domain := ExtractDomainFromURLString(node.HTTPURL.String()); domain != "" {
						httpDomains = append(httpDomains, domain)
					}
				}
			}
			if len(wsDomains) > 0 {
				evmNodeURLs["WSURL_"+chainID] = strings.Join(wsDomains, ",")
			}
			if len(httpDomains) > 0 {
				evmNodeURLs["HTTPURL_"+chainID] = strings.Join(httpDomains, ",")
			}
		}
	}

	return HeartbeatConfig{
		Beat:         cfg.Config.Telemetry().HeartbeatInterval(),
		Lggr:         cfg.Logger,
		P2P:          cfg.Config.P2P().PeerID().String(),
		AppID:        cfg.Config.AppID().String(),
		CSAPublicKey: csaKey,
		EVMNodeURLs:  evmNodeURLs,
	}
}

// Update the constructor to accept optional emitter and meter
func NewHeartbeat(cfg HeartbeatConfig, opts ...HeartbeatOpt) Heartbeat {
	// setup default emitter and meter
	cme := custmsg.NewLabeler()
	labels := map[string]string{
		"system":      "Application",
		"version":     static.Version,
		"commit":      static.Sha,
		"version_tag": static.VersionTag,
	}
	if cfg.P2P != "" {
		labels["peer_id"] = cfg.P2P
	}
	if cfg.AppID != "" {
		labels["appID"] = cfg.AppID
	}
	if cfg.CSAPublicKey != "" {
		labels["csa_key"] = cfg.CSAPublicKey
	}
	for k, v := range cfg.EVMNodeURLs {
		labels[k] = v
	}

	cme.WithMapLabels(labels)
	h := Heartbeat{
		beat:    cfg.Beat,
		opts:    cfg,
		emitter: cme.WithMapLabels(labels),
		meter:   beholder.GetMeter(),
	}

	// Apply test options if any
	for _, opt := range opts {
		opt(&h)
	}

	h.Service, h.eng = commonservices.Config{
		Name:  "Heartbeat",
		Start: h.start,
	}.NewServiceEngine(cfg.Lggr)
	return h
}

// Define options for testing
type HeartbeatOpt func(*Heartbeat)

func WithEmitter(emitter custmsg.MessageEmitter) HeartbeatOpt {
	return func(h *Heartbeat) {
		h.emitter = emitter
	}
}

func WithMeter(meter metric.Meter) HeartbeatOpt {
	return func(h *Heartbeat) {
		h.meter = meter
	}
}

// Update the start method to use the injected emitter if provided
func (h *Heartbeat) start(_ context.Context) error {
	// Setup the heartbeat gauge and count

	gauge, err := h.meter.Int64Gauge("heartbeat")
	if err != nil {
		return fmt.Errorf("failed to create heartbeat gauge: %w", err)
	}
	count, err := h.meter.Int64Gauge("heartbeat_count")
	if err != nil {
		return fmt.Errorf("failed to create heartbeat count gauge: %w", err)
	}

	// Define tick functions
	beatFn := func(ctx context.Context) {
		// TODO allow override of tracer provider into engine for beholder
		_, innerSpan := beholder.GetTracer().Start(ctx, "heartbeat.beat")
		defer innerSpan.End()

		gauge.Record(ctx, 1)
		count.Record(ctx, 1)

		err = h.emitter.Emit(ctx, "heartbeat")
		if err != nil {
			h.eng.Errorw("heartbeat emit failed", "err", err)
		}
	}

	h.eng.GoTick(timeutil.NewTicker(h.GetBeat), beatFn)
	return nil
}

func (h *Heartbeat) GetBeat() time.Duration {
	return h.beat
}

// ExtractDomainFromURLString parses a URL string and returns only the hostname (domain).
// This ensures sensitive information like API keys in URL paths are not exposed.
// The port is stripped to avoid formatting issues.
func ExtractDomainFromURLString(rawURL string) string {
	if rawURL == "" {
		return ""
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return u.Hostname()
}
