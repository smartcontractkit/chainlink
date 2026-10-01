package cresettings

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

func TestValidatedCRESettingsSpec(t *testing.T) {
	t.Parallel()
	settingsString := `Foo = "bar"
`
	noHash := fmt.Sprintf(`type = "cresettings"
schemaVersion = 1
externalJobID = "7dcfa33b-8ed9-4e9f-9216-5b4d3f5c7887"
settings = '''
%s'''`, settingsString)
	hashFn := func(hash string) string {
		return fmt.Sprintf(`type = "cresettings"
schemaVersion = 1
externalJobID = "7dcfa33b-8ed9-4e9f-9216-5b4d3f5c7887"
hash = "%s"
settings = '''
%s'''`, hash, settingsString)
	}
	tests := []struct {
		name    string
		toml    string
		want    job.Job
		wantErr string
	}{
		{name: "empty", toml: `type = "cresettings"
schemaVersion = 1
externalJobID = "7dcfa33b-8ed9-4e9f-9216-5b4d3f5c7887"`, want: job.Job{
			SchemaVersion: 1,
			ExternalJobID: uuid.MustParse("7dcfa33b-8ed9-4e9f-9216-5b4d3f5c7887"),
			CRESettingsSpec: &job.CRESettingsSpec{
				Settings: "",
				Hash:     "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			},
		}},
		{name: "hash", toml: hashFn(`dca1decf4530b4aaa0b6cf8f11ed62f3af75b4ecc4915e4a35292c6918e9160a`), want: job.Job{
			SchemaVersion: 1,
			ExternalJobID: uuid.MustParse("7dcfa33b-8ed9-4e9f-9216-5b4d3f5c7887"),
			CRESettingsSpec: &job.CRESettingsSpec{
				Settings: settingsString,
				Hash:     "dca1decf4530b4aaa0b6cf8f11ed62f3af75b4ecc4915e4a35292c6918e9160a",
			},
		}},
		{name: "no-hash", toml: noHash, want: job.Job{
			SchemaVersion: 1,
			ExternalJobID: uuid.MustParse("7dcfa33b-8ed9-4e9f-9216-5b4d3f5c7887"),
			CRESettingsSpec: &job.CRESettingsSpec{
				Settings: settingsString,
				Hash:     "dca1decf4530b4aaa0b6cf8f11ed62f3af75b4ecc4915e4a35292c6918e9160a",
			},
		}},
		{name: "wrong-type", toml: `type = "asdf"`, wantErr: "unsupported type"},
		{name: "wrong-hash", toml: hashFn(`asdf`), wantErr: "invalid sha256 hash"},
		{name: "shard-assignment", toml: `type = "cresettings"
schemaVersion = 1
externalJobID = "7dcfa33b-8ed9-4e9f-9216-5b4d3f5c7887"
settings = '''
config_type = "shard_assignment"
static_default_assignment = [0, 1]
hashed_default_assignment = false

[per_owner_assignment]
  "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266" = [1]
'''`, want: job.Job{
			SchemaVersion: 1,
			ExternalJobID: uuid.MustParse("7dcfa33b-8ed9-4e9f-9216-5b4d3f5c7887"),
			CRESettingsSpec: &job.CRESettingsSpec{
				Settings: "config_type = \"shard_assignment\"\nstatic_default_assignment = [0, 1]\nhashed_default_assignment = false\n\n[per_owner_assignment]\n  \"0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266\" = [1]\n",
				Hash:     "b82dfa99b83c995c155ee7f067f5469eec95b60c02311292da995e898c74b3c9",
			},
		}},
		{name: "shard-assignment-invalid", toml: `type = "cresettings"
schemaVersion = 1
externalJobID = "7dcfa33b-8ed9-4e9f-9216-5b4d3f5c7887"
settings = '''
config_type = "shard_assignment"
static_default_assignment = [-1]
'''`, wantErr: "invalid shard_assignment config"},
		{name: "unknown-config-type", toml: `type = "cresettings"
schemaVersion = 1
externalJobID = "7dcfa33b-8ed9-4e9f-9216-5b4d3f5c7887"
settings = '''
config_type = "unknown"
Foo = "bar"
'''`, wantErr: "unknown config_type"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.want.Type = job.CRESettings

			got, err := ValidatedCRESettingsSpec(tt.toml)
			if tt.wantErr != "" {
				require.ErrorContains(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestValidatedCRESettingsSpec_CapabilitiesRegistry(t *testing.T) {
	t.Parallel()

	specFn := func(offchain string) string {
		return fmt.Sprintf(`type = "cresettings"
schemaVersion = 1
externalJobID = "7dcfa33b-8ed9-4e9f-9216-5b4d3f5c7887"
config_type = "capabilities_registry"
offchain_config = '''%s'''`, offchain)
	}

	t.Run("valid", func(t *testing.T) {
		t.Parallel()
		got, err := ValidatedCRESettingsSpec(specFn(`{"version":1,"dons":{}}`))
		require.NoError(t, err)
		require.NotNil(t, got.CRESettingsSpec)
		assert.Equal(t, ConfigTypeCapRegistry, got.CRESettingsSpec.ConfigType)
		assert.Equal(t, `{"version":1,"dons":{}}`, got.CRESettingsSpec.OffchainConfig)
		// Hash is computed over OffchainConfig (not Settings).
		assert.NotEmpty(t, got.CRESettingsSpec.Hash)
		assert.Empty(t, got.CRESettingsSpec.Settings)
	})

	t.Run("invalid payload", func(t *testing.T) {
		t.Parallel()
		_, err := ValidatedCRESettingsSpec(specFn(`not json`))
		require.ErrorContains(t, err, "invalid capabilities_registry config")
	})

	t.Run("zero version rejected", func(t *testing.T) {
		t.Parallel()
		_, err := ValidatedCRESettingsSpec(specFn(`{"dons":{}}`))
		require.ErrorContains(t, err, "version must be >= 1")
	})

	t.Run("hash is over offchain_config and must match when provided", func(t *testing.T) {
		t.Parallel()
		payload := `{"version":2}`
		sum := sha256.Sum256([]byte(payload))
		want := hex.EncodeToString(sum[:])
		got, err := ValidatedCRESettingsSpec(specFn(payload))
		require.NoError(t, err)
		assert.Equal(t, want, got.CRESettingsSpec.Hash)

		_, err = ValidatedCRESettingsSpec(specFn(payload) + "\nhash = \"" + want + "\"")
		require.NoError(t, err)
		_, err = ValidatedCRESettingsSpec(specFn(payload) + "\nhash = \"deadbeef\"")
		require.ErrorContains(t, err, "invalid sha256 hash")
	})

	t.Run("settings with capabilities_registry rejected", func(t *testing.T) {
		t.Parallel()
		_, err := ValidatedCRESettingsSpec(specFn(`{"version":1}`) + "\nsettings = '''Foo = \"bar\"'''")
		require.ErrorContains(t, err, "settings must be empty")
	})

	t.Run("offchain_config with another config_type rejected", func(t *testing.T) {
		t.Parallel()
		_, err := ValidatedCRESettingsSpec(`type = "cresettings"
schemaVersion = 1
settings = '''Foo = "bar"'''
offchain_config = '''{"version":1}'''`)
		require.ErrorContains(t, err, "offchain_config is only valid")
	})

	t.Run("routing agrees with delegate", func(t *testing.T) {
		t.Parallel()
		got, err := ValidatedCRESettingsSpec(specFn(`{"version":1}`))
		require.NoError(t, err)
		assert.Equal(t, ConfigTypeCapRegistry, (&delegate{}).configType(got.CRESettingsSpec))
	})
}

// TestValidatedCRESettingsSpec_CapRegistryDiscriminatorForms enumerates the ways a spec could
// try to select capabilities_registry. Only the exact top-level config_type field is accepted,
// and it is persisted verbatim in cre_settings_specs.config_type, the column covered by the
// single-job unique index and the high-water check.
func TestValidatedCRESettingsSpec_CapRegistryDiscriminatorForms(t *testing.T) {
	t.Parallel()

	const header = "type = \"cresettings\"\nschemaVersion = 1\n"
	for _, tc := range []struct {
		name    string
		body    string
		wantErr string
	}{
		{"top-level field", "config_type = \"capabilities_registry\"\noffchain_config = '''{\"version\":1}'''", ""},
		{"embedded in settings", "settings = '''config_type = \"capabilities_registry\"'''", "top-level config_type"},
		{"embedded in settings with offchain_config", "settings = '''config_type = \"capabilities_registry\"'''\noffchain_config = '''{\"version\":1}'''", "top-level config_type"},
		{"top-level plus embedded settings", "config_type = \"capabilities_registry\"\nsettings = '''config_type = \"capabilities_registry\"'''\noffchain_config = '''{\"version\":1}'''", "settings must be empty"},
		{"upper case", "config_type = \"Capabilities_Registry\"\noffchain_config = '''{\"version\":1}'''", "unknown config_type"},
		{"padded", "config_type = \" capabilities_registry\"\noffchain_config = '''{\"version\":1}'''", "unknown config_type"},
		{"settings type with payload", "config_type = \"settings\"\nsettings = '''Foo = \"bar\"'''\noffchain_config = '''{\"version\":1}'''", "offchain_config is only valid"},
		{"no type with payload", "offchain_config = '''{\"version\":1}'''", "offchain_config is only valid"},
		{"top-level without payload", "config_type = \"capabilities_registry\"", "empty payload"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			jb, err := ValidatedCRESettingsSpec(header + tc.body)
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, job.CRESettingsConfigTypeCapRegistry, jb.CRESettingsSpec.ConfigType)
		})
	}
}
