package wsrpc

import (
	"context"
	"crypto/x509"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
	grpc_connectivity "google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"

	// "google.golang.org/grpc/metadata"

	"github.com/smartcontractkit/wsrpc"
	"github.com/smartcontractkit/wsrpc/connectivity"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/cache"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb_grpc"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// MaxConsecutiveRequestFailures controls how many consecutive requests are
// allowed to time out before we reset the connection
const MaxConsecutiveRequestFailures = 10

var (
	timeoutCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mercury_transmit_timeout_count",
		Help: "Running count of transmit timeouts",
	},
		[]string{"serverURL"},
	)
	dialCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mercury_dial_count",
		Help: "Running count of dials to mercury server",
	},
		[]string{"serverURL"},
	)
	dialSuccessCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mercury_dial_success_count",
		Help: "Running count of successful dials to mercury server",
	},
		[]string{"serverURL"},
	)
	dialErrorCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mercury_dial_error_count",
		Help: "Running count of errored dials to mercury server",
	},
		[]string{"serverURL"},
	)
	connectionResetCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mercury_connection_reset_count",
		Help: fmt.Sprintf("Running count of times connection to mercury server has been reset (connection reset happens automatically after %d consecutive request failures)", MaxConsecutiveRequestFailures),
	},
		[]string{"serverURL"},
	)
)

// TODO: Not useful here, and should be moved to the caller
type Client interface {
	services.Service
	pb.MercuryClient
	ServerURL() string
	RawClient() pb.MercuryClient
}

// TODO: Abstraction has already leaked. Worth deleting unless mocks?
type Conn interface {
	WaitForReady(ctx context.Context) bool
	GetState() connectivity.State
	Close()
}

// baseClient is a common base for both GRPC and WSRPC clients
type baseClient struct {
	services.StateMachine

	csaKey       csakey.KeyV2
	serverPubKey []byte
	serverURL    string

	logger logger.Logger
	conn   Conn

	consecutiveTimeoutCnt atomic.Int32
	wg                    sync.WaitGroup
	chStop                services.StopChan
	chResetTransport      chan struct{}

	timeoutCountMetric         prometheus.Counter
	dialCountMetric            prometheus.Counter
	dialSuccessCountMetric     prometheus.Counter
	dialErrorCountMetric       prometheus.Counter
	connectionResetCountMetric prometheus.Counter
}

type GrpcClient struct {
	baseClient
	rawClient   pb_grpc.MercuryClient
	conn        *grpc.ClientConn
	cache       cache.GrpcFetcher
	cacheSet    cache.GrpcCacheSet
	tlsCertFile *string
}

type WsrpcClient struct {
	baseClient
	rawClient pb.MercuryClient
	conn      *wsrpc.ClientConn
	cache     cache.WsrpcFetcher
	cacheSet  cache.WsrpcCacheSet
}

// Consumers of wsrpc package should not usually call NewWSRPCClient directly, but instead use the Pool
func NewWSRPCClient(lggr logger.Logger, clientPrivKey csakey.KeyV2, serverPubKey []byte, serverURL string, cacheSet cache.WsrpcCacheSet) *WsrpcClient {
	return &WsrpcClient{
		baseClient: baseClient{
			csaKey:                     clientPrivKey,
			serverPubKey:               serverPubKey,
			serverURL:                  serverURL,
			logger:                     lggr.Named("WSRPC").With("mercuryServerURL", serverURL),
			chResetTransport:           make(chan struct{}, 1),
			chStop:                     make(services.StopChan),
			timeoutCountMetric:         timeoutCount.WithLabelValues(serverURL),
			dialCountMetric:            dialCount.WithLabelValues(serverURL),
			dialSuccessCountMetric:     dialSuccessCount.WithLabelValues(serverURL),
			dialErrorCountMetric:       dialErrorCount.WithLabelValues(serverURL),
			connectionResetCountMetric: connectionResetCount.WithLabelValues(serverURL),
		},
		cacheSet: cacheSet,
	}
}

