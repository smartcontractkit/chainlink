package keyvalue

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

var _ core.KeyValueStore = (*Client)(nil)

type Client struct {
	grpc pb.KeyValueStoreClient
}

func (k Client) Store(ctx context.Context, key string, val []byte) error {
	_, err := k.grpc.StoreKeyValue(ctx, &pb.StoreKeyValueRequest{Key: key, Value: val})
	if err != nil {
		return fmt.Errorf("failed to store value: %s for key: %s: %w", string(val), key, err)
	}

	return nil
}

func (k Client) Get(ctx context.Context, key string) ([]byte, error) {
	resp, err := k.grpc.GetValueForKey(ctx, &pb.GetValueForKeyRequest{Key: key})
	if err != nil {
		return nil, fmt.Errorf("failed to get value for key: %s: %w", key, err)
	}

	return resp.Value, nil
}

func NewClient(cc grpc.ClientConnInterface) *Client {
	return &Client{pb.NewKeyValueStoreClient(cc)}
}

var _ pb.KeyValueStoreServer = (*Server)(nil)

type Server struct {
	pb.UnimplementedKeyValueStoreServer
	impl core.KeyValueStore
}

func NewServer(impl core.KeyValueStore) *Server {
	return &Server{impl: impl}
}

func (s Server) StoreKeyValue(ctx context.Context, req *pb.StoreKeyValueRequest) (*emptypb.Empty, error) {
	if err := s.impl.Store(ctx, req.Key, req.Value); err != nil {
		return nil, fmt.Errorf("failed to store bytes for key: %s: %w", req.Key, err)
	}
	return &emptypb.Empty{}, nil
}

func (s Server) GetValueForKey(ctx context.Context, req *pb.GetValueForKeyRequest) (*pb.GetValueForKeyResponse, error) {
	bytes, err := s.impl.Get(ctx, req.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to get bytes for key: %s: %w ", req.Key, err)
	}

	return &pb.GetValueForKeyResponse{Value: bytes}, nil
}
