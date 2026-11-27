package standardcapabilities

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/connectivity"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

// TestStandardCapabilities_ForwardsPluginEnvFile tests that the standard capabilities
// LOOPP launcher reads the env file referenced by CL_CAPABILITIES_ENV and forwards
// its contents to the plugin. This is how operator-supplied env vars (e.g. chain-selector's
// EXTRA_SELECTORS_FILE) make it into the LOOPP subprocess despite plugins.NewCmdFactory
// not inheriting os.Environ().
func TestStandardCapabilities_ForwardsPluginEnvFile(t *testing.T) {
	writeEnvFile := func(t *testing.T, contents string) string {
		t.Helper()
		path := filepath.Join(t.TempDir(), "capabilities.env")
		require.NoError(t, os.WriteFile(path, []byte(contents), 0o600))
		return path
	}

	startAndCapture := func(t *testing.T) (plugins.CmdConfig, error) {
		t.Helper()
		var captured plugins.CmdConfig
		registrar := &capturingRegistrar{
			onRegister: func(cfg plugins.CmdConfig) { captured = cfg },
		}
		std := NewStandardCapabilities(
			logger.TestLogger(t),
			"not/found/path/to/binary",
			"{}",
			registrar,
			&staticCapDonIDResolver{},
			core.StandardCapabilitiesDependencies{},
		)
		// Start runs RegisterLOOP synchronously; the capturingRegistrar
		// returns an error so Start unwinds before any LOOPP lifecycle code
		// runs. The returned error from Start is propagated to the caller so
		// each sub-test can assert on it as appropriate.
		return captured, std.Start(t.Context())
	}

	t.Run("env file is parsed and forwarded via CmdConfig.Env", func(t *testing.T) {
		envFile := writeEnvFile(t, "EXTRA_SELECTORS_FILE=/extraConfig/extra-selectors.yaml\nFOO=bar\n")
		t.Setenv(string(env.CapabilitiesPlugin.Env), envFile)

		cfg, err := startAndCapture(t)
		require.Error(t, err, "expected synthetic Start failure from capturingRegistrar")
		require.Contains(t, err.Error(), capturingRegistrarErr)

		require.Contains(t, cfg.Env, "EXTRA_SELECTORS_FILE=/extraConfig/extra-selectors.yaml",
			"standard-capability LOOPP launcher should forward EXTRA_SELECTORS_FILE from the operator-supplied env file so chain-selectors' init() in the LOOPP can read operator-provided selectors")
		require.Contains(t, cfg.Env, "FOO=bar",
			"all entries in the operator-supplied env file should be forwarded, not just well-known ones")
	})

	t.Run("env file unset results in empty CmdConfig.Env", func(t *testing.T) {
		t.Setenv(string(env.CapabilitiesPlugin.Env), "")

		cfg, err := startAndCapture(t)
		require.Error(t, err, "expected synthetic Start failure from capturingRegistrar")
		require.Contains(t, err.Error(), capturingRegistrarErr)

		require.Empty(t, cfg.Env,
			"no operator-provided env vars should be forwarded when CL_CAPABILITIES_ENV is unset")
	})

	t.Run("missing env file fails Start before RegisterLOOP", func(t *testing.T) {
		missingPath := filepath.Join(t.TempDir(), "does-not-exist.env")
		t.Setenv(string(env.CapabilitiesPlugin.Env), missingPath)

		var registerCalled bool
		registrar := &capturingRegistrar{
			onRegister: func(plugins.CmdConfig) { registerCalled = true },
		}
		std := NewStandardCapabilities(
			logger.TestLogger(t),
			"not/found/path/to/binary",
			"{}",
			registrar,
			&staticCapDonIDResolver{},
			core.StandardCapabilitiesDependencies{},
		)

		err := std.Start(t.Context())
		require.Error(t, err, "Start should fail when the env file is missing")
		require.Contains(t, err.Error(), "failed to parse capabilities env file")
		require.False(t, registerCalled, "RegisterLOOP must not be called when env-file parsing fails")
	})
}

