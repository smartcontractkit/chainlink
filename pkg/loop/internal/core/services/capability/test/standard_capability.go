package test

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

type StandardCapabilityService struct {
}

func (t StandardCapabilityService) Info(ctx context.Context) (capabilities.CapabilityInfo, error) {
	return capabilities.CapabilityInfo{
		ID:             "1",
		CapabilityType: capabilities.CapabilityTypeAction,
		Description:    "",
		Version:        "",
		DON:            nil,
	}, nil
}

func (t StandardCapabilityService) Start(ctx context.Context) error { return nil }

func (t StandardCapabilityService) Close() error {
	//TODO implement me
	return nil
}

func (t StandardCapabilityService) Ready() error {
	//TODO implement me
	return nil
}

func (t StandardCapabilityService) HealthReport() map[string]error {
	//TODO implement me
	return nil
}

func (t StandardCapabilityService) Name() string { return "StandardCapability" }

func (t StandardCapabilityService) Initialise(ctx context.Context, config string, telemetryService core.TelemetryService, store core.KeyValueStore,
	capabilityRegistry core.CapabilitiesRegistry, errorLog core.ErrorLog,
	pipelineRunner core.PipelineRunnerService, relayerSet core.RelayerSet) error {
	return nil
}
