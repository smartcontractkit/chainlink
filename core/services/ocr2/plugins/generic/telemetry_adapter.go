package generic

import (
	"context"
	"errors"

	"github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"

	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

var _ core.TelemetryService = (*TelemetryAdapter)(nil)

type TelemetryAdapter struct {
	endpointGenerator telemetry.MonitoringEndpointGenerator
	endpoints         map[[4]string]commontypes.MonitoringEndpoint
}

func NewTelemetryAdapter(endpointGen telemetry.MonitoringEndpointGenerator) *TelemetryAdapter {
	return &TelemetryAdapter{
		endpoints:         make(map[[4]string]commontypes.MonitoringEndpoint),
		endpointGenerator: endpointGen,
	}
}

func (t *TelemetryAdapter) Send(ctx context.Context, network string, chainID string, contractID string, telemetryType string, payload []byte) error {
	e, err := t.getOrCreateEndpoint(network, chainID, contractID, telemetryType)
	if err != nil {
		return err
	}
	e.SendLog(payload)
	return nil
}

func (t *TelemetryAdapter) getOrCreateEndpoint(network string, chainID string, contractID string, telemetryType string) (commontypes.MonitoringEndpoint, error) {
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
		e = t.endpointGenerator.GenMonitoringEndpoint(network, chainID, contractID, synchronization.TelemetryType(telemetryType))
		t.endpoints[key] = e
	}
	return e, nil
}
