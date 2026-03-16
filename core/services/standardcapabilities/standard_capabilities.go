package standardcapabilities

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/services/orgresolver"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
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
	command              string
	config               string
	pluginRegistrar      plugins.RegistrarConfig
	store                core.KeyValueStore
	CapabilitiesRegistry core.CapabilitiesRegistry
	relayerSet           core.RelayerSet
	keystore             core.Keystore
	oracleFactory        core.OracleFactory
	gatewayConnector     core.GatewayConnector
	orgResolver          orgresolver.OrgResolver
	creSettings          core.SettingsBroadcaster
	triggerEventStore    capabilities.EventStore

	capabilitiesLoop *loop.StandardCapabilitiesService

	wg           sync.WaitGroup
	readyChan    chan struct{}
	stopChan     services.StopChan
	startTimeout time.Duration
}

func NewStandardCapabilities(
	log logger.Logger,
	command string,
	configJSON string,
	pluginRegistrar plugins.RegistrarConfig,
	dependencies core.StandardCapabilitiesDependencies,
) *StandardCapabilities {
	return &StandardCapabilities{
		log:                  log,
		command:              command,
		config:               configJSON,
		pluginRegistrar:      pluginRegistrar,
		store:                dependencies.Store,
		CapabilitiesRegistry: dependencies.CapabilityRegistry,
		relayerSet:           dependencies.RelayerSet,
		oracleFactory:        dependencies.OracleFactory,
		gatewayConnector:     dependencies.GatewayConnector,
		keystore:             dependencies.P2PKeystore,
		orgResolver:          dependencies.OrgResolver,
		creSettings:          dependencies.CRESettings,
		triggerEventStore:    dependencies.TriggerEventStore,
		stopChan:             make(chan struct{}),
		readyChan:            make(chan struct{}),
	}
}

func (s *StandardCapabilities) Start(ctx context.Context) error {
	return s.StartOnce("StandardCapabilities", func() error {
		cmdFn, opts, err := s.pluginRegistrar.RegisterLOOP(plugins.CmdConfig{
			ID:  s.log.Name(),
			Cmd: s.command,
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

			dependencies := core.StandardCapabilitiesDependencies{
				Config:             s.config,
				Store:              s.store,
				CapabilityRegistry: s.CapabilitiesRegistry,
				RelayerSet:         s.relayerSet,
				OracleFactory:      s.oracleFactory,
				GatewayConnector:   s.gatewayConnector,
				P2PKeystore:        s.keystore,
				OrgResolver:        s.orgResolver,
				CRESettings:        s.creSettings,
				TriggerEventStore:  s.triggerEventStore,
			}
			if err = s.capabilitiesLoop.Service.Initialise(cctx, dependencies); err != nil {
				s.log.Errorf("error initialising standard capabilities service: %v", err)
				return
			}

			capabilityInfos, err := s.capabilitiesLoop.Service.Infos(cctx)
			if err != nil {
				s.log.Errorf("error getting standard capabilities service info: %v", err)
				return
			}

			s.log.Info("Started standard capabilities", "command", s.command, "capabilities", capabilityInfos)
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
