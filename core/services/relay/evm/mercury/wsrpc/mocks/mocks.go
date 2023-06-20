package mocks

import (
	"context"

	"github.com/smartcontractkit/wsrpc/connectivity"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
)

type MockWSRPCClient struct {
	TransmitF     func(ctx context.Context, in *pb.TransmitRequest) (*pb.TransmitResponse, error)
	LatestReportF func(ctx context.Context, req *pb.LatestReportRequest) (resp *pb.LatestReportResponse, err error)
}

func (m MockWSRPCClient) Name() string                   { return "" }
func (m MockWSRPCClient) Start(context.Context) error    { return nil }
func (m MockWSRPCClient) Close() error                   { return nil }
func (m MockWSRPCClient) HealthReport() map[string]error { return map[string]error{} }
func (m MockWSRPCClient) Ready() error                   { return nil }
func (m MockWSRPCClient) Transmit(ctx context.Context, in *pb.TransmitRequest) (*pb.TransmitResponse, error) {
	return m.TransmitF(ctx, in)
}
func (m MockWSRPCClient) LatestReport(ctx context.Context, in *pb.LatestReportRequest) (*pb.LatestReportResponse, error) {
	return m.LatestReportF(ctx, in)
}

type MockConn struct {
	State  connectivity.State
	Ready  bool
	Closed bool
}

func (m *MockConn) Close() {
	m.Closed = true
}
func (m MockConn) WaitForReady(ctx context.Context) bool {
	return m.Ready
}
func (m MockConn) GetState() connectivity.State { return m.State }
