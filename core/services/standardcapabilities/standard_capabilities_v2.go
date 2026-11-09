package standardcapabilities

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"unicode"

	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry/remote"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	settingsremote "github.com/smartcontractkit/chainlink-common/pkg/settings/remote"
	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

func listenLocal(ctx context.Context) (net.Listener, error) {
	return (&net.ListenConfig{}).Listen(ctx, "tcp", "127.0.0.1:0")
}

// V2 runs a capability built with capability.Run. The host supervises the
// process as an empty LOOP, and gives it its dependencies in env vars. The process reaches the
// capabilities registry and CRE settings through a gRPC server that the host serves for this LOOP only.
type V2 struct {
	services.StateMachine
	log              logger.Logger
	command          string
	app              string
	configJSON       string
	capDonIDResolver CapDonIDResolver
	pluginRegistrar  plugins.RegistrarConfig
	registry         registry.CapabilitiesRegistry
	creSettings      settings.Broadcaster

	emptyLoop atomic.Pointer[loop.EmptyService]
	loopID    string
	wg        sync.WaitGroup
	stopChan  services.StopChan

	// The process is launched in the background once the capability DON ID is
	// resolved. It has its own cancellation, separate from stopChan, because Close
	// must stop the launch before closing the process, and the process before
	// stopping the server it holds streams on.
	launchWG     sync.WaitGroup
	cancelLaunch context.CancelFunc
	launchErr    atomic.Pointer[error]
}

type V2Params struct {
	Log              logger.Logger
	Command          string
	App              string
	ConfigJSON       string
	CapDonIDResolver CapDonIDResolver
	PluginRegistrar  plugins.RegistrarConfig
	// Registry is served to the process. Nil serves no registry.
	Registry registry.CapabilitiesRegistry
	// CRESettings is served to the process. Nil serves no settings.
	CRESettings settings.Broadcaster
}

func NewV2(p V2Params) (*V2, error) {
	return &V2{
		log:              p.Log,
		command:          p.Command,
		app:              p.App,
		configJSON:       p.ConfigJSON,
		capDonIDResolver: p.CapDonIDResolver,
		pluginRegistrar:  p.PluginRegistrar,
		registry:         p.Registry,
		creSettings:      p.CRESettings,
		stopChan:         make(services.StopChan),
	}, nil
}

func (s *V2) Start(ctx context.Context) error {
	return s.StartOnce("StandardCapabilitiesV2", func() error {
		envVars, err := s.prepareEnvVars()
		if err != nil {
			return fmt.Errorf("failed to prepare env vars: %w", err)
		}

		srv, err := s.newGRPCServer(ctx)
		if err != nil {
			return fmt.Errorf("failed to start server: %w", err)
		}

		if s.registry != nil {
			transport, err := remote.NewListenerTransport(remote.ListenerTransportParams{
				Lggr: s.log,
				Listen: func() (net.Listener, error) {
					return listenLocal(context.Background())
				},
			})
			if err != nil {
				srv.close()
				return fmt.Errorf("failed to create registry transport: %w", err)
			}

			remote.RegisterCapabilitiesRegistryServer(srv.Server, s.log, s.registry, transport)
			envVars = append(envVars, "CL_CAPABILITIES_PROXY_URL="+srv.addr())
		}

		if s.creSettings != nil {
			settingsremote.RegisterSettingsServer(srv.Server, s.log, s.creSettings)
			envVars = append(envVars, "CL_SETTINGS_URL="+srv.addr())
		}

		s.startGRPCServer(srv)

		// Launch the actual LOOPP in the background. This needs to resolve the capabilityDonID first
		// if it's not present, which requires the capabilities registyr to sync first.
		launchCtx, cancel := context.WithCancel(context.WithoutCancel(ctx))
		s.cancelLaunch = cancel
		s.launchWG.Go(func() {
			if err := s.launch(launchCtx, envVars); err != nil {
				s.log.Errorw("Failed to launch standard capability", "err", err)
				s.launchErr.Store(&err)
				srv.close()
			}
		})
		return nil
	})
}

