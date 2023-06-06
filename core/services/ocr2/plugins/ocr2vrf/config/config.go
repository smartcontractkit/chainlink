package config

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	dkgconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/dkg/config"
)

// PluginConfig contains custom arguments for the OCR2VRF plugin.
//
// The OCR2VRF plugin runs a DKG under the hood, so it will need both
// DKG and OCR2VRF configuration fields.
//
// The DKG contract address is provided in the plugin configuration,
// however the OCR2VRF contract address is provided in the OCR2 job spec
// under the 'contractID' key.
type PluginConfig struct {
	// DKG configuration fields.
	DKGEncryptionPublicKey string `json:"dkgEncryptionPublicKey"`
	DKGSigningPublicKey    string `json:"dkgSigningPublicKey"`
	DKGKeyID               string `json:"dkgKeyID"`
	DKGContractAddress     string `json:"dkgContractAddress"`

	// VRF configuration fields
	VRFCoordinatorAddress string `json:"vrfCoordinatorAddress"`
	LinkEthFeedAddress    string `json:"linkEthFeedAddress"`
}

// ValidatePluginConfig validates that the given OCR2VRF plugin configuration is correct.
func ValidatePluginConfig(config PluginConfig, dkgSignKs keystore.DKGSign, dkgEncryptKs keystore.DKGEncrypt) error {
	err := dkgconfig.ValidatePluginConfig(dkgconfig.PluginConfig{
		EncryptionPublicKey: config.DKGEncryptionPublicKey,
		SigningPublicKey:    config.DKGSigningPublicKey,
		KeyID:               config.DKGKeyID,
	}, dkgSignKs, dkgEncryptKs)
	if err != nil {
		return err
	}

	// NOTE: a better validation would be to call a method on the on-chain contract pointed to by this
	// address.
	if config.DKGContractAddress == "" {
		return errors.New("dkgContractAddress field must be provided")
	}

	if config.VRFCoordinatorAddress == "" {
		return errors.New("vrfCoordinatorAddress field must be provided")
	}

	// NOTE: similar to the above.
	if config.LinkEthFeedAddress == "" {
		return errors.New("linkEthFieldAddress field must be provided")
	}

	return nil
}
