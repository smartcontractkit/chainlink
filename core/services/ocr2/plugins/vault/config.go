package vault

import (
	"encoding/hex"
	"errors"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
)

type DKGConfig struct {
	MasterPublicKey          string `json:"masterPublicKey"`
	EncryptedPrivateKeyShare string `json:"encryptedPrivateKeyShare"`
}

type Config struct {
	RequestExpiryDuration commonconfig.Duration `json:"requestExpiryDuration"`
	DKG                   *DKGConfig            `json:"dkg,omitempty"`
}

func (c *Config) Validate() error {
	if c.RequestExpiryDuration.Duration() <= 0 {
		return errors.New("request expiry duration cannot be 0")
	}
	if c.DKG == nil {
		return errors.New("DKG configuration is required")
	}
	if _, err := hex.DecodeString(c.DKG.MasterPublicKey); err != nil {
		return errors.New("invalid master public key in DKG configuration: not hex decodable")
	}
	if _, err := hex.DecodeString(c.DKG.EncryptedPrivateKeyShare); err != nil {
		return errors.New("invalid encrypted private key share in DKG configuration: not hex decodable")
	}
	return nil
}