func NewGRPCClient(lggr logger.Logger, clientPrivKey csakey.KeyV2, serverPubKey []byte, serverURL string, cacheSet cache.GrpcCacheSet, tlsCertFile *string) *GrpcClient {
	return &GrpcClient{
		baseClient: baseClient{
			csaKey:                     clientPrivKey,
			serverPubKey:               serverPubKey,
			serverURL:                  serverURL,
			logger:                     lggr.Named("WSRPC").With("mercuryServerURL", serverURL),
			chResetTransport:           make(chan struct{}, 1),
			chStop:                     make(services.StopChan),
			timeoutCountMetric:         timeoutCount.WithLabelValues(serverURL),
			dialCountMetric:            dialCount.WithLabelValues(serverURL),
			dialSuccessCountMetric:     dialSuccessCount.WithLabelValues(serverURL),
			dialErrorCountMetric:       dialErrorCount.WithLabelValues(serverURL),
			connectionResetCountMetric: connectionResetCount.WithLabelValues(serverURL),
		},
		cacheSet: cacheSet,
		tlsCertFile: tlsCertFile,
	}
}

// TODO: cleanup. Defining runLoop on baseClient doesn't work b/c
// resetTransport has clientConn specific details
// func (c *baseClient) runloop() {
// 	defer c.wg.Done()
// 	for {
// 		select {
// 		case <-c.chStop:
// 			return
// 		case <-c.chResetTransport:
// 			// Using channel here ensures we only have one reset in process at
// 			// any given time
// 			c.resetTransport()
// 		}
// 	}
// }

func (c *GrpcClient) Start(ctx context.Context) error {
	return c.StartOnce("GRPC Client", func() (err error) {
		// NOTE: This is not a mistake, dial is non-blocking so it should use a
		// background context, not the Start context
		if err = c.dial(context.Background()); err != nil {
			return err
		}
		c.cache, err = c.cacheSet.Get(ctx, c)
		if err != nil {
			return err
		}
		c.wg.Add(1)
		go c.runloop()
		return nil
	})
}

// NOTE: Dial is non-blocking, and will retry on an exponential backoff
// in the background until close is called, or context is cancelled.
// This is why we use the background context, not the start context here.
//
// Any transmits made while client is still trying to dial will fail
// with error.
func (c *GrpcClient) dial(ctx context.Context, opts ...grpc.DialOption) error {

	// TODO: move this block to TOML config validation
	b, err := os.ReadFile(*c.tlsCertFile)
	if err != nil {
		return err
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		return fmt.Errorf("credentials: failed to append certificates")
	}

	c.dialCountMetric.Inc()

	conn, err := grpc.DialContext(ctx, c.serverURL,
		append(opts,
			grpc.WithTransportCredentials(
				credentials.NewClientTLSFromCert(cp, ""),
			),
		)...,
	)

	if err != nil {
		c.dialErrorCountMetric.Inc()
		setLivenessMetric(false)
		return errors.Wrap(err, "failed to dial wsrpc client")
	}
	c.dialSuccessCountMetric.Inc()
	setLivenessMetric(true)
	c.conn = conn
	c.rawClient = pb_grpc.NewMercuryClient(conn)
	return nil
}

func (c *GrpcClient) runloop() {
	defer c.wg.Done()
	for {
		select {
		case <-c.chStop:
			return
		case <-c.chResetTransport:
			// Using channel here ensures we only have one reset in process at
			// any given time
			c.resetTransport()
		}
	}
}

// resetTransport disconnects and reconnects to the mercury server
func (c *GrpcClient) resetTransport() {
	c.connectionResetCountMetric.Inc()
	ok := c.IfStarted(func() {
		c.conn.Close() // Close is safe to call multiple times
	})
	if !ok {
		panic("resetTransport should never be called unless client is in 'started' state")
	}
	ctx, cancel := c.chStop.Ctx(context.Background())
	defer cancel()
	b := utils.NewRedialBackoff()
	for {
		// Will block until successful dial, or context is canceled (i.e. on close)
		err := c.dial(ctx, grpc.WithBlock())
		if err == nil {
			break
		}
		if ctx.Err() != nil {
			c.logger.Debugw("ResetTransport exiting due to client Close", "err", err)
			return
		}
		c.logger.Errorw("ResetTransport failed to redial", "err", err)
		time.Sleep(b.Duration())
	}
	c.logger.Info("ResetTransport successfully redialled")
}

func (c *GrpcClient) Close() error {
	return c.StopOnce("GRPC Client", func() error {
		close(c.chStop)
		c.conn.Close()
		c.wg.Wait()
		return nil
	})
}

