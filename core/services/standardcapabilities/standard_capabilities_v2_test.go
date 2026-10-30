package standardcapabilities

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry/remote"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	settingsgrpc "github.com/smartcontractkit/chainlink-common/pkg/settings/grpc"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

func TestStandardCapabilitiesV2_Start(t *testing.T) {
	lggr := logger.Test(t)
	local := registry.NewRegistry(lggr)
	settings := &settingsBroadcaster{update: core.SettingsUpdate{Settings: "a = 1", Hash: "h"}}

	var checked bool
	registrar := &capturingRegistrar{onRegister: func(cfg plugins.CmdConfig) {
		env := envMap(cfg.Env)
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

		t.Run("the process gets the settings", func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			updates, err := settingsgrpc.NewClient(lggr, dial(t, env["CL_SETTINGS_URL"])).Subscribe(ctx)
			require.NoError(t, err)
			assert.Equal(t, settings.update, <-updates)
			cancel()
			for range updates {
			}
		})
		checked = true
	}}

	s, err := NewStandardCapabilitiesV2(StandardCapabilitiesV2Params{
		Log:             lggr,
		Command:         "cron-v2",
		App:             "cron",
		ConfigJSON:      `{"fastestScheduleIntervalSeconds": 1}`,
		CapabilityDonID: 7,
		PluginRegistrar: registrar,
		Registry:        local,
		CRESettings:     settings,
	})
	require.NoError(t, err)
	require.ErrorContains(t, s.Start(t.Context()), capturingRegistrarErr)
	require.True(t, checked)
	assert.Empty(t, s.stopServers, "a failed start must stop the servers")
}

func TestStandardCapabilitiesV2_NoServers(t *testing.T) {
	var env map[string]string
	registrar := &capturingRegistrar{onRegister: func(cfg plugins.CmdConfig) { env = envMap(cfg.Env) }}
	s, err := NewStandardCapabilitiesV2(StandardCapabilitiesV2Params{Log: logger.Test(t), Command: "cron-v2", App: "cron", PluginRegistrar: registrar})
	require.NoError(t, err)
	require.ErrorContains(t, s.Start(t.Context()), capturingRegistrarErr)

	assert.NotContains(t, env, "CL_CAPABILITIES_PROXY_URL")
	assert.NotContains(t, env, "CL_SETTINGS_URL")
	assert.NotContains(t, env, "CL_CAPABILITIES_CAPABILITY_DON_ID")
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
		Lggr:   lggr,
		Listen: listenLocal,
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

type settingsBroadcaster struct {
	update core.SettingsUpdate
}

func (b *settingsBroadcaster) Subscribe(context.Context) (<-chan core.SettingsUpdate, error) {
	ch := make(chan core.SettingsUpdate, 1)
	ch <- b.update
	return ch, nil
}

func TestV2App(t *testing.T) {
	app, ok := v2App("/usr/local/bin/cron-v2")
	assert.True(t, ok)
	assert.Equal(t, "cron", app)

	_, ok = v2App("cron")
	assert.False(t, ok, "the old cron binary must keep the standard capabilities path")
}

func TestConfigEnv(t *testing.T) {
	t.Run("maps each key to CL_<APP>_<KEY>", func(t *testing.T) {
		got, err := configEnv("cron", `{"fastestScheduleIntervalSeconds": 1, "chainID": "a", "enabled": true, "big": 12345678901234567890}`)
		require.NoError(t, err)
		assert.Equal(t, []string{
			"CL_CRON_BIG=12345678901234567890",
			"CL_CRON_CHAIN_ID=a",
			"CL_CRON_ENABLED=true",
			"CL_CRON_FASTEST_SCHEDULE_INTERVAL_SECONDS=1",
		}, got)
	})

	t.Run("an empty config sets no env vars", func(t *testing.T) {
		for _, cfg := range []string{"", " ", "{}"} {
			got, err := configEnv("cron", cfg)
			require.NoError(t, err)
			assert.Empty(t, got)
		}
	})

	t.Run("rejects values that are not scalars", func(t *testing.T) {
		for _, cfg := range []string{`{"a": {"b": 1}}`, `{"a": [1]}`, `{"a": null}`} {
			_, err := configEnv("cron", cfg)
			assert.ErrorContains(t, err, `config key "a"`, cfg)
		}
	})

	t.Run("rejects config that is not a JSON object", func(t *testing.T) {
		_, err := configEnv("cron", `[1]`)
		assert.Error(t, err)
	})
}

func TestEnvName(t *testing.T) {
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
