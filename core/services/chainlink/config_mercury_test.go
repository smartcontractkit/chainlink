package chainlink

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

const (
	secretsMercury = `
[Mercury.Credentials.cred1]
URL = "https://chain1.link"
Username = "username1"
Password = "password1"

[Mercury.Credentials.cred2]
URL = "https://chain2.link"
Username = "username2"
Password = "password2"
`
)

func TestMercuryConfig(t *testing.T) {
	opts := GeneralConfigOpts{
		SecretsStrings: []string{secretsMercury},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	m := cfg.Mercury()
	assert.Equal(t, &types.MercuryCredentials{URL: "https://chain1.link", Username: "username1", Password: "password1"}, m.Credentials("cred1"))
	assert.Equal(t, &types.MercuryCredentials{URL: "https://chain2.link", Username: "username2", Password: "password2"}, m.Credentials("cred2"))
}

func TestMercuryTLS(t *testing.T) {
	certPath := "/path/to/cert.pem"
	transmission := toml.Mercury{
		TLS: toml.MercuryTLS{
			CertFile: &certPath,
		},
	}
	cfg := mercuryConfig{c: transmission}

	assert.Equal(t, certPath, cfg.TLS().CertFile())
}

func TestMercuryDataSourceConfig(t *testing.T) {
	t.Parallel()
	t.Run("defaults", func(t *testing.T) {
		t.Parallel()
		opts := GeneralConfigOpts{
			ConfigStrings: []string{`[Feature]
LogPoller = false`},
		}
		cfg, err := opts.New()
		require.NoError(t, err)

		assert.Equal(t, time.Duration(0), cfg.Mercury().DataSource().ObservationTimingBase())
	})

	t.Run("from full fixture", func(t *testing.T) {
		t.Parallel()
		opts := GeneralConfigOpts{
			ConfigStrings: []string{fullTOML},
		}
		cfg, err := opts.New()
		require.NoError(t, err)

		assert.Equal(t, 50*time.Millisecond, cfg.Mercury().DataSource().ObservationTimingBase())
	})
}
