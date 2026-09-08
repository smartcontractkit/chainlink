package chainlink

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

func TestCapabilitiesConfig(t *testing.T) {
	opts := GeneralConfigOpts{
		ConfigStrings: []string{fullTOML},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	p2p := cfg.Capabilities().Peering()
	assert.Equal(t, "p2p_12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw", p2p.PeerID().String())
	assert.Equal(t, 13, p2p.IncomingMessageBufferSize())
	assert.Equal(t, 17, p2p.OutgoingMessageBufferSize())
	assert.True(t, p2p.TraceLogging())

	v2 := p2p.V2()
	assert.False(t, v2.Enabled())
	assert.Equal(t, []string{"a", "b", "c"}, v2.AnnounceAddresses())
	assert.ElementsMatch(
		t,
		[]commontypes.BootstrapperLocator{
			{
				PeerID: "12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw",
				Addrs:  []string{"test:99"},
			},
			{
				PeerID: "12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw",
				Addrs:  []string{"foo:42", "bar:10"},
			},
		},
		v2.DefaultBootstrappers(),
	)
	assert.Equal(t, time.Minute, v2.DeltaDial().Duration())
	assert.Equal(t, 2*time.Second, v2.DeltaReconcile().Duration())
	assert.Equal(t, []string{"foo", "bar"}, v2.ListenAddresses())
}

func TestCapabilitiesLocalConfig(t *testing.T) {
	opts := GeneralConfigOpts{
		ConfigStrings: []string{fullTOML},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	local := cfg.Capabilities().Local()

	// Test RegistryBasedLaunchAllowlist - now contains regex patterns
	allowlist := local.RegistryBasedLaunchAllowlist()
	assert.Equal(t, []string{"^cron@1\\.0\\.0$", "^http-action@.*$"}, allowlist)

	// Test IsAllowlisted with regex matching
	assert.True(t, local.IsAllowlisted("cron@1.0.0"))        // exact match via regex
	assert.False(t, local.IsAllowlisted("cron@2.0.0"))       // version mismatch
	assert.True(t, local.IsAllowlisted("http-action@1.0.0")) // matches any version
	assert.True(t, local.IsAllowlisted("http-action@2.0.0")) // matches any version
	assert.False(t, local.IsAllowlisted("unknown@1.0.0"))

	// Test Capabilities map
	capabilities := local.Capabilities()
	require.NotNil(t, capabilities)
	assert.Len(t, capabilities, 2)

	// Test http-action config
	httpAction := local.GetCapabilityConfig("http-action@1.0.0")
	require.NotNil(t, httpAction)
	assert.Equal(t, "/opt/chainlink/binaries/http_action", httpAction.BinaryPathOverride())
	assert.Equal(t, "gateway", httpAction.Config()["proxyMode"])
	assert.Equal(t, "443,8443", httpAction.Config()["allowedPorts"])

	// Test cron config
	cronConfig := local.GetCapabilityConfig("cron@1.0.0")
	require.NotNil(t, cronConfig)
	assert.Equal(t, "/opt/chainlink/binaries/cron", cronConfig.BinaryPathOverride())
	assert.Equal(t, "60", cronConfig.Config()["fastestScheduleIntervalSeconds"])

	// Test non-existent capability
	unknownConfig := local.GetCapabilityConfig("unknown@1.0.0")
	assert.Nil(t, unknownConfig)
}

func TestCapabilitiesLocalConfigEmpty(t *testing.T) {
	tomlStr := `
[Capabilities.Local]
`
	opts := GeneralConfigOpts{
		ConfigStrings: []string{tomlStr},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	local := cfg.Capabilities().Local()
	assert.Empty(t, local.RegistryBasedLaunchAllowlist())
	assert.Nil(t, local.Capabilities())
	assert.False(t, local.IsAllowlisted("any@1.0.0"))
	assert.Nil(t, local.GetCapabilityConfig("any@1.0.0"))
}

func TestHTTPCapabilitiesConfig(t *testing.T) {
	t.Parallel()

	opts := GeneralConfigOpts{
		ConfigStrings: []string{fullTOML},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	ht := cfg.Capabilities().HTTPTrigger()
	assert.Equal(t, uint16(25), ht.MetadataBatchSize())
	assert.Equal(t, uint16(500), ht.SendChannelBufferSize())
	assert.Equal(t, uint16(10), ht.MaxAuthorizedKeysPerWorkflow())
	assert.Equal(t, uint32(3600), ht.RequestCacheTTL())

	htGw := ht.GatewayConnection()
	assert.Equal(t, uint32(45000), htGw.MaxPushMetadataDurationMs())
	assert.Equal(t, uint32(46000), htGw.MaxPullMetadataDurationMs())

	htRetry := htGw.RetryConfig()
	assert.Equal(t, 200, htRetry.InitialIntervalMs())
	assert.Equal(t, 60000, htRetry.MaxIntervalTimeMs())
	assert.Equal(t, 3.0, htRetry.Multiplier())

	ha := cfg.Capabilities().HTTPAction()
	assert.Equal(t, "direct", ha.ProxyMode())

	haGw := ha.GatewayConnection()
	assert.Equal(t, uint32(200), haGw.InitialIntervalMs())
	assert.Equal(t, uint32(60000), haGw.MaxElapsedTimeMs())
	assert.Equal(t, 3.0, haGw.Multiplier())

	haHc := ha.HTTPClient()
	assert.Equal(t, []string{"10.0.0.1"}, haHc.BlockedIPs())
	assert.Equal(t, []string{"10.0.0.0/8"}, haHc.BlockedIPsCIDR())
	assert.Equal(t, []int{8443, 9443}, haHc.AllowedPorts())
	assert.Equal(t, []string{"http", "https"}, haHc.AllowedSchemes())
	assert.Equal(t, []string{"1.2.3.4"}, haHc.AllowedIPs())
	assert.Equal(t, []string{"1.2.3.0/24"}, haHc.AllowedIPsCIDR())
}

func TestHTTPCapabilitiesConfig_DefaultsWhenOmitted(t *testing.T) {
	t.Parallel()

	// Empty config: the TOML layer applies the documented defaults, matching
	// the capability binaries' built-in defaults exactly. These are the same
	// values the capabilities previously fell back to when job specs omitted
	// the fields, so runtime behavior is unchanged.
	tomlStr := ``
	opts := GeneralConfigOpts{
		ConfigStrings: []string{tomlStr},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	ht := cfg.Capabilities().HTTPTrigger()
	assert.Equal(t, uint16(50), ht.MetadataBatchSize())
	assert.Equal(t, uint16(1000), ht.SendChannelBufferSize())
	assert.Equal(t, uint16(100), ht.MaxAuthorizedKeysPerWorkflow())
	assert.Equal(t, uint32(86400), ht.RequestCacheTTL())
	assert.Equal(t, uint32(30000), ht.GatewayConnection().MaxPushMetadataDurationMs())
	assert.Equal(t, uint32(30000), ht.GatewayConnection().MaxPullMetadataDurationMs())
	assert.Equal(t, 100, ht.GatewayConnection().RetryConfig().InitialIntervalMs())
	assert.Equal(t, 30000, ht.GatewayConnection().RetryConfig().MaxIntervalTimeMs())
	assert.Equal(t, 2.0, ht.GatewayConnection().RetryConfig().Multiplier())

	ha := cfg.Capabilities().HTTPAction()
	assert.Equal(t, "gateway", ha.ProxyMode())
	assert.Equal(t, uint32(100), ha.GatewayConnection().InitialIntervalMs())
	assert.Equal(t, uint32(30000), ha.GatewayConnection().MaxElapsedTimeMs())
	assert.Equal(t, 2.0, ha.GatewayConnection().Multiplier())
	assert.Empty(t, ha.HTTPClient().BlockedIPs())
	assert.Empty(t, ha.HTTPClient().BlockedIPsCIDR())
	assert.Equal(t, []int{443}, ha.HTTPClient().AllowedPorts())
	assert.Equal(t, []string{"https"}, ha.HTTPClient().AllowedSchemes())
	assert.Empty(t, ha.HTTPClient().AllowedIPs())
	assert.Empty(t, ha.HTTPClient().AllowedIPsCIDR())
}

func TestHTTPActionConfig_InvalidProxyMode(t *testing.T) {
	t.Parallel()

	tomlStr := `
[Capabilities.HTTPAction]
ProxyMode = 'invalid-mode'
`
	opts := GeneralConfigOpts{
		ConfigStrings: []string{tomlStr},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	err = cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Capabilities.HTTPAction.ProxyMode")
	assert.Contains(t, err.Error(), "must be either 'gateway' or 'direct'")
}

func TestValidateCapabilityID(t *testing.T) {
	tests := []struct {
		name    string
		capID   string
		wantErr bool
	}{
		{"valid simple", "cron@1.0.0", false},
		{"valid with hyphen", "http-action@1.0.0", false},
		{"valid with prerelease", "http-action@1.0.0-alpha", false},
		{"valid complex version", "my-capability@10.20.30", false},
		{"invalid missing version", "cron", true},
		{"invalid missing name", "@1.0.0", true},
		{"invalid uppercase", "Cron@1.0.0", true},
		{"invalid version format", "cron@1.0", true},
		{"invalid special chars", "cron!@1.0.0", true},
		{"invalid spaces", "cron @1.0.0", true},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := toml.ValidateCapabilityID(tt.capID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLocalCapabilitiesValidation(t *testing.T) {
	t.Run("valid config with regex patterns", func(t *testing.T) {
		cfg := toml.LocalCapabilities{
			RegistryBasedLaunchAllowlist: []string{"^cron@1\\.0\\.0$", "^http-action@.*$", ".*"},
			Capabilities: map[string]toml.CapabilityNodeConfig{
				"cron@1.0.0": {},
			},
		}
		err := cfg.ValidateConfig()
		assert.NoError(t, err)
	})

	t.Run("invalid regex pattern", func(t *testing.T) {
		cfg := toml.LocalCapabilities{
			RegistryBasedLaunchAllowlist: []string{"[invalid"},
		}
		err := cfg.ValidateConfig()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "RegistryBasedLaunchAllowlist")
		assert.Contains(t, err.Error(), "invalid regex pattern")
	})

	t.Run("invalid capabilities key", func(t *testing.T) {
		cfg := toml.LocalCapabilities{
			Capabilities: map[string]toml.CapabilityNodeConfig{
				"invalid": {},
			},
		}
		err := cfg.ValidateConfig()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Capabilities.Local.Capabilities")
	})

	t.Run("multiple errors", func(t *testing.T) {
		cfg := toml.LocalCapabilities{
			RegistryBasedLaunchAllowlist: []string{"[invalid1", "[invalid2"},
			Capabilities: map[string]toml.CapabilityNodeConfig{
				"also-invalid": {},
			},
		}
		err := cfg.ValidateConfig()
		assert.Error(t, err)
	})
}
