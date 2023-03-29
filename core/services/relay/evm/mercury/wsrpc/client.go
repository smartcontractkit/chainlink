package wsrpc

import (
	"context"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/wsrpc"
	"github.com/smartcontractkit/wsrpc/connectivity"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type Client interface {
	services.ServiceCtx
	pb.MercuryClient
}

type client struct {
	utils.StartStopOnce

	csaKey       csakey.KeyV2
	serverPubKey []byte
	serverURL    string

	logger logger.Logger
	conn   *wsrpc.ClientConn
	client pb.MercuryClient
}

// Consumers of wsrpc package should not usually call NewClient directly, but instead use the Pool
func NewClient(lggr logger.Logger, clientPrivKey csakey.KeyV2, serverPubKey []byte, serverURL string) Client {
	return newClient(lggr, clientPrivKey, serverPubKey, serverURL)
}

func newClient(lggr logger.Logger, clientPrivKey csakey.KeyV2, serverPubKey []byte, serverURL string) *client {
	return &client{
		csaKey:       clientPrivKey,
		serverPubKey: serverPubKey,
		serverURL:    serverURL,
		logger:       lggr.Named("WSRPC"),
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
		conn, err := wsrpc.DialWithContext(context.Background(), w.serverURL,
			wsrpc.WithTransportCreds(w.csaKey.Raw().Bytes(), w.serverPubKey),
			wsrpc.WithLogger(w.logger),
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
	lggr := w.logger.With("req.Payload", hexutil.Encode(req.Payload))
	lggr.Debug("Transmit")
	ok := w.IfStarted(func() {
		if ready := w.conn.WaitForReady(ctx); !ready {
			err = errors.Errorf("websocket client not ready; got state: %v", w.conn.GetState())
			return
		}

		resp, err = w.client.Transmit(ctx, req)
	})
	if !ok {
		return nil, errors.New("client is not started")
	}
	if err != nil {
		lggr.Errorw("Transmit failed", "err", err, "req", req, "resp", resp)
	} else if resp.Error != "" {
		lggr.Errorw("Transmit failed; mercury server returned error", "err", resp.Error, "req", req, "resp", resp)
	} else {
		lggr.Debugw("Transmit succeeded", "resp", resp)
	}
	return
}

func (w *client) LatestReport(ctx context.Context, req *pb.LatestReportRequest) (resp *pb.LatestReportResponse, err error) {
	lggr := w.logger.With("req.FeedId", hexutil.Encode(req.FeedId))
	lggr.Debug("LatestReport")
	ok := w.IfStarted(func() {
		if ready := w.conn.WaitForReady(ctx); !ready {
			err = errors.Errorf("websocket client not ready; got state: %v", w.conn.GetState())
			return
		}

		resp, err = w.client.LatestReport(ctx, req)
	})
	if !ok {
		return nil, errors.Errorf("client is not started; state=%v", w.StartStopOnce.State())
	}
	if err != nil {
		lggr.Errorw("LatestReport failed", "err", err, "req", req, "resp", resp)
	} else if resp.Error != "" {
		lggr.Errorw("LatestReport failed; mercury server returned error", "err", resp.Error, "req", req, "resp", resp)
	} else {
		lggr.Debugw("LatestReport succeeded", "resp", resp)
	}
	return
}
