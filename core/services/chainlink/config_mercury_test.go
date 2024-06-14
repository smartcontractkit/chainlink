package chainlink

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/models"
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
	assert.Equal(t, &models.MercuryCredentials{URL: "https://chain1.link", Username: "username1", Password: "password1"}, m.Credentials("cred1"))
	assert.Equal(t, &models.MercuryCredentials{URL: "https://chain2.link", Username: "username2", Password: "password2"}, m.Credentials("cred2"))
}
