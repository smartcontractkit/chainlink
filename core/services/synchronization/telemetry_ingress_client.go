package synchronization

import (
	"context"
	"errors"
	"net/url"
	"sync"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/service"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	telemPb "github.com/smartcontractkit/chainlink/core/services/synchronization/telem"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/smartcontractkit/wsrpc"
	"github.com/smartcontractkit/wsrpc/examples/simple/keys"
)

//go:generate mockery --dir ./telem --name TelemClient --output ./mocks/ --case=underscore

// SendIngressBufferSize is the number of messages to keep in the buffer before dropping additional ones
const SendIngressBufferSize = 100

// TelemetryIngressClient encapsulates all the functionality needed to
// send telemetry to the ingress server using wsrpc
type TelemetryIngressClient interface {
	service.Service
	Start() error
	Close() error
	Send(TelemPayload)
	Unsafe_SetTelemClient(telemPb.TelemClient) bool
}

type NoopTelemetryIngressClient struct{}

func (NoopTelemetryIngressClient) Start() error                                   { return nil }
func (NoopTelemetryIngressClient) Close() error                                   { return nil }
func (NoopTelemetryIngressClient) Send(TelemPayload)                              {}
func (NoopTelemetryIngressClient) Healthy() error                                 { return nil }
func (NoopTelemetryIngressClient) Ready() error                                   { return nil }
func (NoopTelemetryIngressClient) Unsafe_SetTelemClient(telemPb.TelemClient) bool { return true }

type telemetryIngressClient struct {
	utils.StartStopOnce
	url             *url.URL
	ks              keystore.CSA
	serverPubKeyHex string

	telemClient telemPb.TelemClient
	logging     bool

	wgDone           sync.WaitGroup
	chDone           chan struct{}
	dropMessageCount uint32
	chTelemetry      chan TelemPayload
}

type TelemPayload struct {
	Ctx             context.Context
	Telemetry       []byte
	ContractAddress common.Address
}

// NewTelemetryIngressClient returns a client backed by wsrpc that
// can send telemetry to the telemetry ingress server
func NewTelemetryIngressClient(url *url.URL, serverPubKeyHex string, ks keystore.CSA, logging bool) TelemetryIngressClient {
	return &telemetryIngressClient{
		url:             url,
		ks:              ks,
		serverPubKeyHex: serverPubKeyHex,
		logging:         logging,
		chTelemetry:     make(chan TelemPayload, SendIngressBufferSize),
		chDone:          make(chan struct{}),
	}
}

// Start connects the wsrpc client to the telemetry ingress server
func (tc *telemetryIngressClient) Start() error {
	return tc.StartOnce("TelemetryIngressClient", func() error {
		privkey, err := tc.getCSAPrivateKey()
		if err != nil {
			return err
		}

		tc.connect(privkey)

		return nil
	})
}

// Close disconnects the wsrpc client from the ingress server
func (tc *telemetryIngressClient) Close() error {
	return tc.StopOnce("TelemetryIngressClient", func() error {
		close(tc.chDone)
		tc.wgDone.Wait()
		return nil
	})
}

func (tc *telemetryIngressClient) connect(clientPrivKey []byte) {
	tc.wgDone.Add(1)

	go func() {
		defer tc.wgDone.Done()

		serverPubKey := keys.FromHex(tc.serverPubKeyHex)

		conn, err := wsrpc.Dial(tc.url.String(), wsrpc.WithTransportCreds(clientPrivKey, serverPubKey))
		if err != nil {
			logger.Errorf("Error connecting to telemetry ingress server: %v", err)
			return
		}
		defer conn.Close()

		// Initialize a new wsrpc client caller
		// This is used to call RPC methods on the server
		tc.telemClient = telemPb.NewTelemClient(conn)

		// Start handler for telemetry
		tc.handleTelemetry()

		// Wait for close
		<-tc.chDone

	}()
}

func (tc *telemetryIngressClient) handleTelemetry() {
	go func() {
		for {
			select {
			case p := <-tc.chTelemetry:
				// Send telemetry to the ingress server, log any errors
				telemReq := &telemPb.TelemRequest{Telemetry: p.Telemetry, Address: p.ContractAddress.String()}
				_, err := tc.telemClient.Telem(p.Ctx, telemReq)
				if err != nil {
					logger.Errorf("Could not send telemetry: %v", err)
					continue
				}
				if tc.logging {
					logger.Debugw("successfully sent telemetry to ingress server", "contractAddress", p.ContractAddress.String(), "telemetry", p.Telemetry)
				}
			case <-tc.chDone:
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
func (tc *telemetryIngressClient) logBufferFullWithExpBackoff(payload TelemPayload) {
	count := atomic.AddUint32(&tc.dropMessageCount, 1)
	if count > 0 && (count%100 == 0 || count&(count-1) == 0) {
		logger.Warnw("telemetry ingress client buffer full, dropping message", "telemetry", payload.Telemetry, "droppedCount", count)
	}
}

// getCSAPrivateKey gets the client's CSA private key
func (tc *telemetryIngressClient) getCSAPrivateKey() (privkey []byte, err error) {
	// Fetch the client's public key
	keys, err := tc.ks.GetAll()
	if err != nil {
		return privkey, err
	}
	if len(keys) < 1 {
		return privkey, errors.New("CSA key does not exist")
	}

	return keys[0].Raw(), nil
}

// Send sends telemetry to the ingress server using wsrpc if the client is ready.
// Also stores telemetry in a small buffer in case of backpressure from wsrpc,
// throwing away messages once buffer is full
func (tc *telemetryIngressClient) Send(payload TelemPayload) {
	select {
	case tc.chTelemetry <- payload:
		atomic.StoreUint32(&tc.dropMessageCount, 0)
	case <-payload.Ctx.Done():
		return
	default:
		tc.logBufferFullWithExpBackoff(payload)
	}
}

// Unsafe_SetTelemClient sets the TelemClient on the service.
//
// We need to be able to inject a mock for the client to facilitate integration
// tests.
//
// ONLY TO BE USED FOR TESTING.
func (tc *telemetryIngressClient) Unsafe_SetTelemClient(client telemPb.TelemClient) bool {
	if tc.telemClient == nil {
		return false
	}

	tc.telemClient = client
	return true
}