func (s *V2) launch(ctx context.Context, envVars []string) error {
	donID, err := s.capDonIDResolver.awaitCapDonID(ctx)
	if err != nil {
		return fmt.Errorf("failed to resolve capability DON ID: %w", err)
	}
	if donID != 0 {
		envVars = append(envVars, fmt.Sprintf("CL_CAPABILITIES_CAPABILITY_DON_ID=%d", donID))
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

	emptyLoop := loop.NewEmptyService(s.log, opts, cmdFn)
	s.emptyLoop.Store(emptyLoop)
	if err = emptyLoop.Start(ctx); err != nil {
		s.unregisterLOOP()
		return fmt.Errorf("error starting empty loop: %w", err)
	}
	return nil
}

func (s *V2) Close() error {
	defer s.unregisterLOOP()
	return s.StopOnce("StandardCapabilitiesV2", func() error {
		s.cancelLaunch()
		s.launchWG.Wait()

		var err error
		if emptyLoop := s.emptyLoop.Load(); emptyLoop != nil {
			err = emptyLoop.Close()
		}
		close(s.stopChan)
		s.wg.Wait()
		return err
	})
}

type server struct {
	listener net.Listener
	*grpc.Server
}

func (s *server) addr() string {
	return s.listener.Addr().String()
}

func (s *server) close() {
	s.GracefulStop()
	_ = s.listener.Close()
}

func (s *V2) newGRPCServer(ctx context.Context) (*server, error) {
	listener, err := listenLocal(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to listen for the server: %w", err)
	}
	srv := grpc.NewServer()
	return &server{listener: listener, Server: srv}, nil
}

func (s *V2) startGRPCServer(srv *server) {
	ctx, cancel := s.stopChan.NewCtx()
	s.wg.Go(func() {
		defer cancel()
		if err := srv.Serve(srv.listener); err != nil {
			s.log.Errorw("failed to serve", "err", err)
		}
	})
	s.wg.Go(func() {
		<-ctx.Done()
		srv.GracefulStop()
	})
}

func (s *V2) Name() string {
	return s.log.Name()
}

func (s *V2) HealthReport() map[string]error {
	healthErr := s.Healthy()
	if launchErr := s.launchErr.Load(); launchErr != nil {
		healthErr = errors.Join(healthErr, *launchErr)
	}
	report := map[string]error{s.Name(): healthErr}
	if emptyLoop := s.emptyLoop.Load(); emptyLoop != nil {
		maps.Copy(report, emptyLoop.HealthReport())
	}
	return report
}

func (s *V2) unregisterLOOP() {
	if s.loopID == "" {
		return
	}
	s.pluginRegistrar.UnregisterLOOP(s.loopID)
	s.loopID = ""
}

func (s *V2) prepareEnvVars() ([]string, error) {
	cfgEnv, err := envVarsFromConfigJSON(s.app, s.configJSON)
	if err != nil {
		return nil, err
	}
	envVars, err := plugins.ParseEnvFile(env.CapabilitiesPlugin.Env.Get())
	if err != nil {
		return nil, fmt.Errorf("failed to parse capabilities env file: %w", err)
	}
	return slices.Concat(envVars, cfgEnv), nil
}

// envVarsFromConfigJSON turns each key of the job config into the env var CL_<APP>_<KEY>, the name under which
// capability.Run reads it. For example, fastestScheduleIntervalSeconds of the cron app becomes
// CL_CRON_FASTEST_SCHEDULE_INTERVAL_SECONDS.
func envVarsFromConfigJSON(app, configJSON string) ([]string, error) {
	if strings.TrimSpace(configJSON) == "" {
		return nil, nil
	}
	dec := json.NewDecoder(strings.NewReader(configJSON))
	dec.UseNumber()
	var cfg map[string]any
	if err := dec.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	cfgEnv := make([]string, 0, len(cfg))
	for k, v := range cfg {
		switch v.(type) {
		case string, json.Number, bool:
		default:
			return nil, fmt.Errorf("config key %q: only strings, numbers and booleans are supported, got %T", k, v)
		}
		cfgEnv = append(cfgEnv, fmt.Sprintf("CL_%s_%s=%v", envName(app), envName(k), v))
	}
	slices.Sort(cfgEnv)
	return cfgEnv, nil
}

func envName(key string) string {
	var b strings.Builder
	var prev rune
	for i, r := range key {
		switch {
		case r == '-' || r == '.':
			r = '_'
		case unicode.IsUpper(r) && i > 0 && (unicode.IsLower(prev) || unicode.IsDigit(prev)):
			b.WriteRune('_')
		}
		b.WriteRune(unicode.ToUpper(r))
		prev = r
	}
	return b.String()
}
