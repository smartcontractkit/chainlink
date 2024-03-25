package mercury

import (
	"context"
	"fmt"
	"math/big"

	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	mercury_pb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/mercury"
	mercury_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury"
)

var _ mercury_types.ServerFetcher = (*ServerFetcherClient)(nil)

type ServerFetcherClient struct {
	grpc mercury_pb.ServerFetcherClient
}

func NewServerFetcherClient(cc grpc.ClientConnInterface) *ServerFetcherClient {
	return &ServerFetcherClient{grpc: mercury_pb.NewServerFetcherClient(cc)}
}

func (s *ServerFetcherClient) FetchInitialMaxFinalizedBlockNumber(ctx context.Context) (*int64, error) {
	reply, err := s.grpc.FetchInitialMaxFinalizedBlockNumber(ctx, &mercury_pb.FetchInitialMaxFinalizedBlockNumberRequest{})
	if err != nil {
		return nil, err
	}
	return &reply.InitialMaxFinalizedBlockNumber, nil
}

func (s *ServerFetcherClient) LatestPrice(ctx context.Context, feedID [32]byte) (*big.Int, error) {
	reply, err := s.grpc.LatestPrice(ctx, &mercury_pb.LatestPriceRequest{})
	if err != nil {
		return nil, err
	}
	return reply.LatestPrice.Int(), nil
}

func (s *ServerFetcherClient) LatestTimestamp(ctx context.Context) (int64, error) {
	reply, err := s.grpc.LatestTimestamp(ctx, &mercury_pb.LatestTimestampRequest{})
	if err != nil {
		return 0, err
	}
	return reply.LatestTimestamp, nil
}

var _ mercury_pb.ServerFetcherServer = (*ServerFetcherServer)(nil)

type ServerFetcherServer struct {
	mercury_pb.UnimplementedServerFetcherServer

	impl mercury_types.ServerFetcher
}

func NewServerFetcherServer(impl mercury_types.ServerFetcher) *ServerFetcherServer {
	return &ServerFetcherServer{impl: impl}
}

func (s *ServerFetcherServer) FetchInitialMaxFinalizedBlockNumber(ctx context.Context, request *mercury_pb.FetchInitialMaxFinalizedBlockNumberRequest) (*mercury_pb.FetchInitialMaxFinalizedBlockNumberReply, error) {
	val, err := s.impl.FetchInitialMaxFinalizedBlockNumber(ctx)
	if err != nil {
		return nil, err
	}
	// play defense against a nil dereference below. it's a bit weird that we're returning a pointer to an int64.
	if val == nil {
		return nil, fmt.Errorf("FetchInitialMaxFinalizedBlockNumber returned nil")
	}
	return &mercury_pb.FetchInitialMaxFinalizedBlockNumberReply{InitialMaxFinalizedBlockNumber: *val}, nil
}

func (s *ServerFetcherServer) LatestPrice(ctx context.Context, request *mercury_pb.LatestPriceRequest) (*mercury_pb.LatestPriceReply, error) {
	if len(request.FeedID) != 32 {
		return nil, fmt.Errorf("expected feed ID to be 32 bytes, got %d", len(request.FeedID))
	}
	val, err := s.impl.LatestPrice(ctx, ([32]byte(request.FeedID[:32])))
	if err != nil {
		return nil, err
	}
	return &mercury_pb.LatestPriceReply{LatestPrice: pb.NewBigIntFromInt(val)}, nil
}

func (s *ServerFetcherServer) LatestTimestamp(ctx context.Context, request *mercury_pb.LatestTimestampRequest) (*mercury_pb.LatestTimestampReply, error) {
	val, err := s.impl.LatestTimestamp(ctx)
	if err != nil {
		return nil, err
	}
	return &mercury_pb.LatestTimestampReply{LatestTimestamp: val}, nil
}