func (c *GrpcClient) Name() string {
	// TODO: can EVM and Mercury be inferred from up the stack?
	return "EVM.Mercury.GRPCClient"
}

func (c *GrpcClient) HealthReport() map[string]error {
	return map[string]error{c.Name(): c.Healthy()}
}

// Healthy if connected
func (c *GrpcClient) Healthy() (err error) {
	if err = c.StateMachine.Healthy(); err != nil {
		return err
	}
	state := c.conn.GetState()
	// TODO: do WSRPC connectivity states map to GRPC states?
	if state != grpc_connectivity.Ready {
		return errors.Errorf("client state should be %s; got %s", connectivity.Ready, state)
	}
	return nil
}

func (c *GrpcClient) waitForReady(ctx context.Context) (err error) {
	ok := c.IfStarted(func() {
		if ready := WaitForReady(ctx, c); !ready {
			err = errors.Errorf("websocket client not ready; got state: %v", c.conn.GetState())
			return
		}
	})
	if !ok {
		return errors.New("client is not started")
	}
	return
}

// WaitForReady clones wsrpc.WaitForReady
// Blocks on context and waits until the grpc Client Conn state becomes Ready
// It returns true when that happens
// It returns false if the context is cancelled, or the conn is shut down
func WaitForReady(ctx context.Context, c *GrpcClient) bool {
	curState := c.conn.GetState()
	switch curState {
	case grpc_connectivity.Ready:
		return true
	case grpc_connectivity.Shutdown:
		return false
	case grpc_connectivity.Idle, grpc_connectivity.Connecting, grpc_connectivity.TransientFailure:
		break
	}

	c.baseClient.logger.Debugf("Waiting for connection to be ready, current state: %s", curState)

	if ok := c.conn.WaitForStateChange(ctx, curState); !ok {
		return false
	}
	return WaitForReady(ctx, c)
}

func (c *GrpcClient) Transmit(ctx context.Context, req *pb_grpc.TransmitRequest) (resp *pb_grpc.TransmitResponse, err error) {
	c.logger.Trace("Transmit")
	start := time.Now()
	if err = c.waitForReady(ctx); err != nil {
		return nil, errors.Wrap(err, "Transmit call failed")
	}

	// md := metadata.Pairs(
	// 	"key", string(pubKeyStr), // this binary data will be encoded (base64) before sending
	// 							// and will be decoded after being transferred.
	// )
	// ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err = c.rawClient.Transmit(ctx, req)
	c.handleTimeout(err)
	if err != nil {
		c.logger.Warnw("Transmit call failed due to networking error", "err", err, "resp", resp)
		incRequestStatusMetric(statusFailed)
	} else {
		c.logger.Tracew("Transmit call succeeded", "resp", resp)
		incRequestStatusMetric(statusSuccess)
		setRequestLatencyMetric(float64(time.Since(start).Milliseconds()))
	}
	return
}

func (c *GrpcClient) handleTimeout(err error) {
	if errors.Is(err, context.DeadlineExceeded) {
		c.timeoutCountMetric.Inc()
		cnt := c.consecutiveTimeoutCnt.Add(1)
		if cnt == MaxConsecutiveRequestFailures {
			c.logger.Errorf("Timed out on %d consecutive transmits, resetting transport", cnt)
			// NOTE: If we get at least MaxConsecutiveRequestFailures request
			// timeouts in a row, close and re-open the websocket connection.
			//
			// This *shouldn't* be necessary in theory (ideally, wsrpc would
			// handle it for us) but it acts as a "belts and braces" approach
			// to ensure we get a websocket connection back up and running
			// again if it gets itself into a bad state.
			select {
			case c.chResetTransport <- struct{}{}:
			default:
				// This can happen if we had MaxConsecutiveRequestFailures
				// consecutive timeouts, already sent a reset signal, then the
				// connection started working again (resetting the count) then
				// we got MaxConsecutiveRequestFailures additional failures
				// before the runloop was able to close the bad connection.
				//
				// It should be safe to just ignore in this case.
				//
				// Debug log in case my reasoning is wrong.
				c.logger.Debugf("Transport is resetting, cnt=%d", cnt)
			}
		}
	} else {
		c.consecutiveTimeoutCnt.Store(0)
	}
}

