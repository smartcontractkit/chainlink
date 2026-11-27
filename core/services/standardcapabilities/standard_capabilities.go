package standardcapabilities

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/jpillora/backoff"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/services/orgresolver"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/services/standardcapabilities/conversions"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

const (
	defaultStartTimeout = 3 * time.Minute

	defaultDonIDResolveTimeout = 3 * time.Minute
	defaultDonIDRetryInterval  = time.Second

	defaultStartRetryMinDelay = 5 * time.Second
	defaultStartRetryMaxDelay = time.Minute
)

var (
	ErrServiceStopped  = errors.New("service stopped")
	ErrServiceNotReady = errors.New("service not ready")
)

// pluginLoop is the supervised plugin process of a standard capability: the service
// that chainlink-common defines for a standard capability, and the wait for its process
// to be up, which no chainlink-common interface has.
type pluginLoop interface {
	loop.StandardCapabilities
	WaitCtx(context.Context) error
}

var _ pluginLoop = standardCapabilitiesLoop{}

// standardCapabilitiesLoop is the plugin supervisor of chainlink-common with the two
// calls of the plugin added to it. The process must be up before they are called.
type standardCapabilitiesLoop struct {
	*loop.StandardCapabilitiesService
}

func (l standardCapabilitiesLoop) Initialise(ctx context.Context, dependencies core.StandardCapabilitiesDependencies) error {
	return l.Service.Initialise(ctx, dependencies)
}

func (l standardCapabilitiesLoop) Infos(ctx context.Context) ([]capabilities.CapabilityInfo, error) {
	return l.Service.Infos(ctx)
}

type StandardCapabilities struct {
	services.StateMachine
	log                  logger.Logger
	command              string
	capabilityID         string
	config               string
	pluginRegistrar      plugins.RegistrarConfig
	store                core.KeyValueStore
	CapabilitiesRegistry registry.CapabilitiesRegistry
	relayerSet           core.RelayerSet
	keystore             core.Keystore
	oracleFactory        core.OracleFactory
	gatewayConnector     core.GatewayConnector
	orgResolver          orgresolver.OrgResolver
	creSettings          settings.Broadcaster
	triggerEventStore    capabilities.EventStore
	capDonIDResolver     CapDonIDResolver
	// shardIndex is this node's shard from the [Sharding] config, plumbed to the LOOP at Initialise.
	shardIndex uint32

	// capabilitiesLoop is nil between a failed start attempt and the next one.
	// It is written by the startup goroutine, and by Close only after that goroutine has exited.
	capabilitiesLoop pluginLoop
	newPluginLoop    func() pluginLoop
	// loopID is the ID this service holds in the plugin registrar. It is set
	// once RegisterLOOP succeeds and cleared by unregisterLOOP, so the same ID can be
	// registered again when the capability is restarted (e.g. on a config change).
	loopID string

	// registered is CapabilitiesRegistry, with a record of what the plugin added to
	// it. Nil when there is no registry.
	registered            *registeredCapabilities
	registryRemoveTimeout time.Duration

	wg sync.WaitGroup
	// readyChan is closed when the first start attempt has finished.
	readyChan          chan struct{}
	stopChan           services.StopChan
	startTimeout       time.Duration
	startRetryMinDelay time.Duration
	startRetryMaxDelay time.Duration
	startFailures      metric.Int64Counter

	// readyErr records the outcome of the latest start attempt of the
	// asynchronous startup goroutine (WaitCtx/Initialise/Infos). It is nil once
	// startup has succeeded, and holds the last error otherwise. Guarded by
	// readyMu because it is written by the startup goroutine and read by the
	// health checker via Ready/HealthReport on a different goroutine.
	readyMu  sync.RWMutex
	readyErr error
}

func (s *StandardCapabilities) setReadyErr(err error) {
	s.readyMu.Lock()
	defer s.readyMu.Unlock()
	s.readyErr = err
}

func (s *StandardCapabilities) getReadyErr() error {
	s.readyMu.RLock()
	defer s.readyMu.RUnlock()
	return s.readyErr
}

