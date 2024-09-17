package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-plugin"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/triggers/logeventtrigger"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

const (
	serviceName = "LogEventTriggerCapability"
)

type LogEventTriggerGRPCService struct {
	trigger capabilities.TriggerCapability
	s       *loop.Server
	config  logeventtrigger.LogEventConfig
}

func main() {
	s := loop.MustNewStartedServer(serviceName)
	defer s.Stop()

	s.Logger.Infof("Starting %s", serviceName)

	stopCh := make(chan struct{})
	defer close(stopCh)

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: loop.StandardCapabilitiesHandshakeConfig(),
		Plugins: map[string]plugin.Plugin{
			loop.PluginStandardCapabilitiesName: &loop.StandardCapabilitiesLoop{
				PluginServer: &LogEventTriggerGRPCService{
					s: s,
				},
				BrokerConfig: loop.BrokerConfig{Logger: s.Logger, StopCh: stopCh, GRPCOpts: s.GRPCOpts},
			},
		},
		GRPCServer: s.GRPCOpts.NewServer,
	})
}

func (cs *LogEventTriggerGRPCService) Start(ctx context.Context) error {
	return nil
}

func (cs *LogEventTriggerGRPCService) Close() error {
	return nil
}

func (cs *LogEventTriggerGRPCService) Ready() error {
	return nil
}

func (cs *LogEventTriggerGRPCService) HealthReport() map[string]error {
	return nil
}

func (cs *LogEventTriggerGRPCService) Name() string {
	return serviceName
}

func (cs *LogEventTriggerGRPCService) Infos(ctx context.Context) ([]capabilities.CapabilityInfo, error) {
	triggerInfo, err := cs.trigger.Info(ctx)
	if err != nil {
		return nil, err
	}

	return []capabilities.CapabilityInfo{
		triggerInfo,
	}, nil
}

func (cs *LogEventTriggerGRPCService) Initialise(
	ctx context.Context,
	config string,
	telemetryService core.TelemetryService,
	store core.KeyValueStore,
	capabilityRegistry core.CapabilitiesRegistry,
	errorLog core.ErrorLog,
	pipelineRunner core.PipelineRunnerService,
	relayerSet core.RelayerSet,
) error {
	cs.s.Logger.Debugf("Initialising %s", serviceName)

	var logEventConfig logeventtrigger.LogEventConfig
	err := json.Unmarshal([]byte(config), &logEventConfig)
	if err != nil {
		return fmt.Errorf("error decoding log_event_trigger config: %v", err)
	}

	relayID := types.NewRelayID(logEventConfig.Network, fmt.Sprintf("%d", logEventConfig.ChainId))
	relayer, err := relayerSet.Get(ctx, relayID)
	if err != nil {
		return fmt.Errorf("error fetching relayer for chainID %d from relayerSet: %v", logEventConfig.ChainId, err)
	}

	// Set relayer and trigger in LogEventTriggerGRPCService
	cs.config = logEventConfig
	cs.trigger = logeventtrigger.NewLogEventTriggerService(logeventtrigger.Params{
		Logger:         cs.s.Logger,
		Relayer:        relayer,
		LogEventConfig: logEventConfig,
	})

	if err := capabilityRegistry.Add(ctx, cs.trigger); err != nil {
		return fmt.Errorf("error when adding cron trigger to the registry: %w", err)
	}

	return nil
}