func TestStandardCapabilities_InitialiseDependenciesRoundTrip(t *testing.T) {
	t.Parallel()
	want := core.StandardCapabilitiesDependencies{
		Config: "test-config",
	}

	std := NewStandardCapabilities(
		logger.TestLogger(t),
		"not/found/path/to/binary",
		want.Config,
		&capturingRegistrar{},
		&staticCapDonIDResolver{},
		want,
	)

	got := std.initialiseDependencies(0)

	require.Equal(t, want.Config, got.Config, "config should be re-delivered to LOOP Initialise dependencies")
}

func TestStandardCapabilities_ShardIndexDeliveredToLOOP(t *testing.T) {
	t.Parallel()
	std := NewStandardCapabilities(
		logger.TestLogger(t),
		"not/found/path/to/binary",
		"{}",
		&capturingRegistrar{},
		&staticCapDonIDResolver{},
		core.StandardCapabilitiesDependencies{ShardIndex: 3},
	)

	require.Equal(t, uint32(3), std.initialiseDependencies(0).ShardIndex)
}

// Without this, a trigger producer silently falls back to the consumer workflow's
// DON for metering identity and event labels even when the node knows which DON
// it is serving.
func TestStandardCapabilities_CapabilityDonIDDeliveredToLOOP(t *testing.T) {
	t.Parallel()
	std := NewStandardCapabilities(logger.TestLogger(t), "not/found/path/to/binary", "{}", &capturingRegistrar{}, &staticCapDonIDResolver{}, core.StandardCapabilitiesDependencies{})

	require.Equal(t, uint32(42), std.initialiseDependencies(42).CapabilityDonID)
}

const capturingRegistrarErr = "capturingRegistrar: stop after capture"

// capturingRegistrar records the CmdConfig handed to RegisterLOOP and then
// returns a sentinel error so the caller's Start flow exits before any
// LOOPP lifecycle code runs.
type capturingRegistrar struct {
	onRegister func(plugins.CmdConfig)
}

func (c *capturingRegistrar) RegisterLOOP(cfg plugins.CmdConfig) (func() *exec.Cmd, loop.GRPCOpts, error) {
	if c.onRegister != nil {
		c.onRegister(cfg)
	}
	return nil, loop.GRPCOpts{}, errors.New(capturingRegistrarErr)
}

func (c *capturingRegistrar) UnregisterLOOP(string) {}

// TestStandardCapabilities_CloseUnregistersLOOP checks that a capability can be
// restarted with the same LOOP ID, as LocalCapabilityManager does on a config
// change. Before the fix, Close left the ID registered and the second Start
// failed with plugins.ErrExists ("plugin already registered").
func TestStandardCapabilities_CloseUnregistersLOOP(t *testing.T) {
	t.Parallel()
	registered := map[string]bool{}
	pluginRegistrar := plugins.NewRegistrarConfig(loop.GRPCOpts{},
		func(id string) (*plugins.RegisteredLoop, error) {
			if registered[id] {
				return nil, plugins.ErrExists
			}
			registered[id] = true
			return &plugins.RegisteredLoop{}, nil
		},
		func(id string) { delete(registered, id) },
	)
	lggr := logger.TestLogger(t).Named("consensus@1.0.0-alpha")

	newStd := func() *StandardCapabilities {
		std := NewStandardCapabilities(lggr, "not/found/path/to/binary", "{}", pluginRegistrar, &staticCapDonIDResolver{}, core.StandardCapabilitiesDependencies{})
		std.startTimeout = time.Second
		return std
	}

	first := newStd()
	require.NoError(t, first.Start(t.Context()))
	require.True(t, registered[lggr.Name()])
	require.NoError(t, first.Close())
	require.Empty(t, registered, "Close must unregister the LOOP")

	second := newStd()
	require.NoError(t, second.Start(t.Context()), "restart with the same LOOP ID must succeed")
	require.NoError(t, second.Close())
	require.Empty(t, registered)
}

