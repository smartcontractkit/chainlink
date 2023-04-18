package loop

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/smartcontractkit/chainlink-relay/pkg/loop/internal/pb"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
)

var _ types.Service = (*serviceClient)(nil)

type serviceClient struct {
	*lggrBroker
	cc   *grpc.ClientConn
	grpc pb.ServiceClient
}

func newServiceClient(lggrBroker *lggrBroker, cc *grpc.ClientConn) *serviceClient {
	return &serviceClient{lggrBroker, cc, pb.NewServiceClient(cc)}
}

func (s *serviceClient) Start(ctx context.Context) error {
	_, err := s.grpc.Start(ctx, &emptypb.Empty{})
	return err
}

func (s *serviceClient) Close() error {
	_, err := s.grpc.Close(context.TODO(), &emptypb.Empty{})
	return err
}

func (s *serviceClient) Ready() error {
	_, err := s.grpc.Ready(context.TODO(), &emptypb.Empty{})
	return err
}

func (s *serviceClient) Name() string { return s.lggr.Name() }

func (s *serviceClient) HealthReport() map[string]error {
	reply, err := s.grpc.HealthReport(context.TODO(), &emptypb.Empty{})
	if err != nil {
		return map[string]error{s.lggr.Name(): err}
	}
	hr := healthReport(reply.HealthReport)
	hr[s.lggr.Name()] = nil
	return hr
}

var _ pb.ServiceServer = (*serviceServer)(nil)

type serviceServer struct {
	pb.UnimplementedServiceServer
	srv  types.Service
	stop func()
}

func (s *serviceServer) Start(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.srv.Start(ctx)
}

func (s *serviceServer) Close(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	s.stop()
	return &emptypb.Empty{}, s.srv.Close()
}

func (s *serviceServer) Ready(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.srv.Ready()
}
