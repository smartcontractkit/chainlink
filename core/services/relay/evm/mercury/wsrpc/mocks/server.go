package mocks

import (
	"context"
	"crypto/ed25519"
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/smartcontractkit/wsrpc"
	"github.com/smartcontractkit/wsrpc/credentials"
	"github.com/smartcontractkit/wsrpc/peer"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
)

var _ pb.MercuryServer = &Server{}

type Request struct {
	PK  credentials.StaticSizedPublicKey
	Req *pb.TransmitRequest
}

type Server struct {
	privKey ed25519.PrivateKey
	reqsCh  chan Request
	t       *testing.T
}

func NewServer(t *testing.T, privKey ed25519.PrivateKey, reqsCh chan Request) *Server {
	return &Server{privKey, reqsCh, t}
}

func (s *Server) Transmit(ctx context.Context, req *pb.TransmitRequest) (*pb.TransmitResponse, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, errors.New("could not extract public key")
	}
	r := Request{p.PublicKey, req}
	s.reqsCh <- r

	return &pb.TransmitResponse{
		Code:  1,
		Error: "",
	}, nil
}

func (s *Server) LatestReport(ctx context.Context, lrr *pb.LatestReportRequest) (*pb.LatestReportResponse, error) {
	// not implemented in test
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, errors.New("could not extract public key")
	}
	s.t.Logf("mercury server got latest report from %x for feed id 0x%x", p.PublicKey, lrr.FeedId)
	return nil, nil
}

func (s *Server) Start(t *testing.T, pubKeys []ed25519.PublicKey) (serverURL string) {
	// Set up the wsrpc server
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("[MAIN] failed to listen: %v", err)
	}
	serverURL = fmt.Sprintf("%s", lis.Addr().String())
	srv := wsrpc.NewServer(wsrpc.Creds(s.privKey, pubKeys))

	// Register mercury implementation with the wsrpc server
	pb.RegisterMercuryServer(srv, s)

	// Start serving
	go srv.Serve(lis)
	t.Cleanup(srv.Stop)

	return
}