// A registry entry that outlives its plugin blocks the registration of a
// replacement, so Close must remove what the plugin added while the plugin is
// still up.
func TestStandardCapabilities_CloseRemovesCapabilitiesFromRegistry(t *testing.T) {
	t.Parallel()
	const capabilityID = "evm:ChainSelector:16015286601757825753@1.0.0"
	newInstance := func() capabilities.ExecutableCapability {
		return readyCapability{CapabilityInfo: capabilities.MustNewCapabilityInfo(capabilityID, capabilities.CapabilityTypeAction, "test")}
	}

	newStarted := func(t *testing.T, reg core.CapabilitiesRegistry) *StandardCapabilities {
		t.Helper()
		pluginRegistrar := plugins.NewRegistrarConfig(loop.GRPCOpts{},
			func(string) (*plugins.RegisteredLoop, error) { return &plugins.RegisteredLoop{}, nil },
			func(string) {},
		)
		std := NewStandardCapabilities(logger.TestLogger(t), "not/found/path/to/binary", "{}", pluginRegistrar,
			&staticCapDonIDResolver{}, core.StandardCapabilitiesDependencies{CapabilityRegistry: reg})
		std.startTimeout = time.Second
		require.NoError(t, std.Start(t.Context()))
		return std
	}
	// pluginState is the state of the plugin supervisor. These tests use the real one.
	pluginState := func(std *StandardCapabilities) string {
		return std.capabilitiesLoop.(standardCapabilitiesLoop).State()
	}
	// pluginAdds adds a capability the way the plugin does: through the registry it is
	// given in Initialise. The binary does not exist, so no plugin runs in these tests.
	pluginRegistry := func(std *StandardCapabilities) core.CapabilitiesRegistry {
		return std.initialiseDependencies(0).CapabilityRegistry
	}

	t.Run("before the plugin is stopped", func(t *testing.T) {
		t.Parallel()
		reg := mocks.NewCapabilitiesRegistry(t)
		std := newStarted(t, reg)
		reg.EXPECT().Add(mock.Anything, mock.Anything).Return(nil).Once()
		require.NoError(t, pluginRegistry(std).Add(t.Context(), newInstance()))

		var pluginStateAtRemove string
		reg.EXPECT().Get(mock.Anything, capabilityID).Return(nil, nil).Once()
		reg.EXPECT().Remove(mock.Anything, capabilityID).RunAndReturn(func(context.Context, string) error {
			pluginStateAtRemove = pluginState(std)
			return nil
		}).Once()

		require.NoError(t, std.Close())
		require.Equal(t, "Started", pluginStateAtRemove)
		require.Equal(t, "Stopped", pluginState(std))
	})

	// The registry the node uses does not replace an entry that is still live, so
	// without the removal the next instance of the capability cannot be added.
	t.Run("a new instance can be added to the registry after Close", func(t *testing.T) {
		t.Parallel()
		reg := registry.NewRegistry(logger.TestLogger(t))
		std := newStarted(t, reg)
		require.NoError(t, pluginRegistry(std).Add(t.Context(), newInstance()))
		require.ErrorIs(t, reg.Add(t.Context(), newInstance()), registry.ErrCapabilityAlreadyExists)

		require.NoError(t, std.Close())

		require.NoError(t, reg.Add(t.Context(), newInstance()))
	})

	// In the next two cases the mock fails the test if Close sends a Get or a Remove.
	t.Run("not what the registry refused to add", func(t *testing.T) {
		t.Parallel()
		reg := mocks.NewCapabilitiesRegistry(t)
		std := newStarted(t, reg)
		refused := errors.New("refused")
		reg.EXPECT().Add(mock.Anything, mock.Anything).Return(refused).Once()
		require.ErrorIs(t, pluginRegistry(std).Add(t.Context(), newInstance()), refused)

		require.NoError(t, std.Close())
	})

	t.Run("not what the plugin removed itself", func(t *testing.T) {
		t.Parallel()
		reg := mocks.NewCapabilitiesRegistry(t)
		std := newStarted(t, reg)
		reg.EXPECT().Add(mock.Anything, mock.Anything).Return(nil).Once()
		require.NoError(t, pluginRegistry(std).Add(t.Context(), newInstance()))
		reg.EXPECT().Remove(mock.Anything, capabilityID).Return(nil).Once()
		require.NoError(t, pluginRegistry(std).Remove(t.Context(), capabilityID))

		require.NoError(t, std.Close())
	})

	t.Run("a capability that is not in the registry does not fail Close", func(t *testing.T) {
		t.Parallel()
		reg := mocks.NewCapabilitiesRegistry(t)
		std := newStarted(t, reg)
		reg.EXPECT().Add(mock.Anything, mock.Anything).Return(nil).Once()
		require.NoError(t, pluginRegistry(std).Add(t.Context(), newInstance()))

		notFound := errors.New("capability not found")
		reg.EXPECT().Get(mock.Anything, capabilityID).Return(nil, notFound).Once()
		reg.EXPECT().Remove(mock.Anything, capabilityID).Return(notFound).Once()

		require.NoError(t, std.Close())
		require.Equal(t, "Stopped", pluginState(std))
	})

	// A registry that does not observe a context cannot abandon a Remove that waits
	// for a lock, and a Remove that completes after Close could remove the
	// replacement. So when the registry does not answer, Close must go on without a
	// Remove.
	t.Run("a registry that does not answer does not block Close and gets no Remove", func(t *testing.T) {
		t.Parallel()
		reg := mocks.NewCapabilitiesRegistry(t)
		std := newStarted(t, reg)
		std.registryRemoveTimeout = 10 * time.Millisecond
		reg.EXPECT().Add(mock.Anything, mock.Anything).Return(nil).Once()
		require.NoError(t, pluginRegistry(std).Add(t.Context(), newInstance()))

		unblock := make(chan struct{})
		t.Cleanup(func() { close(unblock) })
		reg.EXPECT().Get(mock.Anything, capabilityID).RunAndReturn(func(context.Context, string) (capabilities.BaseCapability, error) {
			<-unblock
			return nil, nil
		}).Once()

		closeErr := make(chan error, 1)
		go func() { closeErr <- std.Close() }()
		select {
		case err := <-closeErr:
			require.NoError(t, err)
		case <-time.After(tests.WaitTimeout(t)):
			require.FailNow(t, "Close waited for the registry")
		}
		require.Equal(t, "Stopped", pluginState(std))
	})

	t.Run("no registry", func(t *testing.T) {
		t.Parallel()
		std := newStarted(t, nil)
		require.Nil(t, pluginRegistry(std))
		require.NoError(t, std.Close())
	})
}

