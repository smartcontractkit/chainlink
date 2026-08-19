package bridgeconn

import (
	"context"
	"encoding/hex"
	"net/url"
	"sync"
	"time"

	"github.com/goccy/go-json"
	"github.com/jonboulle/clockwork"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline/bridgeconn/streamspb"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

// Operational timing and reconnect behavior are hardcoded: this manager has no
// configuration section and EAConn takes no configuration interface.
const (
	subscriptionInterval       = 10 * time.Second
	assetIdleTimeout           = 30 * time.Second
	reconnectBackoffInitial    = time.Second
	reconnectBackoffMultiplier = 2.0
	reconnectBackoffMax        = time.Minute
)

// promEAConnObservationsTotal counts accepted gRPC observation messages per bridge
// and asset pair (the hex-encoded bridgeObservationCacheKey). Rate per second is
// derived externally, e.g. via PromQL rate(bridge_eaconn_observations_total[1m]).
var promEAConnObservationsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "bridge_eaconn_observations_total",
	Help: "Count of gRPC observation messages accepted per bridge and asset pair",
}, []string{"bridgeName", "assetKey"})

// promEAConnTransmitDuration measures the delay between the external adapter receiving
// data from its upstream provider (timestamps.providerDataReceivedUnixMs in the
// observation payload) and this node handling the gRPC message.
var promEAConnTransmitDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "bridge_eaconn_transmit_duration_seconds",
	Help: "Delay from the adapter receiving provider data to the node handling the gRPC observation",
	Buckets: []float64{
		0.001, // 1 ms
		0.002, // 2 ms
		0.005, // 5 ms
		0.010, // 10 ms
		0.015, // 15 ms
		0.020, // 20 ms
		0.030, // 30 ms
		0.040, // 40 ms
		0.050, // 50 ms
		0.075, // 75 ms
		0.100, // 100 ms
		0.150, // 150 ms
		0.200, // 200 ms
		0.500, // 500 ms
	},
}, []string{"bridgeName"})

// promEAConnTransmitSkewTotal counts observations excluded from the transmit-duration
// histogram because now minus providerDataReceivedUnixMs came out negative, meaning the
// adapter's clock runs ahead of this node's. Skew of tens of milliseconds is normal
// between unsynchronized hosts and exceeds the latency the histogram tries to measure,
// so such samples are dropped rather than folded into the lowest bucket. Compare this
// against bridge_eaconn_observations_total to judge how much of a bridge's traffic the
// histogram actually covers: if this tracks the observation rate, the histogram is empty
// for that bridge and its clocks need synchronizing before the latency data means
// anything.
var promEAConnTransmitSkewTotal = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "bridge_eaconn_transmit_skew_total",
	Help: "Count of observations excluded from the transmit-duration histogram due to adapter clock skew",
}, []string{"bridgeName"})

// observationTimestamps is the subset of an observation payload EAConn reads for the
// transmit-duration metric; all other payload fields are ignored.
type observationTimestamps struct {
	Timestamps struct {
		ProviderDataReceivedUnixMs *int64 `json:"providerDataReceivedUnixMs"`
	} `json:"timestamps"`
}

// eaStreamClient is the protobuf-independent contract EAConn depends on for its
// bidirectional observation stream, implemented by grpcStreamClient in production
// and by fakes in tests.
type eaStreamClient interface {
	Send(*streamspb.SubscribeRequest) error
	Recv() (*streamspb.SubscribeResponse, error)
	Close() error
}

// eaStreamDialer opens a new stream to target (host:port derived from the bridge's
// URL authority). useTLS selects TLS (bridge URL scheme https) versus plaintext gRPC
// (http).
type eaStreamDialer func(ctx context.Context, target string, useTLS bool) (eaStreamClient, error)

// eaAsset is a bridge's registered asset: an immutable subscription payload plus the
// timestamp it was last requested through GetObservation.
type eaAsset struct {
	payload  *structpb.Struct
	lastUsed time.Time
}

// eaConn is the single persistent connection owned by BridgeConnManager for one
// bridge name. It sends complete active-asset snapshots on a fixed interval,
// applies indirect unsubscribe by omitting idle assets from the next snapshot, and
// reconnects with fixed exponential backoff on any dial/send/receive failure.
type eaConn struct {
	bridgeName string
	target     string
	useTLS     bool
	dial       eaStreamDialer
	lggr       logger.Logger
	manager    *bridgeConnManager
	clock      clockwork.Clock

	mu     sync.Mutex
	assets map[[32]byte]*eaAsset

	startOnce sync.Once
}

func newEAConn(bridgeName string, bridgeURL models.WebURL, manager *bridgeConnManager) *eaConn {
	u := url.URL(bridgeURL)
	return &eaConn{
		bridgeName: bridgeName,
		target:     u.Host,
		useTLS:     u.Scheme == "https",
		dial:       manager.dial,
		lggr:       logger.With(logger.Named(manager.lggr, "EAConn"), "bridgeName", bridgeName),
		manager:    manager,
		clock:      manager.clock,
		assets:     make(map[[32]byte]*eaAsset),
	}
}

func (c *eaConn) start() {
	c.startOnce.Do(func() {
		go c.run(context.Background())
	})
}

