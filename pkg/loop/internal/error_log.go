package internal

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

var _ types.ErrorLog = (*errorLogClient)(nil)

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

	impl types.ErrorLog
}

func (e *errorLogServer) SaveError(ctx context.Context, request *pb.SaveErrorRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, e.impl.SaveError(ctx, request.Message)
}