// readyCapability is a capability whose connection is up, as it is for a
// running plugin.
type readyCapability struct {
	capabilities.CapabilityInfo
	capabilities.Executable
}

func (readyCapability) GetState() connectivity.State { return connectivity.Ready }

func TestStandardCapabilityStart(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	t.Run("NOK-not_found_binary_does_not_block", func(t *testing.T) {
		ctx := t.Context()
		lggr := logger.TestLogger(t)

		pluginRegistrar := plugins.NewRegistrarConfig(loop.GRPCOpts{}, func(name string) (*plugins.RegisteredLoop, error) { return &plugins.RegisteredLoop{}, nil }, func(loopId string) {})
		registry := mocks.NewCapabilitiesRegistry(t)

		spec := &job.StandardCapabilitiesSpec{
			Command: "not/found/path/to/binary",
			OracleFactory: job.OracleFactoryConfig{
				Enabled: true,
				BootstrapPeers: []string{
					"12D3KooWEBVwbfdhKnicois7FTYVsBFGFcoMhMCKXQC57BQyZMhz@localhost:6690",
				},
				OCRContractAddress: "0x2279B7A0a67DB372996a5FaB50D91eAA73d2eBe6",
				ChainID:            "31337",
			}}

		dependencies := core.StandardCapabilitiesDependencies{
			Config:             spec.Config,
			TelemetryService:   &telemetryServiceMock{},
			Store:              &kvstoreMock{},
			CapabilityRegistry: registry,
			ErrorLog:           &errorLogMock{},
			PipelineRunner:     &pipelineRunnerServiceMock{},
			RelayerSet:         &relayerSetMock{},
			OracleFactory:      &oracleFactoryMock{},
			GatewayConnector:   &gatewayConnectorMock{},
			P2PKeystore:        &keystoreMock{},
		}
		standardCapability := NewStandardCapabilities(lggr, spec.Command, spec.Config, pluginRegistrar, &staticCapDonIDResolver{}, dependencies)
		standardCapability.startTimeout = 1 * time.Second
		err := standardCapability.Start(ctx)
		require.NoError(t, err)

		require.Error(t, standardCapability.Await(ctx), "the first start attempt must fail without the binary")
		// Stop the plugin supervisor. It logs each time it tries to launch the plugin again,
		// and a log after the test has ended fails the whole package.
		require.NoError(t, standardCapability.Close())
	})
}

