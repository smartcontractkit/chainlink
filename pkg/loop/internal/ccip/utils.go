package ccip

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type grpcResourceCloser interface {
	Close(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

// shutdownGRPCServer is a helper function to release server resources
// created by a grpc client.
func shutdownGRPCServer(ctx context.Context, rc grpcResourceCloser) error {
	_, err := rc.Close(ctx, &emptypb.Empty{})
	// due to the handler in the server, it may shutdown before it sends a response to client
	// in that case, we expect the client to receive an Unavailable or Internal error
	if status.Code(err) == codes.Unavailable || status.Code(err) == codes.Internal {
		return nil
	}
	return err
}