func (c *GrpcClient) LatestReport(ctx context.Context, req *pb_grpc.LatestReportRequest) (resp *pb_grpc.LatestReportResponse, err error) {
	lggr := c.logger.With("req.FeedId", hexutil.Encode(req.FeedId))
	lggr.Trace("LatestReport")
	if err = c.waitForReady(ctx); err != nil {
		return nil, errors.Wrap(err, "LatestReport failed")
	}
	var cached bool
	if c.cache == nil {
		resp, err = c.rawClient.LatestReport(ctx, req)
		c.handleTimeout(err)
	} else {
		cached = true
		resp, err = c.cache.LatestReport(ctx, req)
	}
	if err != nil {
		lggr.Errorw("LatestReport failed", "err", err, "resp", resp, "cached", cached)
	} else if resp.Error != "" {
		lggr.Errorw("LatestReport failed; mercury server returned error", "err", resp.Error, "resp", resp, "cached", cached)
	} else if !cached {
		lggr.Debugw("LatestReport succeeded", "resp", resp, "cached", cached)
	} else {
		lggr.Tracew("LatestReport succeeded", "resp", resp, "cached", cached)
	}
	return
}

func (c *GrpcClient) ServerURL() string {
	return c.serverURL
}

func (c *GrpcClient) RawClient() pb_grpc.MercuryClient {
	return c.rawClient
}

func (w *WsrpcClient) Start(ctx context.Context) error {
	return w.StartOnce("WSRPC Client", func() (err error) {
		// NOTE: This is not a mistake, dial is non-blocking so it should use a
		// background context, not the Start context
		if err = w.dial(context.Background()); err != nil {
			return err
		}
		w.cache, err = w.cacheSet.Get(ctx, w)
		if err != nil {
			return err
		}
		w.wg.Add(1)
		go w.runloop()
		return nil
	})
}

// NOTE: Dial is non-blocking, and will retry on an exponential backoff
// in the background until close is called, or context is cancelled.
// This is why we use the background context, not the start context here.
//
// Any transmits made while client is still trying to dial will fail
// with error.
func (w *WsrpcClient) dial(ctx context.Context, opts ...wsrpc.DialOption) error {
	w.dialCountMetric.Inc()
	conn, err := wsrpc.DialWithContext(ctx, w.serverURL,
		append(opts,
			wsrpc.WithTransportCreds(w.csaKey.Raw().Bytes(), w.serverPubKey),
			wsrpc.WithLogger(w.logger),
		)...,
	)
	if err != nil {
		w.dialErrorCountMetric.Inc()
		setLivenessMetric(false)
		return errors.Wrap(err, "failed to dial wsrpc client")
	}
	w.dialSuccessCountMetric.Inc()
	setLivenessMetric(true)
	w.conn = conn
	w.rawClient = pb.NewMercuryClient(conn)
	return nil
}

func (w *WsrpcClient) runloop() {
	defer w.wg.Done()
	for {
		select {
		case <-w.chStop:
			return
		case <-w.chResetTransport:
			// Using channel here ensures we only have one reset in process at
			// any given time
			w.resetTransport()
		}
	}
}

// resetTransport disconnects and reconnects to the mercury server
func (w *WsrpcClient) resetTransport() {
	w.connectionResetCountMetric.Inc()
	ok := w.IfStarted(func() {
		w.conn.Close() // Close is safe to call multiple times
	})
	if !ok {
		panic("resetTransport should never be called unless client is in 'started' state")
	}
	ctx, cancel := w.chStop.Ctx(context.Background())
	defer cancel()
	b := utils.NewRedialBackoff()
	for {
		// Will block until successful dial, or context is canceled (i.e. on close)
		err := w.dial(ctx, wsrpc.WithBlock())
		if err == nil {
			break
		}
		if ctx.Err() != nil {
			w.logger.Debugw("ResetTransport exiting due to client Close", "err", err)
			return
		}
		w.logger.Errorw("ResetTransport failed to redial", "err", err)
		time.Sleep(b.Duration())
	}
	w.logger.Info("ResetTransport successfully redialled")
}

func (w *WsrpcClient) Close() error {
	return w.StopOnce("WSRPC Client", func() error {
		close(w.chStop)
		w.conn.Close()
		w.wg.Wait()
		return nil
	})
}

func (w *WsrpcClient) Name() string {
	return "EVM.Mercury.WSRPCClient"
}

