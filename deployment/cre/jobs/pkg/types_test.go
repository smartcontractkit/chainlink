package pkg

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestUint64_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    Uint64
		wantErr bool
	}{
		{name: "numeric", input: `3`, want: 3},
		{name: "quoted string", input: `"3"`, want: 3},
		{name: "large numeric", input: `3379446385462418246`, want: 3379446385462418246},
		{name: "large quoted string", input: `"3379446385462418246"`, want: 3379446385462418246},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got Uint64
			err := json.Unmarshal([]byte(tt.input), &got)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestAuth0Config_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	t.Run("numeric tenantID", func(t *testing.T) {
		t.Parallel()

		var cfg Auth0Config
		require.NoError(t, json.Unmarshal([]byte(`{
			"issuerURL": "https://example.auth0.com/",
			"audience": "https://vault.example.com",
			"tenantID": 3
		}`), &cfg))
		require.Equal(t, Uint64(3), cfg.TenantID)
	})

	t.Run("string tenantID", func(t *testing.T) {
		t.Parallel()

		var cfg Auth0Config
		require.NoError(t, json.Unmarshal([]byte(`{
			"issuerURL": "https://example.auth0.com/",
			"audience": "https://vault.example.com",
			"tenantID": "3"
		}`), &cfg))
		require.Equal(t, Uint64(3), cfg.TenantID)
	})
}

func TestUint64_UnmarshalYAML(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    Uint64
		wantErr bool
	}{
		{name: "numeric", input: "value: 3\n", want: 3},
		{name: "quoted string", input: "value: \"3\"\n", want: 3},
		{name: "large numeric", input: "value: 3379446385462418246\n", want: 3379446385462418246},
		{name: "large quoted string", input: "value: \"3379446385462418246\"\n", want: 3379446385462418246},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got struct {
				Value Uint64 `yaml:"value"`
			}
			err := yaml.Unmarshal([]byte(tt.input), &got)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got.Value)
		})
	}
}

func TestAuth0Config_UnmarshalYAML(t *testing.T) {
	t.Parallel()

	t.Run("numeric tenantID", func(t *testing.T) {
		t.Parallel()

		var cfg Auth0Config
		require.NoError(t, yaml.Unmarshal([]byte(`
issuerURL: https://example.auth0.com/
audience: https://vault.example.com
tenantID: 3
`), &cfg))
		require.Equal(t, Uint64(3), cfg.TenantID)
	})

	t.Run("string tenantID", func(t *testing.T) {
		t.Parallel()

		var cfg Auth0Config
		require.NoError(t, yaml.Unmarshal([]byte(`
issuerURL: https://example.auth0.com/
audience: https://vault.example.com
tenantID: "3"
`), &cfg))
		require.Equal(t, Uint64(3), cfg.TenantID)
	})
}
