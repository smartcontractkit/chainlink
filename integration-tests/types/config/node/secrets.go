package node

import (
	"fmt"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

var SecretsConf = &chainlink.Secrets{
	Secrets: toml.Secrets{
		Database: toml.DatabaseSecrets{
			AllowSimplePasswords: ptr(true),
		},
		Password: toml.Passwords{
			Keystore: models.NewSecret("................"),
		},
	},
}

type SecretsConfigOpt = func(c *chainlink.Secrets)

func NewSecretsConfig(baseConf *chainlink.Secrets, opts ...SecretsConfigOpt) *chainlink.Secrets {
	for _, opt := range opts {
		opt(baseConf)
	}
	return baseConf
}

func WithDBURL(host, port, dbname string) SecretsConfigOpt {
	return func(c *chainlink.Secrets) {
		c.Secrets.Database.URL = models.MustSecretURL(
			fmt.Sprintf("postgresql://postgres:test@%s:%s/%s?sslmode=disable", host, port, dbname),
		)
	}
}

func WithMercurySecrets(credMapKey, url, username, password string) SecretsConfigOpt {
	return func(c *chainlink.Secrets) {
		c.Secrets.Mercury = toml.MercurySecrets{
			Credentials: map[string]toml.MercuryCredentials{
				credMapKey: {
					URL:      models.MustSecretURL(url),
					Username: models.NewSecret(username),
					Password: models.NewSecret(password),
				},
			},
		}
	}
}
