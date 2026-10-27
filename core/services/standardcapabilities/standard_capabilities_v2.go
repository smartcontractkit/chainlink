package standardcapabilities

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net"
	"path/filepath"
	"slices"
	"strings"
	"unicode"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry/remote"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	settingsgrpc "github.com/smartcontractkit/chainlink-common/pkg/settings/grpc"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

// v2Apps maps the command of a capability built with capability.Run to its app name. The app name
// prefixes the env vars that carry the capability's config.
var v2Apps = map[string]string{
	"cron-v2": "cron",
}

func v2App(command string) (string, bool) {
	app, ok := v2Apps[filepath.Base(command)]
	return app, ok
}

// StandardCapabilitiesV2 runs a capability built with capability.Run. The host supervises the
// process as an empty LOOP, and gives it its dependencies in env vars. The process reaches the
// capabilities registry and the CRE settings through gRPC servers that the host serves for this
// LOOP only.
type StandardCapabilitiesV2 struct {
	services.StateMachine
	log             logger.Logger
	command         string
	env             []string
	pluginRegistrar plugins.RegistrarConfig
	registry        registry.CapabilitiesRegistry
	creSettings     core.SettingsBroadcaster

	emptyLoop   *loop.EmptyService
	loopID      string
	stopServers []func()
}

type StandardCapabilitiesV2Params struct {
	Log             logger.Logger
	Command         string
	App             string
	ConfigJSON      string
	CapabilityDonID uint32
	PluginRegistrar plugins.RegistrarConfig
	// Registry is served to the process. Nil serves no registry.
	Registry registry.CapabilitiesRegistry
	// CRESettings is served to the process. Nil serves no settings, and the process uses its defaults.
	CRESettings core.SettingsBroadcaster
}

func NewStandardCapabilitiesV2(p StandardCapabilitiesV2Params) (*StandardCapabilitiesV2, error) {
	cfgEnv, err := configEnv(p.App, p.ConfigJSON)
	if err != nil {
		return nil, err
	}
	if p.CapabilityDonID != 0 {
		cfgEnv = append(cfgEnv, fmt.Sprintf("CL_CAPABILITIES_CAPABILITY_DON_ID=%d", p.CapabilityDonID))
	}
	return &StandardCapabilitiesV2{
		log:             p.Log,
		command:         p.Command,
		env:             cfgEnv,
		pluginRegistrar: p.PluginRegistrar,
		registry:        p.Registry,
		creSettings:     p.CRESettings,
	}, nil
}

func (s *StandardCapabilitiesV2) Start(ctx context.Context) error {
	return s.StartOnce("StandardCapabilitiesV2", func() error {
		envVars, err := plugins.ParseEnvFile(env.CapabilitiesPlugin.Env.Get())
		if err != nil {
			return fmt.Errorf("failed to parse capabilities env file: %w", err)
		}
		serverEnv, err := s.startServers()
		if err != nil {
			s.closeServers()
			return err
		}
		loopID := s.log.Name()
		cmdFn, opts, err := s.pluginRegistrar.RegisterLOOP(plugins.CmdConfig{
			ID:  loopID,
			Cmd: s.command,
			Env: slices.Concat(envVars, s.env, serverEnv),
		})
		if err != nil {
			s.closeServers()
			return fmt.Errorf("error registering loop: %w", err)
		}
		s.loopID = loopID

		s.emptyLoop = loop.NewEmptyService(s.log, opts, cmdFn)
		if err = s.emptyLoop.Start(ctx); err != nil {
			s.unregisterLOOP()
			s.closeServers()
			return fmt.Errorf("error starting empty loop: %w", err)
		}
		return nil
	})
}

func (s *StandardCapabilitiesV2) Close() error {
	defer s.unregisterLOOP()
	return s.StopOnce("StandardCapabilitiesV2", func() error {
		var err error
		if s.emptyLoop != nil {
			err = s.emptyLoop.Close()
		}
		s.closeServers()
		return err
	})
}

// startServers serves the registry and the CRE settings to the process, and returns the env vars
// that tell the process where they are.
func (s *StandardCapabilitiesV2) startServers() ([]string, error) {
	var serverEnv []string
	if s.registry != nil {
		transport, err := remote.NewListenerTransport(remote.ListenerTransportParams{
			Lggr:   s.log,
			Listen: listenLocal,
			Connect: func(target string) (*grpc.ClientConn, error) {
				return grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create registry transport: %w", err)
		}
		addr, err := s.serve("registry", func(srv *grpc.Server) {
			remote.RegisterCapabilitiesRegistryServer(srv, s.log, s.registry, transport)
		})
		if err != nil {
			return nil, err
		}
		serverEnv = append(serverEnv, "CL_CAPABILITIES_PROXY_URL="+addr)
	}
	if s.creSettings != nil {
		addr, err := s.serve("settings", func(srv *grpc.Server) {
			capabilitiespb.RegisterSettingsServer(srv, settingsgrpc.NewServer(s.creSettings))
		})
		if err != nil {
			return nil, err
		}
		serverEnv = append(serverEnv, "CL_SETTINGS_URL="+addr)
	}
	return serverEnv, nil
}

func (s *StandardCapabilitiesV2) serve(name string, register func(*grpc.Server)) (string, error) {
	lis, err := listenLocal()
	if err != nil {
		return "", fmt.Errorf("failed to listen for the %s server: %w", name, err)
	}
	srv := grpc.NewServer()
	register(srv)
	done := make(chan struct{})
	go func() {
		defer close(done)
		if err := srv.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			s.log.Errorw("Failed to serve", "server", name, "err", err)
		}
	}()
	s.stopServers = append(s.stopServers, func() {
		srv.Stop()
		<-done
	})
	return lis.Addr().String(), nil
}

func (s *StandardCapabilitiesV2) closeServers() {
	for _, stop := range s.stopServers {
		stop()
	}
	s.stopServers = nil
}

func listenLocal() (net.Listener, error) {
	return net.Listen("tcp", "127.0.0.1:0")
}

func (s *StandardCapabilitiesV2) Name() string {
	return s.log.Name()
}

func (s *StandardCapabilitiesV2) HealthReport() map[string]error {
	report := map[string]error{s.Name(): s.Healthy()}
	if s.emptyLoop != nil {
		maps.Copy(report, s.emptyLoop.HealthReport())
	}
	return report
}

func (s *StandardCapabilitiesV2) unregisterLOOP() {
	if s.loopID == "" {
		return
	}
	s.pluginRegistrar.UnregisterLOOP(s.loopID)
	s.loopID = ""
}

// configEnv turns each key of the job config into the env var CL_<APP>_<KEY>, the name under which
// capability.Run reads it. For example, fastestScheduleIntervalSeconds of the cron app becomes
// CL_CRON_FASTEST_SCHEDULE_INTERVAL_SECONDS.
func configEnv(app, configJSON string) ([]string, error) {
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
