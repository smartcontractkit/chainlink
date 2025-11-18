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
			verifierID: secret.VerifierID,
			apiKey:     string(*secret.APIKey),
			apiSecret:  string(*secret.APISecret),
		}
	}
	return secrets
}

func (c *ccvConfig) IndexerSecret() config.IndexerSecret {
	return &indexerSecretConfig{
		apiKey:    string(*c.s.IndexerSecret.APIKey),
		apiSecret: string(*c.s.IndexerSecret.APISecret),
	}
}

type indexerSecretConfig struct {
	apiKey    string
	apiSecret string
}

func (i *indexerSecretConfig) APIKey() string {
	return i.apiKey
}

func (i *indexerSecretConfig) APISecret() string {
	return i.apiSecret
}

type aggregatorSecretConfig struct {
	verifierID string
	apiKey     string
	apiSecret  string
}

func (a *aggregatorSecretConfig) VerifierID() string {
	return a.verifierID
}

func (a *aggregatorSecretConfig) APIKey() string {
	return a.apiKey
}

func (a *aggregatorSecretConfig) APISecret() string {
	return a.apiSecret
}