// A failed Initialise can leave services running in the plugin process, so every
// retry after it must run on a new process, with the capability DON ID resolved again.
func TestStandardCapabilities_RetriesFailedStartOnNewPluginProcess(t *testing.T) {
	t.Parallel()
	errNotReady := errors.New("registry not ready")

	for _, tc := range []struct {
		name       string
		failing    fakePluginLoop
		wantEvents []string
	}{
		{
			name:    "Initialise fails",
			failing: fakePluginLoop{initialiseErr: errNotReady},
			wantEvents: []string{
				"loop1 start", "loop1 wait", "loop1 initialise donID=1", "loop1 close",
				"loop2 start", "loop2 wait", "loop2 initialise donID=2", "loop2 close",
				"loop3 start", "loop3 wait", "loop3 initialise donID=3", "loop3 infos",
			},
		},
		{
			name:    "Infos fails",
			failing: fakePluginLoop{infosErr: errNotReady},
			wantEvents: []string{
				"loop1 start", "loop1 wait", "loop1 initialise donID=1", "loop1 infos", "loop1 close",
				"loop2 start", "loop2 wait", "loop2 initialise donID=2", "loop2 infos", "loop2 close",
				"loop3 start", "loop3 wait", "loop3 initialise donID=3", "loop3 infos",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			loops := &fakePluginLoops{behaviours: []fakePluginLoop{tc.failing, tc.failing, {}}}
			resolver := &gatedCapDonIDResolver{donIDs: make(chan uint32)}
			std := newStandardCapabilitiesWithFakeLoops(t, loops, resolver)

			require.NoError(t, std.Start(t.Context()))
			require.ErrorIs(t, std.Ready(), ErrServiceNotReady)

			resolver.release(t, 1)
			require.ErrorIs(t, std.Await(t.Context()), errNotReady)
			require.ErrorIs(t, std.Ready(), errNotReady, "the service must be unhealthy while the start is retried")

			resolver.release(t, 2)
			resolver.release(t, 3)
			waitForStartToFinish(t, std)

			require.NoError(t, std.Ready(), "the service must turn healthy when a retried start succeeds")
			require.Equal(t, tc.wantEvents, loops.recorded())

			require.NoError(t, std.Close())
			require.Equal(t, slices.Concat(tc.wantEvents, []string{"loop3 close"}), loops.recorded())
		})
	}
}

