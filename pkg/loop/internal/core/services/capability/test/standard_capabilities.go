package test

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

type StandardCapabilitiesService struct {
}

func (t StandardCapabilitiesService) Infos(ctx context.Context) ([]capabilities.CapabilityInfo, error) {
	return []capabilities.CapabilityInfo{
		{
			ID:             "1",
			CapabilityType: capabilities.CapabilityTypeAction,
			Description:    "",
			DON:            nil,
		},
		{
			ID:             "2",
			CapabilityType: capabilities.CapabilityTypeTarget,
			Description:    "",
			DON:            nil,
		},
	}, nil
}

func (t StandardCapabilitiesService) Start(ctx context.Context) error { return nil }

func (t StandardCapabilitiesService) Close() error {
	//TODO implement me
	return nil
}

func (t StandardCapabilitiesService) Ready() error {
	//TODO implement me
	return nil
}

func (t StandardCapabilitiesService) HealthReport() map[string]error {
	//TODO implement me
	return nil
}

func (t StandardCapabilitiesService) Name() string { return "StandardCapabilities" }

func (t StandardCapabilitiesService) Initialise(ctx context.Context, config string, telemetryService core.TelemetryService, store core.KeyValueStore,
	capabilityRegistry core.CapabilitiesRegistry, errorLog core.ErrorLog,
	pipelineRunner core.PipelineRunnerService, relayerSet core.RelayerSet) error {
	return nil
}
