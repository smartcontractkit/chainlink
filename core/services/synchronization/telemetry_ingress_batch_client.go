package synchronization

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/smartcontractkit/wsrpc"
	"github.com/smartcontractkit/wsrpc/examples/simple/keys"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/timeutil"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	telemPb "github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
)

// NoopTelemetryIngressBatchClient is a no-op interface for TelemetryIngressBatchClient
type NoopTelemetryIngressBatchClient struct{}

// Start is a no-op
func (NoopTelemetryIngressBatchClient) Start(context.Context) error { return nil }

// Close is a no-op
func (NoopTelemetryIngressBatchClient) Close() error { return nil }

// Send is a no-op
func (NoopTelemetryIngressBatchClient) Send(TelemPayload) {}

func (NoopTelemetryIngressBatchClient) HealthReport() map[string]error { return map[string]error{} }
func (NoopTelemetryIngressBatchClient) Name() string                   { return "NoopTelemetryIngressBatchClient" }

// Ready is a no-op
func (NoopTelemetryIngressBatchClient) Ready() error { return nil }

type telemetryIngressBatchClient struct {
	services.Service
	eng *services.Engine

	url             *url.URL
	ks              keystore.CSA
	serverPubKeyHex string

	connected   atomic.Bool
	telemClient telemPb.TelemClient
	closeFn     func() error

	logging bool

	telemBufferSize   uint
	telemMaxBatchSize uint
	telemSendInterval time.Duration
	telemSendTimeout  time.Duration

	workers      map[string]*telemetryIngressBatchWorker
	workersMutex sync.Mutex

	useUniConn bool
}

// NewTelemetryIngressBatchClient returns a client backed by wsrpc that
// can send telemetry to the telemetry ingress server
func NewTelemetryIngressBatchClient(url *url.URL, serverPubKeyHex string, ks keystore.CSA, logging bool, lggr logger.Logger, telemBufferSize uint, telemMaxBatchSize uint, telemSendInterval time.Duration, telemSendTimeout time.Duration, useUniconn bool) TelemetryService {
	c := &telemetryIngressBatchClient{
		telemBufferSize:   telemBufferSize,
		telemMaxBatchSize: telemMaxBatchSize,
		telemSendInterval: telemSendInterval,
		telemSendTimeout:  telemSendTimeout,
		url:               url,
		ks:                ks,
		serverPubKeyHex:   serverPubKeyHex,
		logging:           logging,
		workers:           make(map[string]*telemetryIngressBatchWorker),
		useUniConn:        useUniconn,
	}
	c.Service, c.eng = services.Config{
		Name:  "TelemetryIngressBatchClient",
		Start: c.start,
		Close: c.close,
	}.NewServiceEngine(lggr)

	return c
}

// Start connects the wsrpc client to the telemetry ingress server
//
// If a connection cannot be established with the ingress server, Dial will return without
// an error and wsrpc will continue to retry the connection. Eventually when the ingress
// server does come back up, wsrpc will establish the connection without any interaction
// on behalf of the node operator.
func (tc *telemetryIngressBatchClient) start(ctx context.Context) error {
	clientPrivKey, err := tc.getCSAPrivateKey()
	if err != nil {
		return err
	}

	serverPubKey := keys.FromHex(tc.serverPubKeyHex)

	// Initialize a new wsrpc client caller
	// This is used to call RPC methods on the server
	if tc.telemClient == nil { // only preset for tests
		if tc.useUniConn {
			tc.eng.Go(func(ctx context.Context) {
				conn, err := wsrpc.DialUniWithContext(ctx, tc.eng, tc.url.String(), clientPrivKey, serverPubKey)
				if err != nil {
					if ctx.Err() != nil {
						tc.eng.Warnw("gave up connecting to telemetry endpoint", "err", err)
					} else {
						tc.eng.Criticalw("telemetry endpoint dial errored unexpectedly", "err", err, "server pubkey", tc.serverPubKeyHex)
						tc.eng.EmitHealthErr(err)
					}
					return
				}
				tc.telemClient = telemPb.NewTelemClient(conn)
				tc.closeFn = conn.Close
				tc.connected.Store(true)
			})
		} else {
			// Spawns a goroutine that will eventually connect
			conn, err := wsrpc.DialWithContext(ctx, tc.url.String(), wsrpc.WithTransportCreds(clientPrivKey, serverPubKey), wsrpc.WithLogger(tc.eng))
			if err != nil {
				return fmt.Errorf("could not start TelemIngressBatchClient, Dial returned error: %v", err)
			}
			tc.telemClient = telemPb.NewTelemClient(conn)
			tc.closeFn = func() error { conn.Close(); return nil }
		}
	}

	return nil
}

// Close disconnects the wsrpc client from the ingress server and waits for all workers to exit
func (tc *telemetryIngressBatchClient) close() error {
	if (tc.useUniConn && tc.connected.Load()) || !tc.useUniConn {
		return tc.closeFn()
	}
	return nil
}

// getCSAPrivateKey gets the client's CSA private key
func (tc *telemetryIngressBatchClient) getCSAPrivateKey() (privkey []byte, err error) {
	keys, err := tc.ks.GetAll()
	if err != nil {
		return privkey, err
	}
	if len(keys) < 1 {
		return privkey, errors.New("CSA key does not exist")
	}

	return keys[0].Raw(), nil
}

// Send directs incoming telmetry messages to the worker responsible for pushing it to
// the ingress server. If the worker telemetry buffer is full, messages are dropped
// and a warning is logged.
func (tc *telemetryIngressBatchClient) Send(ctx context.Context, telemData []byte, contractID string, telemType TelemetryType) {
	if tc.useUniConn && !tc.connected.Load() {
		tc.eng.Warnw("not connected to telemetry endpoint", "endpoint", tc.url.String())
		return
	}
	payload := TelemPayload{
		Telemetry:  telemData,
		TelemType:  telemType,
		ContractID: contractID,
	}
	worker := tc.findOrCreateWorker(payload)

	select {
	case worker.chTelemetry <- payload:
		worker.dropMessageCount.Store(0)
	case <-ctx.Done():
		return
	default:
		worker.logBufferFullWithExpBackoff(payload)
	}
}

// findOrCreateWorker finds a worker by ContractID or creates a new one if none exists
func (tc *telemetryIngressBatchClient) findOrCreateWorker(payload TelemPayload) *telemetryIngressBatchWorker {
	tc.workersMutex.Lock()
	defer tc.workersMutex.Unlock()

	workerKey := fmt.Sprintf("%s_%s", payload.ContractID, payload.TelemType)
	worker, found := tc.workers[workerKey]

	if !found {
		worker = NewTelemetryIngressBatchWorker(
			tc.telemMaxBatchSize,
			tc.telemSendTimeout,
			tc.telemClient,
			make(chan TelemPayload, tc.telemBufferSize),
			payload.ContractID,
			payload.TelemType,
			tc.eng,
			tc.logging,
		)
		tc.eng.GoTick(timeutil.NewTicker(func() time.Duration {
			return tc.telemSendInterval
		}), worker.Send)
		tc.workers[workerKey] = worker
	}

	return worker
}
