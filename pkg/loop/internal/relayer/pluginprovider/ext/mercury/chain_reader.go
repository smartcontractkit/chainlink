package mercury

import (
	"context"

	"google.golang.org/grpc"

	mercury_pb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/mercury"
	mercury_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury"
)

var _ mercury_types.ChainReader = (*ChainReaderClient)(nil)

type ChainReaderClient struct {
	grpc mercury_pb.MercuryChainReaderClient
}

func NewChainReaderClient(cc grpc.ClientConnInterface) *ChainReaderClient {
	return &ChainReaderClient{grpc: mercury_pb.NewMercuryChainReaderClient(cc)}
}

func (c *ChainReaderClient) LatestHeads(ctx context.Context, n int) ([]mercury_types.Head, error) {
	reply, err := c.grpc.LatestHeads(ctx, &mercury_pb.LatestHeadsRequest{
		NumHeads: int64(n),
	})
	if err != nil {
		return nil, err
	}
	return heads(reply.Heads), nil
}

func heads(heads []*mercury_pb.Head) []mercury_types.Head {
	res := make([]mercury_types.Head, len(heads))
	for i, head := range heads {
		res[i] = headFromPb(head)
	}
	return res
}

func headFromPb(head *mercury_pb.Head) mercury_types.Head {
	return mercury_types.Head{
		Number: head.Number,
		Hash:   head.Hash,
	}
}

var _ mercury_pb.MercuryChainReaderServer = (*ChainReaderServer)(nil)

type ChainReaderServer struct {
	mercury_pb.UnimplementedMercuryChainReaderServer

	impl mercury_types.ChainReader
}

func NewChainReaderServer(impl mercury_types.ChainReader) *ChainReaderServer {
	return &ChainReaderServer{impl: impl}
}

func (c *ChainReaderServer) LatestHeads(ctx context.Context, request *mercury_pb.LatestHeadsRequest) (*mercury_pb.LatestHeadsReply, error) {
	heads, err := c.impl.LatestHeads(ctx, int(request.NumHeads))
	if err != nil {
		return nil, err
	}
	return &mercury_pb.LatestHeadsReply{
		Heads: pbHeads(heads),
	}, nil
}

func pbHeads(heads []mercury_types.Head) []*mercury_pb.Head {
	res := make([]*mercury_pb.Head, len(heads))
	for i, head := range heads {
		res[i] = pbHead(head)
	}
	return res
}

func pbHead(head mercury_types.Head) *mercury_pb.Head {
	return &mercury_pb.Head{
		Number: head.Number,
		Hash:   head.Hash,
	}
}
