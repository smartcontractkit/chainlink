package standardcapabilities

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

const defaultStartTimeout = 3 * time.Minute

type standardCapabilities struct {
	services.StateMachine
	log                  logger.Logger
	spec                 *job.StandardCapabilitiesSpec
	pluginRegistrar      plugins.RegistrarConfig
	telemetryService     core.TelemetryService
	store                core.KeyValueStore
	CapabilitiesRegistry core.CapabilitiesRegistry
	errorLog             core.ErrorLog
	pipelineRunner       core.PipelineRunnerService
	relayerSet           core.RelayerSet
	oracleFactory        core.OracleFactory

	capabilitiesLoop *loop.StandardCapabilitiesService

	wg           sync.WaitGroup
	stopChan     services.StopChan
	startTimeout time.Duration
}

func newStandardCapabilities(
	log logger.Logger,
	spec *job.StandardCapabilitiesSpec,
	pluginRegistrar plugins.RegistrarConfig,
	telemetryService core.TelemetryService,
	store core.KeyValueStore,
	CapabilitiesRegistry core.CapabilitiesRegistry,
	errorLog core.ErrorLog,
	pipelineRunner core.PipelineRunnerService,
	relayerSet core.RelayerSet,
	oracleFactory core.OracleFactory,
) *standardCapabilities {
	return &standardCapabilities{
		log:                  log,
		spec:                 spec,
		pluginRegistrar:      pluginRegistrar,
		telemetryService:     telemetryService,
		store:                store,
		CapabilitiesRegistry: CapabilitiesRegistry,
		errorLog:             errorLog,
		pipelineRunner:       pipelineRunner,
		relayerSet:           relayerSet,
		oracleFactory:        oracleFactory,
		stopChan:             make(chan struct{}),
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

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()

			if s.startTimeout == 0 {
				s.startTimeout = defaultStartTimeout
			}

			cctx, cancel := s.stopChan.CtxWithTimeout(s.startTimeout)
			defer cancel()

			if err = s.capabilitiesLoop.WaitCtx(cctx); err != nil {
				s.log.Errorf("error waiting for standard capabilities service to start: %v", err)
				return
			}

			if err = s.capabilitiesLoop.Service.Initialise(cctx, s.spec.Config, s.telemetryService, s.store, s.CapabilitiesRegistry, s.errorLog,
				s.pipelineRunner, s.relayerSet, s.oracleFactory); err != nil {
				s.log.Errorf("error initialising standard capabilities service: %v", err)
				return
			}

			capabilityInfos, err := s.capabilitiesLoop.Service.Infos(cctx)
			if err != nil {
				s.log.Errorf("error getting standard capabilities service info: %v", err)
				return
			}

			s.log.Info("Started standard capabilities for job spec", "spec", s.spec, "capabilities", capabilityInfos)
		}()

		return nil
	})
}

func (s *standardCapabilities) Close() error {
	close(s.stopChan)
	s.wg.Wait()
	return s.StopOnce("StandardCapabilities", func() error {
		if s.capabilitiesLoop != nil {
			return s.capabilitiesLoop.Close()
		}

		return nil
	})
}
