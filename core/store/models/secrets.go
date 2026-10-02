package models

import (
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
)

// Secret is a string that formats and encodes redacted, as "xxxxx".
//
// Deprecated: use commonconfig.SecretString instead.
type Secret = commonconfig.SecretString

// Deprecated: use commonconfig.NewSecretString instead.
func NewSecret(s string) *Secret { return commonconfig.NewSecretString(s) }

// SecretURL is a URL that formats and encodes redacted, as "xxxxx".
//
// Deprecated: use commonconfig.SecretURL instead.
type SecretURL = commonconfig.SecretURL

// Deprecated: use commonconfig.NewSecretURL instead.
func NewSecretURL(u *commonconfig.URL) *commonconfig.SecretURL { return commonconfig.NewSecretURL(u) }

// Deprecated: use commonconfig.MustSecretURL instead.
func MustSecretURL(u string) *commonconfig.SecretURL {
	return NewSecretURL(commonconfig.MustParseURL(u))
}