// A plugin whose start fails after it has added its capability leaves a registry entry.
// The entry must be removed before the process is replaced: the registry does not
// replace an entry that is live, so the next process could not add its capability.
func TestStandardCapabilities_RemovesCapabilitiesBeforeItReplacesThePluginProcess(t *testing.T) {
	t.Parallel()
	const capabilityID = "evm:ChainSelector:16015286601757825753@1.0.0"
	added := readyCapability{CapabilityInfo: capabilities.MustNewCapabilityInfo(capabilityID, capabilities.CapabilityTypeAction, "test")}
	loops := &fakePluginLoops{behaviours: []fakePluginLoop{
		{adds: added, infosErr: errors.New("infos failed")},
		{adds: added},
	}}
	reg := &recordingRegistry{CapabilitiesRegistry: registry.NewRegistry(logger.TestLogger(t)), loops: loops}
	std := newStandardCapabilitiesWithFakeLoopsAndRegistry(t, loops, &staticCapDonIDResolver{}, reg)

	require.NoError(t, std.Start(t.Context()))
	waitForStartToFinish(t, std)
	require.NoError(t, std.Ready())

	wantEvents := []string{
		"loop1 start", "loop1 wait", "loop1 initialise donID=0", "registry add err=false", "loop1 infos",
		"registry remove err=false", "loop1 close",
		"loop2 start", "loop2 wait", "loop2 initialise donID=0", "registry add err=false", "loop2 infos",
	}
	require.Equal(t, wantEvents, loops.recorded())

	require.NoError(t, std.Close())
	require.Equal(t, slices.Concat(wantEvents, []string{"registry remove err=false", "loop2 close"}), loops.recorded())
}

// A plugin process that is not up within the start timeout has not run Initialise,
// so the start is retried on the same process.
func TestStandardCapabilities_RetriesStartOnSamePluginProcessWhenItIsLate(t *testing.T) {
	t.Parallel()
	loops := &fakePluginLoops{behaviours: []fakePluginLoop{{waitErrs: []error{context.DeadlineExceeded}}}}
	std := newStandardCapabilitiesWithFakeLoops(t, loops, &staticCapDonIDResolver{donID: 7})

	require.NoError(t, std.Start(t.Context()))
	waitForStartToFinish(t, std)

	require.NoError(t, std.Ready())
	require.Equal(t, []string{"loop1 start", "loop1 wait", "loop1 wait", "loop1 initialise donID=7", "loop1 infos"}, loops.recorded())
	require.NoError(t, std.Close())
}

func TestStandardCapabilities_CloseStopsRetryingTheStart(t *testing.T) {
	t.Parallel()
	loops := &fakePluginLoops{behaviours: []fakePluginLoop{{initialiseErr: errors.New("database timeout")}}}
	std := newStandardCapabilitiesWithFakeLoops(t, loops, &staticCapDonIDResolver{})
	std.startRetryMinDelay = time.Hour
	std.startRetryMaxDelay = time.Hour

	require.NoError(t, std.Start(t.Context()))
	require.Error(t, std.Await(t.Context()))

	closeErr := make(chan error, 1)
	go func() { closeErr <- std.Close() }()
	select {
	case err := <-closeErr:
		require.NoError(t, err)
	case <-time.After(tests.WaitTimeout(t)):
		require.FailNow(t, "Close did not interrupt the wait before the next start attempt")
	}

	require.Equal(t, []string{"loop1 start", "loop1 wait", "loop1 initialise donID=0", "loop1 close"}, loops.recorded())
}

func newStandardCapabilitiesWithFakeLoops(t *testing.T, loops *fakePluginLoops, resolver CapDonIDResolver) *StandardCapabilities {
	t.Helper()
	return newStandardCapabilitiesWithFakeLoopsAndRegistry(t, loops, resolver, nil)
}

