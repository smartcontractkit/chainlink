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
	if c.Auth0 != nil {
		if c.Auth0.IssuerURL == "" {
			return errors.New("auth0 issuerURL is required when auth0 is configured")
		}
		if c.Auth0.Audience == "" {
			return errors.New("auth0 audience is required when auth0 is configured")
		}
		if c.Auth0.TenantID == 0 {
			return errors.New("auth0 tenantID is required when auth0 is configured")
		}
	}
	return nil
}
