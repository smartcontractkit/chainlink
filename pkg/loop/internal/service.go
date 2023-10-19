package internal

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop/internal/pb"
	"github.com/smartcontractkit/chainlink-relay/pkg/services"
)

var ErrPluginUnavailable = errors.New("plugin unavailable")

var _ services.Service = (*serviceClient)(nil)

type serviceClient struct {
	b    *brokerExt
	cc   grpc.ClientConnInterface
	grpc pb.ServiceClient
}

func newServiceClient(b *brokerExt, cc grpc.ClientConnInterface) *serviceClient {
	return &serviceClient{b, cc, pb.NewServiceClient(cc)}
}

func (s *serviceClient) Start(ctx context.Context) error {
	return nil // no-op: server side starts automatically
}

func (s *serviceClient) Close() error {
	ctx, cancel := s.b.stopCtx()
	defer cancel()

	_, err := s.grpc.Close(ctx, &emptypb.Empty{})
	return err
}

func (s *serviceClient) Ready() error {
	ctx, cancel := s.b.stopCtx()
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, time.Second)
	defer cancel()

	_, err := s.grpc.Ready(ctx, &emptypb.Empty{})
	return err
}

func (s *serviceClient) Name() string { return s.b.Logger.Name() }

func (s *serviceClient) HealthReport() map[string]error {
	ctx, cancel := s.b.stopCtx()
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, time.Second)
	defer cancel()

	reply, err := s.grpc.HealthReport(ctx, &emptypb.Empty{})
	if err != nil {
		return map[string]error{s.b.Logger.Name(): err}
	}
	hr := healthReport(reply.HealthReport)
	hr[s.b.Logger.Name()] = nil
	return hr
}

var _ pb.ServiceServer = (*serviceServer)(nil)

type serviceServer struct {
	pb.UnimplementedServiceServer
	srv services.Service
}

func (s *serviceServer) Close(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.srv.Close()
}

func (s *serviceServer) Ready(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.srv.Ready()
}

func (s *serviceServer) HealthReport(ctx context.Context, empty *emptypb.Empty) (*pb.HealthReportReply, error) {
	var r pb.HealthReportReply
	r.HealthReport = make(map[string]string)
	for n, err := range s.srv.HealthReport() {
		var serr string
		if err != nil {
			serr = err.Error()
		}
		r.HealthReport[n] = serr
	}
	return &r, nil
}