func newStandardCapabilitiesWithFakeLoopsAndRegistry(t *testing.T, loops *fakePluginLoops, resolver CapDonIDResolver, reg core.CapabilitiesRegistry) *StandardCapabilities {
	t.Helper()
	registrar := plugins.NewRegistrarConfig(loop.GRPCOpts{},
		func(string) (*plugins.RegisteredLoop, error) { return &plugins.RegisteredLoop{}, nil },
		func(string) {},
	)
	std := NewStandardCapabilities(logger.TestLogger(t), "not/found/path/to/binary", "{}", registrar, resolver,
		core.StandardCapabilitiesDependencies{CapabilityRegistry: reg})
	std.newPluginLoop = loops.launch
	std.startRetryMinDelay = time.Millisecond
	std.startRetryMaxDelay = time.Millisecond
	return std
}

func waitForStartToFinish(t *testing.T, std *StandardCapabilities) {
	t.Helper()
	finished := make(chan struct{})
	go func() {
		defer close(finished)
		std.wg.Wait()
	}()
	tests.RequireSignal(t, finished, "the start did not finish")
}

// gatedCapDonIDResolver lets a test release one start attempt at a time, with the
// capability DON ID that the attempt resolves.
type gatedCapDonIDResolver struct {
	donIDs chan uint32
}

func (r *gatedCapDonIDResolver) awaitCapDonID(ctx context.Context) (uint32, error) {
	select {
	case donID := <-r.donIDs:
		return donID, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

func (r *gatedCapDonIDResolver) release(t *testing.T, donID uint32) {
	t.Helper()
	select {
	case r.donIDs <- donID:
	case <-time.After(tests.WaitTimeout(t)):
		require.FailNow(t, "no start attempt asked for the capability DON ID")
	}
}

// fakePluginLoops launches one fakePluginLoop for each plugin process the service
// starts, and records in order what the service does with each of them.
type fakePluginLoops struct {
	// behaviours has one entry for each launch. The last entry repeats.
	behaviours []fakePluginLoop

	mu       sync.Mutex
	launched int
	events   []string
}

func (f *fakePluginLoops) launch() pluginLoop {
	f.mu.Lock()
	defer f.mu.Unlock()
	launched := f.behaviours[min(f.launched, len(f.behaviours)-1)]
	f.launched++
	launched.id = f.launched
	launched.loops = f
	return &launched
}

func (f *fakePluginLoops) record(id int, event string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.events = append(f.events, fmt.Sprintf("loop%d %s", id, event))
}

func (f *fakePluginLoops) recordEvent(event string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.events = append(f.events, event)
}

func (f *fakePluginLoops) recorded() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return slices.Clone(f.events)
}

type fakePluginLoop struct {
	loops *fakePluginLoops
	id    int

	// waitErrs are returned by successive WaitCtx calls. After them WaitCtx succeeds.
	waitErrs      []error
	initialiseErr error
	infosErr      error
	// adds is added to the registry in Initialise, as a plugin does.
	adds capabilities.BaseCapability
}

func (l *fakePluginLoop) Start(context.Context) error {
	l.loops.record(l.id, "start")
	return nil
}

func (l *fakePluginLoop) Close() error {
	l.loops.record(l.id, "close")
	return nil
}

func (l *fakePluginLoop) Ready() error { return nil }

func (l *fakePluginLoop) Name() string { return fmt.Sprintf("loop%d", l.id) }

func (l *fakePluginLoop) HealthReport() map[string]error {
	return map[string]error{l.Name(): nil}
}

func (l *fakePluginLoop) WaitCtx(context.Context) error {
	l.loops.record(l.id, "wait")
	if len(l.waitErrs) == 0 {
		return nil
	}
	err := l.waitErrs[0]
	l.waitErrs = l.waitErrs[1:]
	return err
}

func (l *fakePluginLoop) Initialise(_ context.Context, dependencies core.StandardCapabilitiesDependencies) error {
	l.loops.record(l.id, fmt.Sprintf("initialise donID=%d", dependencies.CapabilityDonID))
	if l.adds != nil {
		if err := dependencies.CapabilityRegistry.Add(context.Background(), l.adds); err != nil {
			return err
		}
	}
	return l.initialiseErr
}

func (l *fakePluginLoop) Infos(context.Context) ([]capabilities.CapabilityInfo, error) {
	l.loops.record(l.id, "infos")
	return nil, l.infosErr
}

// recordingRegistry records the Add and Remove calls that reach the registry, in the
// same list as the events of the fake plugin processes.
type recordingRegistry struct {
	core.CapabilitiesRegistry
	loops *fakePluginLoops
}

func (r *recordingRegistry) Add(ctx context.Context, c capabilities.BaseCapability) error {
	err := r.CapabilitiesRegistry.Add(ctx, c)
	r.loops.recordEvent(fmt.Sprintf("registry add err=%v", err != nil))
	return err
}

func (r *recordingRegistry) Remove(ctx context.Context, id string) error {
	err := r.CapabilitiesRegistry.Remove(ctx, id)
	r.loops.recordEvent(fmt.Sprintf("registry remove err=%v", err != nil))
	return err
}

type telemetryServiceMock struct{}

func (t *telemetryServiceMock) Send(ctx context.Context, network string, chainID string, contractID string, telemetryType string, payload []byte) error {
	return nil
}

type kvstoreMock struct{}

func (k *kvstoreMock) Store(ctx context.Context, key string, val []byte) error {
	return nil
}
func (k *kvstoreMock) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, nil
}
func (k *kvstoreMock) PruneExpiredEntries(ctx context.Context, maxAge time.Duration) (int64, error) {
	return 0, nil
}

