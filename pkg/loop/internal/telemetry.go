package internal

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

var _ types.TelemetryService = (*telemetryServiceClient)(nil)
var _ types.TelemetryClient = (*telemetryClient)(nil)

type telemetryEndpoint struct {
	client        types.TelemetryService
	network       string
	chainID       string
	telemetryType string
	contractID    string
}

func (t *telemetryEndpoint) SendLog(ctx context.Context, log []byte) error {
	return t.client.Send(ctx, t.network, t.chainID, t.contractID, t.telemetryType, log)
}

func NewTelemetryClient(client types.TelemetryService) *telemetryClient {
	return &telemetryClient{TelemetryService: client}
}

type telemetryClient struct {
	types.TelemetryService
}

// NewEndpoint generates a new monitoring endpoint, returns nil if one cannot be generated
func (t *telemetryClient) NewEndpoint(ctx context.Context, network string, chainID string, contractID string, telemetryType string) (types.TelemetryClientEndpoint, error) {
	if contractID == "" {
		return nil, errors.New("contractID cannot be empty")
	}
	if telemetryType == "" {
		return nil, errors.New("telemetryType cannot be empty")
	}
	if network == "" {
		return nil, errors.New("network cannot be empty")
	}
	if chainID == "" {
		return nil, errors.New("chainId cannot be empty")
	}

	return &telemetryEndpoint{
		client:        t.TelemetryService,
		network:       network,
		chainID:       chainID,
		contractID:    contractID,
		telemetryType: telemetryType,
	}, nil
}

type telemetryServiceClient struct {
	grpc pb.TelemetryClient
}

// Send sends payload to the desired endpoint based on network and chainID
func (t *telemetryServiceClient) Send(ctx context.Context, network string, chainID string, contractID string, telemetryType string, payload []byte) error {
	if contractID == "" {
		return errors.New("contractID cannot be empty")
	}
	if telemetryType == "" {
		return errors.New("telemetryType cannot be empty")
	}
	if network == "" {
		return errors.New("network cannot be empty")
	}
	if chainID == "" {
		return errors.New("chainId cannot be empty")
	}
	if len(payload) == 0 {
		return errors.New("payload cannot be empty")
	}
	_, err := t.grpc.Send(ctx, &pb.TelemetryMessage{
		RelayID: &pb.RelayID{
			Network: network,
			ChainId: chainID,
		},
		ContractID:    contractID,
		TelemetryType: telemetryType,
		Payload:       payload,
	})
	return err
}

func NewTelemetryServiceClient(cc grpc.ClientConnInterface) *telemetryServiceClient {
	return &telemetryServiceClient{grpc: pb.NewTelemetryClient(cc)}
}

var _ pb.TelemetryServer = (*telemetryServer)(nil)

type telemetryServer struct {
	pb.UnimplementedTelemetryServer

	impl types.TelemetryService
}

func (t *telemetryServer) Send(ctx context.Context, message *pb.TelemetryMessage) (*emptypb.Empty, error) {
	var network, chainID string
	if message.RelayID != nil {
		network = message.RelayID.Network
		chainID = message.RelayID.ChainId
	}
	err := t.impl.Send(ctx, network, chainID, message.ContractID, message.TelemetryType, message.Payload)
	return &emptypb.Empty{}, err
}

func NewTelemetryServer(impl types.TelemetryService) *telemetryServer {
	return &telemetryServer{impl: impl}
}
