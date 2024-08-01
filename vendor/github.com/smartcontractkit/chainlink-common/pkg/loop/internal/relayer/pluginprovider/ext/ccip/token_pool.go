package ccip

import (
	"context"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	ccippb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

// TokenPoolBatchedReaderGRPCClient implements [cciptypes.TokenPoolBatchedReader] by wrapping a
// [ccippb.TokenPoolBatchedReaderGRPCClient] grpc client.
// It is used by a ReportingPlugin to call the TokenPoolBatchedReader service, which
// is hosted by the relayer

type TokenPoolBatchedReaderGRPCClient struct {
	client ccippb.TokenPoolBatcherReaderClient
	conn   grpc.ClientConnInterface
}

func NewTokenPoolBatchedReaderGRPCClient(cc grpc.ClientConnInterface) *TokenPoolBatchedReaderGRPCClient {
	return &TokenPoolBatchedReaderGRPCClient{client: ccippb.NewTokenPoolBatcherReaderClient(cc), conn: cc}
}

// TokenPoolBatchedReaderGRPCServer implements [ccippb.TokenPoolBatchedReaderServer] by wrapping a
// [cciptypes.TokenPoolBatchedReader] implementation.
// This server is hosted by the relayer and is called ReportingPlugin via
// the [TokenPoolBatchedReaderGRPCClient]
type TokenPoolBatchedReaderGRPCServer struct {
	ccippb.UnimplementedTokenPoolBatcherReaderServer

	impl cciptypes.TokenPoolBatchedReader
	deps []io.Closer
}

func NewTokenPoolBatchedReaderGRPCServer(impl cciptypes.TokenPoolBatchedReader) *TokenPoolBatchedReaderGRPCServer {
	return &TokenPoolBatchedReaderGRPCServer{impl: impl, deps: []io.Closer{impl}}
}

// ensure interface is implemented
var _ ccippb.TokenPoolBatcherReaderServer = (*TokenPoolBatchedReaderGRPCServer)(nil)
var _ cciptypes.TokenPoolBatchedReader = (*TokenPoolBatchedReaderGRPCClient)(nil)

func (t *TokenPoolBatchedReaderGRPCClient) ClientConn() grpc.ClientConnInterface {
	return t.conn
}

// Close implements ccip.TokenPoolBatchedReader.
func (t *TokenPoolBatchedReaderGRPCClient) Close() error {
	return shutdownGRPCServer(context.Background(), t.client)
}

// GetInboundTokenPoolRateLimits implements ccip.TokenPoolBatchedReader.
func (t *TokenPoolBatchedReaderGRPCClient) GetInboundTokenPoolRateLimits(ctx context.Context, tokenPoolReaders []cciptypes.Address) ([]cciptypes.TokenBucketRateLimit, error) {
	req := &ccippb.GetInboundTokenPoolRateLimitsRequest{
		TokenPoolReaders: cciptypes.Addresses(tokenPoolReaders).Strings(),
	}
	resp, err := t.client.GetInboundTokenPoolRateLimits(ctx, req)
	if err != nil {
		return nil, err
	}
	rateLimits := make([]cciptypes.TokenBucketRateLimit, len(resp.TokenPoolRateLimits))
	for i, r := range resp.TokenPoolRateLimits {
		rateLimits[i] = tokenBucketRateLimit(r)
	}
	return rateLimits, nil
}

// Server methods

func (t *TokenPoolBatchedReaderGRPCServer) AddDep(dep io.Closer) *TokenPoolBatchedReaderGRPCServer {
	t.deps = append(t.deps, dep)
	return t
}

// GetInboundTokenPoolRateLimits implements ccippb.TokenPoolBatcherReaderServer.
func (t *TokenPoolBatchedReaderGRPCServer) GetInboundTokenPoolRateLimits(ctx context.Context, req *ccippb.GetInboundTokenPoolRateLimitsRequest) (*ccippb.GetInboundTokenPoolRateLimitsResponse, error) {
	rateLimts, err := t.impl.GetInboundTokenPoolRateLimits(ctx, cciptypes.MakeAddresses(req.TokenPoolReaders))
	if err != nil {
		return nil, err
	}
	pbRateLimits := make([]*ccippb.TokenPoolRateLimit, len(rateLimts))
	for i, r := range rateLimts {
		pbRateLimits[i] = tokenBucketRateLimitPB(r)
	}
	return &ccippb.GetInboundTokenPoolRateLimitsResponse{
		TokenPoolRateLimits: pbRateLimits,
	}, nil
}

func (t *TokenPoolBatchedReaderGRPCServer) Close(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, services.MultiCloser(t.deps).Close()
}
