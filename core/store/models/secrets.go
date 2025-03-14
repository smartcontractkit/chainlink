package models

import (
	"github.com/smartcontractkit/chainlink-common/pkg/config"
)

// Secret is a string that formats and encodes redacted, as "xxxxx".
// Deprecated
type Secret = config.SecretString

// Deprecated
func NewSecret(s string) *Secret { return config.NewSecretString(s) }

// SecretURL is a URL that formats and encodes redacted, as "xxxxx".
// Deprecated
type SecretURL = config.SecretURL

// Deprecated
func NewSecretURL(u *config.URL) *config.SecretURL { return (*config.SecretURL)(u) }

// Deprecated
func MustSecretURL(u string) *config.SecretURL {
	return NewSecretURL(config.MustParseURL(u))
}
