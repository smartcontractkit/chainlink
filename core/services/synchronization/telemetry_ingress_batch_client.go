package synchronization

import (
	"context"
	"errors"
	"net/url"
	"sync"
	"time"

	"go.uber.org/atomic"

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
	Start() error
	Close() error
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

// telemetryIngressBatchWorker pushes telemetry in batches to the ingress server via wsrpc.
// A worker is created per ContractID.
type telemetryIngressBatchWorker struct {
	telemMaxBatchSize uint
	telemSendInterval time.Duration
	telemClient       telemPb.TelemClient
	wgDone            *sync.WaitGroup
	chDone            chan struct{}
	chTelemetry       chan TelemPayload
	contractID        string
	logging           bool
	lggr              logger.Logger
	dropMessageCount  atomic.Uint32
}

// NewTelemetryIngressBatchWorker returns a worker for a given contractID that can send
// telemetry to the ingress server via WSRPC
func NewTelemetryIngressBatchWorker(
	telemMaxBatchSize uint,
	telemSendInterval time.Duration,
	telemClient telemPb.TelemClient,
	wgDone *sync.WaitGroup,
	chDone chan struct{},
	chTelemetry chan TelemPayload,
	contractID string,
	globalLogger logger.Logger,
	logging bool,
) *telemetryIngressBatchWorker {
	return &telemetryIngressBatchWorker{
		telemSendInterval: telemSendInterval,
		telemMaxBatchSize: telemMaxBatchSize,
		telemClient:       telemClient,
		wgDone:            wgDone,
		chDone:            chDone,
		chTelemetry:       chTelemetry,
		contractID:        contractID,
		logging:           logging,
		lggr:              globalLogger.Named("TelemetryIngressBatchWorker"),
	}
}

// Start sends batched telemetry to the ingress server on an interval
func (tw *telemetryIngressBatchWorker) Start() {
	tw.wgDone.Add(1)
	sendTicker := time.NewTicker(tw.telemSendInterval)

	go func() {
		defer tw.wgDone.Done()

		for {
			select {
			case <-sendTicker.C:
				if len(tw.chTelemetry) == 0 {
					continue
				}

				// Send batched telemetry to the ingress server, log any errors
				telemBatchReq := tw.BuildTelemBatchReq()
				_, err := tw.telemClient.TelemBatch(context.Background(), telemBatchReq)

				if err != nil {
					tw.lggr.Warnf("Could not send telemetry: %v", err)
					continue
				}
				if tw.logging {
					tw.lggr.Debugw("Successfully sent telemetry to ingress server", "contractID", telemBatchReq.ContractId, "telemetry", telemBatchReq.Telemetry)
				}
			case <-tw.chDone:
				return
			}
		}
	}()
}

// logBufferFullWithExpBackoff logs messages at
// 1
// 2
// 4
// 8
// 16
// 32
// 64
// 100
// 200
// 300
// etc...
func (tw *telemetryIngressBatchWorker) logBufferFullWithExpBackoff(payload TelemPayload) {
	count := tw.dropMessageCount.Inc()
	if count > 0 && (count%100 == 0 || count&(count-1) == 0) {
		tw.lggr.Warnw("telemetry ingress client buffer full, dropping message", "telemetry", payload.Telemetry, "droppedCount", count)
	}
}

// BuildTelemBatchReq reads telemetry off the worker channel and packages it into a batch request
func (tw *telemetryIngressBatchWorker) BuildTelemBatchReq() *telemPb.TelemBatchRequest {
	var telemBatch [][]byte

	// Read telemetry off the channel up to the max batch size
	for len(tw.chTelemetry) > 0 && len(telemBatch) < int(tw.telemMaxBatchSize) {
		telemPayload := <-tw.chTelemetry
		telemBatch = append(telemBatch, telemPayload.Telemetry)
	}

	return &telemPb.TelemBatchRequest{
		ContractId: tw.contractID,
		Telemetry:  telemBatch,
	}
}
