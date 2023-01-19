package wsrpc

import (
	"context"
	"net/url"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/wsrpc"
	"github.com/smartcontractkit/wsrpc/connectivity"

	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/core/services/relay/evm/mercury/wsrpc/pb"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type Client interface {
	services.ServiceCtx
	pb.MercuryClient
}

type client struct {
	utils.StartStopOnce

	csaKey       csakey.KeyV2
	serverPubKey []byte
	serverURL    *url.URL

	conn   *wsrpc.ClientConn
	client pb.MercuryClient
}

func NewClient(privKey csakey.KeyV2, serverPubKey []byte, serverURL *url.URL) Client {
	return &client{
		csaKey:       privKey,
		serverPubKey: serverPubKey,
		serverURL:    serverURL,
	}
}

func (w *client) Start(_ context.Context) error {
	return w.StartOnce("WSRPC Client", func() error {
		// NOTE: Dial is non-blocking, and will retry on an exponential backoff
		// in the background until close is called, or context is cancelled.
		// This is why we use the background context, not the start context here.
		//
		// Any transmits made while client is still trying to dial will fail
		// with error.
		conn, err := wsrpc.DialWithContext(context.Background(), w.serverURL.String(),
			wsrpc.WithTransportCreds(w.csaKey.Raw().Bytes(), w.serverPubKey),
		)
		if err != nil {
			return errors.Wrap(err, "failed to dial wsrpc client")
		}
		w.conn = conn
		w.client = pb.NewMercuryClient(conn)
		return nil
	})
}

func (w *client) Close() error {
	return w.StopOnce("WSRPC Client", func() error {
		w.conn.Close()
		return nil
	})
}

func (w *client) Name() string {
	return "EVM.Mercury.WSRPCClient"
}

func (w *client) HealthReport() map[string]error {
	return map[string]error{w.Name(): w.Healthy()}
}

func (w *client) Transmit(ctx context.Context, in *report.ReportRequest) (rr *report.ReportResponse, err error) {
	ok := w.IfStarted(func() {
		rr, err = w.client.Transmit(ctx, in)
	})
	if !ok {
		return nil, errors.New("client is not started")
	}
	return
}

// Healthy if connected
func (w *client) Healthy() (err error) {
	if err = w.StartStopOnce.Healthy(); err != nil {
		return err
	}
	state := w.conn.GetState()
	if state != connectivity.Ready {
		return errors.Errorf("client state should be %s; got %s", connectivity.Ready, state)
	}
	return nil
}

func (w *client) Transmit(ctx context.Context, req *pb.TransmitRequest) (resp *pb.TransmitResponse, err error) {
	ok := w.IfStarted(func() {
		resp, err = w.client.Transmit(ctx, req)
	})
	if !ok {
		return nil, errors.New("client is not started")
	}
	return
}

func (w *client) LatestReport(ctx context.Context, req *pb.LatestReportRequest) (resp *pb.LatestReportResponse, err error) {
	ok := w.IfStarted(func() {
		resp, err = w.client.LatestReport(ctx, req)
	})
	if !ok {
		return nil, errors.New("client is not started")
	}
	return
}