type keystoreMock struct{ core.UnimplementedKeystore }

func (k *keystoreMock) Accounts(ctx context.Context) (accounts []string, err error) {
	return nil, nil
}
func (k *keystoreMock) Sign(ctx context.Context, account string, data []byte) (signed []byte, err error) {
	return nil, nil
}

type errorLogMock struct{}

func (e *errorLogMock) SaveError(ctx context.Context, msg string) error {
	return nil
}

type relayerSetMock struct{}

func (r *relayerSetMock) Get(ctx context.Context, relayID types.RelayID) (core.Relayer, error) {
	return nil, nil
}
func (r *relayerSetMock) List(ctx context.Context, relayIDs ...types.RelayID) (map[types.RelayID]core.Relayer, error) {
	return nil, nil
}

type pipelineRunnerServiceMock struct{}

func (p *pipelineRunnerServiceMock) ExecuteRun(ctx context.Context, spec string, vars core.Vars, options core.Options) (core.TaskResults, error) {
	return nil, nil
}

type oracleFactoryMock struct{}

func (o *oracleFactoryMock) NewOracle(ctx context.Context, args core.OracleArgs) (core.Oracle, error) {
	return &oracleMock{}, nil
}

type oracleMock struct{}

func (o *oracleMock) Start(ctx context.Context) error {
	return nil
}
func (o *oracleMock) Close(ctx context.Context) error {
	return nil
}

type gatewayConnectorMock struct {
	core.UnimplementedGatewayConnector
}

func (g *gatewayConnectorMock) Start(context.Context) error {
	return nil
}

func (g *gatewayConnectorMock) Close() error {
	return nil
}

func (g *gatewayConnectorMock) AddHandler(ctx context.Context, methods []string, handler core.GatewayConnectorHandler) error {
	return nil
}

func (g *gatewayConnectorMock) RemoveHandler(ctx context.Context, methods []string) error {
	return nil
}

func (g *gatewayConnectorMock) SendToGateway(ctx context.Context, gatewayID string, resp *jsonrpc.Response[json.RawMessage]) error {
	return nil
}

func (g *gatewayConnectorMock) SignMessage(ctx context.Context, msg []byte) ([]byte, error) {
	return nil, nil
}

func (g *gatewayConnectorMock) GatewayIDs(ctx context.Context) ([]string, error) {
	return nil, nil
}

func (g *gatewayConnectorMock) DonID(ctx context.Context) (string, error) {
	return "", nil
}

func (g *gatewayConnectorMock) AwaitConnection(ctx context.Context, gatewayID string) error {
	return nil
}
