package models

import (
	"github.com/smartcontractkit/chainlink-common/pkg/config"
)

// Secret is a string that formats and encodes redacted, as "xxxxx".
//
// Deprecated: use config.SecretString instead.
type Secret = config.SecretString

// Deprecated: use config.NewSecretString instead.
func NewSecret(s string) *Secret { return config.NewSecretString(s) }

// SecretURL is a URL that formats and encodes redacted, as "xxxxx".
//
// Deprecated: use config.SecretURL instead.
type SecretURL = config.SecretURL

// Deprecated: use config.NewSecretURL instead.
func NewSecretURL(u *config.URL) *config.SecretURL { return config.NewSecretURL(u) }

// Deprecated: use config.MustSecretURL instead.
func MustSecretURL(u string) *config.SecretURL {
	return NewSecretURL(config.MustParseURL(u))
}
