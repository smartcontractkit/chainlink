package wsrpc

import (
	"context"
	"crypto"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	grpc_connectivity "google.golang.org/grpc/connectivity"

	"github.com/smartcontractkit/wsrpc"
	"github.com/smartcontractkit/wsrpc/connectivity"

	"github.com/smartcontractkit/chainlink-data-streams/rpc"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/core/services/llo/grpc"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/cache"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
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

type Client interface {
	services.Service
	pb.MercuryClient
	ServerURL() string
	RawClient() pb.MercuryClient
}

type Conn interface {
	wsrpc.ClientInterface
	WaitForReady(ctx context.Context) bool
	GetState() grpc_connectivity.State
	Close() error
}

type DialWithContextFunc func(ctxCaller context.Context, target string, opts ...wsrpc.DialOption) (Conn, error)

type client struct {
	services.StateMachine

	csaSigner    crypto.Signer
	serverPubKey []byte
	serverURL    string

	dialWithContext DialWithContextFunc

	logger    logger.SugaredLogger
	conn      Conn
	rawClient pb.MercuryClient
	mu        sync.RWMutex

	consecutiveTimeoutCnt atomic.Int32
	wg                    sync.WaitGroup
	chStop                services.StopChan
	chResetTransport      chan struct{}

	cacheSet cache.CacheSet
	cache    cache.Fetcher

	timeoutCountMetric         prometheus.Counter
	dialCountMetric            prometheus.Counter
	dialSuccessCountMetric     prometheus.Counter
	dialErrorCountMetric       prometheus.Counter
	connectionResetCountMetric prometheus.Counter
}

type ClientOpts struct {
	Logger       logger.SugaredLogger
	CSASigner    crypto.Signer
	ServerPubKey []byte
	ServerURL    string
	CacheSet     cache.CacheSet

	// DialWithContext allows optional dependency injection for testing
	DialWithContext DialWithContextFunc
}

// Consumers of wsrpc package should not usually call NewClient directly, but instead use the Pool
func NewClient(opts ClientOpts) Client {
	return newClient(opts)
}

func newClient(opts ClientOpts) *client {
	var dialWithContext DialWithContextFunc
	if opts.DialWithContext != nil {
		dialWithContext = opts.DialWithContext
	} else {
		// NOTE: Wrap here since wsrpc.DialWithContext returns a concrete *wsrpc.Conn, not an interface
		dialWithContext = func(ctx context.Context, target string, opts ...wsrpc.DialOption) (Conn, error) {
			conn, err := wsrpc.DialWithContext(ctx, target, opts...)
			return conn, err
		}
	}
	return &client{
		dialWithContext:            dialWithContext,
		csaSigner:                  opts.CSASigner,
		serverPubKey:               opts.ServerPubKey,
		serverURL:                  opts.ServerURL,
		logger:                     opts.Logger.Named("WSRPC").Named(opts.ServerURL).With("serverURL", opts.ServerURL),
		chResetTransport:           make(chan struct{}, 1),
		cacheSet:                   opts.CacheSet,
		chStop:                     make(services.StopChan),
		timeoutCountMetric:         timeoutCount.WithLabelValues(opts.ServerURL),
		dialCountMetric:            dialCount.WithLabelValues(opts.ServerURL),
		dialSuccessCountMetric:     dialSuccessCount.WithLabelValues(opts.ServerURL),
		dialErrorCountMetric:       dialErrorCount.WithLabelValues(opts.ServerURL),
		connectionResetCountMetric: connectionResetCount.WithLabelValues(opts.ServerURL),
	}
}

func (w *client) Start(ctx context.Context) error {
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
func (w *client) dial(ctx context.Context, opts ...wsrpc.DialOption) error {
	w.dialCountMetric.Inc()
	conn, err := w.dialWithContext(ctx, w.serverURL,
		append(opts,
			wsrpc.WithTransportSigner(w.csaSigner, w.serverPubKey),
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
	w.mu.Lock()
	w.conn = conn
	w.rawClient = pb.NewMercuryClient(conn)
	w.mu.Unlock()
	return nil
}

func (w *client) runloop() {
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
func (w *client) resetTransport() {
	w.connectionResetCountMetric.Inc()
	ok := w.IfStarted(func() {
		w.mu.RLock()
		defer w.mu.RUnlock()
		w.conn.Close() // Close is safe to call multiple times
	})
	if !ok {
		panic("resetTransport should never be called unless client is in 'started' state")
	}
	ctx, cancel := w.chStop.NewCtx()
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

func (w *client) Close() error {
	return w.StopOnce("WSRPC Client", func() error {
		close(w.chStop)
		w.mu.RLock()
		w.conn.Close()
		w.mu.RUnlock()
		w.wg.Wait()
		return nil
	})
}

func (w *client) Name() string {
	return w.logger.Name()
}

func (w *client) HealthReport() map[string]error {
	return map[string]error{w.Name(): w.Healthy()}
}

// Healthy if connected
func (w *client) Healthy() (err error) {
	if err = w.StateMachine.Healthy(); err != nil {
		return err
	}
	state := w.conn.GetState()
	if state != grpc_connectivity.Ready {
		return errors.Errorf("client state should be %s; got %s", connectivity.Ready, state)
	}
	return nil
}

func (w *client) waitForReady(ctx context.Context) (err error) {
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

func (w *client) Transmit(ctx context.Context, req *pb.TransmitRequest) (resp *pb.TransmitResponse, err error) {
	ok := w.IfStarted(func() {
		w.logger.Trace("Transmit")
		start := time.Now()
		if err = w.waitForReady(ctx); err != nil {
			err = errors.Wrap(err, "Transmit call failed")
			return
		}
		w.mu.RLock()
		rc := w.rawClient
		w.mu.RUnlock()
		resp, err = rc.Transmit(ctx, req)
		w.handleTimeout(err)
		if err != nil {
			w.logger.Warnw("Transmit call failed due to networking error", "err", err, "resp", resp)
			incRequestStatusMetric(statusFailed)
		} else {
			w.logger.Tracew("Transmit call succeeded", "resp", resp)
			incRequestStatusMetric(statusSuccess)
			setRequestLatencyMetric(float64(time.Since(start).Milliseconds()))
		}
	})
	if !ok {
		err = errors.New("client is not started")
	}
	return
}

// hacky workaround to trap panics from buggy underlying wsrpc lib and restart
// the connection from a known good state
func (w *client) handlePanic(r interface{}) {
	w.chResetTransport <- struct{}{}
}

func (w *client) handleTimeout(err error) {
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

func (w *client) LatestReport(ctx context.Context, req *pb.LatestReportRequest) (resp *pb.LatestReportResponse, err error) {
	ok := w.IfStarted(func() {
		lggr := w.logger.With("req.FeedId", hexutil.Encode(req.FeedId))
		lggr.Trace("LatestReport")
		if err = w.waitForReady(ctx); err != nil {
			err = errors.Wrap(err, "LatestReport failed")
			return
		}
		var cached bool
		if w.cache == nil {
			w.mu.RLock()
			rc := w.rawClient
			w.mu.RUnlock()
			resp, err = rc.LatestReport(ctx, req)
			w.handleTimeout(err)
		} else {
			cached = true
			resp, err = w.cache.LatestReport(ctx, req)
		}
		switch {
		case err != nil:
			lggr.Errorw("LatestReport failed", "err", err, "resp", resp, "cached", cached)
		case resp.Error != "":
			lggr.Errorw("LatestReport failed; mercury server returned error", "err", resp.Error, "resp", resp, "cached", cached)
		case !cached:
			lggr.Debugw("LatestReport succeeded", "resp", resp, "cached", cached)
		default:
			lggr.Tracew("LatestReport succeeded", "resp", resp, "cached", cached)
		}
	})
	if !ok {
		err = errors.New("client is not started")
	}
	return
}

func (w *client) ServerURL() string {
	return w.serverURL
}

func (w *client) RawClient() pb.MercuryClient {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.rawClient
}

var _ grpc.Client = GRPCCompatibilityWrapper{}

type GRPCCompatibilityWrapper struct {
	Client
}

func (w GRPCCompatibilityWrapper) Transmit(ctx context.Context, in *rpc.TransmitRequest) (*rpc.TransmitResponse, error) {
	req := &pb.TransmitRequest{
		Payload:      in.Payload,
		ReportFormat: in.ReportFormat,
	}
	resp, err := w.Client.Transmit(ctx, req)
	if err != nil {
		return nil, err
	}
	return &rpc.TransmitResponse{
		Code:  resp.Code,
		Error: resp.Error,
	}, nil
}
