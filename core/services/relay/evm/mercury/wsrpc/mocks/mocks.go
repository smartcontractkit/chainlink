package mocks

import (
	"context"

	grpc_connectivity "google.golang.org/grpc/connectivity"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
)

type MockWSRPCClient struct {
	TransmitF     func(ctx context.Context, in *pb.TransmitRequest) (*pb.TransmitResponse, error)
	LatestReportF func(ctx context.Context, req *pb.LatestReportRequest) (resp *pb.LatestReportResponse, err error)
}

func (m *MockWSRPCClient) Name() string                   { return "" }
func (m *MockWSRPCClient) Start(context.Context) error    { return nil }
func (m *MockWSRPCClient) Close() error                   { return nil }
func (m *MockWSRPCClient) HealthReport() map[string]error { return map[string]error{} }
func (m *MockWSRPCClient) Ready() error                   { return nil }
func (m *MockWSRPCClient) Transmit(ctx context.Context, in *pb.TransmitRequest) (*pb.TransmitResponse, error) {
	return m.TransmitF(ctx, in)
}
func (m *MockWSRPCClient) LatestReport(ctx context.Context, in *pb.LatestReportRequest) (*pb.LatestReportResponse, error) {
	return m.LatestReportF(ctx, in)
}
func (m *MockWSRPCClient) ServerURL() string { return "mock server url" }

func (m *MockWSRPCClient) RawClient() pb.MercuryClient { return nil }

type MockConn struct {
	State  grpc_connectivity.State
	Ready  bool
	Closed bool
}

func (m *MockConn) Close() error {
	m.Closed = true
	return nil
}
func (m MockConn) WaitForReady(ctx context.Context) bool {
	return m.Ready
}
func (m MockConn) GetState() grpc_connectivity.State { return m.State }
