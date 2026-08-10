package cresettings

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

func TestValidatedCRESettingsSpec(t *testing.T) {
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
config_type = "shard_assignment"
shard_assignment = '''
static_default_assignment = [0, 1]
hashed_default_assignment = false

[per_owner_assignment]
  "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266" = [1]
'''`, want: job.Job{
			SchemaVersion: 1,
			ExternalJobID: uuid.MustParse("7dcfa33b-8ed9-4e9f-9216-5b4d3f5c7887"),
			CRESettingsSpec: &job.CRESettingsSpec{
				ConfigType:      "shard_assignment",
				ShardAssignment: "static_default_assignment = [0, 1]\nhashed_default_assignment = false\n\n[per_owner_assignment]\n  \"0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266\" = [1]\n",
				Hash:            "fd86c9beb2526983703dd165002529f30aec8c38ec9eef91de3813513b0cf46b",
			},
		}},
		{name: "shard-assignment-invalid", toml: `type = "cresettings"
schemaVersion = 1
externalJobID = "7dcfa33b-8ed9-4e9f-9216-5b4d3f5c7887"
config_type = "shard_assignment"
shard_assignment = '''
static_default_assignment = [-1]
'''`, wantErr: "invalid shard_assignment config"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
