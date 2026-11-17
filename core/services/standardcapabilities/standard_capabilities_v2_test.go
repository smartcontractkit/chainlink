package standardcapabilities

import (
	"context"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry/remote"
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	settingsremote "github.com/smartcontractkit/chainlink-common/pkg/settings/remote"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

func TestStandardCapabilitiesV2_Start(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	local := registry.NewRegistry(lggr)
	creSettings := &settings.Atomic{}
	require.NoError(t, creSettings.Store(settings.Snapshot{Settings: `Foo = "initial"`, Hash: "initial"}))

	var checked bool
	var proxyURL string
	registrar := &capturingRegistrar{onRegister: func(cfg plugins.CmdConfig) {
		env := envMap(cfg.Env)
		proxyURL = env["CL_CAPABILITIES_PROXY_URL"]
		assert.Equal(t, "1", env["CL_CRON_FASTEST_SCHEDULE_INTERVAL_SECONDS"])
		assert.Equal(t, "7", env["CL_CAPABILITIES_CAPABILITY_DON_ID"])

		t.Run("the process can add a capability to the registry", func(t *testing.T) {
			proxy := remote.NewCapabilitiesRegistryClient(lggr, dial(t, env["CL_CAPABILITIES_PROXY_URL"]), clientTransport(t, lggr))
			id := "cron-trigger@1.0.0"
			require.NoError(t, proxy.Add(t.Context(), &trigger{info: capabilities.CapabilityInfo{ID: id, CapabilityType: capabilities.CapabilityTypeTrigger}}))

			stored, err := local.GetTrigger(t.Context(), id)
			require.NoError(t, err)
			info, err := stored.Info(t.Context())
			require.NoError(t, err)
			assert.Equal(t, id, info.ID)
		})

		// The subtest's context and connection end with it, so the stream does not hold up the
		// server stopping after the launch fails.
		t.Run("the process receives CRE settings", func(t *testing.T) {
			assert.Equal(t, proxyURL, env["CL_SETTINGS_URL"], "settings are served on the registry proxy's server")
			updates, err := settingsremote.NewSettingsClient(lggr, dial(t, env["CL_SETTINGS_URL"])).Subscribe(t.Context())
			require.NoError(t, err)
			requireSnapshot(t, updates, settings.Snapshot{Settings: `Foo = "initial"`, Hash: "initial"})

			next := settings.Snapshot{Settings: `Foo = "next"`, Hash: "next"}
			require.NoError(t, creSettings.Store(next))
			requireSnapshot(t, updates, next)
		})
		checked = true
	}}

	s, err := NewV2(V2Params{
		Log:              lggr,
		Command:          "cronv2",
		App:              "cron",
		ConfigJSON:       `{"fastestScheduleIntervalSeconds": 1}`,
		CapDonIDResolver: &staticCapDonIDResolver{donID: 7},
		PluginRegistrar:  registrar,
		Registry:         local,
		CRESettings:      creSettings,
	})
	require.NoError(t, err)
	require.NoError(t, s.Start(t.Context()))
	s.launchWG.Wait()
	require.True(t, checked)
	require.ErrorContains(t, s.HealthReport()[s.Name()], capturingRegistrarErr)

	_, err = (&net.Dialer{}).DialContext(t.Context(), "tcp", proxyURL)
	require.Error(t, err, "a failed launch must stop the server")
	require.NoError(t, s.Close())
}

func TestStandardCapabilitiesV2_NoServers(t *testing.T) {
	t.Parallel()
	var env map[string]string
	registrar := &capturingRegistrar{onRegister: func(cfg plugins.CmdConfig) { env = envMap(cfg.Env) }}
	s, err := NewV2(V2Params{Log: logger.Test(t), Command: "cronv2", App: "cron", PluginRegistrar: registrar, CapDonIDResolver: &staticCapDonIDResolver{}})
	require.NoError(t, err)
	require.NoError(t, s.Start(t.Context()))
	s.launchWG.Wait()
	require.NoError(t, s.Close())

	assert.NotContains(t, env, "CL_CAPABILITIES_PROXY_URL")
	assert.NotContains(t, env, "CL_SETTINGS_URL")
	assert.NotContains(t, env, "CL_CAPABILITIES_CAPABILITY_DON_ID")
}

func TestDelegate_NewServicesV2_ForwardsCRESettings(t *testing.T) {
	t.Parallel()
	creSettings := &settings.Atomic{}
	d := &Delegate{logger: logger.Test(t), creSettings: creSettings}
	svcs, err := d.newServicesV2(t.Context(), logger.Test(t), "cronv2", "cron", "", 7)
	require.NoError(t, err)
	require.Len(t, svcs, 1)
	v2, ok := svcs[0].(*V2)
	require.True(t, ok)
	assert.Same(t, creSettings, v2.creSettings)
}

func requireSnapshot(t *testing.T, updates <-chan settings.Snapshot, want settings.Snapshot) {
	t.Helper()
	select {
	case got := <-updates:
		require.Equal(t, want, got)
	case <-time.After(5 * time.Second):
		t.Fatalf("timed out waiting for settings %q", want.Hash)
	}
}

// chanCapDonIDResolver blocks until a DON ID is sent, so tests control when it resolves.
type chanCapDonIDResolver chan uint32

func (c chanCapDonIDResolver) awaitCapDonID(ctx context.Context) (uint32, error) {
	select {
	case id := <-c:
		return id, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

func TestStandardCapabilitiesV2_LaunchWaitsForDonID(t *testing.T) {
	t.Parallel()
	registered := make(chan map[string]string, 1)
	registrar := &capturingRegistrar{onRegister: func(cfg plugins.CmdConfig) { registered <- envMap(cfg.Env) }}
	resolver := make(chanCapDonIDResolver)
	s, err := NewV2(V2Params{
		Log: logger.Test(t), Command: "cronv2", App: "cron", PluginRegistrar: registrar, CapDonIDResolver: resolver,
	})
	require.NoError(t, err)
	require.NoError(t, s.Start(t.Context()))

	select {
	case <-registered:
		t.Fatal("the process must not be launched before the DON ID is resolved")
	case <-time.After(50 * time.Millisecond):
	}

	resolver <- 9
	select {
	case env := <-registered:
		assert.Equal(t, "9", env["CL_CAPABILITIES_CAPABILITY_DON_ID"])
	case <-time.After(5 * time.Second):
		t.Fatal("the process was not launched after the DON ID was resolved")
	}
	require.NoError(t, s.Close())
}

func TestStandardCapabilitiesV2_CloseDuringDonIDResolution(t *testing.T) {
	t.Parallel()
	var launched bool
	registrar := &capturingRegistrar{onRegister: func(plugins.CmdConfig) { launched = true }}
	s, err := NewV2(V2Params{
		Log: logger.Test(t), Command: "cronv2", App: "cron", PluginRegistrar: registrar, CapDonIDResolver: make(chanCapDonIDResolver),
	})
	require.NoError(t, err)
	require.NoError(t, s.Start(t.Context()))

	require.NoError(t, s.Close())
	assert.False(t, launched, "a capability closed while resolving its DON ID must not launch its process")
}

func envMap(env []string) map[string]string {
	m := make(map[string]string, len(env))
	for _, kv := range env {
		k, v, _ := strings.Cut(kv, "=")
		m[k] = v
	}
	return m
}

func dial(t *testing.T, target string) *grpc.ClientConn {
	t.Helper()
	require.NotEmpty(t, target)
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })
	return conn
}

// clientTransport is the transport of the process side: it publishes the capabilities the
// process adds, for the host to dial.
func clientTransport(t *testing.T, lggr logger.Logger) remote.Transport {
	t.Helper()
	transport, err := remote.NewListenerTransport(remote.ListenerTransportParams{
		Lggr: lggr,
		Listen: func() (net.Listener, error) {
			return listenLocal(t.Context())
		},
		Connect: func(target string) (*grpc.ClientConn, error) {
			return grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
		},
	})
	require.NoError(t, err)
	return transport
}

type trigger struct {
	info capabilities.CapabilityInfo
}

func (c *trigger) Info(context.Context) (capabilities.CapabilityInfo, error) { return c.info, nil }

func (c *trigger) RegisterTrigger(context.Context, capabilities.TriggerRegistrationRequest) (<-chan capabilities.TriggerResponse, error) {
	return make(chan capabilities.TriggerResponse), nil
}

func (c *trigger) UnregisterTrigger(context.Context, capabilities.TriggerRegistrationRequest) error {
	return nil
}

func (c *trigger) AckEvent(context.Context, string, string, string) error { return nil }

func TestDelegate_CronV2Active(t *testing.T) {
	t.Parallel()
	now := time.Now()
	for name, tc := range map[string]struct {
		settings settings.Getter
		want     bool
	}{
		"default": {},
		"inactive": {settings: activePeriodGetter{
			Lower: commonconfig.Timestamp(now.Add(time.Hour).Unix()),
			Upper: commonconfig.Timestamp(now.Add(2 * time.Hour).Unix()),
		}},
		"active": {settings: activePeriodGetter{
			Lower: commonconfig.Timestamp(now.Add(-time.Hour).Unix()),
			Upper: commonconfig.Timestamp(now.Add(time.Hour).Unix()),
		}, want: true},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			d := &Delegate{logger: logger.Test(t), limitsFactory: limits.Factory{Settings: tc.settings}}
			assert.Equal(t, tc.want, d.cronV2Active(t.Context()))
		})
	}
}

