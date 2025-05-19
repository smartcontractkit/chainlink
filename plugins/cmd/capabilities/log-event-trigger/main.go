package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-plugin"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/triggers/logevent"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

const (
	serviceName = "LogEventTriggerCapability"
)

type LogEventTriggerGRPCService struct {
	triggerService *logevent.TriggerService
	s              *loop.Server
	config         logevent.Config
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
	err := cs.triggerService.Close()
	if err != nil {
		return fmt.Errorf("error closing trigger service for chainID %s: %w", cs.config.ChainID, err)
	}
	return nil
}

func (cs *LogEventTriggerGRPCService) Ready() error {
	return nil
}

func (cs *LogEventTriggerGRPCService) HealthReport() map[string]error {
	return map[string]error{cs.Name(): nil}
}

func (cs *LogEventTriggerGRPCService) Name() string {
	return cs.s.Logger.Name()
}

func (cs *LogEventTriggerGRPCService) Infos(ctx context.Context) ([]capabilities.CapabilityInfo, error) {
	triggerInfo, err := cs.triggerService.Info(ctx)
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
	oracleFactory core.OracleFactory,
) error {
	cs.s.Logger.Debugf("Initialising %s", serviceName)

	var logEventConfig logevent.Config
	err := json.Unmarshal([]byte(config), &logEventConfig)
	if err != nil {
		return fmt.Errorf("error decoding log_event_trigger config: %w", err)
	}

	relayID := types.NewRelayID(logEventConfig.Network, logEventConfig.ChainID)
	relayer, err := relayerSet.Get(ctx, relayID)
	if err != nil {
		return fmt.Errorf("error fetching relayer for chainID %s from relayerSet: %w", logEventConfig.ChainID, err)
	}

	// Set relayer and trigger in LogEventTriggerGRPCService
	cs.config = logEventConfig
	triggerService, err := logevent.NewTriggerService(ctx, cs.s.Logger, relayer, logEventConfig)
	if err != nil {
		return fmt.Errorf("error creating trigger service for chainID %s: %w", logEventConfig.ChainID, err)
	}
	err = triggerService.Start(ctx)
	if err != nil {
		return fmt.Errorf("error starting trigger service for chainID %s: %w", logEventConfig.ChainID, err)
	}
	cs.triggerService = triggerService

	if err := capabilityRegistry.Add(ctx, cs.triggerService); err != nil {
		return fmt.Errorf("error when adding cron trigger to the registry: %w", err)
	}

	return nil
}
