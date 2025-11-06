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
	base := fmt.Sprintf(`type = "cresettings"
externalJobID = "7dcfa33b-8ed9-4e9f-9216-5b4d3f5c7887"
Settings = '''
%s'''`, settingsString)
	tests := []struct {
		name    string
		toml    string
		want    job.Job
		wantErr string
	}{
		{name: "empty", toml: `type = "cresettings"
externalJobID = "7dcfa33b-8ed9-4e9f-9216-5b4d3f5c7887"`, want: job.Job{
			ID:            0,
			ExternalJobID: uuid.MustParse("7dcfa33b-8ed9-4e9f-9216-5b4d3f5c7887"),
			CRESettingsSpec: &job.CRESettingsSpec{
				Settings: "",
				Hash:     "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			},
		}},
		{name: "hash", toml: base + `
Hash = "dca1decf4530b4aaa0b6cf8f11ed62f3af75b4ecc4915e4a35292c6918e9160a"
`, want: job.Job{
			ID:            0,
			ExternalJobID: uuid.MustParse("7dcfa33b-8ed9-4e9f-9216-5b4d3f5c7887"),
			CRESettingsSpec: &job.CRESettingsSpec{
				Settings: settingsString,
				Hash:     "dca1decf4530b4aaa0b6cf8f11ed62f3af75b4ecc4915e4a35292c6918e9160a",
			},
		}},
		{name: "no-hash", toml: base, want: job.Job{
			ID:            0,
			ExternalJobID: uuid.MustParse("7dcfa33b-8ed9-4e9f-9216-5b4d3f5c7887"),
			CRESettingsSpec: &job.CRESettingsSpec{
				Settings: settingsString,
				Hash:     "dca1decf4530b4aaa0b6cf8f11ed62f3af75b4ecc4915e4a35292c6918e9160a",
			},
		}},
		{name: "wrong-type", toml: `type = "asdf"`, wantErr: "unsupported type"},
		{name: "wrong-hash", toml: base + `
Hash = "asdf"`, wantErr: "invalid sha256 hash"},
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
