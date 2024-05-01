package standardcapability

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

type StandardCapability struct {
	log                logger.Logger
	spec               *job.StandardCapabilitySpec
	pluginRegistrar    plugins.RegistrarConfig
	telemetryService   core.TelemetryService
	store              core.KeyValueStore
	capabilityRegistry core.CapabilitiesRegistry
	errorLog           core.ErrorLog
	pipelineRunner     core.PipelineRunnerService
	relayerSet         core.RelayerSet

	capabilityLoop *loop.StandardCapabilityService
}

func NewStandardCapability(log logger.Logger, spec *job.StandardCapabilitySpec,
	pluginRegistrar plugins.RegistrarConfig,
	telemetryService core.TelemetryService,
	store core.KeyValueStore,
	capabilityRegistry core.CapabilitiesRegistry,
	errorLog core.ErrorLog,
	pipelineRunner core.PipelineRunnerService,
	relayerSet core.RelayerSet) *StandardCapability {
	return &StandardCapability{
		log:                log,
		spec:               spec,
		pluginRegistrar:    pluginRegistrar,
		telemetryService:   telemetryService,
		store:              store,
		capabilityRegistry: capabilityRegistry,
		errorLog:           errorLog,
		pipelineRunner:     pipelineRunner,
		relayerSet:         relayerSet,
	}
}

func (s *StandardCapability) Start(ctx context.Context) error {
	cmdName := s.spec.Command

	cmdFn, opts, err := s.pluginRegistrar.RegisterLOOP(plugins.CmdConfig{
		ID:  s.log.Name(),
		Cmd: cmdName,
		Env: nil,
	})

	if err != nil {
		return fmt.Errorf("error registering loop: %v", err)
	}

	s.capabilityLoop = loop.NewStandardCapability(s.log, opts, cmdFn)

	if err = s.capabilityLoop.Start(ctx); err != nil {
		return fmt.Errorf("error starting standard capability service: %v", err)
	}

	if err = s.capabilityLoop.WaitCtx(ctx); err != nil {
		return fmt.Errorf("error waiting for standard capability service to start: %v", err)
	}

	if err = s.capabilityLoop.Service.Initialise(ctx, s.spec.Config, s.telemetryService, s.store, s.capabilityRegistry, s.errorLog,
		s.pipelineRunner, s.relayerSet); err != nil {
		return fmt.Errorf("error initialising standard capability service: %v", err)
	}

	capabilityInfo, err := s.capabilityLoop.Service.Info(ctx)
	if err != nil {
		return fmt.Errorf("error getting standard capability service info: %v", err)
	}

	s.log.Info("Started standard capability", "info", capabilityInfo)

	return nil
}

func (s *StandardCapability) Close() error {
	if s.capabilityLoop != nil {
		return s.capabilityLoop.Close()
	}

	return nil
}
