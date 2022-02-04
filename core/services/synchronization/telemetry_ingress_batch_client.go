package synchronization

import (
	"errors"
	"net/url"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	telemPb "github.com/smartcontractkit/chainlink/core/services/synchronization/telem"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/wsrpc"
	"github.com/smartcontractkit/wsrpc/examples/simple/keys"
)

//go:generate mockery --dir ./telem --name TelemClient --output ./mocks/ --case=underscore

// TelemetryIngressBatchClient encapsulates all the functionality needed to
// send telemetry to the ingress server using wsrpc
type TelemetryIngressBatchClient interface {
	services.Service
	Send(TelemPayload)
}

// NoopTelemetryIngressBatchClient is a no-op interface for TelemetryIngressBatchClient
type NoopTelemetryIngressBatchClient struct{}

// Start is a no-op
func (NoopTelemetryIngressBatchClient) Start() error { return nil }

// Close is a no-op
func (NoopTelemetryIngressBatchClient) Close() error { return nil }

// Send is a no-op
func (NoopTelemetryIngressBatchClient) Send(TelemPayload) {}

// Healthy is a no-op
func (NoopTelemetryIngressBatchClient) Healthy() error { return nil }

// Ready is a no-op
func (NoopTelemetryIngressBatchClient) Ready() error { return nil }

type telemetryIngressBatchClient struct {
	utils.StartStopOnce
	url             *url.URL
	ks              keystore.CSA
	serverPubKeyHex string

	telemClient  telemPb.TelemClient
	globalLogger logger.Logger
	logging      bool
	lggr         logger.Logger

	wgDone sync.WaitGroup
	chDone chan struct{}

	telemBufferSize   uint
	telemMaxBatchSize uint
	telemSendInterval time.Duration

	workers      map[string]*telemetryIngressBatchWorker
	workersMutex sync.Mutex
}

// NewTelemetryIngressBatchClient returns a client backed by wsrpc that
// can send telemetry to the telemetry ingress server
func NewTelemetryIngressBatchClient(url *url.URL, serverPubKeyHex string, ks keystore.CSA, logging bool, lggr logger.Logger, telemBufferSize uint, telemMaxBatchSize uint, telemSendInterval time.Duration) TelemetryIngressBatchClient {
	return &telemetryIngressBatchClient{
		telemBufferSize:   telemBufferSize,
		telemMaxBatchSize: telemMaxBatchSize,
		telemSendInterval: telemSendInterval,
		url:               url,
		ks:                ks,
		serverPubKeyHex:   serverPubKeyHex,
		globalLogger:      lggr,
		logging:           logging,
		lggr:              lggr.Named("TelemetryIngressBatchClient"),
		chDone:            make(chan struct{}),
		workers:           make(map[string]*telemetryIngressBatchWorker),
	}
}

// Start connects the wsrpc client to the telemetry ingress server
func (tc *telemetryIngressBatchClient) Start() error {
	return tc.StartOnce("TelemetryIngressBatchClient", func() error {
		privkey, err := tc.getCSAPrivateKey()
		if err != nil {
			return err
		}

		tc.connect(privkey)

		return nil
	})
}

// Close disconnects the wsrpc client from the ingress server and waits for all workers to exit
func (tc *telemetryIngressBatchClient) Close() error {
	return tc.StopOnce("TelemetryIngressBatchClient", func() error {
		close(tc.chDone)
		tc.wgDone.Wait()
		return nil
	})
}

// Connects to the telemetry ingress server
//
// Connection is handled in a goroutine because Dial will block
// until it can establish a connection. This is important during startup because
// we do not want to block other services from starting.
//
// Eventually when the ingress server does come back up, wsrpc will establish the connection
// without any interaction on behalf of the node operator.
func (tc *telemetryIngressBatchClient) connect(clientPrivKey []byte) {
	tc.wgDone.Add(1)

	go func() {
		defer tc.wgDone.Done()

		serverPubKey := keys.FromHex(tc.serverPubKeyHex)

		conn, err := wsrpc.Dial(tc.url.String(), wsrpc.WithTransportCreds(clientPrivKey, serverPubKey))
		if err != nil {
			tc.lggr.Errorf("Error connecting to telemetry ingress server: %v", err)
			return
		}
		defer conn.Close()

		// Initialize a new wsrpc client caller
		// This is used to call RPC methods on the server
		if tc.telemClient == nil { // only preset for tests
			tc.telemClient = telemPb.NewTelemClient(conn)
		}

		// Wait for close
		<-tc.chDone

	}()
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
func (tc *telemetryIngressBatchClient) Send(payload TelemPayload) {
	worker := tc.findOrCreateWorker(payload)
	select {
	case worker.chTelemetry <- payload:
		worker.dropMessageCount.Store(0)
	case <-payload.Ctx.Done():
		return
	default:
		worker.logBufferFullWithExpBackoff(payload)
	}
}

// findOrCreateWorker finds a worker by ContractID or creates a new one if none exists
func (tc *telemetryIngressBatchClient) findOrCreateWorker(payload TelemPayload) *telemetryIngressBatchWorker {
	tc.workersMutex.Lock()
	defer tc.workersMutex.Unlock()

	worker, found := tc.workers[payload.ContractID]

	if !found {
		worker = NewTelemetryIngressBatchWorker(
			tc.telemMaxBatchSize,
			tc.telemSendInterval,
			tc.telemClient,
			&tc.wgDone,
			tc.chDone,
			make(chan TelemPayload, tc.telemBufferSize),
			payload.ContractID,
			tc.globalLogger,
			tc.logging,
		)
		worker.Start()
		tc.workers[payload.ContractID] = worker
	}

	return worker
}
