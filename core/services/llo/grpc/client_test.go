package grpc

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-data-streams/rpc"
	"github.com/smartcontractkit/chainlink-data-streams/rpc/mtls"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test_Client(t *testing.T) {
	clientPrivKey := ed25519.NewKeyFromSeed(randomBytes(t, 32))
	serverPrivKey := ed25519.NewKeyFromSeed(randomBytes(t, 32))

	t.Run("Transmit errors if not started", func(t *testing.T) {
		c := NewClient(ClientOpts{
			Logger:       logger.TestLogger(t),
			ClientSigner: clientPrivKey,
			ServerPubKey: serverPrivKey.Public().(ed25519.PublicKey),
			ServerURL:    "example.com",
		})

		resp, err := c.Transmit(t.Context(), &rpc.TransmitRequest{})
		assert.Nil(t, resp)
		require.EqualError(t, err, "service is Unstarted, not started")
	})
	t.Run("Transmits report including client public key metadata", func(t *testing.T) {
		ch := make(chan packet, 100)
		srv := newMercuryServer(t, serverPrivKey, ch)
		serverURL := srv.start(t, []ed25519.PublicKey{clientPrivKey.Public().(ed25519.PublicKey)})

		c := NewClient(ClientOpts{
			Logger:       logger.TestLogger(t),
			ClientSigner: clientPrivKey,
			ServerPubKey: serverPrivKey.Public().(ed25519.PublicKey),
			ServerURL:    serverURL,
		})

		servicetest.Run(t, c)

		req := &rpc.TransmitRequest{
			Payload:      []byte("report"),
			ReportFormat: 42,
		}
		resp, err := c.Transmit(t.Context(), req)
		require.NoError(t, err)

		assert.Equal(t, "", resp.Error)
		assert.Equal(t, int32(1), resp.Code)

		select {
		case p := <-ch:
			assert.Equal(t, req.Payload, p.req.Payload)
			assert.Equal(t, req.ReportFormat, p.req.ReportFormat)
			m, ok := metadata.FromIncomingContext(p.ctx)
			require.True(t, ok)
			require.Len(t, m["client_public_key"], 1)
			assert.Equal(t, hex.EncodeToString(clientPrivKey.Public().(ed25519.PublicKey)), m["client_public_key"][0])
		default:
			t.Fatal("expected request to be received")
		}
	})
}

func randomBytes(t *testing.T, n int) (r []byte) {
	r = make([]byte, n)
	_, err := rand.Read(r)
	require.NoError(t, err)
	return
}

type packet struct {
	ctx context.Context //nolint:containedctx // this is used solely for test purposes
	req *rpc.TransmitRequest
}

type mercuryServer struct {
	rpc.UnimplementedTransmitterServer
	privKey   ed25519.PrivateKey
	packetsCh chan packet
	t         *testing.T
}

func newMercuryServer(t *testing.T, privKey ed25519.PrivateKey, packetsCh chan packet) *mercuryServer {
	return &mercuryServer{rpc.UnimplementedTransmitterServer{}, privKey, packetsCh, t}
}

func (srv *mercuryServer) start(t *testing.T, clientPubKeys []ed25519.PublicKey) (serverURL string) {
	// Set up the grpc server
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("[MAIN] failed to listen: %v", err)
	}
	serverURL = lis.Addr().String()
	sMtls, err := mtls.NewTransportCredentials(srv.privKey, clientPubKeys)
	require.NoError(t, err)
	s := grpc.NewServer(grpc.Creds(sMtls))

	// Register mercury implementation with the wsrpc server
	rpc.RegisterTransmitterServer(s, srv)

	// Start serving
	go func() {
		s.Serve(lis) //nolint:errcheck // don't care about errors in tests
	}()

	t.Cleanup(s.Stop)

	return
}

func (srv *mercuryServer) Transmit(ctx context.Context, req *rpc.TransmitRequest) (*rpc.TransmitResponse, error) {
	srv.packetsCh <- packet{ctx, req}

	return &rpc.TransmitResponse{
		Code:  1,
		Error: "",
	}, nil
}

func (srv *mercuryServer) LatestReport(ctx context.Context, lrr *rpc.LatestReportRequest) (*rpc.LatestReportResponse, error) {
	panic("should not be called")
}