func (w *WsrpcClient) HealthReport() map[string]error {
	return map[string]error{w.Name(): w.Healthy()}
}

// Healthy if connected
func (w *WsrpcClient) Healthy() (err error) {
	if err = w.StateMachine.Healthy(); err != nil {
		return err
	}
	state := w.conn.GetState()
	if state != connectivity.Ready {
		return errors.Errorf("client state should be %s; got %s", connectivity.Ready, state)
	}
	return nil
}

func (w *WsrpcClient) waitForReady(ctx context.Context) (err error) {
	ok := w.IfStarted(func() {
		if ready := w.conn.WaitForReady(ctx); !ready {
			err = errors.Errorf("websocket client not ready; got state: %v", w.conn.GetState())
			return
		}
	})
	if !ok {
		return errors.New("client is not started")
	}
	return
}

func (w *WsrpcClient) Transmit(ctx context.Context, req *pb.TransmitRequest) (resp *pb.TransmitResponse, err error) {
	w.logger.Trace("Transmit")
	start := time.Now()
	if err = w.waitForReady(ctx); err != nil {
		return nil, errors.Wrap(err, "Transmit call failed")
	}
	resp, err = w.rawClient.Transmit(ctx, req)
	w.handleTimeout(err)
	if err != nil {
		w.logger.Warnw("Transmit call failed due to networking error", "err", err, "resp", resp)
		incRequestStatusMetric(statusFailed)
	} else {
		w.logger.Tracew("Transmit call succeeded", "resp", resp)
		incRequestStatusMetric(statusSuccess)
		setRequestLatencyMetric(float64(time.Since(start).Milliseconds()))
	}
	return
}

func (w *WsrpcClient) handleTimeout(err error) {
	if errors.Is(err, context.DeadlineExceeded) {
		w.timeoutCountMetric.Inc()
		cnt := w.consecutiveTimeoutCnt.Add(1)
		if cnt == MaxConsecutiveRequestFailures {
			w.logger.Errorf("Timed out on %d consecutive transmits, resetting transport", cnt)
			// NOTE: If we get at least MaxConsecutiveRequestFailures request
			// timeouts in a row, close and re-open the websocket connection.
			//
			// This *shouldn't* be necessary in theory (ideally, wsrpc would
			// handle it for us) but it acts as a "belts and braces" approach
			// to ensure we get a websocket connection back up and running
			// again if it gets itself into a bad state.
			select {
			case w.chResetTransport <- struct{}{}:
			default:
				// This can happen if we had MaxConsecutiveRequestFailures
				// consecutive timeouts, already sent a reset signal, then the
				// connection started working again (resetting the count) then
				// we got MaxConsecutiveRequestFailures additional failures
				// before the runloop was able to close the bad connection.
				//
				// It should be safe to just ignore in this case.
				//
				// Debug log in case my reasoning is wrong.
				w.logger.Debugf("Transport is resetting, cnt=%d", cnt)
			}
		}
	} else {
		w.consecutiveTimeoutCnt.Store(0)
	}
}

func (w *WsrpcClient) LatestReport(ctx context.Context, req *pb.LatestReportRequest) (resp *pb.LatestReportResponse, err error) {
	lggr := w.logger.With("req.FeedId", hexutil.Encode(req.FeedId))
	lggr.Trace("LatestReport")
	if err = w.waitForReady(ctx); err != nil {
		return nil, errors.Wrap(err, "LatestReport failed")
	}
	var cached bool
	if w.cache == nil {
		resp, err = w.rawClient.LatestReport(ctx, req)
		w.handleTimeout(err)
	} else {
		cached = true
		resp, err = w.cache.LatestReport(ctx, req)
	}
	if err != nil {
		lggr.Errorw("LatestReport failed", "err", err, "resp", resp, "cached", cached)
	} else if resp.Error != "" {
		lggr.Errorw("LatestReport failed; mercury server returned error", "err", resp.Error, "resp", resp, "cached", cached)
	} else if !cached {
		lggr.Debugw("LatestReport succeeded", "resp", resp, "cached", cached)
	} else {
		lggr.Tracew("LatestReport succeeded", "resp", resp, "cached", cached)
	}
	return
}

func (w *WsrpcClient) ServerURL() string {
	return w.serverURL
}

func (w *WsrpcClient) RawClient() pb.MercuryClient {
	return w.rawClient
}