func NewStandardCapabilities(
	log logger.Logger,
	command string,
	configJSON string,
	pluginRegistrar plugins.RegistrarConfig,
	capDonIDResolver CapDonIDResolver,
	dependencies core.StandardCapabilitiesDependencies,
) *StandardCapabilities {
	// The plugin is given the registry through a wrapper that records what it adds, so
	// that Close can remove it.
	capabilitiesRegistry := dependencies.CapabilityRegistry
	var registered *registeredCapabilities
	if capabilitiesRegistry != nil {
		registered = trackRegisteredCapabilities(capabilitiesRegistry)
		capabilitiesRegistry = registered
	}
	return &StandardCapabilities{
		log:                  log,
		command:              command,
		capabilityID:         conversions.GetCapabilityIDFromCommand(command, configJSON),
		config:               configJSON,
		pluginRegistrar:      pluginRegistrar,
		store:                dependencies.Store,
		CapabilitiesRegistry: capabilitiesRegistry,
		registered:           registered,
		relayerSet:           dependencies.RelayerSet,
		oracleFactory:        dependencies.OracleFactory,
		gatewayConnector:     dependencies.GatewayConnector,
		keystore:             dependencies.P2PKeystore,
		orgResolver:          dependencies.OrgResolver,
		creSettings:          dependencies.CRESettings,
		triggerEventStore:    dependencies.TriggerEventStore,
		capDonIDResolver:     capDonIDResolver,
		shardIndex:           dependencies.ShardIndex,
		stopChan:             make(chan struct{}),
		readyChan:            make(chan struct{}),
	}
}

// initialiseDependencies builds the StandardCapabilitiesDependencies delivered to
// the capability LOOP via Initialise.
func (s *StandardCapabilities) initialiseDependencies(capDonID uint32) core.StandardCapabilitiesDependencies {
	return core.StandardCapabilitiesDependencies{
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
		CapabilityDonID:    capDonID,
		ShardIndex:         s.shardIndex,
	}
}

func (s *StandardCapabilities) Start(ctx context.Context) error {
	return s.StartOnce("StandardCapabilities", func() error {
		startFailures, err := newStartFailuresCounter()
		if err != nil {
			return fmt.Errorf("failed to create start failures counter: %w", err)
		}
		s.startFailures = startFailures

		envVars, err := plugins.ParseEnvFile(env.CapabilitiesPlugin.Env.Get())
		if err != nil {
			return fmt.Errorf("failed to parse capabilities env file: %w", err)
		}
		loopID := s.log.Name()
		cmdFn, opts, err := s.pluginRegistrar.RegisterLOOP(plugins.CmdConfig{
			ID:  loopID,
			Cmd: s.command,
			Env: envVars,
		})
		if err != nil {
			return fmt.Errorf("error registering loop: %w", err)
		}
		s.loopID = loopID

		if s.newPluginLoop == nil {
			s.newPluginLoop = func() pluginLoop {
				return standardCapabilitiesLoop{loop.NewStandardCapabilitiesService(s.log, opts, cmdFn)}
			}
		}
		if err = s.launchPluginLoop(ctx); err != nil {
			s.unregisterLOOP()
			return fmt.Errorf("error starting standard capabilities service: %w", err)
		}

		s.wg.Go(s.startUntilReady)

		return nil
	})
}

// startUntilReady retries the asynchronous part of the start until it succeeds
// or the service is closed. A registry entry, the database or an RPC can be
// unavailable for longer than any fixed wait, so there is no deadline.
func (s *StandardCapabilities) startUntilReady() {
	var firstAttempt sync.Once
	finishFirstAttempt := func() { firstAttempt.Do(func() { close(s.readyChan) }) }
	defer finishFirstAttempt()

	if s.startTimeout == 0 {
		s.startTimeout = defaultStartTimeout
	}
	if s.startRetryMinDelay == 0 {
		s.startRetryMinDelay = defaultStartRetryMinDelay
	}
	if s.startRetryMaxDelay == 0 {
		s.startRetryMaxDelay = defaultStartRetryMaxDelay
	}
	retryDelay := backoff.Backoff{Min: s.startRetryMinDelay, Max: s.startRetryMaxDelay, Factor: 2, Jitter: true}

	ctx, cancel := s.stopChan.NewCtx()
	defer cancel()

	for attempt := 1; ; attempt++ {
		initialiseCalled, err := s.attemptStart(ctx)
		s.setReadyErr(err)
		finishFirstAttempt()
		if err == nil || ctx.Err() != nil {
			return
		}

		s.startFailures.Add(ctx, 1, metric.WithAttributes(
			attribute.String("capability_id", s.capabilityID),
			attribute.String("command", filepath.Base(s.command)),
		))
		delay := retryDelay.Duration()
		s.log.Warnw("Standard capabilities start failed, will retry", "command", s.command, "attempt", attempt, "retryIn", delay, "newPluginProcess", initialiseCalled)

		select {
		case <-ctx.Done():
			return
		case <-time.After(delay):
		}

		// A failed Initialise can leave services running in the plugin process,
		// so Initialise must not run twice on one process. The process is kept
		// until now so that its logs and metrics for the failure are exported.
		if initialiseCalled {
			s.closePluginLoop()
		}
	}
}

