package internal

import (
	"context"
	"errors"
	"fmt"

	"github.com/smartcontractkit/libocr/commontypes"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/loop/internal/pb"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
)

var _ types.Telemetry = (*telemetryClient)(nil)
var _ types.MonitoringEndpointGenerator = (*telemetryClient)(nil)
var _ commontypes.MonitoringEndpoint = (*telemetryEndpoint)(nil)

type TelemetryClient struct {
	*telemetryClient
}

type telemetryClient struct {
	grpc pb.TelemetryClient

	lggr logger.Logger
}

type telemetryEndpoint struct {
	lggr logger.Logger

	grpc          pb.TelemetryClient
	relayID       pb.RelayID
	contractID    string
	telemetryType string
}

func (t *telemetryEndpoint) SendLog(log []byte) {
	_, err := t.grpc.Send(context.Background(), &pb.TelemetryMessage{
		RelayID:       &t.relayID,
		ContractID:    t.contractID,
		TelemetryType: t.telemetryType,
		Payload:       log,
	})
	if err != nil {
		t.lggr.Errorw("cannot send telemetry", "err", err)
	}
}

// GenMonitoringEndpoint generates a new monitoring endpoint, returns nil if one cannot be generated
func (t *telemetryClient) GenMonitoringEndpoint(network string, chainID string, contractID string, telemetryType string) commontypes.MonitoringEndpoint {
	if contractID == "" {
		t.lggr.Errorw("cannot generate monitoring endpoint, contractID is empty", "contractID", contractID, "telemetryType", telemetryType, "network", network, "chainID", chainID)
		return nil
	}
	if telemetryType == "" {
		t.lggr.Errorw("cannot generate monitoring endpoint, telemetryType is empty", "contractID", contractID, "telemetryType", telemetryType, "network", network, "chainID", chainID)
		return nil
	}
	if network == "" {
		t.lggr.Errorw("cannot generate monitoring endpoint, network is empty", "contractID", contractID, "telemetryType", telemetryType, "network", network, "chainID", chainID)
		return nil
	}
	if chainID == "" {
		t.lggr.Errorw("cannot generate monitoring endpoint, chainID is empty", "contractID", contractID, "telemetryType", telemetryType, "network", network, "chainID", chainID)
		return nil
	}
	return &telemetryEndpoint{
		grpc: t.grpc,
		relayID: pb.RelayID{
			Network: network,
			ChainId: chainID,
		},
		contractID:    contractID,
		telemetryType: telemetryType,
		lggr:          t.lggr,
	}
}

// Send sends payload to the desired endpoint based on network and chainID
func (t *telemetryClient) Send(ctx context.Context, network string, chainID string, contractID string, telemetryType string, payload []byte) error {
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
	if err != nil {
		return err
	}
	return nil
}

func NewTelemetryClient(cc grpc.ClientConnInterface, lggr logger.Logger) *telemetryClient {
	return &telemetryClient{grpc: pb.NewTelemetryClient(cc), lggr: lggr}
}

var _ pb.TelemetryServer = (*telemetryServer)(nil)

type telemetryServer struct {
	pb.UnimplementedTelemetryServer

	impl      types.MonitoringEndpointGenerator
	endpoints map[string]commontypes.MonitoringEndpoint
}

func (t *telemetryServer) Send(ctx context.Context, message *pb.TelemetryMessage) (*emptypb.Empty, error) {
	e, err := t.getOrCreateEndpoint(message)
	if err != nil {
		return nil, err
	}
	e.SendLog(message.Payload)

	return nil, nil
}

func (t *telemetryServer) getOrCreateEndpoint(m *pb.TelemetryMessage) (commontypes.MonitoringEndpoint, error) {
	if m.ContractID == "" {
		return nil, errors.New("contractID cannot be empty")
	}
	if m.TelemetryType == "" {
		return nil, errors.New("telemetryType cannot be empty")
	}
	if m.RelayID == nil {
		return nil, errors.New("RelayID cannot be nil")
	}
	if m.RelayID.Network == "" {
		return nil, errors.New("RelayID.Network cannot be empty")
	}
	if m.RelayID.ChainId == "" {
		return nil, errors.New("RelayID.ChainId cannot be empty")
	}

	key := makeKey(m)
	e, ok := t.endpoints[key]
	if !ok {
		e = t.impl.GenMonitoringEndpoint(m.RelayID.Network, m.RelayID.ChainId, m.ContractID, m.TelemetryType)
	}
	return e, nil
}

func makeKey(m *pb.TelemetryMessage) string {
	return fmt.Sprintf("%s_%s_%s_%s", m.RelayID.Network, m.RelayID.ChainId, m.ContractID, m.TelemetryType)
}

func NewTelemetryServer(impl types.MonitoringEndpointGenerator) *telemetryServer {
	return &telemetryServer{impl: impl}
}
