package pkg

import (
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/cresettings"
)

func TestCRESettingsJob_ResolveJob(t *testing.T) {
	tests := []struct {
		name     string
		settings string
	}{
		{"empty", ""},
		{"simple", "Foo = 42"},
		{"multi-line", `Foo = 42
[Bar]
Baz = "test"
`},
		{"single-quotes", `Foo = 'bar'
`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := CRESettingsJob{Settings: tt.settings}
			spec, err := j.ResolveJob()
			require.NoError(t, err)

			t.Log("Spec:", spec)

			got, err := cresettings.ValidatedCRESettingsSpec(spec)
			require.NoError(t, err)
			require.Equal(t, "config_type = \"settings\"\n"+tt.settings, got.CRESettingsSpec.Settings)

			var decoded struct {
				ConfigType string `toml:"config_type"`
			}
			require.NoError(t, toml.Unmarshal([]byte(got.CRESettingsSpec.Settings), &decoded))
			require.Equal(t, cresettings.ConfigTypeSettings, decoded.ConfigType)
		})
	}
}

func TestCapabilitiesRegistryJob_ResolveJob(t *testing.T) {
	tests := []struct {
		name           string
		offchainConfig string
	}{
		{
			name:           "simple",
			offchainConfig: `{"domain":"cre","env":"staging","version":1,"dons":{"workflow_1_zone-a":{"capabilities":{"cron@1.0.0":{"specConfig":{}}}}}}`,
		},
		{
			// A JSON string value containing quotes, single quotes, and backslashes: proves the
			// TOML encoding survives round-trip instead of corrupting the outer
			// settings = '''...''' wrapper cre-settings.tmpl applies (see ResolveJob's comment).
			name:           "special characters in a spec_config value",
			offchainConfig: `{"domain":"cre","env":"staging","version":1,"dons":{"workflow_1_zone-a":{"capabilities":{"cron@1.0.0":{"specConfig":{"fields":{"note":{"stringValue":"it's a \"quoted\" value, with a \\ backslash"}}}}}}}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := CapabilitiesRegistryJob{OffchainConfig: tt.offchainConfig}
			spec, err := j.ResolveJob()
			require.NoError(t, err)

			t.Log("Spec:", spec)

			got, err := cresettings.ValidatedCRESettingsSpec(spec)
			require.NoError(t, err)

			var decoded struct {
				ConfigType     string `toml:"config_type"`
				OffchainConfig string `toml:"offchain_config"`
			}
			require.NoError(t, toml.Unmarshal([]byte(got.CRESettingsSpec.Settings), &decoded))
			require.Equal(t, "capabilities_registry", decoded.ConfigType)
			require.Equal(t, tt.offchainConfig, decoded.OffchainConfig)
		})
	}
}

func TestCapabilitiesRegistryJob_ResolveJob_RejectsTripleQuote(t *testing.T) {
	j := CapabilitiesRegistryJob{OffchainConfig: `{"note":"it's a '''triple''' quote value"}`}
	_, err := j.ResolveJob()
	require.ErrorContains(t, err, "'''")
}
