package keyvalue

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
)

func Test_KeyValueStoreClient(t *testing.T) {
	// Setup
	client := Client{grpc: &testGrpcClient{store: make(map[string][]byte)}}
	key := "key"
	insertedVal := "aval"

	err := client.Store(context.Background(), key, []byte(insertedVal))
	assert.NoError(t, err)

	retrievedVal, err := client.Get(context.Background(), key)
	assert.NoError(t, err)
	assert.Equal(t, insertedVal, string(retrievedVal))
}

func Test_KeyValueStoreServer(t *testing.T) {
	// Setup
	server := Server{impl: &testKeyValueStore{store: make(map[string][]byte)}}

	_, err := server.StoreKeyValue(context.Background(), &pb.StoreKeyValueRequest{Key: "key", Value: []byte(`{"A":"a","B":1}`)})
	assert.NoError(t, err)
	resp, err := server.GetValueForKey(context.Background(), &pb.GetValueForKeyRequest{Key: "key"})
	assert.NoError(t, err)
	assert.Equal(t, []byte(`{"A":"a","B":1}`), resp.Value)
}

type testGrpcClient struct {
	store map[string][]byte
}

func (t *testGrpcClient) StoreKeyValue(ctx context.Context, in *pb.StoreKeyValueRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	t.store[in.Key] = in.Value
	return &emptypb.Empty{}, nil
}

func (t *testGrpcClient) GetValueForKey(ctx context.Context, in *pb.GetValueForKeyRequest, opts ...grpc.CallOption) (*pb.GetValueForKeyResponse, error) {
	return &pb.GetValueForKeyResponse{Value: t.store[in.Key]}, nil
}

type testKeyValueStore struct {
	store map[string][]byte
}

func (t *testKeyValueStore) Store(ctx context.Context, key string, val []byte) error {
	t.store[key] = val
	return nil
}

func (t *testKeyValueStore) Get(ctx context.Context, key string) ([]byte, error) {
	return t.store[key], nil
}
