package validation

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/goplugin"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type validationServiceClient struct {
	*net.BrokerExt
	*goplugin.ServiceClient
	grpc pb.ValidationServiceClient
}

func (v *validationServiceClient) ValidateConfig(ctx context.Context, config map[string]interface{}) error {
	pbConfig, err := structpb.NewStruct(config)
	if err != nil {
		return err
	}
	_, err = v.grpc.ValidateConfig(ctx, &pb.ValidateConfigRequest{
		Config: pbConfig,
	})
	return err
}

func NewValidationServiceClient(b *net.BrokerExt, cc grpc.ClientConnInterface) *validationServiceClient {
	return &validationServiceClient{b.WithName("ReportingPluginProviderClient"), goplugin.NewServiceClient(b, cc), pb.NewValidationServiceClient(cc)}
}

type validationServiceServer struct {
	pb.UnimplementedValidationServiceServer

	*net.BrokerExt

	impl types.ValidationServiceServer
}

func (v *validationServiceServer) ValidateConfig(ctx context.Context, c *pb.ValidateConfigRequest) (*pb.ValidateConfigResponse, error) {
	err := v.impl.ValidateConfig(ctx, c.Config.AsMap())
	if err != nil {
		return nil, err
	}
	return &pb.ValidateConfigResponse{}, nil
}

func NewValidationServiceServer(impl types.ValidationServiceServer, b *net.BrokerExt) *validationServiceServer {
	return &validationServiceServer{impl: impl, BrokerExt: b.WithName("ReportingPluginFactoryServer")}
}