type activePeriodGetter settings.Range[commonconfig.Timestamp]

func (g activePeriodGetter) GetScoped(context.Context, settings.Scope, string) (string, error) {
	return settings.Range[commonconfig.Timestamp](g).String(), nil
}

func TestEnvVarsFromConfigJSON(t *testing.T) {
	t.Parallel()
	t.Run("maps each key to CL_<APP>_<KEY>", func(t *testing.T) {
		t.Parallel()
		got, err := envVarsFromConfigJSON("cron", `{"fastestScheduleIntervalSeconds": 1, "chainID": "a", "enabled": true, "big": 12345678901234567890}`)
		require.NoError(t, err)
		assert.Equal(t, []string{
			"CL_CRON_BIG=12345678901234567890",
			"CL_CRON_CHAIN_ID=a",
			"CL_CRON_ENABLED=true",
			"CL_CRON_FASTEST_SCHEDULE_INTERVAL_SECONDS=1",
		}, got)
	})

	t.Run("an empty config sets no env vars", func(t *testing.T) {
		t.Parallel()
		for _, cfg := range []string{"", " ", "{}"} {
			got, err := envVarsFromConfigJSON("cron", cfg)
			require.NoError(t, err)
			assert.Empty(t, got)
		}
	})

	t.Run("rejects values that are not scalars", func(t *testing.T) {
		t.Parallel()
		for _, cfg := range []string{`{"a": {"b": 1}}`, `{"a": [1]}`, `{"a": null}`} {
			_, err := envVarsFromConfigJSON("cron", cfg)
			assert.ErrorContains(t, err, `config key "a"`, cfg)
		}
	})

	t.Run("rejects config that is not a JSON object", func(t *testing.T) {
		t.Parallel()
		_, err := envVarsFromConfigJSON("cron", `[1]`)
		assert.Error(t, err)
	})
}

func TestEnvName(t *testing.T) {
	t.Parallel()
	for in, want := range map[string]string{
		"cron":                           "CRON",
		"fastestScheduleIntervalSeconds": "FASTEST_SCHEDULE_INTERVAL_SECONDS",
		"chainID":                        "CHAIN_ID",
		"v2Value":                        "V2_VALUE",
		"fastest-schedule.interval":      "FASTEST_SCHEDULE_INTERVAL",
	} {
		assert.Equal(t, want, envName(in), in)
	}
}
