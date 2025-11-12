package chainlink

import (
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

type ccvConfig struct {
	s toml.CCVSecrets
}

func (c *ccvConfig) AggregatorSecrets() []config.AggregatorSecret {
	secrets := make([]config.AggregatorSecret, len(c.s.AggregatorSecrets))
	for i, secret := range c.s.AggregatorSecrets {
		secrets[i] = &aggregatorSecretConfig{
			committeeID: secret.CommitteeID,
			apiKey:      string(*secret.APIKey),
			apiSecret:   string(*secret.APISecret),
		}
	}
	return secrets
}

type aggregatorSecretConfig struct {
	committeeID string
	apiKey      string
	apiSecret   string
}

func (a *aggregatorSecretConfig) CommitteeID() string {
	return a.committeeID
}

func (a *aggregatorSecretConfig) APIKey() string {
	return a.apiKey
}

func (a *aggregatorSecretConfig) APISecret() string {
	return a.apiSecret
}
