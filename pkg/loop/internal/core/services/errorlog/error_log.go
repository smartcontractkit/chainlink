package errorlog

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

func NewClient(cc grpc.ClientConnInterface) *errorLogClient {
	return &errorLogClient{pb.NewErrorLogClient(cc)}
}

var _ pb.ErrorLogServer = (*Server)(nil)

type Server struct {
	pb.UnimplementedErrorLogServer

	Impl types.ErrorLog
}

func (e *Server) SaveError(ctx context.Context, request *pb.SaveErrorRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, e.Impl.SaveError(ctx, request.Message)
}