// registerAsset records or refreshes the lease for key, lazily adding it to the next
// snapshot. Asset changes never trigger an immediate send.
func (c *eaConn) registerAsset(key [32]byte, payload *structpb.Struct) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if asset, ok := c.assets[key]; ok {
		asset.lastUsed = c.clock.Now()
		return
	}
	c.assets[key] = &eaAsset{payload: payload, lastUsed: c.clock.Now()}
	c.lggr.Debugw("EAConn: registered asset", "key", hex.EncodeToString(key[:]))
}

// pruneAndSnapshot removes assets idle for at least assetIdleTimeout (indirect
// unsubscribe by omission) and returns a full-snapshot request for what remains.
func (c *eaConn) pruneAndSnapshot() (map[[32]byte]struct{}, *streamspb.SubscribeRequest) {
	now := c.clock.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	for key, asset := range c.assets {
		if now.Sub(asset.lastUsed) >= assetIdleTimeout {
			delete(c.assets, key)
		}
	}
	keys := make(map[[32]byte]struct{}, len(c.assets))
	req := &streamspb.SubscribeRequest{Subscriptions: make([]*streamspb.Subscription, 0, len(c.assets))}
	for key, asset := range c.assets {
		keys[key] = struct{}{}
		req.Subscriptions = append(req.Subscriptions, &streamspb.Subscription{Data: asset.payload})
	}
	return keys, req
}

// run owns the connect/send/receive/reconnect lifecycle for the process lifetime;
// there is no explicit stop, since this manager has no runner-owned lifecycle.
func (c *eaConn) run(ctx context.Context) {
	backoff := reconnectBackoffInitial
	for {
		stream, err := c.dial(ctx, c.target, c.useTLS)
		if err != nil {
			c.lggr.Errorw("EAConn: dial failed", "target", c.target, "err", err)
			backoff = c.sleepBackoff(backoff)
			continue
		}
		c.lggr.Infow("EAConn: stream opened", "target", c.target)
		backoff = reconnectBackoffInitial

		streamErr := c.serve(stream)
		_ = stream.Close()
		c.lggr.Warnw("EAConn: stream closed, reconnecting", "target", c.target, "err", streamErr)
		backoff = c.sleepBackoff(backoff)
	}
}

// serve runs the sender ticker and receiver loop against one open stream until
// either fails, returning the failure.
func (c *eaConn) serve(stream eaStreamClient) error {
	errCh := make(chan error, 2)
	done := make(chan struct{})
	defer close(done)

	go func() {
		errCh <- c.recvLoop(stream, done)
	}()
	go func() {
		errCh <- c.sendLoop(stream, done)
	}()
	return <-errCh
}

func (c *eaConn) sendLoop(stream eaStreamClient, done <-chan struct{}) error {
	ticker := c.clock.NewTicker(subscriptionInterval)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			return nil
		case <-ticker.Chan():
			_, req := c.pruneAndSnapshot()
			if err := stream.Send(req); err != nil {
				return err
			}
		}
	}
}

func (c *eaConn) recvLoop(stream eaStreamClient, done <-chan struct{}) error {
	for {
		resp, err := stream.Recv()
		if err != nil {
			return err
		}
		c.handleObservation(resp)
		select {
		case <-done:
			return nil
		default:
		}
	}
}

// handleObservation accepts an observation only for a key currently registered to
// this EAConn (payload_hash is expected to equal bridgeObservationCacheKey of the
// subscription payload we sent); unregistered or malformed keys are discarded.
func (c *eaConn) handleObservation(resp *streamspb.SubscribeResponse) {
	if len(resp.PayloadHash) != 32 {
		c.lggr.Warnw("EAConn: discarding observation with malformed payload_hash", "len", len(resp.PayloadHash))
		return
	}
	var key [32]byte
	copy(key[:], resp.PayloadHash)

	c.mu.Lock()
	_, registered := c.assets[key]
	c.mu.Unlock()
	if !registered {
		c.lggr.Debugw("EAConn: discarding observation for unregistered key", "payloadHash", hex.EncodeToString(key[:]))
		return
	}
	promEAConnObservationsTotal.WithLabelValues(c.bridgeName, hex.EncodeToString(key[:])).Inc()
	c.observeTransmitDuration(resp.ObservationJson)
	c.manager.PutObservation(key, resp.ObservationJson)
}

// observeTransmitDuration records now minus providerDataReceivedUnixMs. Payloads that
// omit the timestamp are skipped silently; negative durations are counted in
// promEAConnTransmitSkewTotal instead, since they reflect adapter clock skew rather than
// transmission time. Both branches are per-message hot paths at a few thousand
// observations per second, so neither logs: the counters carry the signal.
func (c *eaConn) observeTransmitDuration(observationJSON []byte) {
	var obs observationTimestamps
	if err := json.Unmarshal(observationJSON, &obs); err != nil {
		return
	}
	if obs.Timestamps.ProviderDataReceivedUnixMs == nil {
		return
	}
	received := time.UnixMilli(*obs.Timestamps.ProviderDataReceivedUnixMs)
	d := c.clock.Now().Sub(received)
	if d < 0 {
		promEAConnTransmitSkewTotal.WithLabelValues(c.bridgeName).Inc()
		return
	}
	promEAConnTransmitDuration.WithLabelValues(c.bridgeName).Observe(d.Seconds())
}

func (c *eaConn) sleepBackoff(current time.Duration) time.Duration {
	c.clock.Sleep(current)
	return nextBackoff(current)
}

func nextBackoff(current time.Duration) time.Duration {
	next := time.Duration(float64(current) * reconnectBackoffMultiplier)
	return min(next, reconnectBackoffMax)
}
