package synchronization

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc/connectivity"

	"github.com/smartcontractkit/wsrpc"
	"github.com/smartcontractkit/wsrpc/credentials"
	"github.com/smartcontractkit/wsrpc/examples/simple/keys"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/timeutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
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
	csaKeyStore     keystore.CSA
	csaSigner       *core.Ed25519Signer
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
	workersMutex sync.RWMutex

	useUniConn bool

	healthMonitorCancel context.CancelFunc
}

// NewTelemetryIngressBatchClient returns a client backed by wsrpc that
// can send telemetry to the telemetry ingress server
func NewTelemetryIngressBatchClient(url *url.URL, serverPubKeyHex string, csaKeyStore keystore.CSA, logging bool, lggr logger.Logger, telemBufferSize uint, telemMaxBatchSize uint, telemSendInterval time.Duration, telemSendTimeout time.Duration, useUniconn bool) TelemetryService {
	c := &telemetryIngressBatchClient{
		telemBufferSize:   telemBufferSize,
		telemMaxBatchSize: telemMaxBatchSize,
		telemSendInterval: telemSendInterval,
		telemSendTimeout:  telemSendTimeout,
		url:               url,
		csaKeyStore:       csaKeyStore,
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
	serverPubKey := keys.FromHex(tc.serverPubKeyHex)

	// Initialize a new wsrpc client caller
	// This is used to call RPC methods on the server
	if tc.telemClient == nil { // only preset for tests
		key, err := keystore.GetDefault(ctx, tc.csaKeyStore)
		if err != nil {
			return err
		}
		tc.csaSigner, err = core.NewEd25519Signer(key.ID(), keystore.CSASigner{CSA: tc.csaKeyStore}.Sign)
		if err != nil {
			return err
		}
		if tc.useUniConn {
			tc.eng.Go(func(ctx context.Context) {
				conn, err := func() (*wsrpc.UniClientConn, error) {
					pubs, err := credentials.ValidPublicKeysFromEd25519(serverPubKey)
					if err != nil {
						return nil, err
					}
					tlsConfig, err := credentials.NewClientTLSSigner(tc.csaSigner, pubs)
					if err != nil {
						return nil, err
					}
					conn := wsrpc.NewTLSUniClientConn(tc.eng, tc.url.String(), tlsConfig)
					return conn, conn.Dial(ctx)
				}()
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
			// Spawns a goroutine that will eventually connect. Don't pass ctx, which is cancelled after returning.
			conn, err := wsrpc.Dial(tc.url.String(), wsrpc.WithTransportSigner(tc.csaSigner, serverPubKey), wsrpc.WithLogger(tc.eng))
			if err != nil {
				return fmt.Errorf("could not start TelemIngressBatchClient, Dial returned error: %w", err)
			}
			tc.telemClient = telemPb.NewTelemClient(conn)
			tc.closeFn = func() error { conn.Close(); return nil }
			tc.startHealthMonitoring(ctx, conn)
		}
	}

	return nil
}

// startHealthMonitoring starts a goroutine to monitor the connection state and update other relevant metrics every 5 seconds
func (tc *telemetryIngressBatchClient) startHealthMonitoring(ctx context.Context, conn *wsrpc.ClientConn) {
	_, cancel := context.WithCancel(ctx)
	tc.healthMonitorCancel = cancel

	tc.eng.Go(func(ctx context.Context) {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Check the connection state
				connected := float64(0)
				if conn.GetState() == connectivity.Ready {
					connected = float64(1)
				}
				TelemetryClientConnectionStatus.WithLabelValues(tc.url.String()).Set(connected)
			case <-ctx.Done():
				return
			}
		}
	})
}

// Close disconnects the wsrpc client from the ingress server and waits for all workers to exit
func (tc *telemetryIngressBatchClient) close() (err error) {
	if tc.healthMonitorCancel != nil {
		tc.healthMonitorCancel()
	}
	if (tc.useUniConn && tc.connected.Load()) || !tc.useUniConn {
		err = errors.Join(err, tc.closeFn())
	}
	if tc.csaSigner != nil {
		err = errors.Join(err, tc.csaSigner.Close())
	}
	return
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
			tc.url.String(),
		)
		tc.eng.GoTick(timeutil.NewTicker(func() time.Duration {
			return tc.telemSendInterval
		}), worker.Send)
		tc.workers[workerKey] = worker

		TelemetryClientWorkers.WithLabelValues(tc.url.String(), string(payload.TelemType)).Inc()
	}

	return worker
}
