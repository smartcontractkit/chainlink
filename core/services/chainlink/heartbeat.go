package chainlink

import (
	"context"
	"fmt"
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

	lggr    logger.Logger
	opts    HeartbeatConfig
	beat    time.Duration
	emitter custmsg.MessageEmitter // Add emitter field
	meter   metric.Meter           // Add meter field
}

type HeartbeatConfig struct {
	Beat  time.Duration
	Lggr  logger.Logger
	P2P   string
	AppID string
}

func NewHeartbeatConfig(cfg ApplicationOpts) HeartbeatConfig {
	return HeartbeatConfig{
		Beat:  cfg.Config.Telemetry().HeartbeatInterval(),
		Lggr:  cfg.Logger,
		P2P:   cfg.Config.P2P().PeerID().String(),
		AppID: cfg.Config.AppID().String(),
	}
}

// Update the constructor to accept optional emitter and meter
func NewHeartbeat2(opts HeartbeatConfig, testOpts ...HeartbeatOpt) Heartbeat {
	lggr := logger.Sugared(opts.Lggr).Named("Heartbeat")
	h := Heartbeat{
		beat:    opts.Beat,
		opts:    opts,
		emitter: nil, // Will be set in start() if nil
		meter:   nil, // Will be set in start() if nil
	}

	// Apply test options if any
	for _, opt := range testOpts {
		opt(&h)
	}

	h.Service, h.eng = commonservices.Config{
		Name:  "Heartbeat",
		Start: h.start,
	}.NewServiceEngine(lggr)
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
	// Setup beholder resources
	var gauge, count metric.Int64Gauge
	var err error

	// Use injected gauger or get default
	gauger := h.meter
	if gauger == nil {
		gauge, err = beholder.GetMeter().Int64Gauge("heartbeat")
		if err != nil {
			return err
		}
		count, err = beholder.GetMeter().Int64Gauge("heartbeat_count")
		if err != nil {
			return err
		}
	} else {
		gauge, err = gauger.Int64Gauge("heartbeat")
		if err != nil {
			return fmt.Errorf("failed to create heartbeat gauge: %w", err)
		}
		count, err = gauger.Int64Gauge("heartbeat_count")
		if err != nil {
			return fmt.Errorf("failed to create heartbeat count gauge: %w", err)
		}
	}

	// Use injected emitter or create a new one
	cme := h.emitter
	if cme == nil {
		cme = custmsg.NewLabeler()
		labels := map[string]string{"system": "Application", "version": static.Version, "commit": static.Sha}
		if h.opts.P2P != "" {
			labels["peer_id"] = h.opts.P2P
		}
		if h.opts.AppID != "" {
			labels["appID"] = h.opts.AppID
		}
		cme.WithMapLabels(labels)
	}

	// Define tick functions
	beatFn := func(ctx context.Context) {
		// TODO allow override of tracer provider into engine for beholder
		_, innerSpan := beholder.GetTracer().Start(ctx, "heartbeat.beat")
		defer innerSpan.End()

		gauge.Record(ctx, 1)
		count.Record(ctx, 1)

		err = cme.Emit(ctx, "heartbeat")
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
