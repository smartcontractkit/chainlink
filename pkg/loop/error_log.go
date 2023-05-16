package loop

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/smartcontractkit/chainlink-relay/pkg/loop/internal/pb"
)

var _ ErrorLog = (*errorLogClient)(nil)

type errorLogClient struct {
	grpc pb.ErrorLogClient
}

func (e errorLogClient) SaveError(ctx context.Context, msg string) error {
	_, err := e.grpc.SaveError(ctx, &pb.SaveErrorRequest{Message: msg})
	return err
}

func newErrorLogClient(cc grpc.ClientConnInterface) *errorLogClient {
	return &errorLogClient{pb.NewErrorLogClient(cc)}
}

var _ pb.ErrorLogServer = (*errorLogServer)(nil)

type errorLogServer struct {
	pb.UnimplementedErrorLogServer

	impl ErrorLog
}

func (e *errorLogServer) SaveError(ctx context.Context, request *pb.SaveErrorRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, e.impl.SaveError(ctx, request.Message)
}
