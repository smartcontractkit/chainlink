package config

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
)

// PluginConfig contains custom arguments for the DKG plugin.
type PluginConfig struct {
	EncryptionPublicKey string `json:"encryptionPublicKey"`
	SigningPublicKey    string `json:"signingPublicKey"`
	KeyID               string `json:"keyID"`
}

// ValidatePluginConfig validates that the given DKG plugin configuration is correct.
func ValidatePluginConfig(config PluginConfig, dkgSignKs keystore.DKGSign, dkgEncryptKs keystore.DKGEncrypt) error {
	_, err := dkgEncryptKs.Get(config.EncryptionPublicKey)
	if err != nil {
		return errors.Wrapf(err, "DKG encryption key: %s not found in key store", config.EncryptionPublicKey)
	}
	_, err = dkgSignKs.Get(config.SigningPublicKey)
	if err != nil {
		return errors.Wrapf(err, "DKG sign key: %s not found in key store", config.SigningPublicKey)
	}

	return nil
}
