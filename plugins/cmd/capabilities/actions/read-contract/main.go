package read_contract

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-plugin"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/actions"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/logeventtrigger"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

const (
	serviceName = "ReadContractCapability"
)

type ReadContractGRPCService struct {
	action capabilities.ActionCapability
	s      *loop.Server
	config actions.ReadContractConfig
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
				PluginServer: &ReadContractGRPCService{
					s: s,
				},
				BrokerConfig: loop.BrokerConfig{Logger: s.Logger, StopCh: stopCh, GRPCOpts: s.GRPCOpts},
			},
		},
		GRPCServer: s.GRPCOpts.NewServer,
	})
}

func (cs *ReadContractGRPCService) Start(ctx context.Context) error {
	return nil
}

func (cs *ReadContractGRPCService) Close() error {
	return nil
}

func (cs *ReadContractGRPCService) Ready() error {
	return nil
}

func (cs *ReadContractGRPCService) HealthReport() map[string]error {
	return nil
}

func (cs *ReadContractGRPCService) Name() string {
	return serviceName
}

func (cs *ReadContractGRPCService) Infos(ctx context.Context) ([]capabilities.CapabilityInfo, error) {
	triggerInfo, err := cs.action.Info(ctx)
	if err != nil {
		return nil, err
	}

	return []capabilities.CapabilityInfo{
		triggerInfo,
	}, nil
}

func (cs *ReadContractGRPCService) Initialise(
	ctx context.Context,
	config string,
	telemetryService core.TelemetryService,
	store core.KeyValueStore,
	capabilityRegistry core.CapabilitiesRegistry,
	errorLog core.ErrorLog,
	pipelineRunner core.PipelineRunnerService,
	relayerSet core.RelayerSet,
) error {
	cs.s.Logger.Infof("Initialising %s", serviceName)

	var readContractConfig actions.ReadContractConfig
	err := json.Unmarshal([]byte(config), &readContractConfig)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	relayID := types.NewRelayID(readContractConfig.Network, fmt.Sprintf("%d", readContractConfig.ChainId))
	relayer, err := relayerSet.Get(ctx, relayID)
	if err != nil {
		return fmt.Errorf("failed to fetch relayer for chainID %d from relayerSet: %w", readContractConfig.ChainId, err)
	}

	// Set relayer and trigger in ReadContractGRPCService
	cs.config = readContractConfig

	relayer.NewContractReader()

	cs.action = actions.NewReadContractAction(logeventtrigger.Params{
		Logger:         cs.s.Logger,
		Relayer:        relayer,
		LogEventConfig: readContractConfig,
	})

	if err := capabilityRegistry.Add(ctx, cs.trigger); err != nil {
		return fmt.Errorf("error when adding cron trigger to the registry: %w", err)
	}

	return nil
}
