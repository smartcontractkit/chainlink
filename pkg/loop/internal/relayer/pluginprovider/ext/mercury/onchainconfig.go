package mercury

import (
	"context"

	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	mercury_pb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/mercury"
	mercury_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury"
)

var _ mercury_types.OnchainConfigCodec = (*onchainConfigCodecClient)(nil)

type onchainConfigCodecClient struct {
	grpc mercury_pb.OnchainConfigCodecClient
}

func newOnchainConfigCodecClient(cc grpc.ClientConnInterface) *onchainConfigCodecClient {
	return &onchainConfigCodecClient{grpc: mercury_pb.NewOnchainConfigCodecClient(cc)}
}

func (o *onchainConfigCodecClient) Encode(ctx context.Context, config mercury_types.OnchainConfig) ([]byte, error) {
	reply, err := o.grpc.Encode(ctx, &mercury_pb.EncodeOnchainConfigRequest{
		OnchainConfig: pbOnchainConfig(config),
	})
	if err != nil {
		return nil, err
	}
	return reply.OnchainConfig, nil
}

func (o *onchainConfigCodecClient) Decode(ctx context.Context, data []byte) (mercury_types.OnchainConfig, error) {
	reply, err := o.grpc.Decode(ctx, &mercury_pb.DecodeOnchainConfigRequest{
		OnchainConfig: data,
	})
	if err != nil {
		return mercury_types.OnchainConfig{}, err
	}
	return onchainConfig(reply.OnchainConfig), nil
}

func pbOnchainConfig(config mercury_types.OnchainConfig) *mercury_pb.OnchainConfig {
	return &mercury_pb.OnchainConfig{
		Min: pb.NewBigIntFromInt(config.Min),
		Max: pb.NewBigIntFromInt(config.Max),
	}
}

func onchainConfig(config *mercury_pb.OnchainConfig) mercury_types.OnchainConfig {
	return mercury_types.OnchainConfig{
		Min: config.Min.Int(),
		Max: config.Max.Int(),
	}
}

var _ mercury_pb.OnchainConfigCodecServer = (*onchainConfigCodecServer)(nil)

type onchainConfigCodecServer struct {
	mercury_pb.UnimplementedOnchainConfigCodecServer

	impl mercury_types.OnchainConfigCodec
}

func newOnchainConfigCodecServer(impl mercury_types.OnchainConfigCodec) *onchainConfigCodecServer {
	return &onchainConfigCodecServer{impl: impl}
}

func (o *onchainConfigCodecServer) Encode(ctx context.Context, request *mercury_pb.EncodeOnchainConfigRequest) (*mercury_pb.EncodeOnchainConfigReply, error) {
	val, err := o.impl.Encode(ctx, onchainConfig(request.OnchainConfig))
	if err != nil {
		return nil, err
	}
	return &mercury_pb.EncodeOnchainConfigReply{OnchainConfig: val}, nil
}

func (o *onchainConfigCodecServer) Decode(ctx context.Context, request *mercury_pb.DecodeOnchainConfigRequest) (*mercury_pb.DecodeOnchainConfigReply, error) {
	val, err := o.impl.Decode(ctx, request.OnchainConfig)
	if err != nil {
		return nil, err
	}
	return &mercury_pb.DecodeOnchainConfigReply{OnchainConfig: pbOnchainConfig(val)}, nil
}