// attemptStart runs the asynchronous part of the start once. The returned bool
// reports whether Initialise ran on the current plugin process.
func (s *StandardCapabilities) attemptStart(ctx context.Context) (bool, error) {
	if s.capabilitiesLoop == nil {
		if err := s.launchPluginLoop(ctx); err != nil {
			s.log.Errorf("error starting standard capabilities service: %v", err)
			return false, fmt.Errorf("error starting standard capabilities service: %w", err)
		}
	}

	capDonID, err := s.capDonIDResolver.awaitCapDonID(ctx)
	if err != nil {
		s.log.Errorf("error waiting to resolve the capability Don ID: %v", err)
		return false, fmt.Errorf("error waiting to resolve the capability Don ID: %w", err)
	}

	cctx, cancel := context.WithTimeout(ctx, s.startTimeout)
	defer cancel()

	if err = s.capabilitiesLoop.WaitCtx(cctx); err != nil {
		s.log.Errorf("error waiting for standard capabilities service to start: %v", err)
		return false, fmt.Errorf("waiting for standard capabilities service to start: %w", err)
	}

	dependencies := s.initialiseDependencies(capDonID)
	if err = s.capabilitiesLoop.Initialise(cctx, dependencies); err != nil {
		s.log.Errorf("error initialising standard capabilities service: %v", err)
		return true, fmt.Errorf("initialising standard capabilities service: %w", err)
	}

	capabilityInfos, err := s.capabilitiesLoop.Infos(cctx)
	if err != nil {
		s.log.Errorf("error getting standard capabilities service info: %v", err)
		return true, fmt.Errorf("getting standard capabilities service info: %w", err)
	}

	s.log.Info("Started standard capabilities", "command", s.command, "capabilities", capabilityInfos)
	return true, nil
}

func (s *StandardCapabilities) launchPluginLoop(ctx context.Context) error {
	launched := s.newPluginLoop()
	if err := launched.Start(ctx); err != nil {
		return err
	}
	s.capabilitiesLoop = launched
	return nil
}

func (s *StandardCapabilities) closePluginLoop() {
	if s.capabilitiesLoop == nil {
		return
	}
	// A plugin whose start failed may have added its capabilities already. Remove
	// them before the process is stopped: see registeredCapabilities.
	s.registered.removeAll(s.log, s.registryRemoveTimeout)
	if err := s.capabilitiesLoop.Close(); err != nil {
		s.log.Warnw("error closing standard capabilities service before retrying the start", "err", err)
	}
	s.capabilitiesLoop = nil
}

// Ready is a non-blocking check for the service's ready state.  Errors if not
// ready when called.
func (s *StandardCapabilities) Ready() error {
	if err := s.StateMachine.Ready(); err != nil {
		return err
	}
	select {
	case <-s.readyChan:
		// The first start attempt has finished; surface the outcome of the
		// latest attempt (nil on success, the WaitCtx/Initialise/Infos error
		// otherwise).
		return s.getReadyErr()
	case <-s.stopChan:
		return ErrServiceStopped
	default:
		return ErrServiceNotReady
	}
}

// Name implements services.HealthReporter. It returns the per-capability
// logger name so health metrics are attributed to this specific capability.
func (s *StandardCapabilities) Name() string {
	return s.log.Name()
}

// HealthReport implements services.HealthReporter. Unlike the embedded
// StateMachine's default (which reports healthy as soon as Start returns), this
// reflects the asynchronous startup outcome so a capability that failed to
// initialise is reported unhealthy instead of falsely ready. It turns healthy
// when a retried start succeeds.
func (s *StandardCapabilities) HealthReport() map[string]error {
	return map[string]error{s.Name(): s.Ready()}
}

// Await waits for the first start attempt to finish or for the context to be
// cancelled. A failed start is retried in the background, see Ready.
func (s *StandardCapabilities) Await(ctx context.Context) error {
	select {
	case <-s.readyChan:
		return s.getReadyErr()
	case <-s.stopChan:
		return ErrServiceStopped
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *StandardCapabilities) Close() error {
	close(s.stopChan)
	s.wg.Wait()
	// Release the LOOP registration after the LOOP is stopped, even if StopOnce
	// fails, so that a new instance with the same ID can register it again.
	defer s.unregisterLOOP()
	return s.StopOnce("StandardCapabilities", func() error {
		// Before the plugin is stopped: see registeredCapabilities.
		s.registered.removeAll(s.log, s.registryRemoveTimeout)

		if s.capabilitiesLoop != nil {
			return s.capabilitiesLoop.Close()
		}

		return nil
	})
}

func (s *StandardCapabilities) unregisterLOOP() {
	if s.loopID == "" {
		return
	}
	s.pluginRegistrar.UnregisterLOOP(s.loopID)
	s.loopID = ""
}
