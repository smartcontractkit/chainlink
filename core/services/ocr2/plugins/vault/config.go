package vault

import (
	"errors"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	vaultcap "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
)

type DKGConfig struct {
	ContractID string `json:"dkgContractID"`
}

type Config struct {
	RequestExpiryDuration commonconfig.Duration `json:"requestExpiryDuration"`
	DKG                   *DKGConfig            `json:"dkg,omitempty"`
	Auth0                 *vaultcap.Auth0Config `json:"auth0,omitempty"`
}

func (c *Config) Validate() error {
	if c.RequestExpiryDuration.Duration() <= 0 {
		return errors.New("request expiry duration cannot be 0")
	}
	return nil
}
