package standardcapabilities

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

const defaultStartTimeout = 3 * time.Minute

var (
	ErrServiceStopped  = errors.New("service stopped")
	ErrServiceNotReady = errors.New("service not ready")
)

type StandardCapabilities struct {
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
	keystore             core.Keystore
	oracleFactory        core.OracleFactory
	gatewayConnector     core.GatewayConnector

	capabilitiesLoop *loop.StandardCapabilitiesService

	wg           sync.WaitGroup
	readyChan    chan struct{}
	stopChan     services.StopChan
	startTimeout time.Duration
}

func NewStandardCapabilities(
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
	gatewayConnector core.GatewayConnector,
	keystore core.Keystore,
) *StandardCapabilities {
	return &StandardCapabilities{
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
		gatewayConnector:     gatewayConnector,
		keystore:             keystore,
		stopChan:             make(chan struct{}),
		readyChan:            make(chan struct{}),
	}
}

func (s *StandardCapabilities) Start(ctx context.Context) error {
	return s.StartOnce("StandardCapabilities", func() error {
		cmdName := s.spec.Command

		cmdFn, opts, err := s.pluginRegistrar.RegisterLOOP(plugins.CmdConfig{
			ID:  s.log.Name(),
			Cmd: cmdName,
			Env: nil,
		})
		if err != nil {
			return fmt.Errorf("error registering loop: %w", err)
		}

		s.capabilitiesLoop = loop.NewStandardCapabilitiesService(s.log, opts, cmdFn)
		if err = s.capabilitiesLoop.Start(ctx); err != nil {
			return fmt.Errorf("error starting standard capabilities service: %w", err)
		}

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			defer close(s.readyChan)

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
				s.pipelineRunner, s.relayerSet, s.oracleFactory, s.gatewayConnector, s.keystore); err != nil {
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

// Ready is a non-blocking check for the service's ready state.  Errors if not
// ready when called.
func (s *StandardCapabilities) Ready() error {
	if err := s.StateMachine.Ready(); err != nil {
		return err
	}
	select {
	case <-s.readyChan:
		return nil
	case <-s.stopChan:
		return ErrServiceStopped
	default:
		return ErrServiceNotReady
	}
}

// Await waits for the service to be ready or for the context to be cancelled.
func (s *StandardCapabilities) Await(ctx context.Context) error {
	select {
	case <-s.readyChan:
		return nil
	case <-s.stopChan:
		return ErrServiceStopped
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *StandardCapabilities) Close() error {
	close(s.stopChan)
	s.wg.Wait()
	return s.StopOnce("StandardCapabilities", func() error {
		if s.capabilitiesLoop != nil {
			return s.capabilitiesLoop.Close()
		}

		return nil
	})
}
