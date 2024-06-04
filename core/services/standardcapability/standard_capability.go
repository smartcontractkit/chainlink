package standardcapability

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

type standardCapability struct {
	services.StateMachine
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

func newStandardCapability(log logger.Logger, spec *job.StandardCapabilitySpec,
	pluginRegistrar plugins.RegistrarConfig,
	telemetryService core.TelemetryService,
	store core.KeyValueStore,
	capabilityRegistry core.CapabilitiesRegistry,
	errorLog core.ErrorLog,
	pipelineRunner core.PipelineRunnerService,
	relayerSet core.RelayerSet) *standardCapability {
	return &standardCapability{
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

func (s *standardCapability) Start(ctx context.Context) error {
	return s.StartOnce("StandardCapability", func() error {
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
	})
}

func (s *standardCapability) Close() error {
	return s.StopOnce("StandardCapability", func() error {
		if s.capabilityLoop != nil {
			return s.capabilityLoop.Close()
		}

		return nil
	})
}
