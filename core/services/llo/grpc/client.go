package grpc

import (
	"context"
	"crypto"
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-data-streams/rpc"
	"github.com/smartcontractkit/chainlink-data-streams/rpc/mtls"
)

type Client interface {
	services.Service
	Transmit(ctx context.Context, in *rpc.TransmitRequest) (*rpc.TransmitResponse, error)
	ServerURL() string
}

var _ Client = (*client)(nil)

type client struct {
	services.Service
	eng *services.Engine

	clientSigner    crypto.Signer
	clientPubKeyHex string
	serverPubKey    ed25519.PublicKey
	serverURL       string

	conn   *grpc.ClientConn
	client rpc.TransmitterClient
}

type ClientOpts struct {
	Logger       logger.Logger
	ClientSigner crypto.Signer
	ServerPubKey ed25519.PublicKey
	ServerURL    string
}

func NewClient(opts ClientOpts) Client {
	return newClient(opts)
}

func newClient(opts ClientOpts) Client {
	c := &client{
		clientSigner:    opts.ClientSigner,
		clientPubKeyHex: hex.EncodeToString(opts.ClientSigner.Public().(ed25519.PublicKey)),
		serverPubKey:    opts.ServerPubKey,
		serverURL:       opts.ServerURL,
	}
	c.Service, c.eng = services.Config{
		Name:  "GRPCClient",
		Start: c.start,
		Close: c.close,
	}.NewServiceEngine(opts.Logger)
	return c
}

func (c *client) start(context.Context) error {
	cMtls, err := mtls.NewTransportSigner(c.clientSigner, []ed25519.PublicKey{c.serverPubKey})
	if err != nil {
		return fmt.Errorf("failed to create client mTLS credentials: %w", err)
	}
	// Latency is critical so configure aggressively for fast
	// redial attempts and short keepalive
	clientConn, err := grpc.NewClient(
		c.serverURL,
		grpc.WithTransportCredentials(cMtls),
		grpc.WithConnectParams(
			grpc.ConnectParams{
				Backoff: backoff.Config{
					BaseDelay:  1 * time.Second,
					Multiplier: 2,
					Jitter:     0.2,
					MaxDelay:   30 * time.Second,
				},
				MinConnectTimeout: time.Second,
			},
		),
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				Time:                time.Second * 10,
				Timeout:             time.Second * 20,
				PermitWithoutStream: true,
			}),
	)
	if err != nil {
		return fmt.Errorf("failed to create client connection: %w", err)
	}
	c.conn = clientConn
	c.client = rpc.NewTransmitterClient(c.conn)
	return nil
}

func (c *client) close() error {
	return c.conn.Close()
}

func (c *client) Transmit(ctx context.Context, req *rpc.TransmitRequest) (resp *rpc.TransmitResponse, err error) {
	err = c.eng.IfStarted(func() error {
		// This is a self-identified client ID
		// It is not cryptographically verified
		transmitCtx := metadata.AppendToOutgoingContext(ctx, "client_public_key", c.clientPubKeyHex)
		resp, err = c.client.Transmit(transmitCtx, req)
		return err
	})
	return
}

func (c *client) LatestReport(ctx context.Context, req *rpc.LatestReportRequest) (resp *rpc.LatestReportResponse, err error) {
	return nil, errors.New("LatestReport is not supported in grpc mode")
}

func (c *client) ServerURL() string {
	return c.serverURL
}
