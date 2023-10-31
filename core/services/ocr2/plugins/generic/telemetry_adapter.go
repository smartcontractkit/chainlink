package generic

import (
	"context"
	"errors"

	"github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink-relay/pkg/types"
)

var _ types.TelemetryService = (*TelemetryAdapter)(nil)

type TelemetryAdapter struct {
	endpointGenerator types.MonitoringEndpointGenerator
	endpoints         map[[4]string]commontypes.MonitoringEndpoint
}

func NewTelemetryAdapter(endpointGen types.MonitoringEndpointGenerator) *TelemetryAdapter {
	return &TelemetryAdapter{
		endpoints:         make(map[[4]string]commontypes.MonitoringEndpoint),
		endpointGenerator: endpointGen,
	}
}

func (t *TelemetryAdapter) Send(ctx context.Context, network string, chainID string, contractID string, telemetryType string, payload []byte) error {
	e, err := t.getOrCreateEndpoint(contractID, telemetryType, network, chainID)
	if err != nil {
		return err
	}
	e.SendLog(payload)
	return nil
}

func (t *TelemetryAdapter) getOrCreateEndpoint(contractID string, telemetryType string, network string, chainID string) (commontypes.MonitoringEndpoint, error) {
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
		return nil, errors.New("chainID cannot be empty")
	}

	key := [4]string{network, chainID, contractID, telemetryType}
	e, ok := t.endpoints[key]
	if !ok {
		e = t.endpointGenerator.GenMonitoringEndpoint(network, chainID, contractID, telemetryType)
		t.endpoints[key] = e
	}
	return e, nil
}
