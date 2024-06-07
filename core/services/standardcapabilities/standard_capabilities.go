package standardcapabilities

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

type standardCapabilities struct {
	services.StateMachine
	log                logger.Logger
	spec               *job.StandardCapabilitiesSpec
	pluginRegistrar    plugins.RegistrarConfig
	telemetryService   core.TelemetryService
	store              core.KeyValueStore
	capabilityRegistry core.CapabilitiesRegistry
	errorLog           core.ErrorLog
	pipelineRunner     core.PipelineRunnerService
	relayerSet         core.RelayerSet

	capabilitiesLoop *loop.StandardCapabilitiesService
}

func newStandardCapabilities(log logger.Logger, spec *job.StandardCapabilitiesSpec,
	pluginRegistrar plugins.RegistrarConfig,
	telemetryService core.TelemetryService,
	store core.KeyValueStore,
	capabilityRegistry core.CapabilitiesRegistry,
	errorLog core.ErrorLog,
	pipelineRunner core.PipelineRunnerService,
	relayerSet core.RelayerSet) *standardCapabilities {
	return &standardCapabilities{
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

func (s *standardCapabilities) Start(ctx context.Context) error {
	return s.StartOnce("StandardCapabilities", func() error {
		cmdName := s.spec.Command

		cmdFn, opts, err := s.pluginRegistrar.RegisterLOOP(plugins.CmdConfig{
			ID:  s.log.Name(),
			Cmd: cmdName,
			Env: nil,
		})

		if err != nil {
			return fmt.Errorf("error registering loop: %v", err)
		}

		s.capabilitiesLoop = loop.NewStandardCapabilitiesService(s.log, opts, cmdFn)

		if err = s.capabilitiesLoop.Start(ctx); err != nil {
			return fmt.Errorf("error starting standard capabilities service: %v", err)
		}

		if err = s.capabilitiesLoop.WaitCtx(ctx); err != nil {
			return fmt.Errorf("error waiting for standard capabilities service to start: %v", err)
		}

		if err = s.capabilitiesLoop.Service.Initialise(ctx, s.spec.Config, s.telemetryService, s.store, s.capabilityRegistry, s.errorLog,
			s.pipelineRunner, s.relayerSet); err != nil {
			return fmt.Errorf("error initialising standard capabilities service: %v", err)
		}

		capabilityInfos, err := s.capabilitiesLoop.Service.Infos(ctx)
		if err != nil {
			return fmt.Errorf("error getting standard capabilities service info: %v", err)
		}

		s.log.Info("Started standard capabilities for job spec", "spec", s.spec, "capabilities", capabilityInfos)

		return nil
	})
}

func (s *standardCapabilities) Close() error {
	return s.StopOnce("StandardCapabilities", func() error {
		if s.capabilitiesLoop != nil {
			return s.capabilitiesLoop.Close()
		}

		return nil
	})
}
